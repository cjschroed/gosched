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

// an event is an instance of an activity
// it contains the time and attendence for the instance
type Event_entity struct {
  Id string `json:"id" datastore:"-"`
  Title  string `json:"title"`
  Description string `json:"description"`
  Owner string  `json:"owner"`
	Activity_id string `json:"activity_id"`
  Start_time time.Time `json:"start_time"`
  End_time time.Time `json:"end_time"`
	Interval int `json:"interval"` // crontabesque code for repeating events
	Max_attendees int `json:"max_attendees"` // max attendees for this event
	Bookings_count int `json:"bookings_count"` // count of bookings against
	Available bool `json:"available"`  // true if there are open slots
}

// the event handler switches on the HTTP method to determine
// the function to call
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

// EventGet: retreive an event by it's ID
// Form parameters expected:
//   id: string representing the datastore ID of the event to be retreived 
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
// into a go structure. This is the place where limits on the data are enforced
// to prevent malicious intent by API users
func InhaleEvent(r *http.Request) Event_entity {
	var e Event_entity
	const longForm = "Jan 2, 2006 at 3:04pm (MST)"
	e.Activity_id = r.FormValue("activity_id")
	e.Description = r.FormValue("description")
	e.Title = r.FormValue("title")
	e.Start_time,_ = time.Parse(longForm, r.FormValue("start_time"))
	e.End_time,_ = time.Parse(longForm, r.FormValue("end_time"))
	e.Max_attendees,_ = strconv.Atoi(r.FormValue("max_attendees"))
	return e
}

// EventInsert: create a new event with data supplied by the client
// Form parameters expected:
//	title: string title field 
//  description: string description field 
//  activity_id: integer in string form denoting the datastore ID of the
//		activity this event correpsonds to
//  start_time: string representation for time of the start the event,
//		example format:  "Jan 2, 2006 at 3:04pm (MST)"
//  end_time: string representation for the ending time for the event, 
//  max_attendees: integer in string form denoting the maximum number of 
//		people allowed to attend the event 
func EventInsert(w http.ResponseWriter, r *http.Request) {
	ds := appengine.NewContext(r)
  event := InhaleEvent(r)
  u := user.Current(ds)
  if u != nil {
    event.Owner = u.Email
  } else {
    event.Owner = "Guest"
  }
	if event.Max_attendees > event.Bookings_count {
		event.Available = true
	} else {
		event.Available = false 
	}
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

// EventDelete: remove an event by ID
// Form parameters expected:
//   id: string representing the datastore ID of the event to be deleted
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
	fmt.Fprintf(w, "{\"method\":\"DELETE\",\"id\":\"%v\",\"message\":\"SUCCESS\"}",id64)
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
	// set availabilty to false if there are no more slots available to book
	if event.Max_attendees > event.Bookings_count {
		event.Available = true
	} else {
		event.Available = false 
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

// EventListGet: returns all event objects that correspond to an activity
// Form parameters expected:
// 	activity_id: string representing the datastore ID of the activity 	 
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
	// build query to return events for an activity that have availability
	q = datastore.NewQuery("Event_entity").Filter("Activity_id = ", actid).Filter("Available = ", true)
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

// ListEventByAvailability lists all the events that are not full which 
// are available for an activity in a time range
// Form parameters expected by ListEventByAvailability:
// 	day: Mmm DD, YYYY in string format, starting day for the time range
//  dayend: Mmm DD, YYYY in string format, last day (inclusive) 
//  activity_id: datastore ID of the activity the event belongs to
//  results: If present, returns the set of days for which the activity 
//		has available events. If absent, returns the list of available events 
//		that exist for the time period specified

func ListEventByAvailability(w http.ResponseWriter, r *http.Request) {
	var q *datastore.Query
	var act Event_entity
	var jf []byte
	var err error
	const longForm = "Jan 2, 2006"
	day,_ := time.Parse(longForm, r.FormValue("day"))
	dayend,err := time.Parse(longForm, r.FormValue("dayend"))
	if err != nil {
		dayend = day.Add(time.Duration(24)*time.Hour)
	} else {
		dayend = dayend.Add(time.Duration(24)*time.Hour)
	}
	l := make([]Event_entity,0)
	daylist := make([]string, 0)
	ds := appengine.NewContext(r)
	ds.Infof("day in question: %v",day.Format(longForm))
	ds.Infof("dayend: %v",dayend.Format(longForm))
	actid := r.FormValue("activity_id")
	returns := r.FormValue("results")
	if actid == "" {
    fmt.Fprint(w, "{\"errror\":\"No activity ID specified\"}")
		return
	}
	// build query to return events for an activity that have availability
	q = datastore.NewQuery("Event_entity").Filter("Activity_id = ", actid).Filter("Available = ", true).Filter("Start_time >= " , day).Filter("Start_time <= ", dayend)
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
		years,months,days := act.Start_time.Date()
		d := time.Date(years,months,days,0,0,0,0,time.UTC)
		daylist = append(daylist,d.Format(longForm))
	}
	if returns == "" {
		jf, err = json.Marshal(l)
	} else {
		jf, err = json.Marshal(daylist)
	}
  if err != nil {
    fmt.Fprint(w, "{\"errror\":\"Error marshalling json\"}")
  } else {
    w.Write(jf)
  }
}
