/*
 handle-logout.go - respond to requests for /Logout/

    Copyright 2011 Jeffrey C. Barrett <jcbarret@gmail.com>

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
    "appengine"
    "appengine/user"
    "os"
)

func init(){
	http.HandleFunc("/Logout/", handleLogout)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)	
	returnURL := "/"
	// parse form	
       	err := r.ParseForm()
       	if err != nil {
       	  serveError(c, w, err)
      	  return
	}
      	if r.FormValue("continue") != "" {
      	   returnURL = r.FormValue("continue")
      	}

	if useOpenID {
	  // adjust returnURL to bring us back to a local user login form
	  laterReturnUrl := returnURL
	  returnURL = "/Login/?chooseLogin=1&continue="+http.URLEscape(laterReturnUrl)
	}
	// redirect to google logout (for OpenID as well, or else we won't be locally logged out)
	lourl,err := user.LogoutURL(c, returnURL)
	if err != nil {
	     c.Errorf("handleLogout: error getting LogoutURL")
	}	
	c.Debugf("handleLogout: redirecting to logoutURL=%v", lourl)
	http.Redirect(w, r, lourl, http.StatusFound)
	return
}

func getLogoutURL(c appengine.Context, returnURL string) (url string, err os.Error) {
    // set logout URL
    if useOpenID {
       // use our logout page
       url = "/Logout/"+"?continue="+returnURL
    } else {
       // use google logout
       url, err = user.LogoutURL(c, returnURL)
    }
    return // url, err
}
