package office683_shared

import (
  "os"
  "path/filepath"
  "fmt"
  "github.com/pkg/errors"
  "math/rand"
  "time"
  "strings"
  "html/template"
  "net/http"
  "os/exec"
  "github.com/saenuma/flaarum"
  "github.com/saenuma/zazabul"
)


func GetRootPath() (string, error) {
  hd, err := os.UserHomeDir()
  if err != nil {
    return "", errors.Wrap(err, "os error")
  }
  dd := os.Getenv("SNAP_COMMON")
  if dd == "/var/snap/go/common" || dd == "" {
    dd = filepath.Join(hd, "office683_data")
    os.MkdirAll(dd, 0777)
  }

  return dd, nil
}


func UntestedRandomString(length int) string {
  var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
  const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

  b := make([]byte, length)
  for i := range b {
    b[i] = charset[seededRand.Intn(len(charset))]
  }
  return string(b)
}


func DoesPathExists(p string) bool {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return false
	}
	return true
}


func GetFlaarumClient() flaarum.Client {
  conf, _ := GetInstallationConfig()
  keyStr := conf.Get("flaarum_keystr")
  flaarumClient := flaarum.NewClient("127.0.0.1", keyStr, "first_proj")

  err = flaarumClient.Ping()
  if err != nil {
    panic(err)
  }

  return flaarumClient
}


func GetInstallationConfig() (zazabul.Config, error) {
  rootPath, err := GetRootPath()
  if err != nil {
    return zazabul.Config{}, err
  }

  confPath := filepath.Join(rootPath, "install.zconf")

  conf, err := zazabul.LoadConfigFile(confPath)
  if err != nil {
    return zazabul.Config{}, errors.New(fmt.Sprintf("The file '%s' cannot be loaded.", confPath))
  }

  for _, item := range conf.Items {
    if item.Value == "" {
      return zazabul.Config{}, errors.New("Every field in the launch file is compulsory.")
    }
  }

  return conf, nil
}


func IsLoggedInUser(r *http.Request) (bool, map[string]any) {
  cookie, err := r.Cookie("thingy_thing")
  if err != nil {
    return false, nil
  }

  flaarumClient := GetFlaarumClient()

  count, err := flaarumClient.CountRows(fmt.Sprintf(`
    table: sessions
    where:
      keystr = '%s'
    `, cookie.Value))
  if err != nil {
    fmt.Println(err)
    return false, nil
  }

  if count == 0 {
    return false, nil
  }

  sessionRow, err := flaarumClient.SearchForOne(fmt.Sprintf(`
    table: sessions expand
    where:
      keystr = '%s'
    `, cookie.Value))
  if err != nil {
    fmt.Println(err)
    return false, nil
  }

  return true, map[string]any {
    "id": (*sessionRow)["userid"],
    "firstname": (*sessionRow)["userid.firstname"],
    "surname": (*sessionRow)["userid.surname"],
    "email": (*sessionRow)["userid.email"],
    "registration_dt": (*sessionRow)["userid.registration_dt"],
  }
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
	tmpl := template.Must(template.ParseFS(Content, "templates/error.html"))
	tmpl.Execute(w, Context{template.HTML(msg)})
}


func BelongToThisTeam(teamid, userid int64) bool {
  flaarumClient := GetFlaarumClient()

  count, err := flaarumClient.CountRows(fmt.Sprintf(`
    table: team_members
    where:
      teamid = %d
      and userid = %d
    `, teamid, userid))

  if err != nil {
    fmt.Println(err.Error())
    return false
  }

  if count > 0 {
    return true
  }

  return false
}
