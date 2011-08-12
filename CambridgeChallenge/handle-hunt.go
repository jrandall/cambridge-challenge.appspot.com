/*
 handle-hunt.go - respond to requests for /Hunt/

    Copyright 2011 Joshua C. Randall <jcrandall@alum.mit.edu>

    This file is part of CambridgeChallenge.

    CambridgeChallenge is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    CambridgeChallenge is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with CambridgeChallenge.  If not, see <http://www.gnu.org/licenses/>.
*/
package CambridgeChallenge

import (
    "appengine"
    "http"
    "template"
    "os"
    "strings"
    "time"
    "regexp"
    "fmt"
    "appengine/blobstore"
    "appengine/datastore"
	"net"
)

const (
      huntPath = "/Hunt/"
      huntTemplateFileName = "template/hunt.html.gotmpl"
      huntStateDatastore = "HuntState"
)

// data that gets processed by the html template
type HuntTemplateData struct {
     User string
     Error string
     SuppressAnswerBox bool
     *HuntDirectoryEntry
     DebugHuntData *Hunt
     CurrentState *State
}

type HuntState struct {
     User   string
     HuntName   string
     CurrentStateName  string
     FromState string
     Date   datastore.Time
}

func init() {
    http.HandleFunc(huntPath, handleHunt)
}

var huntTemplate = template.MustParseFile(huntTemplateFileName, template.FormatterMap{"dstime" : dstimeFormatter})

func StateTransition(c appengine.Context, user string, huntName string, toState string, fromState string) (err os.Error) {
     c.Debugf("StateTransition called to %v from %v", toState, fromState)
	     hs := HuntState {
     	     	User: user,
     		HuntName: huntName,
     		CurrentStateName: toState,
     		FromState: fromState,
     		Date: datastore.SecondsToTime(time.Seconds()),
    		}
    	     _, err = datastore.Put(c, datastore.NewIncompleteKey(huntStateDatastore), &hs)
	     return // err
}

