/*
 handle-logout.go - respond to requests for /Logout/

    Copyright 2011 Jeffrey C. Barrett <jcbarret@gmail.com>

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
	"appengine"
	"appengine/user"
)

func init(){
	http.HandleFunc("/Logout/", handleLogout)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)	
	lourl,err := user.LogoutURL(c, "/")
	if err != nil {
	   c.Errorf("handleLogout: error getting LogoutURL")
	}	
	c.Debugf("handleLogout: redirecting to logoutURL=%v", lourl)
	http.Redirect(w, r, lourl, http.StatusFound)
}
