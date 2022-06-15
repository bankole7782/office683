package main

import (
  "net/http"
  "fmt"
  "path/filepath"
  "os"
  "html/template"
  "github.com/gorilla/mux"
  "github.com/bankole7782/office683/office683_shared"
  "github.com/bankole7782/office683/sites/auths"
  "github.com/bankole7782/office683/sites/docs"
  "github.com/bankole7782/office683/sites/events"
  "github.com/bankole7782/office683/sites/cab"
)


func main() {
  rootPath, err := office683_shared.GetRootPath()
  if err != nil {
    panic(err)
  }

  os.MkdirAll(filepath.Join(rootPath, "cab"), 0777)
  os.MkdirAll(filepath.Join(rootPath, "docs"), 0777)
  os.MkdirAll(filepath.Join(rootPath, "docs_images"), 0777)

  conf, err := office683_shared.GetInstallationConfig()
  if err != nil {
    panic(err)
  }

  r := mux.NewRouter()

  r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    type Context struct {
      CompanyName string
      CompanyLogoPath string
      AdminEmail string
      Msg string
    }

    ctx := Context {
      CompanyName: conf.Get("company_name"), CompanyLogoPath: conf.Get("company_logo"),
      AdminEmail: conf.Get("admin_email"),
    }

    tmpl := template.Must(template.ParseFS(office683_shared.Content, "templates/home.html"))
    tmpl.Execute(w, ctx)
  })

  r.HandleFunc("/programs", func(w http.ResponseWriter, r *http.Request) {
    status, userDetails := office683_shared.IsLoggedInUser(r)
    if status == false {
      http.Redirect(w, r, "/", 307)
      return
    }
    type Context struct {
      Fullname string
    }
    tmpl := template.Must(template.ParseFS(office683_shared.Content, "templates/programs.html"))
    tmpl.Execute(w, Context{userDetails["firstname"].(string) + " " + userDetails["surname"].(string)})
  })

  r.HandleFunc("/gs/{obj}", func (w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    rawObj, err := office683_shared.ContentStatics.ReadFile("statics/" + vars["obj"])
    if err != nil {
      panic(err)
    }
    w.Header().Set("Content-Disposition", "attachment; filename=" + vars["obj"])
    ext := vars["obj"][len(vars["obj"])-3:]
    if ext == ".js" {
      w.Header().Set("Content-Type", "text/javascript")
    } else {
      contentType := http.DetectContentType(rawObj)
      w.Header().Set("Content-Type", contentType)
    }
    w.Write(rawObj)
  })

  // attach handlers
  docs.AddHandlers(r)
  events.AddHandlers(r)
  auths.AddHandlers(r)
  cab.AddHandlers(r)

  fmt.Println("Running office683 @ http://127.0.0.1:8080")

  err = http.ListenAndServe(":8080", r)
  if err != nil {
    fmt.Println(err)
    panic(err)
  }
}
