package main

import (
  "fmt"
  "os/exec"
  "github.com/bankole7782/office683/office683_shared"
)

func main() {
  fmt.Println("Starting ssl-proxy.")

  command := "ssl-proxy-linux-amd64"

  conf, err := office683_shared.GetInstallationConfig()
  if err != nil {
    panic(err)
  }

  var c *exec.Cmd
  if conf.Get("domain") != "" {
    fmt.Println("Domain name set. Trying letsencrypt")
    c = exec.Command(command, "-from", "0.0.0.0:443", "-to", "127.0.0.1:8387",
      "-domain", conf.Get("domain"))
  } else {
    fmt.Println("Domain name not set. Using self-signed certificate")
    c = exec.Command(command, "-from", "0.0.0.0:443", "-to", "127.0.0.1:8387")
  }

  c.Run()
}
