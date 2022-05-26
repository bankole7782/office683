package office683_shared

import (
  "os"
  "path/filepath"
  // "strings"
  "github.com/pkg/errors"
  "math/rand"
  "time"
  "github.com/saenuma/flaarum"
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
  const charset = "abcdefghijklmnopqrstuvwxyz1234567890"

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
