/*
 user-acl.go - functions to get user and check access

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
    "appengine/user"
    "http"
    "strings"
)

func requireAnyUser(w http.ResponseWriter, r *http.Request) (User string) {
    c := appengine.NewContext(r)
    u := user.Current(c); 
    if u == nil {
       // user not logged in 
       if useOpenID {
              // redirect to our login page
	      returnUrl := "/Login/"+"?chooseLogin=1&continue="+http.URLEscape(r.URL.RawPath)
	      c.Debugf("handleLogin: redirecting to %v", returnUrl)
	      http.Redirect(w, r, returnUrl, http.StatusFound)
       	      return
       } else {
       	      // redirect to google login page
	      returnUrl := r.URL.RawPath
	      loginUrl, err := user.LoginURL(c, returnUrl)
	      if err != nil {
	      	 c.Errorf("handleLogin: error getting LoginURL for %v %v", returnUrl)
	      }	
	      c.Debugf("handleLogin: redirecting to loginURL=%v", loginUrl)
	      http.Redirect(w, r, loginUrl, http.StatusFound)
       }
    }
    // valid user logged in
    User = u.String()
    url, err := http.ParseURL(User)
    if err != nil {
       // error parsing URL, redirect to Login?
       c.Errorf("error parsing User as URL: %v", User)
       w.Header().Set("Location", "/Login/")
       w.WriteHeader(http.StatusFound)
    }
    User = strings.Split(url.Host, ".",2)[0]
    return // User
}

func requireAdminUser(w http.ResponseWriter, r *http.Request) (User string) {
     c := appengine.NewContext(r)
     User = requireAnyUser(w, r)
     if appengine.IsDevAppServer() { // dev app server always admin
     	return // User
     }
     // TODO this should absolutely not be a hardcoded userid
     if User == "ccjava" {
     	// the admin user!
	return // User
     }
     // not admin user
     c.Debugf("requireAdminUser: %v is not admin, redirecting to %v", User, "/Login/?chooseLogin=1&continue="+http.URLEscape(r.URL.RawPath))
     w.Header().Set("Location", "/Login/?chooseLogin=1&continue="+http.URLEscape(r.URL.RawPath))
     w.WriteHeader(http.StatusFound)
     return // User
}
