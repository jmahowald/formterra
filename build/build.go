package main

import (
	"os"

	"github.com/backpackhealth/formterra/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	var dir string
	if len(os.Args) != 2 {
		dir = "."
	} else {
		dir = os.Args[1]
	}
	doc.GenMarkdownTree(cmd.RootCmd, dir)
}
