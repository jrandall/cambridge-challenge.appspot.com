/*
 user-error.go - handle user errors

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
    "os"
    "fmt"
)

const (
      userErrorTemplateFileName = "template/error.html.gotmpl"
)



var userErrorTemplate = template.MustParseFile(userErrorTemplateFileName, template.FormatterMap{"dstime" : dstimeFormatter})

func recoverUserError(w http.ResponseWriter, r *http.Request) {
  if rec := recover(); rec != nil {
     var err os.Error
     c.Infof("Recovering from panic in recoverUserError: %v", rec)     

     w.Header().Set("Content-Type", "text/html")
     etd := map[string]string{
     	 "Error": fmt.Sprintf("%v",rec),
     }

     if err = userErrorTemplate.Execute(w, etd); err != nil {
        serveError(c, w, err)
     }
  }
}
