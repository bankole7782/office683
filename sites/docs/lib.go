package docs

import (
  "github.com/gorilla/mux"
)


func AddHandlers(r *mux.Router) {
  // documents

  r.HandleFunc("/docs/", allFolders)
  r.HandleFunc("/docs/{id}", docsOfFolder)
  r.HandleFunc("/new_doc", newDocument)
  r.HandleFunc("/new_folder", newFolder)
  r.HandleFunc("/update_doc/{id}", updateDoc)
  r.HandleFunc("/save_doc/{id}", saveDoc)
  r.HandleFunc("/doc/{id}", viewRenderedDoc)

  // gallery
  r.HandleFunc("/docs_images/", gallery)
  r.HandleFunc("/upload_image", uploadImage)
  r.HandleFunc("/gi/{img}", getImage)
}
