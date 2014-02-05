gosched
=======

A simple scheduling API in go for google appengine

The data model consists of two entities:
 1) events describe a time, place, and duration and belong to a single activity
 2) activites describe a class of events (description, vendor, etc)

Events can be searched for by activity + time range and by owner. Booking an event adds the owner's party to the attendee list. 
