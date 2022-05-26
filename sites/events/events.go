package events

import (
  "fmt"
  "time"
  "strconv"
  "strings"
  "net/http"
  "html/template"
  "github.com/gorilla/mux"
  "github.com/pkg/errors"
  "github.com/bankole7782/office683/office683_shared"
)


func allEvents(w http.ResponseWriter, r *http.Request) {
  flaarumClient := office683_shared.GetFlaarumClient()
  rows, err := flaarumClient.Search(`
    table: events
    order_by: begin_date asc
    limit: 100
    `)
  if err != nil {
    ErrorPage(w, errors.Wrap(err, "flaarum search error"))
    return
  }

  evts := make([]map[string]string, 0)
  oldEvts := make([]map[string]string, 0)

  for _, row := range *rows {
    el := map[string]string {
      "id": strconv.FormatInt(row["id"].(int64), 10),
      "title": row["title"].(string),
      "begin_date": row["begin_date"].(time.Time).Format("2006-01-02"),
    }
    beginDate := row["begin_date"].(time.Time)
    if beginDate.After(time.Now()) {
      evts = append(evts, el)
    } else {
      oldEvts = append(oldEvts, el)
    }
  }

  type Context struct {
    Events []map[string]string
    OldEvents []map[string]string
  }

  tmpl := template.Must(template.ParseFS(content, "templates/all_events.html"))
  tmpl.Execute(w, Context{evts, oldEvts})
}


func newEvent(w http.ResponseWriter, r *http.Request) {
  flaarumClient := office683_shared.GetFlaarumClient()

  if r.Method == http.MethodGet {
    tmpl := template.Must(template.ParseFS(content, "templates/new_event.html"))
    tmpl.Execute(w, nil)
  } else {
    _, err := flaarumClient.InsertRowStr("events", map[string]string {
      "title": r.FormValue("title"),
      "begin_date": r.FormValue("begin_date"),
      "begin_time": r.FormValue("begin_time"),
      "end_date": r.FormValue("end_date"),
      "end_time": r.FormValue("end_time"),
      "event_description": r.FormValue("event_description"),
      "event_preparation": r.FormValue("event_preparation"),
    })

    if err != nil {
      ErrorPage(w, errors.Wrap(err, "flaarum insert error"))
      return
    }

    http.Redirect(w, r, "/events/", 307)
  }
}


func aEvent(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  docId := vars["id"]

  flaarumClient := office683_shared.GetFlaarumClient()
  eventRow, err := flaarumClient.SearchForOne(fmt.Sprintf(`
    table: events
    where:
      id = %s
    `, docId))
  if err != nil {
    ErrorPage(w, errors.Wrap(err, "flaarum search error"))
    return
  }

  eventDetails := map[string]string {
    "id": strconv.FormatInt((*eventRow)["id"].(int64), 10),
    "title": (*eventRow)["title"].(string),
    "begin_date": (*eventRow)["begin_date"].(time.Time).Format("2006-01-02"),
    "end_date": (*eventRow)["end_date"].(time.Time).Format("2006-01-02"),
    "begin_time": (*eventRow)["begin_time"].(string),
    "end_time": (*eventRow)["end_time"].(string),
    "desc": "",
    "preps": "",
  }

  if (*eventRow)["event_description"] != nil {
    eventDetails["desc"] = (*eventRow)["event_description"].(string)
  }

  if (*eventRow)["event_preparation"] != nil {
    eventDetails["preps"] = (*eventRow)["event_description"].(string)
  }

  type Context struct {
    Event map[string]string
    IsNotEmpty func(string) bool
  }

  ine := func(s string) bool {
    return strings.TrimSpace(s) != ""
  }

  tmpl := template.Must(template.ParseFS(content, "templates/a_event.html"))
  tmpl.Execute(w, Context{eventDetails, ine})
}


// func editAEvent(w http.ResponseWriter, r *http.Request) {
//   vars := mux.Vars(r)
//   docId := vars["id"]
//
//   flaarumClient := office683_shared.GetFlaarumClient()
//   docRow, err := flaarumClient.SearchForOne(fmt.Sprintf(`
//     table: events
//     where:
//       id = %s
//     `, docId))
//   if err != nil {
//     ErrorPage(w, errors.Wrap(err, "flaarum search error"))
//     return
//   }
//
//
// }
