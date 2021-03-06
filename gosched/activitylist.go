package gosched

import (
  "fmt"
  "appengine"
  "appengine/datastore"
	"appengine/user"
  "encoding/json"
	"strconv"
  "net/http"
)


func ActivityListHandler(w http.ResponseWriter, r *http.Request) {
	ds := appengine.NewContext(r)
  u := user.Current(ds)
  if u == nil {
    fmt.Fprint(w, "{\"errror\":\"User credentials could not be determined.\"}")
    return
	}
  switch {
    case r.Method == "GET":
      ActivityListGet(w,r,ds,u)
    default:
      fmt.Fprint(w, "Event handler.")
  }
}

// returns all activities correpsonding to the logged in user
//  or one specified by the 'owner' form parameter
func ActivityListGet(w http.ResponseWriter, r *http.Request, ds appengine.Context, u *user.User) {
	var q *datastore.Query
	var act Activity_entity
	var owner string
	l := make([]Activity_entity,0)
	owner = r.FormValue("owner")
	// if owner is not specified, use currently logged in user 
	if owner == "" {
		owner = u.Email
	}

	if owner == "all" {
		q = datastore.NewQuery("Activity_entity")
	} else {
		// build query to return activities by owner
		q = datastore.NewQuery("Activity_entity").Filter("Owner = ", owner)
	}
	t := q.Run(ds)
	for t != nil {
		key,err := t.Next(&act)
    if err == datastore.Done {
			break
    }
		if err != nil {
			break
		}
		ds.Infof("activity key: %v", key.IntID())
		act.Id = strconv.FormatInt(key.IntID(),10)
		l = append(l,act)
	}
	jf, err := json.Marshal(l)
  if err != nil {
    fmt.Fprint(w, "{\"errror\":\"Error marshalling json\"}")
  } else {
    w.Write(jf)
  }
}

// Clear all data owned by a user
// dangerous, use with care (would be protected by admin priveleges in
// a production system)
func ClearUserData(w http.ResponseWriter, r *http.Request) {
	var q *datastore.Query
	ds := appengine.NewContext(r)
	owner := r.FormValue("owner")
	if owner == "" {
		u := user.Current(ds)
		if u != nil {
			owner = u.Email
		} else {
			fmt.Fprint(w, "{\"errror\":\"Error obtaining user credentials\"}")
			return
		}
	} else {
		// user specified, must be admin
		if user.IsAdmin(ds) == false {
			fmt.Fprint(w, "{\"errror\":\"Error user must be admin\"}")
			return
		}
	}
	q = datastore.NewQuery("Activity_entity").Filter("Owner = ", owner).KeysOnly()
	keylist,_ := q.GetAll(ds,nil)
	datastore.DeleteMulti(ds, keylist)
	q = datastore.NewQuery("Event_entity").Filter("Owner = ", owner).KeysOnly()
	keylist,_ = q.GetAll(ds,nil)
	datastore.DeleteMulti(ds, keylist)
	q = datastore.NewQuery("Booking_entity").Filter("Owner = ", owner).KeysOnly()
	keylist,_ = q.GetAll(ds,nil)
	datastore.DeleteMulti(ds, keylist)
	fmt.Fprint(w, "{\"method\":\"DELETE\",\"message\":\"SUCCESS\"}")
}


