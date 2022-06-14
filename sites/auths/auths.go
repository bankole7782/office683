package auths

import (
  "net/http"
  "html/template"
)


func registerUser(w http.ResponseWriter, r *http.Request) {
  tmpl := template.Must(template.ParseFS(content, "templates/register_user.html"))
  tmpl.Execute(w, nil)
}
