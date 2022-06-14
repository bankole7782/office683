package auths

import (
  "github.com/gorilla/mux"
)


func AddHandlers(r *mux.Router) {
  // main auths
  r.HandleFunc("/register", registerUser)
  r.HandleFunc("/signin", signInHandler)
  r.HandleFunc("/signout", signout)

  // teams
  r.HandleFunc("/teams", allTeams)
  r.HandleFunc("/new_team", newTeam)
  r.HandleFunc("/update_team_members", updateTeamMembers)
}
