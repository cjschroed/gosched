package gosched

import (
    "fmt"
		"html/template"
		"appengine"
		"appengine/user"
    "net/http"
)

const google_openid_url = "https://www.google.com/accounts/o8/id"
const yahoo_openid_url = "http://open.login.yahooapis.com/openid20/www.yahoo.com/xrds"
type logindata struct {
  GoogleLogin string
  YahooLogin string
}


// main event handler 
func init() {
	http.HandleFunc("/", LoginPage)
	http.HandleFunc("/gosched/v1/hello", HelloHandler)
	http.HandleFunc("/gosched/v1/activity", ActivityHandler)
	http.HandleFunc("/gosched/v1/activity/book", BookingHandler)
	http.HandleFunc("/gosched/v1/activity/list", ActivityListHandler)
	http.HandleFunc("/gosched/v1/activity/events", EventsHandler)
	http.HandleFunc("/gosched/v1/activity/events/search", ListEventByAvailability)
	http.HandleFunc("/gosched/v1/activity/events/list", EventListGet)
	http.HandleFunc("/gosched/v1/activity/clear", ClearUserData)
	http.HandleFunc("/gosched/unittests", UnitTestSection)
}

// I'm awake! I'm awake.
func HelloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hello, world!")
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
  t,_ := template.New("foo").Delims("[[", "]]").Parse(LoginPageTemplate)
  ds := appengine.NewContext(r)
  var m logindata
  u := user.Current(ds)
  if u != nil {
		http.Redirect(w, r, "/gosched/unittests", http.StatusFound)
		return
	}

  m.GoogleLogin,_ = user.LoginURLFederated(ds, "/gosched/v1/unittests", google_openid_url)
  m.YahooLogin,_ = user.LoginURLFederated(ds, "/gosched/v1/unittests", yahoo_openid_url)
  t.Execute(w, m)
}

