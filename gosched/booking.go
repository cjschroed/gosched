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

// a booking is an object describing a party booking reservations for an event
type Booking_entity struct {
    Id string `json:"id" datastore:"-"`
    Event_id  int64 `json:"event_id"`
    Instance int `json:"instance"`
		Count int `json:"count"`
    Owner string  `json:"owner"`
}

// the event handler switches on the HTTP method to determine
// the function to call
func BookingHandler(w http.ResponseWriter, r *http.Request) {
  switch {
    case r.Method == "GET":
      BookingGet(w,r)
    case r.Method == "POST":
      BookingInsert(w,r)
	}
}

func InhaleBooking(r *http.Request) Booking_entity {
	var a Booking_entity
	a.Event_id,_ = strconv.ParseInt(r.FormValue("event_id"),10,0)
	a.Count,_ = strconv.Atoi(r.FormValue("count"))
	// if count is not supplied, assumed to be 1
	if a.Count <= 0 {
		a.Count = 1
	}
	return a
}

// BookingGet: retreive a booking object by it's ID
// Form parameters expected:
//   id: string representing the datastore ID of the booking to be retreived 
func BookingGet(w http.ResponseWriter, r *http.Request) {
  var book Booking_entity
  ds := appengine.NewContext(r)
  id64,err := InhaleID(r,false)
  if err != nil {
    fmt.Fprint(w, "{\"errror\":\"ID could not be parsed\"}")
    return
  }
  key := datastore.NewKey(ds, "Booking_entity", "", id64, nil)
  err = datastore.Get(ds, key, &book)
  if err != nil {
    fmt.Fprintf(w, "{\"errror\":\"Booking %v not found, %v\"}", id64, err)
    return
  }
  book.Id = strconv.FormatInt(id64,10)
  jf, err := json.Marshal(book)
  if err != nil {
    fmt.Fprint(w, "{\"errror\":\"Error marshalling json\"}")
  } else {
    w.Write(jf)
  }
}

// BookingInsert: create a new event with data supplied by the client
// Form parameters expected:
//  event_id: integer in string form denoting the datastore ID of the
//		event this booking is for
//  count: integer denoting the number of attendees for this booking
func BookingInsert(w http.ResponseWriter, r *http.Request) {
	var e Event_entity
  ds := appengine.NewContext(r)
  bkey := datastore.NewIncompleteKey(ds, "Booking_entity", nil)
  book := InhaleBooking(r)
  u := user.Current(ds)
  if u != nil {
    book.Owner = u.Email
  } else {
    book.Owner = "Guest"
  }
	// before putting, get the event
	// verify that it exists
  ekey := datastore.NewKey(ds, "Event_entity", "", book.Event_id, nil)
  err := datastore.Get(ds, ekey, &e)
  if err != nil {
    fmt.Fprintf(w, "{\"errror\":\"Event %v not found, %v\"}", book.Event_id, err)
    return
  }
	// ensure the event is available and has room for the party
	if !e.Available && e.Bookings_count + book.Count <= e.Max_attendees {
    fmt.Fprintf(w, "{\"errror\":\"Event %v not available\"}", book.Event_id)
    return
  }
	// update event with booking count incremented
	e.Bookings_count = e.Bookings_count + book.Count 
	// change status if closed
	if e.Bookings_count >= e.Max_attendees {
		e.Available = false
	}
  bkey,err = datastore.Put(ds, bkey, &book)
  if err != nil {
    fmt.Fprint(w, "{\"errror\":\"Booking insert failed\"}")
    return
  }
  ekey,err = datastore.Put(ds, ekey, &e)
  if err != nil {
    fmt.Fprint(w, "{\"errror\":\"Event update failed\"}")
    return
  }
  book.Id = strconv.FormatInt(bkey.IntID(),10)
  jf, err := json.Marshal(book)
  if err != nil {
    fmt.Fprint(w, "{\"errror\":\"Error marshalling json\"}")
  } else {
    w.Write(jf)
	}
}