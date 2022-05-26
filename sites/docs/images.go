package docs

import (
  "fmt"
  "strconv"
  "time"
  "path/filepath"
  "io"
  "os"
  "net/http"
  "html/template"
  "github.com/pkg/errors"
  "github.com/gorilla/mux"
  "github.com/bankole7782/office683/office683_shared"
)


func gallery(w http.ResponseWriter, r *http.Request) {
  flaarumClient := office683_shared.GetFlaarumClient()
  rows, err := flaarumClient.Search(`
    table: docs_images
    order_by: upload_dt desc
    limit: 100
    `)
  if err != nil {
    ErrorPage(w, errors.Wrap(err, "flaarum search error"))
    return
  }

  type Context struct {
    ImageDescs []map[string]interface{}
  }

  tmpl := template.Must(template.ParseFS(content, "templates/gallery.html"))
  tmpl.Execute(w, Context{*rows})
}


func uploadImage(w http.ResponseWriter, r *http.Request) {
  rootPath, err := office683_shared.GetRootPath()
  if err != nil {
    panic(err)
  }

  flaarumClient := office683_shared.GetFlaarumClient()

  // Maximum upload of 10 MB files
	r.ParseMultipartForm(10000 << 20)

	file, handler, err := r.FormFile("image")
	if err != nil {
		fmt.Fprintf(w, "not_ok")
		fmt.Println(err)
		return
	}
	defer file.Close()

	rawFile, err := io.ReadAll(file)
	if err != nil {
		fmt.Fprintf(w, "not_ok")
		fmt.Println(err)
		return
	}

  ext := filepath.Ext(handler.Filename)

  retId, err := flaarumClient.InsertRowAny("docs_images", map[string]interface{} {
    "original_name": handler.Filename[:len(handler.Filename)-4],
    "upload_dt": time.Now(),
    "format": ext[1:],
  })

  if err != nil {
    fmt.Fprintf(w, "not_ok")
    fmt.Println(err)
    return
  }

	err = os.WriteFile(filepath.Join(rootPath, "docs_images", strconv.FormatInt(retId, 10) + ext), rawFile, 0777)
	if err != nil {
		fmt.Fprintf(w, "not_ok")
		fmt.Println(err)
		return
	}

  http.Redirect(w, r, "/docs_images/", 307)
}


func getImage(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  imageName := vars["img"]

  rootPath, err := office683_shared.GetRootPath()
  if err != nil {
    panic(err)
  }

  http.ServeFile(w, r, filepath.Join(rootPath, "docs_images", imageName))
}
