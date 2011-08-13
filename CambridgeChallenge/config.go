/*
 config.go - configuration constants

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
)

const (
      huntPath = "/Hunt/"
      huntAdminPath = "/HuntAdmin/"
      huntAdminUploadPath = "/HuntAdmin/Upload/"
      huntAdminDownloadPath = "/HuntAdmin/Download/"

      huntStateDatastore = "HuntState"
      huntDirectoryDatastore = "HuntDirectory"

      huntLimit = 100 // number of hunts to load in the admin interface

      huntTemplateFileName = "template/hunt.html.gotmpl"
      huntAdminTemplateFileName = "template/huntadmin.html.gotmpl"
      rootTemplateFileName = "template/root.html.gotmpl"
      loginTemplateFileName = "template/login.html.gotmpl"
      
      useOpenID = true
)



/*
var openIdProviders []string = []string{
      		      "Google.com/accounts/o8/id",
		      "Yahoo.com",
    		      "MySpace.com",
    		      "AOL.com",
    		      "MyOpenID.com",
      }
*/

// package global appengine Context 
var c appengine.Context

