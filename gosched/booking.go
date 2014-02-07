package gosched

import (
  "fmt"
	"strconv"
	"appengine"
	"appengine/user"
  "appengine/datastore"
	"encoding/json"
  "net/http"
	"errors"
)

// a booking is an object describing a party booking reservations for an event
type Booking_entity struct {
    Id string `json:"id" datastore:"-"`
    Event_id  int64 `json:"event_id"`
		Count int `json:"count"`
    Owner string  `json:"owner"`
}

// the event handler switches on the HTTP method to determine
// the function to call
func BookingHandler(w http.ResponseWriter, r *http.Request) {
  ds := appengine.NewContext(r)
  u := user.Current(ds)
	if u == nil {
    fmt.Fprint(w, "{\"errror\":\"User credentials could not be determined.\"}")
		return
	}
  switch {
    case r.Method == "GET":
      BookingGet(w,r,ds,u)
    case r.Method == "POST":
      BookingInsert(w,r,ds,u)
	}
	return
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
func BookingGet(w http.ResponseWriter, r *http.Request, ds appengine.Context, u *user.User) {
  var book Booking_entity
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
func BookingInsert(w http.ResponseWriter, r *http.Request, ds appengine.Context, u *user.User) {
	var e Event_entity
  bkey := datastore.NewIncompleteKey(ds, "Booking_entity", nil)
  book := InhaleBooking(r)
  book.Owner = u.Email

	// before putting, get the event
	// verify that it exists
  ekey := datastore.NewKey(ds, "Event_entity", "", book.Event_id, nil)
	err := datastore.RunInTransaction(ds, func(ds appengine.Context) error {
		err := datastore.Get(ds, ekey, &e)
    if err != nil {
      return err
    }
	  // ensure the event is available and has room for the party
	  if !e.Available && e.Bookings_count + book.Count <= e.Max_attendees {
			err = errors.New("No availability")
      return err
    }
	  // update event with booking count incremented
	  e.Bookings_count = e.Bookings_count + book.Count
	  // change status if closed
	  if e.Bookings_count >= e.Max_attendees {
		  e.Available = false
	  }
    bkey,err = datastore.Put(ds, bkey, &book)
    if err != nil {
      return err
    }
    ekey,err = datastore.Put(ds, ekey, &e)
    if err != nil {
      return err
    }
		return err
	}, nil)
  if err != nil {
		fmt.Fprintf(w, "{\"errror\":\"Unable to book %v\"}", book.Event_id)
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
