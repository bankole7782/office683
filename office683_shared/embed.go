package office683_shared

import (
  "embed"
)

//go:embed statics/*
var ContentStatics embed.FS

//go:embed templates/*
var Content embed.FS

//go:embed flaarum_stmts/*
var FlaarumStmts embed.FS
