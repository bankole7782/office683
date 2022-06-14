package office683_shared

import (
  "os"
  "path/filepath"
  "fmt"
  "github.com/pkg/errors"
  "math/rand"
  "time"
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
