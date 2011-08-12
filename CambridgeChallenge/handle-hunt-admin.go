/*
 huntadmin.go - respond to requests for /huntadmin, an administrative 
 	      interface for adding, deleting, and modifying hunt data.

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

func init() {
    http.HandleFunc("/huntadmin", huntadmin)
}

func huntadmin(w http.ResponseWriter, r *http.Request) {
    accesslog := Access{
        Date:       datastore.SecondsToTime(time.Seconds()),
	RemoteAddress: r.RemoteAddr,
	Method: r.Method,
	URL: r.RawURL,
    }
    accesslog.User = requireAnyUser(w, r)
/*
    c := appengine.NewContext(r)
    q := datastore.NewQuery("Access").Order("-Date").Limit(10)
    greetings := make([]Access, 0, 10)
    if _, err := q.GetAll(c, &greetings); err != nil {
        http.Error(w, err.String(), http.StatusInternalServerError)
        return
    }
    if err := guestbookTemplate.Execute(w, greetings); err != nil {
        http.Error(w, err.String(), http.StatusInternalServerError)
    }
*/
    if err := rootTemplate.Execute(w, accesslog); err != nil {
        http.Error(w, err.String(), http.StatusInternalServerError)
    }
}
