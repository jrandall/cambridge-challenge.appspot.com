/*
 access-log.go - log all accesses

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
    "appengine/datastore"
    "time"
    "os"
    "http"
)

type Access struct {
    User          string
    Date          datastore.Time
    RemoteAddress string
    URL		  string
    Method	  string
}

func LogAccess(r *http.Request, User string) (accessLog Access, err os.Error) {
    accessLog = Access{
        Date:       datastore.SecondsToTime(time.Seconds()),
	RemoteAddress: r.RemoteAddr,
	Method: r.Method,
	URL: r.RawURL,
	User: User,
    }
    c := appengine.NewContext(r)

    _, err = datastore.Put(c, datastore.NewIncompleteKey("Access"), &accessLog)
    return // Access, err
}

