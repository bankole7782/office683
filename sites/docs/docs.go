package main

import (
  "net/http"
  "html/template"
  "time"
  "strconv"
  "fmt"
  "os"
  "path/filepath"
  "github.com/pkg/errors"
  "github.com/gorilla/mux"
  "github.com/bankole7782/office683/office683_shared"
)


func newDocument(w http.ResponseWriter, r *http.Request) {
  if r.Method == http.MethodPost {
    flaarumClient := office683_shared.GetFlaarumClient()
    retId, err := flaarumClient.InsertRowAny("docs", map[string]interface{} {
      "doc_title": r.FormValue("doc_title"),
      "folder": r.FormValue("folder"),
      "update_dt": time.Now(),
      "public": false,
    })

    if err != nil {
      ErrorPage(w, errors.Wrap(err, "flaarum insert error"))
      return
    }

    http.Redirect(w, r, "/update_doc/" + strconv.FormatInt(retId, 10), 307)
  }
}



func allDocs(w http.ResponseWriter, r *http.Request) {
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


func updateDoc(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  docId := vars["id"]

  flaarumClient := office683_shared.GetFlaarumClient()
  docRow, err := flaarumClient.SearchForOne(fmt.Sprintf(`
    table: docs
    where:
      id = %s
    `, docId))
  if err != nil {
    ErrorPage(w, errors.Wrap(err, "flaarum search error"))
    return
  }

  docDetails := map[string]string {
    "doc_title": (*docRow)["doc_title"].(string),
    "folder": (*docRow)["folder"].(string),
    "updated": (*docRow)["update_dt"].(time.Time).String(),
  }

  rootPath, err := office683_shared.GetRootPath()
  if err != nil {
    panic(err)
  }
  rawDoc := ""
  docPath := filepath.Join(rootPath, "docs", docId + ".md")
  if office683_shared.DoesPathExists(docPath) {
    raw, err := os.ReadFile(docPath)
    if err != nil {
      ErrorPage(w, errors.Wrap(err, "os error"))
      return
    }
    rawDoc = string(raw)
  }

  type Context struct {
    DocDetails map[string]string
    RawDoc string
    DocId string
  }
  tmpl := template.Must(template.ParseFS(content, "templates/update_doc.html"))
  tmpl.Execute(w, Context{docDetails, rawDoc, docId})
}


func saveDoc(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  docId := vars["id"]

  rootPath, err := office683_shared.GetRootPath()
  if err != nil {
    panic(err)
  }

  docPath := filepath.Join(rootPath, "docs", docId + ".md")
  os.WriteFile(docPath, []byte(r.FormValue("raw_doc")), 0777)

  flaarumClient := office683_shared.GetFlaarumClient()
  err = flaarumClient.UpdateRowsAny(fmt.Sprintf(`
    table: docs
    where:
     id = %s
    `, docId), map[string]interface{} {
    "update_dt": time.Now(),
  })
  if err != nil {
    fmt.Println(err)
  }
  fmt.Fprintf(w, "ok")
}
