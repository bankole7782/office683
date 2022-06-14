package auths

import (
  "fmt"
  "html/template"
  "net/http"
  "github.com/pkg/errors"
  "strconv"
  "github.com/bankole7782/office683/office683_shared"
)


func allTeams(w http.ResponseWriter, r *http.Request) {
  status, userDetails := office683_shared.IsLoggedInUser(r)
  if status == false {
    http.Redirect(w, r, "/", 307)
    return
  }

  flaarumClient := office683_shared.GetFlaarumClient()
  rows, err := flaarumClient.Search(fmt.Sprintf(`
    table: teams
    order_by: team_name asc
    `))
  if err != nil {
    office683_shared.ErrorPage(w, errors.Wrap(err, "flaarum error"))
    return
  }

  teams := make([]map[string]any, 0)
  teamMembers := make(map[string][]map[string]string)
  for _, row := range *rows {
    teamMembersRows, err := flaarumClient.Search(fmt.Sprintf(`
      table: team_members expand
      order_by: userid.firstname asc
      where:
        teamid = %d
      `, row["id"].(int64)))
    if err != nil {
      office683_shared.ErrorPage(w, errors.Wrap(err, "flaarum error"))
      return
    }

    members := make([]map[string]string, 0)

    for _, teamRow := range *teamMembersRows {
      elem := map[string]string {
        "fullname": teamRow["userid.firstname"].(string) + " " + teamRow["userid.surname"].(string),
        "email": teamRow["userid.email"].(string), "id": strconv.FormatInt(teamRow["userid"].(int64), 10),
      }

      members = append(members, elem)
    }

    teams = append(teams, row)
    teamMembers[row["team_name"].(string)] = members
  }

  allUserRows, err := flaarumClient.Search(`
    table: users
    fields: firstname surname id
    order_by: firstname asc
    `)
  if err != nil {
    office683_shared.ErrorPage(w, errors.Wrap(err, "flaarum error"))
    return
  }

  type Context struct {
    UserDetails map[string]any
    Teams []map[string]any
    TeamMembers map[string][]map[string]string
    AllUsers []map[string]any
    IsEven func(int) bool
  }
  ie := func(i int) bool {
    return ((i+1) %2) == 0
  }
  tmpl := template.Must(template.ParseFS(content, "templates/all_teams.html"))
  tmpl.Execute(w, Context{userDetails, teams, teamMembers, *allUserRows, ie})
}


func newTeam(w http.ResponseWriter, r *http.Request) {
  if r.Method == http.MethodPost {
    status, _ := office683_shared.IsLoggedInUser(r)
    if status == false {
      http.Redirect(w, r, "/", 307)
      return
    }

    flaarumClient := office683_shared.GetFlaarumClient()

    _, err := flaarumClient.InsertRowStr("teams", map[string]string {
      "team_name": r.FormValue("team_name"), "desc": r.FormValue("desc"),
    })

    if err != nil {
      office683_shared.ErrorPage(w, errors.Wrap(err, "flaarum error"))
      return
    }

    r.Method = http.MethodGet
    http.Redirect(w, r, "/teams", 307)
  }

}

func updateTeamMembers(w http.ResponseWriter, r *http.Request) {
  if r.Method == http.MethodPost {
    flaarumClient := office683_shared.GetFlaarumClient()

    r.FormValue("team_id")

    if len(r.Form["member"]) > 0 {
      err := flaarumClient.DeleteRows(fmt.Sprintf(`
        table: team_members
        where:
          teamid = %s
        `, r.FormValue("team_id")))
      if err != nil {
        office683_shared.ErrorPage(w, errors.Wrap(err, "flaarum error"))
        return
      }

      for _, memberId := range r.Form["member"] {
        _, err = flaarumClient.InsertRowStr("team_members", map[string]string {
          "teamid": r.FormValue("team_id"), "userid": memberId,
        })

        if err != nil {
          office683_shared.ErrorPage(w, errors.Wrap(err, "flaarum error"))
          return
        }
      }

    }

    r.Method = http.MethodGet
    http.Redirect(w, r, "/teams", 307)
  }
}
