package docs

import (
  "net/http"
  "html/template"
  "time"
  "strconv"
  "fmt"
  "os"
  "strings"
  "path/filepath"
  "github.com/pkg/errors"
  "github.com/gorilla/mux"
  "github.com/russross/blackfriday/v2"
  "github.com/dustin/go-humanize"
  "github.com/bankole7782/office683/office683_shared"
)


func newDocument(w http.ResponseWriter, r *http.Request) {
  status, _ := office683_shared.IsLoggedInUser(r)
  if status == false {
    http.Redirect(w, r, "/", 307)
    return
  }

  if r.Method == http.MethodPost {
    flaarumClient := office683_shared.GetFlaarumClient()
    parts := strings.Split(r.FormValue("team_folder"), "-")
    folderIdInt, _ := strconv.ParseInt(parts[1], 10, 64)
    retId, err := flaarumClient.InsertRowAny("docs", map[string]interface{} {
      "doc_title": r.FormValue("doc_title"),
      "folderid": folderIdInt,
      "update_dt": time.Now(),
      "public": false,
    })

    if err != nil {
      office683_shared.ErrorPage(w, errors.Wrap(err, "flaarum insert error"))
      return
    }

    http.Redirect(w, r, "/update_doc/" + strconv.FormatInt(retId, 10), 307)
  }
}



func newFolder(w http.ResponseWriter, r *http.Request) {
  status, _ := office683_shared.IsLoggedInUser(r)
  if status == false {
    http.Redirect(w, r, "/", 307)
    return
  }

  if r.Method == http.MethodPost {
    flaarumClient := office683_shared.GetFlaarumClient()
    if r.FormValue("teamid") == "" {
      office683_shared.ErrorPage(w, errors.New("The team must be selected."))
      return
    }
    teamIdInt, _ := strconv.ParseInt(r.FormValue("teamid"), 10, 64)
    _, err := flaarumClient.InsertRowAny("docs_folders", map[string]interface{} {
      "folder_name": r.FormValue("folder_name"),
      "teamid": teamIdInt,
      "desc": r.FormValue("desc"),
    })

    if err != nil {
      office683_shared.ErrorPage(w, errors.Wrap(err, "flaarum insert error"))
      return
    }

    r.Method = http.MethodGet
    http.Redirect(w, r, "/docs/", 307)
  }
}


func allFolders(w http.ResponseWriter, r *http.Request) {
  status, userDetails := office683_shared.IsLoggedInUser(r)
  if status == false {
    http.Redirect(w, r, "/", 307)
    return
  }

  flaarumClient := office683_shared.GetFlaarumClient()

  teamMembersRows, err := flaarumClient.Search(fmt.Sprintf(`
    table: team_members expand
    order_by: teamid.team_name asc
    fields: teamid.team_name teamid
    where:
      userid = %d
    `, userDetails["id"]))
  if err != nil {
    office683_shared.ErrorPage(w, errors.Wrap(err, "flaarum error"))
    return
  }


  yourTeams := make([]map[string]any, 0)
  teamsToFolders := make(map[string][]map[string]any)
  teamsFolders := make([]map[string]any, 0)
  for _, row := range *teamMembersRows {
    elem := map[string]any {
      "teamid": row["teamid"],
      "team_name": row["teamid.team_name"],
    }
    yourTeams = append(yourTeams, elem)

    folderRows, err := flaarumClient.Search(fmt.Sprintf(`
      table: docs_folders expand
      order_by: folder asc
      where:
        teamid = %d
      `, row["teamid"].(int64)))
    if err != nil {
      office683_shared.ErrorPage(w, errors.Wrap(err, "flaarum search error"))
      return
    }

    teamName := row["teamid.team_name"].(string)
    teamsToFolders[teamName] = *folderRows
    for _, frow := range *folderRows {
      teamsFolders = append(teamsFolders, map[string]any {
        "folder_name": frow["folder_name"], "team_name": frow["teamid.team_name"],
        "folderid": frow["id"], "teamid": frow["teamid"],
      })
    }
  }

  type Context struct {
    YourTeams []map[string]any
    HaveTeams bool
    TeamsToFolders map[string][]map[string]any
    Folders []map[string]any
  }

  tmpl := template.Must(template.ParseFS(content, "templates/all_folders.html"))
  tmpl.Execute(w, Context{yourTeams, len(yourTeams) > 0, teamsToFolders, teamsFolders})
}


