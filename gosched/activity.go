package gosched

import (
  "fmt"
	"strconv"
	"appengine"
	"appengine/user"
  "appengine/datastore"
	"encoding/json"
  "net/http"
)

type Activity_entity struct {
    Id string `json:"id" datastore:"-"`
    Title  string `json:"title"`
    Description string `json:"description"`
    Owner string  `json:"owner"`
		Vendor_name string `json:"vendor_name"`
}

// the event handler switches on the HTTP method to determine
// the function to call
func ActivityHandler(w http.ResponseWriter, r *http.Request) {
  switch {
    case r.Method == "GET":
      ActivityGet(w,r)
    case r.Method == "POST":
      ActivityInsert(w,r)
    case r.Method == "DELETE":
      ActivityDelete(w,r)
    case r.Method == "PUT":
      ActivityUpdate(w,r)
    default:
      fmt.Fprint(w, "Activity handler.")
  }
}

// InhaleID reads the id field from a form request and then
// validates that it only conforms to a string representation of an int64
// returns the ID in int64 form
func InhaleID(r *http.Request, l bool) (int64, error) {
	var id string
	if l {
		id = r.PostFormValue("id")
	} else {
		id = r.FormValue("id")
	}
	id64,err := strconv.ParseInt(id ,10, 0)
	return id64,err
}

// This is the function responsible for converting the http parameter data
// into an activity structure. Limits on the data are enforced
// to prevent malicious intent by API users
func InhaleActivity(r *http.Request) Activity_entity {
	var a Activity_entity
	a.Description = r.FormValue("description")
	a.Title = r.FormValue("title")
	a.Vendor_name = r.FormValue("vendor_name")
	return a
}

// ActivityGet: retreive an activity by it's ID
// Form parameters expected:
//   id: string representing the datastore ID of the activity to be retreived 
func ActivityGet(w http.ResponseWriter, r *http.Request) {
	var act Activity_entity
  ds := appengine.NewContext(r)
  id64,err := InhaleID(r,false)
  if err != nil {
    fmt.Fprint(w, "{\"errror\":\"ID could not be parsed\"}")
    return
  }
  key := datastore.NewKey(ds, "Activity_entity", "", id64, nil)
  err = datastore.Get(ds, key, &act)
  if err != nil {
    fmt.Fprintf(w, "{\"errror\":\"Activity %v not found, %v\"}", id64, err)
    return
  }
	act.Id = strconv.FormatInt(id64,10)
	jf, err := json.Marshal(act)
	if err != nil {
		fmt.Fprint(w, "{\"errror\":\"Error marshalling json\"}")
	} else {
		w.Write(jf)
	}
}

// ActivityInsert: create a new activity with data supplied by the client
// Form parameters expected:
//	title: string title field 
//  description: string description field 
//  vendor_name: string describing the vendor supplying the activity 
func ActivityInsert(w http.ResponseWriter, r *http.Request) {
	ds := appengine.NewContext(r)
	r.ParseForm()
  key := datastore.NewIncompleteKey(ds, "Activity_entity", nil)
	act := InhaleActivity(r)
	u := user.Current(ds)
	if u != nil {
		act.Owner = u.Email
	} else {
		act.Owner = "Guest"
	}
  mykey,err := datastore.Put(ds, key, &act)
  if err != nil {
    fmt.Fprint(w, "{\"errror\":\"Activity not found\"}")
    return
  }
	act.Id = strconv.FormatInt(mykey.IntID(),10)
  jf, err := json.Marshal(act)
  if err != nil {
    fmt.Fprint(w, "{\"errror\":\"Error marshalling json\"}")
  } else {
    w.Write(jf)
  }
}

// ActivityDelete: delete an activity object by it's ID
// Form parameters expected:
//   id: string representing the datastore ID of the activity to be removed
func ActivityDelete(w http.ResponseWriter, r *http.Request) {
  ds := appengine.NewContext(r)
  id64,err := InhaleID(r,false)
  if err != nil {
    fmt.Fprint(w, "{\"errror\":\"ID could not be read\"}")
    return
  }
  key := datastore.NewKey(ds, "Activity_entity", "", id64, nil)
  err = datastore.Delete(ds,key)
  if err != nil {
    fmt.Fprint(w, "{\"errror\":\"Activity %v not found\"}", id64)
    return
  }
	fmt.Fprintf(w, "{\"method\":\"DELETE\",\"id\":\"%v\",\"message\":\"SUCCESS\"}",id64)
}

// ActivityUpdate takes the same parameters as insert, but replaces the
//  exisiting object instead of creating a new one
func ActivityUpdate(w http.ResponseWriter, r *http.Request) {
  ds := appengine.NewContext(r)
	var old Activity_entity
  act := InhaleActivity(r)
  id64,err := InhaleID(r,true)
  key := datastore.NewKey(ds, "Activity_entity", "", id64, nil)
  err = datastore.Get(ds, key, &old)
  if err != nil {
    fmt.Fprint(w, "{\"errror\":\"Activity not found\"}")
    return
  }
  _,err = datastore.Put(ds, key, &act)
  if err != nil {
    fmt.Fprint(w, "{\"errror\":\"Error updating activity\"}")
    return
  }
  jf, err := json.Marshal(act)
  if err != nil {
    fmt.Fprint(w, "{\"errror\":\"Error marshalling json\"}")
  } else {
    w.Write(jf)
  }
}

