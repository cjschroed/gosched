package gosched

import (
    "fmt"
    "net/http"
)

func init() {
	http.HandleFunc("/gosched/v1/hello", HelloHandler)
	http.HandleFunc("/gosched/v1/activity", ActivityHandler)
	http.HandleFunc("/gosched/v1/activity/list", ActivityListHandler)
	http.HandleFunc("/gosched/v1/activity/events", EventsHandler)
	http.HandleFunc("/gosched/v1/activity/events/list", EventListGet)
	http.HandleFunc("/gosched/v1/unittests", UnitTestSection)
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hello, world!")
}
