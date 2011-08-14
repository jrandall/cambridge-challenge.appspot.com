/*
 hunt-state.go - hunt state transitions and queries

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
    "os"
    "time"
    "appengine/datastore"
)

type HuntState struct {
     User   string
     HuntName   string
     CurrentStateName  string
     FromState string
     Date   datastore.Time
}

func StateTransition(user string, huntName string, toState string, fromState string) (err os.Error) {
     c.Debugf("StateTransition called for user %v hunt %v to state %v from state %v", user, huntName, toState, fromState)
	     hs := HuntState {
     	     	User: user,
     		HuntName: huntName,
     		CurrentStateName: toState,
     		FromState: fromState,
     		Date: datastore.SecondsToTime(time.Seconds()),
    		}
		c.Debugf("StateTransition: calling datastore.Put to %v with hs: [%v]", huntStateDatastore, hs)
    	     _, err = datastore.Put(c, datastore.NewIncompleteKey(huntStateDatastore), &hs)
	     if err != nil {
	     	c.Errorf("StateTransition: error %v", err)
	     }
	     return // err
}

func GetCurrentHuntState(user string, huntName string) (currentHuntState *HuntState, err os.Error) {
     c.Debugf("GetCurrentHuntState(%v, %v)", user, huntName)
       stateQuery := datastore.NewQuery(huntStateDatastore).Filter("User=", user).Filter("HuntName=", huntName).Order("-Date").Limit(1)
       states := make([]HuntState, 0, 1)
     c.Debugf("GetCurrentHuntState: executing datastore query")
       if _, err := stateQuery.GetAll(c, &states); err != nil {
       	  return // currentHuntState, err
       }
       if len(states) == 1 {
       	  // found current state, set return value to it!
          c.Debugf("GetCurrentHuntState: found current state")
          currentHuntState = &states[0]
       }
       c.Debugf("GetCurrentHuntState: returning currentHuntState %v", currentHuntState)
       return // currentHuntState, err
}

