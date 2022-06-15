package cab

import (
  "net/http"
  "strconv"
  "fmt"
  "time"
  "io"
  "os"
  "strings"
  "path/filepath"
  "html/template"
  "github.com/pkg/errors"
  "github.com/gorilla/mux"
  "github.com/bankole7782/office683/office683_shared"
)


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
    _, err := flaarumClient.InsertRowAny("cab_folders", map[string]interface{} {
      "folder_name": r.FormValue("folder_name"),
      "teamid": teamIdInt,
      "desc": r.FormValue("desc"),
    })

    if err != nil {
      office683_shared.ErrorPage(w, errors.Wrap(err, "flaarum insert error"))
      return
    }

    r.Method = http.MethodGet
    http.Redirect(w, r, "/cab/", 307)
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
      table: cab_folders expand
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


func filesOfFolder(w http.ResponseWriter, r *http.Request) {
  status, userDetails := office683_shared.IsLoggedInUser(r)
  if status == false {
    http.Redirect(w, r, "/", 307)
    return
  }

  vars := mux.Vars(r)
  folderId := vars["id"]

  flaarumClient := office683_shared.GetFlaarumClient()

  // rootPath, err := office683_shared.GetRootPath()
  // if err != nil {
  //   panic(err)
  // }

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

  teamsFolders := make([]map[string]any, 0)
  for _, row := range *teamMembersRows {
    folderRows, err := flaarumClient.Search(fmt.Sprintf(`
      table: cab_folders expand
      order_by: folder asc
      where:
        teamid = %d
      `, row["teamid"].(int64)))
    if err != nil {
      office683_shared.ErrorPage(w, errors.Wrap(err, "flaarum search error"))
      return
    }

    for _, frow := range *folderRows {
      teamsFolders = append(teamsFolders, map[string]any {
        "folder_name": frow["folder_name"], "team_name": frow["teamid.team_name"],
        "folderid": frow["id"], "teamid": frow["teamid"],
      })
    }
  }

  filesRows, err := flaarumClient.Search(fmt.Sprintf(`
    table: cab_files
    order_by: upload_dt desc
    where:
      folderid = %s
    `, folderId))
  if err != nil {
    office683_shared.ErrorPage(w, errors.Wrap(err, "flaarum search error"))
    return
  }

  type Context struct {
    Folders []map[string]any
    FolderName string
    TeamName string
    Files []map[string]any
  }

  tmpl := template.Must(template.ParseFS(content, "templates/files.html"))
  tmpl.Execute(w, Context{teamsFolders, (*folderRow)["folder_name"].(string),
    (*folderRow)["teamid.team_name"].(string), *filesRows})
}


func createWrittenName() string {
  flaarumClient := office683_shared.GetFlaarumClient()

  for {
    rs := office683_shared.UntestedRandomString(100)
    count, err := flaarumClient.CountRows(fmt.Sprintf(`
      table: cab_files
      where:
        written_filename = '%s'
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


func uploadFile(w http.ResponseWriter, r *http.Request) {
  rootPath, err := office683_shared.GetRootPath()
  if err != nil {
    panic(err)
  }

  flaarumClient := office683_shared.GetFlaarumClient()

  // Maximum upload of 10 MB files
	r.ParseMultipartForm(10000 << 20)

	file, handler, err := r.FormFile("file")
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
  parts := strings.Split(r.FormValue("team_folder"), "-")
  folderIdInt, _ := strconv.ParseInt(parts[1], 10, 64)
  writtenFileName := createWrittenName()

  _, err = flaarumClient.InsertRowAny("cab_files", map[string]any {
    "original_name": handler.Filename[:len(handler.Filename)-len(ext)],
    "upload_dt": time.Now(),
    "format": ext[1:],
    "folderid": folderIdInt,
    "written_filename": writtenFileName,
  })

  if err != nil {
    fmt.Fprintf(w, "not_ok")
    fmt.Println(err)
    return
  }

	err = os.WriteFile(filepath.Join(rootPath, "cab", writtenFileName + ext), rawFile, 0777)
	if err != nil {
		fmt.Fprintf(w, "not_ok")
		fmt.Println(err)
		return
	}

  http.Redirect(w, r, "/cab/" + parts[1], 307)
}


func getFile(w http.ResponseWriter, r *http.Request) {
  status, _ := office683_shared.IsLoggedInUser(r)
  if status == false {
    http.Redirect(w, r, "/", 307)
    return
  }

  vars := mux.Vars(r)
  fileName := vars["name"]
  flaarumClient := office683_shared.GetFlaarumClient()
  ext := filepath.Ext(fileName)
  fileNameNoExt := fileName[:len(fileName)-len(ext)]

  fileRow, err := flaarumClient.SearchForOne(fmt.Sprintf(`
    table: cab_files
    where:
      written_filename = '%s'
    `, fileNameNoExt))
  if err != nil {
    office683_shared.ErrorPage(w, errors.Wrap(err, "os error"))
    return
  }

  originalFileName := (*fileRow)["original_name"].(string)

  rootPath, err := office683_shared.GetRootPath()
  if err != nil {
    office683_shared.ErrorPage(w, errors.Wrap(err, "os error"))
    return
  }
  rawObj, err := os.ReadFile(filepath.Join(rootPath, "cab", fileName))
  if err != nil {
    office683_shared.ErrorPage(w, errors.Wrap(err, "os error"))
    return
  }

  w.Header().Set("Content-Disposition", "attachment; filename=" + originalFileName + ext)
  if ext == ".js" {
    w.Header().Set("Content-Type", "text/javascript")
  } else {
    contentType := http.DetectContentType(rawObj)
    w.Header().Set("Content-Type", contentType)
  }
  w.Write(rawObj)
}
