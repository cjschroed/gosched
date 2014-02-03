package gosched

import (
  "fmt"
  "time"
	"strconv"
	"appengine"
  "appengine/datastore"
  "encoding/json"
  "net/http"
)

type Event_entity struct {
  Id string `json:"id" datastore:"-"`
  Title  string `json:"title"`
  Description string `json:"description"`
  Owner string  `json:"owner"`
	Activity_id string `json:"activity_id"`
  Start_time time.Time `json:"start_time"`
  End_time time.Time `json:"end_time"`
	Status string `json:"status"`
	Max_attendees int `json:"max_attendees"`
}

func EventsHandler(w http.ResponseWriter, r *http.Request) {
  switch {
    case r.Method == "GET":
      EventGet(w,r)
    case r.Method == "POST":
      EventInsert(w,r)
    case r.Method == "DELETE":
      EventDelete(w,r)
    case r.Method == "PUT":
      EventUpdate(w,r)
    default:
      fmt.Fprint(w, "Event handler.")
  }
}

// retrieve an event by it's ID
func EventGet(w http.ResponseWriter, r *http.Request) {
  var event Event_entity
	ds := appengine.NewContext(r)
  id64,_ := strconv.ParseInt(r.FormValue("id"), 10, 0)
  key := datastore.NewKey(ds, "Event_entity", "", id64, nil)
  g_err := datastore.Get(ds, key, &event)
  if g_err != nil {
    fmt.Fprint(w, "{\"errror\":\"Event not found\"}")
    return
  }
  jf, err := json.Marshal(event)
  if err != nil {
    fmt.Fprint(w, "{\"errror\":\"Error marshalling json\"}")
  } else {
    w.Write(jf)
  }
}

// This is the function responsible for converting the http parameter data
// into a go structure, this is the place where limits on the data are enforced
// to prevent malicious intent by API users
func InhaleEvent(r *http.Request) Event_entity {
	var e Event_entity
	const longForm = "Jan 2, 2006 at 3:04pm (MST)"
	e.Activity_id = r.FormValue("activity_id")
	e.Description = r.FormValue("description")
	e.Title = r.FormValue("title")
	e.Start_time,_ = time.Parse(longForm, r.FormValue("start_time"))
	e.End_time,_ = time.Parse(longForm, r.FormValue("end_time"))
	e.Status = r.FormValue("status")
	e.Max_attendees,_ = strconv.Atoi(r.FormValue("max_attendees"))
	return e
}

func EventInsert(w http.ResponseWriter, r *http.Request) {
	ds := appengine.NewContext(r)
  event := InhaleEvent(r)
  key := datastore.NewIncompleteKey(ds, "Event_entity", nil)
  mykey,err := datastore.Put(ds, key, &event)
  if err != nil {
    fmt.Fprint(w, "{\"errror\":\"Event not found\"}")
    return
  }
	event.Id = strconv.FormatInt(mykey.IntID(),10)
  jf, err := json.Marshal(event)
  if err != nil {
    fmt.Fprint(w, "{\"errror\":\"Error marshalling json\"}")
  } else {
    w.Write(jf)
  }
}

func EventDelete(w http.ResponseWriter, r *http.Request) {
	ds := appengine.NewContext(r)
	id64,err := strconv.ParseInt(r.FormValue("id"), 10, 0)
  if err != nil {
    fmt.Fprint(w, "{\"errror\":\"ID could not be read\"}")
    return
  }
	key := datastore.NewKey(ds, "Event_entity", "", id64, nil)
	err = datastore.Delete(ds,key)
  if err != nil {
    fmt.Fprint(w, "{\"errror\":\"Event not found\"}")
    return
  }
	fmt.Fprint(w, "{\"method\":\"DELETE\",\"id\":\"%v\",\"message\":\"SUCCESS\"}",id64)
}

func EventUpdate(w http.ResponseWriter, r *http.Request) {
	ds := appengine.NewContext(r)
	var e Event_entity
  event := InhaleEvent(r)
	id64,err := strconv.ParseInt(event.Id, 10, 0)
	key := datastore.NewKey(ds, "Event_entity", "", id64, nil)
	err = datastore.Get(ds, key, &e)
  if err != nil {
    fmt.Fprint(w, "{\"errror\":\"Event not found\"}")
    return
  }
  _,err = datastore.Put(ds, key, &event)
  if err != nil {
    fmt.Fprint(w, "{\"errror\":\"Error updating event\"}")
		return
	}
  jf, err := json.Marshal(event)
  if err != nil {
    fmt.Fprint(w, "{\"errror\":\"Error marshalling json\"}")
  } else {
    w.Write(jf)
  }
}

func EventListGet(w http.ResponseWriter, r *http.Request) {
	var q *datastore.Query
	var act Event_entity
	l := make([]Event_entity,0)
	ds := appengine.NewContext(r)
	actid := r.FormValue("activity_id")
	if actid == "" {
    fmt.Fprint(w, "{\"errror\":\"No activity ID specified\"}")
		return
	}
	// build query to return activities by an owner
	q = datastore.NewQuery("Event_entity").Filter("Activity_id = ", actid)
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

