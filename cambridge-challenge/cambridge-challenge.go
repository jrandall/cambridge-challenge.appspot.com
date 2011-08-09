package hello

import (
    "appengine"
    "appengine/datastore"
    "appengine/user"
    "http"
    "template"
    "time"
    "io"
)

type Log struct {
    User          string
    Date          datastore.Time
    RemoteAddress string
    Content 	  string
}

func init() {
    http.HandleFunc("/", root)
    http.HandleFunc("/sign", sign)
}

func root(w http.ResponseWriter, r *http.Request) {
    requireAnyUser(w, r)
    c := appengine.NewContext(r)
    q := datastore.NewQuery("Log").Order("-Date").Limit(10)
    greetings := make([]Log, 0, 10)
    if _, err := q.GetAll(c, &greetings); err != nil {
        http.Error(w, err.String(), http.StatusInternalServerError)
        return
    }
    if err := guestbookTemplate.Execute(w, greetings); err != nil {
        http.Error(w, err.String(), http.StatusInternalServerError)
    }
}

func dstimeFormatter(wr io.Writer, formatter string, dstime ...interface{}) {
  io.WriteString(wr, time.SecondsToUTC(int64(dstime[0].(datastore.Time))/1000000).Format(time.RFC1123));
}

var guestbookTemplate = template.MustParse(guestbookTemplateHTML, template.FormatterMap{"dstime" : dstimeFormatter})

const guestbookTemplateHTML = `
<html>
  <body>
   <ul>
    {.repeated section @}
     <li>
      {.section User}
        <b>{@|html}</b> accessed at: 
      {.or}
        An anonymous person accessed at: 
      {.end}
      {Date|dstime|html} from {RemoteAddress|html}
    {.end}
   </ul>
    <form action="/sign" method="post">
      <div><textarea name="content" rows="3" cols="60"></textarea></div>
      <div><input type="submit" value="Submit"></div>
    </form>
  </body>
</html>
`

func requireAnyUser(w http.ResponseWriter, r *http.Request) (userString string) {
    c := appengine.NewContext(r)
    u := user.Current(c); 
    if u != nil {
        // valid user logged in
        return(u.String())
    } else {
       // user not logged in, redirect to login page
       url, err := user.LoginURL(c, r.URL.String())
       if err != nil {
       	 http.Error(w, err.String(), http.StatusInternalServerError)
	 return ""
       }
       w.Header().Set("Location", url)
       w.WriteHeader(http.StatusFound)
       return ""
    }
    return ""
}

func sign(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    g := Log{
        Content:    r.FormValue("content"),
        Date:       datastore.SecondsToTime(time.Seconds()),
	RemoteAddress: r.RemoteAddr,
    }
    g.User = requireAnyUser(w, r)


    _, err := datastore.Put(c, datastore.NewIncompleteKey("Log"), &g)
    if err != nil {
        http.Error(w, err.String(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/", http.StatusFound)
}

