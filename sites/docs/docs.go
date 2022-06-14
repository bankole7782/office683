package docs

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
  "github.com/russross/blackfriday/v2"
  "github.com/dustin/go-humanize"
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
  rootPath, err := office683_shared.GetRootPath()
  if err != nil {
    panic(err)
  }

  flaarumClient := office683_shared.GetFlaarumClient()
  folderRows, err := flaarumClient.Search(`
    table: docs distinct
    fields: folder
    order_by: folder asc
    `)
  if err != nil {
    ErrorPage(w, errors.Wrap(err, "flaarum search error"))
    return
  }
  folders := make([]string, 0)
  folderToDocsMap := make(map[string][]map[string]interface{})
  for _, row := range *folderRows {
    folderName := row["folder"].(string)
    folders = append(folders, row["folder"].(string))

    docRows, err := flaarumClient.Search(fmt.Sprintf(`
      table: docs
      where:
        folder = '%s'
      `, folderName))

    if err != nil {
      fmt.Println(err)
    }

    elems := make([]map[string]interface{}, 0)
    for _, docRow := range *docRows {
      docPath := filepath.Join(rootPath, "docs", strconv.FormatInt(docRow["id"].(int64), 10) + ".md")
      docSize := "0B"
      if office683_shared.DoesPathExists(docPath) {

        file, err := os.Open(docPath)
        if err != nil {
          ErrorPage(w, errors.Wrap(err, "os error"))
          return
        }
        defer file.Close()

        stat, err := file.Stat()
        if err != nil {
          ErrorPage(w, errors.Wrap(err, "os error"))
          return
        }
        docSize = humanize.Bytes(uint64(stat.Size()))
      }

      elems = append(elems, map[string]interface{} {
        "doc_title": docRow["doc_title"],
        "updated": docRow["update_dt"],
        "id": docRow["id"],
        "doc_size": docSize,
      })

    }

    folderToDocsMap[folderName] = elems
  }

  type Context struct {
    Folders []string
    FoldersToDocsMap map[string][]map[string]interface{}
  }

  tmpl := template.Must(template.ParseFS(content, "templates/all_docs.html"))
  tmpl.Execute(w, Context{folders, folderToDocsMap})
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


func viewRenderedDoc(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  docId := vars["id"]

  rootPath, err := office683_shared.GetRootPath()
  if err != nil {
    panic(err)
  }

  docPath := filepath.Join(rootPath, "docs", docId + ".md")

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

  var rawDoc []byte
  if office683_shared.DoesPathExists(docPath) {
    rawDoc2, err := os.ReadFile(docPath)
    if err != nil {
      ErrorPage(w, errors.Wrap(err, "os error"))
      return
    }
    rawDoc = rawDoc2
  }

  renderedDocRaw := blackfriday.Run(rawDoc)
  renderedDoc := template.HTML(string(renderedDocRaw))

  type Context struct {
    RenderedDoc template.HTML
    DocDetails map[string]string
    DocId string
  }

  tmpl := template.Must(template.ParseFS(content, "templates/render_doc.html"))
  tmpl.Execute(w, Context{renderedDoc, docDetails, docId})
}


func deleteDoc(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  docId := vars["id"]

  rootPath, err := office683_shared.GetRootPath()
  if err != nil {
    panic(err)
  }

  flaarumClient := office683_shared.GetFlaarumClient()
  err = flaarumClient.DeleteRows(fmt.Sprintf(`
    table: docs
    where:
      id = %s
    `, docId))
  if err != nil {
    ErrorPage(w, errors.Wrap(err, "flaarum search error"))
    return
  }

  docPath := filepath.Join(rootPath, "docs", docId + ".md")
  os.Remove(docPath)

  fmt.Fprintf(w, "ok")
}
