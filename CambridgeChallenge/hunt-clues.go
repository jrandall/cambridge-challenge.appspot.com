/*
 hunt-clues.go - Clue functions

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
       "appengine/datastore"
       "os"
       "fmt"
       "regexp"
       "time"
       "strconv"
)

// implement json.Marshaler interface
func (clue *Clue) UnmarshalJSON(jsondata []byte) (err os.Error) {
     fmt.Printf("Clue.UnmarshalJSON called with jsondata=%v\n", jsondata)
     return // err
}

func (clue *Clue) Answerable() bool {
     switch clue.AnswerType {
     	    case "": return true
     	    case "regexp": return true
     }
     return false
}

func (clue *Clue) AnswerCorrect(answerAttempt string, user string, huntName string) (correct bool) {
     c.Debugf("Clue.AnswerCorrect called with answerAttempt=%v on clue [%v]\n", answerAttempt, clue)

     //FIXME: this is not where the default belongs!
     if clue.AnswerType == "" {
     	c.Infof("Defaulting AnswerType to regexp")
     	clue.AnswerType = "regexp"
     }

     switch clue.AnswerType {
     	    case "regexp": {
	    	   c.Debugf("AnswerCorrect: handling regexp")
		   // refuse to match against blank regexp or blank answer
		   if clue.Answer == "" || answerAttempt == "" {
		      c.Debugf("AnswerCorrect regexp ignoring blank answer &/or answerAttempt")
		      return
		   }
     	     	   match, err := regexp.MatchString(clue.Answer, answerAttempt)
	     	   if err != nil {
	     	      panic(fmt.Sprintf("error attempting to match answerAttempt %v against regexp %v: %v", answerAttempt, clue.Answer, err))
	     	   }
	     	   if match {
		      correct = true
	     	   } 
		   }
	    case "timeofdayafter": { // answer is false until the time specified
	    	 c.Debugf("AnswerCorrect: handling timeofdayafter")
	    	 todThreshold, err := time.Parse(time.RFC822Z, clue.Answer)
		 if err != nil {
		    // could not parse time
		    c.Errorf("could not parse time %v as RFC822Z", clue.Answer)
		 }
		 // todThreshold now set to UTC time of threshold
		 todThresholdUTCSeconds := todThreshold.Seconds()
		 // this is UTC seconds
		 currentTimeUTCSeconds := time.Seconds()
	    	 c.Debugf("AnswerCorrect: checking whether todThreshold %v is less than or equal to time %v", todThresholdUTCSeconds, currentTimeUTCSeconds)
		 if todThresholdUTCSeconds <= currentTimeUTCSeconds {
		    correct = true
		 }
	    }
	    case "timeofdaybefore": { // answer is true until the time specified
	    	 c.Debugf("AnswerCorrect: handling timeofdaybefore")
	    	 todThreshold, err := time.Parse(time.RFC822Z, clue.Answer)
		 if err != nil {
		    // could not parse time
		    c.Errorf("could not parse time %v as RFC822Z", clue.Answer)
		 }
		 // todThreshold now set to UTC time of threshold
		 todThresholdUTCSeconds := todThreshold.Seconds()
		 currentTimeUTCSeconds := time.Seconds()
	    	 c.Debugf("AnswerCorrect: checking whether todThreshold %v is greater than or equal to time %v", todThresholdUTCSeconds, currentTimeUTCSeconds)
		 if todThresholdUTCSeconds >= currentTimeUTCSeconds {
		    correct = true
		 }
	    }
	    case "statetimewithin": {
	    	 c.Debugf("AnswerCorrect: handling statetimewithin")
	    	 // parse "Answer"
		 fields := splitOnWhiteSpaceTwoRegexp.FindStringSubmatch(clue.Answer)
		 state := fields[1]  // statename
		 withinTime, err := strconv.Atoi64(fields[2]) // time (s)
		 if err != nil {
		    panic(fmt.Sprintf("error parsing statetimewithin answer second field (withinTime=%v) to int64: %v", fields[1], err))
		 }

		 var arriveTime int64

		 //get earliest entry into state
		 //FIXME this belongs in hunt-state.go?
       		 stateQuery := datastore.NewQuery(huntStateDatastore).Filter("User=", user).Filter("HuntName=", huntName).Filter("CurrentStateName=",state).Order("Date").Limit(1)
       		 states := make([]HuntState, 0, 1)
       		 if _, err := stateQuery.GetAll(c, &states); err != nil {
		    panic(fmt.Sprintf("error getting earliest transition to state %v: %v", state, err))
       		 }
       		 if len(states) == 1 {
       	  	    // found state, get arrival time
		    stateDate := states[0].Date
		    arriveTime = int64(stateDate/1000000)
		    c.Debugf("stateDate %v yields arriveTime %v", stateDate, arriveTime)
       		 }

		 currentTime := time.Seconds()

		 // compare arriveTime + withinTime to current time
		 c.Debugf("comparing arriveTime %v + withinTime %v >= currentTime %v", arriveTime, withinTime, currentTime)
	    	 if arriveTime + withinTime >= currentTime {
		    c.Debugf("CORRECT!")
		    correct = true
		 }
	    }
	    default: {
	    	 c.Warningf("AnswerCorrect: unhandled AnswerType %v", clue.AnswerType)	    
		 //FIXME should this fail silently?
	    }
     }
     if correct {
     	c.Debugf("AnswerCorrect: CORRECT!")
     }

     return // correct
}


var splitOnWhiteSpaceTwoRegexp = regexp.MustCompile("^(.*)[ \t]+(.*)$")
