/*
 template-formatters.go - custom template formatters 

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
       "io"
       "time"
       "appengine/datastore"
)

func dstimeFormatter(wr io.Writer, formatter string, dstime ...interface{}) {
  io.WriteString(wr, time.SecondsToUTC(int64(dstime[0].(datastore.Time))/1000000).Format(time.RFC1123));
}

