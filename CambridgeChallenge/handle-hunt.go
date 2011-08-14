/*
 handle-hunt.go - respond to requests for /Hunt/

    Copyright 2011 Joshua C. Randall <jcrandall@alum.mit.edu>

    This file is part of CambridgeChallenge.

    CambridgeChallenge is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    CambridgeChallenge is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

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
    "fmt"
    "appengine/blobstore"
    "appengine/datastore"
	"net"
)

// data that gets processed by the html template
type HuntTemplateData struct {
     User string
     Error string
     SuppressAnswerBox bool
     SuppressBackButton bool
     *HuntDirectoryEntry
     DebugHuntData *Hunt
     CurrentState *State
}

func init() {
    http.HandleFunc(huntPath, handleHunt)
}

var huntTemplate = template.MustParseFile(huntTemplateFileName, template.FormatterMap{"dstime" : dstimeFormatter})

func handleHunt(w http.ResponseWriter, r *http.Request) {
    var td HuntTemplateData
    var err os.Error
    defer recoverUserError(w, r)
    c = appengine.NewContext(r)

    td.User = requireAnyUser(w, r)
    LogAccess(r, td.User)
    if td.User == "" {
       panic("requireAnyUser did not return a username")
    }

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

       c.Debugf("handleHunt: calling DecodeHuntData on blobkey %v", td.BlobKey)
       huntData, err := DecodeHuntData(blobstore.NewReader(c, td.BlobKey))
       if err != nil {
       	  serveError(c, w, err)
       }
       // now have huntData
       if appengine.IsDevAppServer() {
       	  td.DebugHuntData = huntData // exposes all huntData fields in td
       }

       // get and/or initialize current state
       var currentHuntState *HuntState
       currentHuntState, err = GetCurrentHuntState(td.User, td.HuntName)
       if err != nil {
       	  serveError(c, w, err)
       }
       if currentHuntState == nil {
       	  // set initial state (add transition from "START" to the hunt EnterStaet
          if huntData.EnterState == "" {
             // hunt didn't specify initial state
             panic("Hunt did not specify EnterState, cannot initialize hunt")
          } 
	  err = StateTransition(td.User, td.HuntName, huntData.EnterState, "")
	  if err != nil {
	     panic("error setting initial state")
	  }
	  c.Debugf("handleHunt: StateTransition complete, getting current hunt state")
       	  currentHuntState, err = GetCurrentHuntState(td.User, td.HuntName)
       	  if err != nil {
       	     serveError(c, w, err)
       	  }
	  if currentHuntState == nil {
	     panic("currentHuntState nil after setting initial state")
	  }
       }
       c.Debugf("handleHunt: Have currentHuntState %v", currentHuntState)
       currentStateName := currentHuntState.CurrentStateName
       td.CurrentState = huntData.States[currentStateName]

       if currentHuntState.FromState == "" {
       	  // suppress back button when previous state is not set
      	 td.SuppressBackButton = true
	}

      // get answer submission, if any      
       err = r.ParseForm()
       if err != nil {
       	  serveError(c, w, err)
      	  return
      }
      answerAttempt := r.FormValue("Answer")

      var correct = false
      var cluesHaveAnswers = false
      var allCluesCorrect = true
	 // check answer against all clues with answers, also noting whether any clues in this state contain an answerable answer (otherwise we are in next/previous state)
	 var clues []Clue = td.CurrentState.Clues
	 for _, clue := range clues {
	     // have a clue
	     c.Debugf("handleHunt: have clue with answer: %v and answertype: %v", clue.Answer, clue.AnswerType)

	     if clue.Answerable() {
	     	// this clues answer requires a form to answer
	       	cluesHaveAnswers = true 
	     }
	     
	     correct = clue.AnswerCorrect(answerAttempt, td.User, td.HuntName)

	     if !correct {
	     	allCluesCorrect = false
	     	// check if there is an IncorrectAnswerState and transition to it
		if clue.IncorrectAnswerState != "" {
	     	   err = StateTransition(td.User, td.HuntName, clue.IncorrectAnswerState, td.CurrentState.StateName)
	     	   if err != nil {
	     	      panic(fmt.Sprintf("error advancing from State %v to IncorrectAnswerState %v: %v", td.CurrentState.StateName, clue.IncorrectAnswerState, err))
	     	   }
		   // redirect to self to get updated state
		   // FIXME: this will miss answers just submitted now for later clues!
	    	    http.Redirect(w, r, huntPath+"/"+td.HuntName, http.StatusFound)
	     	    return
		}
	     }

	     if correct {
	     	// check if there is a CorrectAnswerState and transition to it immediately
		if clue.CorrectAnswerState != "" {
	     	   err = StateTransition(td.User, td.HuntName, clue.CorrectAnswerState, td.CurrentState.StateName)
	     	   if err != nil {
	     	      panic(fmt.Sprintf("error advancing from State %v to CorrectAnswerState %v: %v", td.CurrentState.StateName, clue.CorrectAnswerState, err))
	     	   }
		   // redirect to self now that we have transitioned
	    	    http.Redirect(w, r, huntPath+"/"+td.HuntName, http.StatusFound)
	     	    return
		}
	     }
	 }

      if cluesHaveAnswers && allCluesCorrect { // could be based on points instead
      	 // advance to NextState and redirect to self
	     err = StateTransition(td.User, td.HuntName, td.CurrentState.NextState, td.CurrentState.StateName)
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
				err = StateTransition(td.User, td.HuntName, td.CurrentState.NextState, td.CurrentState.StateName)
					     if err != nil {
					     	panic(fmt.Sprintf("error advancing from State %v to NextState %v: %v", td.CurrentState.StateName, td.CurrentState.NextState, err))
					     }
					 // redirect
					 http.Redirect(w, r, huntPath+"/"+td.HuntName, http.StatusFound)
			}
		}else{
err = StateTransition(td.User, td.HuntName, td.CurrentState.NextState, td.CurrentState.StateName)
	     if err != nil {
	     	panic(fmt.Sprintf("error advancing from State %v to NextState %v: %v", td.CurrentState.StateName, td.CurrentState.NextState, err))
	     }
	 // redirect
	 http.Redirect(w, r, huntPath+"/"+td.HuntName, http.StatusFound)	    
	}
	 } 
	 if r.FormValue("Navigate") == "Back" {
	    err = StateTransition(td.User, td.HuntName, currentHuntState.FromState, td.CurrentState.StateName)
	     if err != nil {
	     	panic(fmt.Sprintf("error advancing from State %v to NextState %v: %v", td.CurrentState.StateName, td.CurrentState.NextState, err))
	     }
	 // redirect
	 http.Redirect(w, r, huntPath+"/"+td.HuntName, http.StatusFound)
	 }
 }


    w.Header().Set("Content-Type", "text/html")
    c.Debugf("handleHunt: calling huntTemplate.Execute on td: %v", td)
    if err = huntTemplate.Execute(w, td); err != nil {
        serveError(c, w, err)
    }
    return
}