func handleHunt(w http.ResponseWriter, r *http.Request) {
    var td HuntTemplateData
    var err os.Error
    defer recoverUserError(w, r)
    c := appengine.NewContext(r)

    td.User = requireAnyUser(w, r)
    LogAccess(r, td.User)

    // get hunt name from URL path
    huntSearchName := strings.Split(strings.Replace(r.URL.Path, huntPath, "", 1), "/", 2)[0]


    // look up hunt in directory
    huntQuery := datastore.NewQuery(huntDirectoryDatastore).Filter("HuntName=", huntSearchName).Limit(1)
    huntentries := make([]HuntDirectoryEntry, 0, 1)
    if _, err := huntQuery.GetAll(c, &huntentries); err != nil {
       serveError(c, w, err)
       return
    }
    if len(huntentries) != 1 {
       panic("Hunt not found")
    }   
       // hunt found, load hunt data from blobstore
       td.HuntDirectoryEntry = &huntentries[0] // sets all HuntDirectoryEntry fields without copying data

       huntData, err := DecodeHuntData(blobstore.NewReader(c, td.BlobKey))
       if err != nil {
       	  serveError(c, w, err)
       }
       // now have huntData
       if appengine.IsDevAppServer() {
       	  td.DebugHuntData = huntData // exposes all huntData fields in td
       }

       // get and/or initialize current state
       var currentHuntState HuntState
       stateQuery := datastore.NewQuery(huntStateDatastore).Filter("User=", td.User).Filter("HuntName=", td.HuntName).Order("-Date").Limit(1)
       states := make([]HuntState, 0, 1)
       if _, err := stateQuery.GetAll(c, &states); err != nil {
       	  serveError(c, w, err)
       	  return
       }
       if len(states) == 1 {
          currentHuntState = states[0]
       	  currentStateName := currentHuntState.CurrentStateName
       	  td.CurrentState = huntData.States[currentStateName]
       } else {
       	  // state doesn't exist, set initial state now
	  if huntData.EnterState == "" {
	     // hunt didn't specify initial state
	     panic("Hunt did not specify EnterState, cannot initialize")
	  } else {
	     // get EnterState
	     td.CurrentState = huntData.States[huntData.EnterState]
	     // add Enter->CurrentState transition to huntStateDatastore
	     err = StateTransition(c, td.User, td.HuntName, td.CurrentState.StateName, "START")
	     if err != nil {
	     	panic("error setting initial state")
	     }
	  }
       }

      // get answer submission, if any      
       err = r.ParseForm()
       if err != nil {
       	  serveError(c, w, err)
      	  return
      }
      answerAttempt := r.FormValue("Answer")

      // TODO: state points (for multi-clue states) 
      var correct = false
      var cluesHaveAnswers = false
	 // check answer against all clues, also checking whether any clues in this state contain an answer
	 var clues []Clue = td.CurrentState.Clues
	 for _, clue := range clues {
	     c.Debugf("checking answer, have clue with answer: %v", clue.Answer)
	     if clue.Answer != "" {
	       cluesHaveAnswers = true 
     	     match, err := regexp.MatchString(clue.Answer, answerAttempt)
	     if err != nil {
	     	panic(fmt.Sprintf("error attempting to match answer %v against regexp %v: %v", answerAttempt, clue.Answer, err))
	     }
	     if match {
	        c.Infof("CORRECT ANSWER!")
		correct = true
	     } 
	     }
	 }

      if correct { // would be based on points instead
      	 // advance to NextState and redirect to self
	     err = StateTransition(c, td.User, td.HuntName, td.CurrentState.NextState, td.CurrentState.StateName)
	     if err != nil {
	     	panic(fmt.Sprintf("error advancing from State %v to NextState %v: %v", td.CurrentState.StateName, td.CurrentState.NextState, err))
	     }
	 // redirect
	 http.Redirect(w, r, huntPath+"/"+td.HuntName, http.StatusFound)
   
	     return
      }

      if !cluesHaveAnswers {
      	 td.SuppressAnswerBox = true
	 // if these clues don't have answers, we may have a "Forward" / "Back" button submission, check for that and act on it
	 if r.FormValue("Navigate") == "Forward" {
		//check for an IP address requirement
		criterion := td.CurrentState.AllowNetMask
		if (criterion != ""){
			maskchunks := strings.Split(criterion,"/",-1)
			mask := net.ParseIP(maskchunks[0])

			ra:= net.ParseIP(r.RemoteAddr)
			//these are the correct bytes from the IPv6 storage...I tried *desperately*
			// to get it to make a mask from an IP object but it seems impossible!
			maskedIP := ra.Mask(net.IPv4Mask(mask[12],mask[13],mask[14],mask[15]))
			if (maskedIP.Equal(net.ParseIP(maskchunks[1]))){
				err = StateTransition(c, td.User, td.HuntName, td.CurrentState.NextState, td.CurrentState.StateName)
					     if err != nil {
					     	panic(fmt.Sprintf("error advancing from State %v to NextState %v: %v", td.CurrentState.StateName, td.CurrentState.NextState, err))
					     }
					 // redirect
					 http.Redirect(w, r, huntPath+"/"+td.HuntName, http.StatusFound)
			}
		}else{
err = StateTransition(c, td.User, td.HuntName, td.CurrentState.NextState, td.CurrentState.StateName)
	     if err != nil {
	     	panic(fmt.Sprintf("error advancing from State %v to NextState %v: %v", td.CurrentState.StateName, td.CurrentState.NextState, err))
	     }
	 // redirect
	 http.Redirect(w, r, huntPath+"/"+td.HuntName, http.StatusFound)	    
	}
	 } 
	 if r.FormValue("Navigate") == "Back" {
	    err = StateTransition(c, td.User, td.HuntName, currentHuntState.FromState, td.CurrentState.StateName)
	     if err != nil {
	     	panic(fmt.Sprintf("error advancing from State %v to NextState %v: %v", td.CurrentState.StateName, td.CurrentState.NextState, err))
	     }
	 // redirect
	 http.Redirect(w, r, huntPath+"/"+td.HuntName, http.StatusFound)
	 }
 }


    w.Header().Set("Content-Type", "text/html")
    c.Debugf("calling huntTemplate.Execute on td: %v", td)
    if err = huntTemplate.Execute(w, td); err != nil {
        serveError(c, w, err)
    }
    return
}


