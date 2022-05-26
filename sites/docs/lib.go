package docs

import (
  "net/http"
  "github.com/gorilla/mux"
  "fmt"
  "html/template"
  "strings"
)


func AddHandlers(r *mux.Router) {
  // documents

  r.HandleFunc("/docs/", allDocs)
  r.HandleFunc("/new_doc", newDocument)
  r.HandleFunc("/update_doc/{id}", updateDoc)
  r.HandleFunc("/save_doc/{id}", saveDoc)
  r.HandleFunc("/doc/{id}", viewRenderedDoc)
  r.HandleFunc("/delete_doc_u89xe/{id}", deleteDoc)

  // gallery
  r.HandleFunc("/docs_images/", gallery)
  r.HandleFunc("/upload_image", uploadImage)
  r.HandleFunc("/gi/{img}", getImage)
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
