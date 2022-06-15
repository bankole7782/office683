package cab

import (
  "github.com/gorilla/mux"
)


func AddHandlers(r *mux.Router) {
  r.HandleFunc("/cab/", allFolders)
  r.HandleFunc("/cab/{id}", filesOfFolder)
  r.HandleFunc("/cab_new_folder", newFolder)
  r.HandleFunc("/cab_upload_file", uploadFile)
  r.HandleFunc("/gcf/{name}", getFile)

}