func docsOfFolder(w http.ResponseWriter, r *http.Request) {
  status, userDetails := office683_shared.IsLoggedInUser(r)
  if status == false {
    http.Redirect(w, r, "/", 307)
    return
  }

  vars := mux.Vars(r)
  folderId := vars["id"]

  rootPath, err := office683_shared.GetRootPath()
  if err != nil {
    panic(err)
  }

  flaarumClient := office683_shared.GetFlaarumClient()

  docRows, err := flaarumClient.Search(fmt.Sprintf(`
    table: docs
    order_by: doc_title asc
    where:
      folderid = '%s'
    `, folderId))

  if err != nil {
    office683_shared.ErrorPage(w, errors.Wrap(err, "os error"))
    return
  }

  elems := make([]map[string]any, 0)
  for _, docRow := range *docRows {
    docPath := filepath.Join(rootPath, "docs", strconv.FormatInt(docRow["id"].(int64), 10) + ".md")
    docSize := "0B"
    if office683_shared.DoesPathExists(docPath) {

      file, err := os.Open(docPath)
      if err != nil {
        office683_shared.ErrorPage(w, errors.Wrap(err, "os error"))
        return
      }
      defer file.Close()

      stat, err := file.Stat()
      if err != nil {
        office683_shared.ErrorPage(w, errors.Wrap(err, "os error"))
        return
      }
      docSize = humanize.Bytes(uint64(stat.Size()))
    }

    elems = append(elems, map[string]any {
      "doc_title": docRow["doc_title"],
      "updated": docRow["update_dt"],
      "id": docRow["id"],
      "doc_size": docSize,
    })

  }

  folderRow, err := flaarumClient.SearchForOne(fmt.Sprintf(`
    table: docs_folders expand
    where:
      id = %s
    `, folderId))
  if err != nil {
    office683_shared.ErrorPage(w, errors.Wrap(err, "flaarum error"))
    return
  }


  teamMembersRows, err := flaarumClient.Search(fmt.Sprintf(`
    table: team_members expand
    order_by: teamid.team_name asc
    fields: teamid.team_name teamid
    where:
      userid = %d
    `, userDetails["id"]))
  if err != nil {
    office683_shared.ErrorPage(w, errors.Wrap(err, "flaarum error"))
    return
  }

  yourTeams := make([]map[string]any, 0)
  teamsToFolders := make(map[string][]map[string]any)
  teamsFolders := make([]map[string]any, 0)
  for _, row := range *teamMembersRows {
    elem := map[string]any {
      "teamid": row["teamid"],
      "team_name": row["teamid.team_name"],
    }
    yourTeams = append(yourTeams, elem)

    folderRows, err := flaarumClient.Search(fmt.Sprintf(`
      table: docs_folders expand
      order_by: folder asc
      where:
        teamid = %d
      `, row["teamid"].(int64)))
    if err != nil {
      office683_shared.ErrorPage(w, errors.Wrap(err, "flaarum search error"))
      return
    }

    teamName := row["teamid.team_name"].(string)
    teamsToFolders[teamName] = *folderRows
    for _, frow := range *folderRows {
      teamsFolders = append(teamsFolders, map[string]any {
        "folder_name": frow["folder_name"], "team_name": frow["teamid.team_name"],
        "folderid": frow["id"], "teamid": frow["teamid"],
      })
    }
  }

  type Context struct {
    Documents []map[string]any
    FolderName string
    TeamName string
    Folders []map[string]any
    YourTeams []map[string]any
    HaveTeams bool
  }

  tmpl := template.Must(template.ParseFS(content, "templates/all_docs.html"))
  tmpl.Execute(w, Context{elems, (*folderRow)["folder_name"].(string),
    (*folderRow)["teamid.team_name"].(string), teamsFolders, yourTeams, len(yourTeams)>0})
}


