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
      office683_shared.ErrorPage(w, errors.Wrap(err, "zazabul error"))
      return
    }

    if r.FormValue("apassword") != conf.Get("admin_pass") {
      office683_shared.ErrorPage(w, errors.New("The admin password is invalid."))
      return
    }

    if r.FormValue("password") != r.FormValue("cpassword") {
      office683_shared.ErrorPage(w, errors.New("The two passwords don't match"))
      return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)
    if err != nil {
      office683_shared.ErrorPage(w, errors.Wrap(err, "bcrypt error"))
      return
    }

    flaarumClient := office683_shared.GetFlaarumClient()

    _, err = flaarumClient.InsertRowAny("users", map[string]any {
      "firstname": r.FormValue("firstname"), "surname": r.FormValue("surname"),
      "email": r.FormValue("email"), "registration_dt": time.Now(),
      "password": string(hashedPassword),
    })

    if err != nil {
      office683_shared.ErrorPage(w, errors.Wrap(err, "flaarum error"))
      return
    }

    tmpl := template.Must(template.ParseFS(content, "templates/done_registration.html"))
    tmpl.Execute(w, nil)
  }
}


func createSessionCode() string {
  flaarumClient := office683_shared.GetFlaarumClient()

  for {
    rs := office683_shared.UntestedRandomString(100)
    count, err := flaarumClient.CountRows(fmt.Sprintf(`
      table: sessions
      where:
        keystr = '%s'
      `, rs))
    if err != nil {
      fmt.Println(err.Error())
      return ""
    }

    if count == 0 {
      return rs
    }
  }
}


func signInHandler(w http.ResponseWriter, r *http.Request) {
  flaarumClient := office683_shared.GetFlaarumClient()

  if r.Method == http.MethodPost {
    email := r.FormValue("email")
    password := r.FormValue("password")

    userRow, err := flaarumClient.SearchForOne(fmt.Sprintf(`
      table: users
      where:
        email = '%s'
      `, email))
    if err != nil {
      office683_shared.ErrorPage(w, errors.Wrap(err, "flaarum error"))
      return
    }

    databasePassword := (*userRow)["password"].(string)
    err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
    if err != nil {
      office683_shared.ErrorPage(w, errors.New("Invalid Credentials"))
      return
    }

    now := time.Now()
    sessionCode := createSessionCode()

    _, err = flaarumClient.InsertRowAny("sessions", map[string]any {
      "keystr": sessionCode, "creation_dt": time.Now(),
      "userid": (*userRow)["id"].(int64),
    })
    if err != nil {
      office683_shared.ErrorPage(w, errors.Wrap(err, "flaarum error"))
      return
    }

    expires := now.Add(time.Hour * 24)
    cookie := &http.Cookie {
      Name: "thingy_thing",
      Value: sessionCode,
      Path: "/",
      Expires: expires,
    }
    http.SetCookie(w, cookie)

    http.Redirect(w, r, "/programs", 307)
  }
}


func signout(w http.ResponseWriter, r *http.Request) {
  cookie := &http.Cookie {
    Name: "nasc_thing",
    Value: "",
    Path: "/",
    MaxAge: -1,
  }
  http.SetCookie(w, cookie)
  http.Redirect(w, r, "/", 302)
}
