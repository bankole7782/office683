package events

import (
  "github.com/gorilla/mux"
)


func AddHandlers(r *mux.Router) {
  r.HandleFunc("/events/", allEvents)
  r.HandleFunc("/new_event", newEvent)
  r.HandleFunc("/event/{id}", aEvent)
  r.HandleFunc("/delete_event/{id}", deleteAEvent)
}
