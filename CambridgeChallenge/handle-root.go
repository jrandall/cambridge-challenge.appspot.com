/*
 handle-root.go - respond to requests for /

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
    "http"
    "template"
    "appengine"
    "os"
)

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
    rd.LogoutURL, err = getLogoutURL(c, "/")
    if err != nil {
       c.Errorf("could not get LogoutURL: %v", err)
       serveError(c, w, err)
       return
    }   
    
    if err := rootTemplate.Execute(w, rd); err != nil {
            serveError(c, w, err)
    }
}

var rootTemplate = template.MustParseFile(rootTemplateFileName, template.FormatterMap{"dstime" : dstimeFormatter})



