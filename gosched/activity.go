package gosched

import (
  "fmt"
	"time"
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
    Creation_date time.Time `json:"creation_date"`
    Last_modified time.Time `json:"last_modified"`
}

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
    case r.Method == "PATCH":
      ActivityPatch(w,r)
    default:
      fmt.Fprint(w, "Activity handler.")
  }
}

func InhaleID(r *http.Request) (int64, error) {
	// read the id field from a form request
	// validate that it only contains numeric characters
	// return the ID in int64 form
	id := r.PostFormValue("id")
	id64,err := strconv.ParseUint(id ,10, 64)
	return id64,err
}

func InhaleActivity(r *http.Request) Activity_entity {
	var a Activity_entity
	ds := appengine.NewContext(r)
	a.Id = r.PostFormValue("id")
	a.Description = r.PostFormValue("description")
	a.Title = r.PostFormValue("title")
	ds.Infof("title:%v", r.PostFormValue("title"))
	return a
}

func ActivityGet(w http.ResponseWriter, r *http.Request) {
	var act Activity_entity
  ds := appengine.NewContext(r)
  id64,_ := InhaleID(r)
  key := datastore.NewKey(ds, "Activity_entity", "", id64, nil)
  err := datastore.Get(ds, key, &act)
  if err != nil {
    fmt.Fprint(w, "{\"errror\":\"Activity %v not found\"}", id64)
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

func ActivityDelete(w http.ResponseWriter, r *http.Request) {
  ds := appengine.NewContext(r)
  id64,err := strconv.ParseInt(r.FormValue("id"), 10, 0)
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
	fmt.Fprint(w, "{\"method\":\"DELETE\",\"id\":\"%v\",\"message\":\"SUCCESS\"}",id64)
}

func ActivityUpdate(w http.ResponseWriter, r *http.Request) {
	var act Activity_entity
	act.Title = "Update test activity"
	jf, err := json.Marshal(act)
	if err != nil {
		fmt.Fprint(w, "{\"errror\":\"Error marshalling json\"}")
	} else {
		w.Write(jf)
	}
}

func ActivityPatch(w http.ResponseWriter, r *http.Request) {
	var act Activity_entity
	act.Title = "Patch test activity"
	jf, err := json.Marshal(act)
	if err != nil {
		fmt.Fprint(w, "{\"errror\":\"Error marshalling json\"}")
	} else {
		w.Write(jf)
	}
}
