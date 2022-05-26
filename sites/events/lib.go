package events


import (
  "net/http"
  "github.com/gorilla/mux"
  "fmt"
  "html/template"
  "strings"
)


func AddHandlers(r *mux.Router) {
  r.HandleFunc("/events/", allEvents)
  r.HandleFunc("/new_event", newEvent)
  r.HandleFunc("/event/{id}", aEvent)
  r.HandleFunc("/delete_event/{id}", deleteAEvent)
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
