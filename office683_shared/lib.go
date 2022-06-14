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
  "github.com/saenuma/flaarum"
  "github.com/saenuma/zazabul"
)


func GetRootPath() (string, error) {
  var rootPath string

	devCheckStr := os.Getenv("OFFICE683_DEV")
  if devCheckStr == "true" {
    hd, err := os.UserHomeDir()
  	if err != nil {
  		return "", errors.Wrap(err, "os error")
  	}
    rootPath = filepath.Join(hd, "office683_data")
  } else {
    rootPath = "/office683"
  }

  err := os.MkdirAll(rootPath, 0777)
  if err != nil {
    panic(err)
  }

	return rootPath, nil
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
  devCheckStr := os.Getenv("OFFICE683_DEV")
  var flaarumClient flaarum.Client
  if devCheckStr == "true" {
    flaarumClient = flaarum.NewClient("127.0.0.1", "not-set", "first_proj")
  } else {
    flaarumClient = flaarum.NewClient("127.0.0.1", "not-set", "first_proj")
  }

  err := flaarumClient.Ping()
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
