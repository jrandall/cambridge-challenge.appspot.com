/*
 handle-root.go - respond to requests for /

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
    "http"
    "template"
    "appengine"
    "appengine/user"
    "os"
)

const (
      rootTemplateFileName = "template/root.html.gotmpl"
)


/*
type StateTransitionLog struct {
    User          string
    Date          datastore.Time
    FromState	  string
    ToState	  string
}

type CurrentState struct {
    User          string
    State	  string
}
*/

func init() {
    http.HandleFunc("/", handleRoot)
}

type RootTemplateData struct {
     User string
     LogoutURL string
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
    var rd RootTemplateData
    rd.User = requireAnyUser(w, r)
    LogAccess(r, rd.User)
    var err os.Error
    c := appengine.NewContext(r)
    rd.LogoutURL, err = user.LogoutURL(c, huntAdminPath)
    if err != nil {
       c.Errorf("could not get LogoutURL: %v", err)
       serveError(c, w, err)
       return
    }   
    
/*
    q := datastore.NewQuery("Access").Order("-Date").Limit(10)
    greetings := make([]Access, 0, 10)
    if _, err := q.GetAll(c, &greetings); err != nil {
        http.Error(w, err.String(), http.StatusInternalServerError)
        return
    }
    if err := guestbookTemplate.Execute(w, greetings); err != nil {
        http.Error(w, err.String(), http.StatusInternalServerError)
    }
*/
    if err := rootTemplate.Execute(w, rd); err != nil {
        http.Error(w, err.String(), http.StatusInternalServerError)
    }
}

var rootTemplate = template.MustParseFile(rootTemplateFileName, template.FormatterMap{"dstime" : dstimeFormatter})



