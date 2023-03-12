package main

import (
	"github.com/cidverse/cid-actions-go/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	status  = "clean"
)

func main() {
	cmd.Execute()
}
