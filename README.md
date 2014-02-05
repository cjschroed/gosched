gosched
=======

A simple scheduling API in go for google appengine

The data model consists of three entities:

 1) events describe a time, place, and duration and belong to a single activity

 2) activites describe a class of events (description, vendor, etc)

 3) bookings describe a party that is attending an event

Events can be searched for by activity + time range and by owner. Booking an event adds the owner's party to the attendee list. 

The live version can be viewed at https://csgosched.appspot.com/

To run the development environment, you need to download and unzip the [b]Google App Engine SDK for Go[/b] found here: 
https://developers.google.com/appengine/downloads 

Tips on getting started with Go on Appengine can be viewed here:
https://developers.google.com/appengine/docs/go/gettingstarted/devenvironment

Clone the repository (produces the gosched directory) and then run <path to Go SDK>/goapp serve gosched 

You should be able to view the app in the development environment at localhost:8080. You can see the development console at localhost:8000

Scheduling API Methods
======================

To create an activity:  
 	url: /gosched/v1/activity  
	method: POST  
	parameters:  
		title: string  
		description: string  
		vendor_name: string  
	returns:  
		JSON of the activity object, including its ID and owner  

To get an activity object:  
	url: /gosched/v1/activity  
	method: GET   
	parameters:  
		id: the ID returned in the JSON object that was created  
	returns:  
		JSON of the activity object	  

To delete an activity:  
	url: /gosched/v1/activity  
	method: DELETE  
	parameters:  
		id: the ID returned in the JSON object that was created  
	returns:  
		JSON object with the ID of the deleted object and message == "SUCCESS"  

To create an event for an activity:  
	url: /gosched/v1/activity/events  
	method: POST  
	parameters:  
		title:  
		description:  
		activity_id:  
		start_time:  
		end_time:  
		max_attendees:  
	returns:  
		JSON of the event object, including the datastore ID  

To delete an event:  
	url: /gosched/v1/activity/events  
	method: DELETE  
	parameters:  
		id: the ID returned in the JSON object that was created  
	returns:  
		JSON object with the ID of the deleted object and message == "SUCCESS"  

To search for activity day availability over a time range:  
	url: /gosched/v1/activity/events/search  
	method: GET  
	parameters:  
		day: date for the first day in the range, ex. Jan 4, 2013  
		dayend: date for the last day in the range inclusive, same format as day  
		activity_id: integer in string form of the activity ID  
		results: string to designate the type of return, must be present to   
			return an array of dates with avaiable events. ex. results=days  
	returns:  
		A JSON array of dates that have event objects with availability for  
		the designated activity  

To search for activity availability over a range of time:  
	url: /gosched/v1/activity/events/search  
	method: GET  
	parameters:  
		day: date for the first day in the range, ex. Jan 4, 2013  
		dayend: date for the last day in the range inclusive, same format as day  
		activity_id: integer in string form of the activity ID  
	returns:  
		A JSON array of event objects for the designated activity with   
		start times between day (start) and dayend (end, inclusive) that   
		have availability  

To search for activity availability during a day:  
	url: /gosched/v1/activity/events/search  
	method: GET  
	parameters:  
		day: date for the first day in the range, ex. Jan 4, 2013  
		activity_id: integer in string form of the activity ID  
	returns:  
		A JSON array of event objects for the designated activity with   
		start times that fall on day with availability  

To create a booking against an activity event:  
- url: /gosched/v1/activity/book  
- method: POST  
- parameters:  
    - event_id: integer in string form for the event to book  
    - count: integer in string form for number of attendees in the booking,  presumed to be 1 if absent  
- returns:  JSON of the booking object, including the ID 	  

To view a list of activities:  
- url: /gosched/v1/activity/list  
- method: GET  
- parameters:  
    - owner: a string reprenting the owner of the activies you wish to   view, or "all" to see all activities (limit 100). Cursors to walk   the entire list are not implemented yet. The list of all   activities across all vendors is for demonstration purposes only.   The vendor's email address used to create the activities is the key  to find all activities supplied by that vendor. In a production env  there would be a vendor registration system that would provide a more appropriate key than the email address.  

To delete all data associated with a user:  
- url: /gosched/v1/activity/clear  
- method: GET  
- parameters:  
    - owner: ownername in string format  
- returns:  JSON object with message = "SUCCESS" or error object  

Improvements  

