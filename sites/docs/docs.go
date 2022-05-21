package main

import (
  "net/http"
  "html/template"
  "time"
  "strconv"
  "github.com/pkg/errors"
  "github.com/bankole7782/office683/office683_shared"
)


func newDocument(w http.ResponseWriter, r *http.Request) {
  if r.Method == http.MethodPost {
    flaarumClient := office683_shared.GetFlaarumClient()
    retId, err := flaarumClient.InsertRowAny("docs", map[string]interface{} {
      "doc_title": r.FormValue("doc_title"),
      "folder": r.FormValue("folder"),
      "update_dt": time.Now(),
    })

    if err != nil {
      ErrorPage(w, errors.Wrap(err, "flaarum insert error"))
      return
    }

    http.Redirect(w, r, "/update_doc/" + strconv.FormatInt(retId, 10), 307)
  }
}



func allDocs(w http.ResponseWriter, r *http.Request) {
  // rootPath, err := office683_shared.GetRootPath()
  // if err != nil {
  //   panic(err)
  // }
  //

  // flaarumClient := office683_shared.GetFlaarumClient()
  // folderRows, err := flaarumClient.Search(`
  //   table: docs distinct
  //   fields: folder
  //   `)
  // folders := make([]string, 0)
  // for _, row := range *rows {
  //   folders = append(folders, row["folder"].(string))
  // }
  //

  tmpl := template.Must(template.ParseFS(content, "templates/all_docs.html"))
  tmpl.Execute(w, nil)
}
