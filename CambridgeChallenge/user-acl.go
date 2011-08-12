/*
 user-acl.go - functions to get user and check access

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
    "appengine/user"
    "http"
)

func requireAnyUser(w http.ResponseWriter, r *http.Request) (User string) {
    c := appengine.NewContext(r)
    u := user.Current(c); 
    if u != nil {
        // valid user logged in
	User = u.String()
    } else {
       // user not logged in, redirect to login page
       url, err := user.LoginURL(c, r.URL.String())
       if err != nil {
       	 http.Error(w, err.String(), http.StatusInternalServerError)
       }
       w.Header().Set("Location", url)
       w.WriteHeader(http.StatusFound)
    }
    return // User
}
