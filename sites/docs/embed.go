package docs

import (
  "embed"
)

//go:embed templates/*
var content embed.FS
