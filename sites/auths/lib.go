package auths

import (
  "github.com/gorilla/mux"
)


func AddHandlers(r *mux.Router) {
  // documents
  r.HandleFunc("/register", registerUser)
  r.HandleFunc("/signin", signInHandler)
}
