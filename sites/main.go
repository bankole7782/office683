package main

import (
  "net/http"
  "fmt"
  "path/filepath"
  "os"
  "html/template"
  "github.com/gorilla/mux"
  "github.com/bankole7782/office683/office683_shared"
  "github.com/bankole7782/office683/sites/docs"
)

func main() {
  rootPath, err := office683_shared.GetRootPath()
  if err != nil {
    panic(err)
  }

  os.MkdirAll(filepath.Join(rootPath, "docs"), 0777)
  os.MkdirAll(filepath.Join(rootPath, "docs_images"), 0777)

  r := mux.NewRouter()

  r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFS(office683_shared.Content, "templates/home.html"))
    tmpl.Execute(w, nil)
  })

  r.HandleFunc("/gs/{obj}", func (w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    rawObj, err := office683_shared.ContentStatics.ReadFile("statics/" + vars["obj"])
    if err != nil {
      panic(err)
    }
    w.Header().Set("Content-Disposition", "attachment; filename=" + vars["obj"])
    contentType := http.DetectContentType(rawObj)
    w.Header().Set("Content-Type", contentType)
    w.Write(rawObj)
  })

  docs.AddHandlers(r)

  fmt.Println("Running docs @ http://127.0.0.1:8080")

  err = http.ListenAndServe(":8080", r)
  if err != nil {
    fmt.Println(err)
    panic(err)
  }
}
