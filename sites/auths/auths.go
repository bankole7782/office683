package auths

import (
  "fmt"
  "time"
  "net/http"
  "html/template"
  "github.com/pkg/errors"
  "golang.org/x/crypto/bcrypt"
  "github.com/bankole7782/office683/office683_shared"
)


func registerUser(w http.ResponseWriter, r *http.Request) {
  if r.Method == http.MethodGet {
    tmpl := template.Must(template.ParseFS(content, "templates/register_user.html"))
    tmpl.Execute(w, nil)
  } else {
    conf, err := office683_shared.GetInstallationConfig()
    if err != nil {
      ErrorPage(w, errors.Wrap(err, "zazabul error"))
      return
    }

    if r.FormValue("apassword") != conf.Get("admin_pass") {
      ErrorPage(w, errors.New("The admin password is invalid."))
      return
    }

    if r.FormValue("password") != r.FormValue("cpassword") {
      ErrorPage(w, errors.New("The two passwords don't match"))
      return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)
    if err != nil {
      ErrorPage(w, errors.Wrap(err, "bcrypt error"))
      return
    }

    fmt.Println(hashedPassword)
    fmt.Println(string(hashedPassword))
    flaarumClient := office683_shared.GetFlaarumClient()

    _, err = flaarumClient.InsertRowAny("users", map[string]any {
      "firstname": r.FormValue("firstname"), "surname": r.FormValue("surname"),
      "email": r.FormValue("email"), "registration_dt": time.Now(),
      "password": string(hashedPassword),
    })

    if err != nil {
      ErrorPage(w, errors.Wrap(err, "flaarum error"))
      return
    }

    tmpl := template.Must(template.ParseFS(content, "templates/done_registration.html"))
    tmpl.Execute(w, nil)
  }
}
