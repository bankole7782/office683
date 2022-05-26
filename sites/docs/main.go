package main

import (
  "net/http"
  "github.com/bankole7782/office683/office683_shared"
  "github.com/gorilla/mux"
  "fmt"
  "html/template"
  "strings"
  "path/filepath"
  "os"
)


func main() {
  rootPath, err := office683_shared.GetRootPath()
  if err != nil {
    panic(err)
  }

  os.MkdirAll(filepath.Join(rootPath, "docs"), 0777)

  r := mux.NewRouter()

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

  r.HandleFunc("/", allDocs)
  r.HandleFunc("/new_doc", newDocument)
  r.HandleFunc("/update_doc/{id}", updateDoc)
  r.HandleFunc("/save_doc/{id}", saveDoc)
  r.HandleFunc("/doc/{id}", viewRenderedDoc)
  r.HandleFunc("/delete_doc_u89xe/{id}", deleteDoc)

  fmt.Printf("Running docs @ http://127.0.0.1:%s\n", office683_shared.DocsPort)

  err = http.ListenAndServe(fmt.Sprintf(":%s", office683_shared.DocsPort), r)
  if err != nil {
    fmt.Println(err)
    panic(err)
  }
}


func ErrorPage(w http.ResponseWriter, err error) {
	type Context struct {
		Msg template.HTML
	}
	msg := fmt.Sprintf("%+v", err)
	fmt.Println(msg)
	msg = strings.ReplaceAll(msg, "\n", "<br>")
	msg = strings.ReplaceAll(msg, " ", "&nbsp;")
	msg = strings.ReplaceAll(msg, "\t", "&nbsp;&nbsp;")
	tmpl := template.Must(template.ParseFS(content, "templates/error.html"))
	tmpl.Execute(w, Context{template.HTML(msg)})
}