func updateDoc(w http.ResponseWriter, r *http.Request) {
  status, _ := office683_shared.IsLoggedInUser(r)
  if status == false {
    http.Redirect(w, r, "/", 307)
    return
  }

  vars := mux.Vars(r)
  docId := vars["id"]

  flaarumClient := office683_shared.GetFlaarumClient()
  docRow, err := flaarumClient.SearchForOne(fmt.Sprintf(`
    table: docs
    where:
      id = %s
    `, docId))
  if err != nil {
    office683_shared.ErrorPage(w, errors.Wrap(err, "flaarum search error"))
    return
  }

  docDetails := map[string]string {
    "doc_title": (*docRow)["doc_title"].(string),
    "updated": (*docRow)["update_dt"].(time.Time).String(),
  }

  rootPath, err := office683_shared.GetRootPath()
  if err != nil {
    panic(err)
  }
  rawDoc := ""
  docPath := filepath.Join(rootPath, "docs", docId + ".md")
  if office683_shared.DoesPathExists(docPath) {
    raw, err := os.ReadFile(docPath)
    if err != nil {
      office683_shared.ErrorPage(w, errors.Wrap(err, "os error"))
      return
    }
    rawDoc = string(raw)
  }

  type Context struct {
    DocDetails map[string]string
    RawDoc string
    DocId string
  }
  tmpl := template.Must(template.ParseFS(content, "templates/update_doc.html"))
  tmpl.Execute(w, Context{docDetails, rawDoc, docId})
}


func saveDoc(w http.ResponseWriter, r *http.Request) {
  status, _ := office683_shared.IsLoggedInUser(r)
  if status == false {
    http.Redirect(w, r, "/", 307)
    return
  }

  vars := mux.Vars(r)
  docId := vars["id"]

  rootPath, err := office683_shared.GetRootPath()
  if err != nil {
    panic(err)
  }

  docPath := filepath.Join(rootPath, "docs", docId + ".md")
  os.WriteFile(docPath, []byte(r.FormValue("raw_doc")), 0777)

  flaarumClient := office683_shared.GetFlaarumClient()
  err = flaarumClient.UpdateRowsAny(fmt.Sprintf(`
    table: docs
    where:
     id = %s
    `, docId), map[string]interface{} {
    "update_dt": time.Now(),
  })
  if err != nil {
    fmt.Println(err)
  }
  fmt.Fprintf(w, "ok")
}


func viewRenderedDoc(w http.ResponseWriter, r *http.Request) {
  status, _ := office683_shared.IsLoggedInUser(r)

  vars := mux.Vars(r)
  docId := vars["id"]

  rootPath, err := office683_shared.GetRootPath()
  if err != nil {
    panic(err)
  }

  docPath := filepath.Join(rootPath, "docs", docId + ".md")

  flaarumClient := office683_shared.GetFlaarumClient()
  docRow, err := flaarumClient.SearchForOne(fmt.Sprintf(`
    table: docs
    where:
      id = %s
    `, docId))
  if err != nil {
    office683_shared.ErrorPage(w, errors.Wrap(err, "flaarum search error"))
    return
  }

  if status == false && (*docRow)["public"] == false {
    http.Redirect(w, r, "/", 307)
    return
  }

  docDetails := map[string]string {
    "doc_title": (*docRow)["doc_title"].(string),
    "updated": (*docRow)["update_dt"].(time.Time).String(),
  }

  var rawDoc []byte
  if office683_shared.DoesPathExists(docPath) {
    rawDoc2, err := os.ReadFile(docPath)
    if err != nil {
      office683_shared.ErrorPage(w, errors.Wrap(err, "os error"))
      return
    }
    rawDoc = rawDoc2
  }

  renderedDocRaw := blackfriday.Run(rawDoc)
  renderedDoc := template.HTML(string(renderedDocRaw))

  type Context struct {
    RenderedDoc template.HTML
    DocDetails map[string]string
    DocId string
  }

  tmpl := template.Must(template.ParseFS(content, "templates/render_doc.html"))
  tmpl.Execute(w, Context{renderedDoc, docDetails, docId})
}


func deleteDoc(w http.ResponseWriter, r *http.Request) {
  status, _ := office683_shared.IsLoggedInUser(r)
  if status == false {
    http.Redirect(w, r, "/", 307)
    return
  }

  vars := mux.Vars(r)
  docId := vars["id"]

  rootPath, err := office683_shared.GetRootPath()
  if err != nil {
    panic(err)
  }

  flaarumClient := office683_shared.GetFlaarumClient()
  err = flaarumClient.DeleteRows(fmt.Sprintf(`
    table: docs
    where:
      id = %s
    `, docId))
  if err != nil {
    office683_shared.ErrorPage(w, errors.Wrap(err, "flaarum search error"))
    return
  }

  docPath := filepath.Join(rootPath, "docs", docId + ".md")
  os.Remove(docPath)

  fmt.Fprintf(w, "ok")
}
