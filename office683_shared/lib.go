package office683_shared

import (
  "os"
  "path/filepath"
  "fmt"
  "strings"
  "html/template"
  "net/http"
  "strconv"
  "runtime"
  "github.com/pkg/errors"
  "github.com/gookit/color"
  "github.com/saenuma/flaarum"
  "github.com/saenuma/zazabul"
  "github.com/saenuma/flaarum/flaarum_shared"
)


const O6_CONF_TEMPLATE = `// name of the company that created this office tools information
company_name: Test1

// logo of the company on your website.
company_logo: https://sae.ng/static/logo.png

// admin_pass is the password used by all admins
admin_pass:

// admin_email is for contacting the admin to get access
admin_email: admin@admin.com

// domain must be set after launching your server
domain:

`


func GetRootPath() (string, error) {
  var dd string
  if runtime.GOOS == "windows" {
    userHomeDir, err := os.UserHomeDir()
    if err != nil {
      return "", err
    }
    dd = filepath.Join(userHomeDir, "Office683")
  } else {
    dd = "/var/lib/office683"
  }

  os.MkdirAll(dd, 0777)
  return dd, nil
}


func DoesPathExists(p string) bool {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return false
	}
	return true
}


func GetFlaarumClient() flaarum.Client {
  var keyStr string
	inProd := flaarum_shared.GetSetting("in_production")
	if inProd == "" {
		color.Red.Println("unexpected error. Have you installed  and launched flaarum?")
		os.Exit(1)
	}
	if inProd == "true"{
		keyStrPath := flaarum_shared.GetKeyStrPath()
		raw, err := os.ReadFile(keyStrPath)
		if err != nil {
			color.Red.Println(err)
			os.Exit(1)
		}
		keyStr = string(raw)
	} else {
		keyStr = "not-yet-set"
	}
	port := flaarum_shared.GetSetting("port")
	if port == "" {
		color.Red.Println("unexpected error. Have you installed  and launched flaarum?")
		os.Exit(1)
	}
	var cl flaarum.Client

	portInt, err := strconv.Atoi(port)
	if err != nil {
		color.Red.Println("Invalid port setting.")
		os.Exit(1)
	}

	if portInt != flaarum_shared.PORT {
		cl = flaarum.NewClientCustomPort("127.0.0.1", keyStr, "first_proj", portInt)
	} else {
		cl = flaarum.NewClient("127.0.0.1", keyStr, "first_proj")
	}

  err = cl.Ping()
  if err != nil {
    panic(err)
  }

  return cl
}


func GetInstallationConfig() (zazabul.Config, error) {
  rootPath, err := GetRootPath()
  if err != nil {
    return zazabul.Config{}, err
  }

  confPath := filepath.Join(rootPath, "install.zconf")
  if ! DoesPathExists(confPath) {
    conf, err := zazabul.ParseConfig(O6_CONF_TEMPLATE)
    if err != nil {
      panic(err)
    }
    conf.Update(map[string]string {
  		"admin_pass": GenerateSecureRandomString(50),
  	})
    conf.Write(confPath)
  }

  conf, err := zazabul.LoadConfigFile(confPath)
  if err != nil {
    return zazabul.Config{}, errors.New(fmt.Sprintf("The file '%s' cannot be loaded.", confPath))
  }

  for _, item := range conf.Items {
    if item.Value == "" {
      if item.Name != "domain" {
        return zazabul.Config{}, errors.New("Every field in the launch file is compulsory.")
      }
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
