/*
 handle-login.go - respond to login_required to dispatch to openid provider

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
	"appengine"
	"appengine/user"
	"os"
	"template"
	"strings"
)

func init(){
	http.HandleFunc("/_ah/login_required", handleLogin)
	http.HandleFunc("/Login/", handleLogin)
}

type LoginTemplateData struct {
     User string
     ContinueURL string
}

var loginTemplate = template.MustParseFile(loginTemplateFileName, template.FormatterMap{"dstime" : dstimeFormatter})

func handleLogin(w http.ResponseWriter, r *http.Request) {
    	c = appengine.NewContext(r)	
    	u := user.Current(c)
	var identity string
	var err os.Error
	returnUrl := "/"

	// parse form	
       err = r.ParseForm()
       if err != nil {
       	  serveError(c, w, err)
      	  return
      }
      if r.FormValue("continue") != "" {
      	 returnUrl = r.FormValue("continue")
      }
      if r.FormValue("chooseLogin") == "1" {
         // display form instead of redirecting to OpenID
      	 c.Debugf("handleLogin: have request to display login form")
    	 var td LoginTemplateData
	 if u != nil {
	    td.User = u.String()
	 }
	 td.ContinueURL = returnUrl
    	 if err := loginTemplate.Execute(w, td); err != nil {
            serveError(c, w, err)
    	 }
	 return
      }
      if r.FormValue("identity") != "" {
      	 identity = r.FormValue("identity")
      }

      // only allow MyOpenID.com for now
      identity = strings.Split(identity, ".", 2)[0]
      if identity != "" {
      	 identity += "."
      }
      identity += "MyOpenID.com"
      c.Debugf("getting LoginURLFederated for %v %v", returnUrl, identity)
      loginUrl,err := user.LoginURLFederated(c, returnUrl, identity)
      if err != nil {
	   c.Errorf("handleLogin: error getting LoginURL for %v %v", returnUrl, identity)
      }	
      c.Debugf("handleLogin: redirecting to loginURL=%v", loginUrl)
      http.Redirect(w, r, loginUrl, http.StatusFound)
}
