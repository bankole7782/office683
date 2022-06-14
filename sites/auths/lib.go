package auths

import (
  "net/http"
  "github.com/gorilla/mux"
  "fmt"
  "html/template"
  "strings"
  "github.com/bankole7782/office683/office683_shared"
)


func AddHandlers(r *mux.Router) {
  // documents
  r.HandleFunc("/register", registerUser)
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
	tmpl := template.Must(template.ParseFS(office683_shared.Content, "templates/error.html"))
	tmpl.Execute(w, Context{template.HTML(msg)})
}
