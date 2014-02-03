package gosched

import (
  "fmt"
  "appengine"
  "appengine/datastore"
	"appengine/user"
  "encoding/json"
  "net/http"
)


func ActivityListHandler(w http.ResponseWriter, r *http.Request) {
  switch {
    case r.Method == "GET":
      ActivityListGet(w,r)
    default:
      fmt.Fprint(w, "Event handler.")
  }
}

func ActivityListGet(w http.ResponseWriter, r *http.Request) {
	var q *datastore.Query
	var act Activity_entity
	var owner string
	l := make([]Activity_entity,0)
	ds := appengine.NewContext(r)
	owner = r.FormValue("owner")
	// if owner is not specified, use currently logged in user 
	// or guest if not logged in 
	if owner == "" {
		u := user.Current(ds)
		if u != nil {
			owner = u.Email
		} else {
			owner = "Guest"
		}
	}
	// build query to return activities by an owner
	q = datastore.NewQuery("Activity_entity").Filter("Owner = ", owner)
	t := q.Run(ds)
	for t != nil {
		_,err := t.Next(&act)
    if err == datastore.Done {
			break
    }
		if err != nil {
			break
		}
		l = append(l,act)
	}
	jf, err := json.Marshal(l)
  if err != nil {
    fmt.Fprint(w, "{\"errror\":\"Error marshalling json\"}")
  } else {
    w.Write(jf)
  }

}
