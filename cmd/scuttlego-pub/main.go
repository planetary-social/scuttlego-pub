package main

import (
	"fmt"
	"os"

	"github.com/boreq/guinea"
)

var rootCommand = guinea.Command{
	Run: nil,
	Subcommands: map[string]*guinea.Command{
		"run":  &runCommand,
		"init": &initCommand,
	},
	Options:          nil,
	Arguments:        nil,
	ShortDescription: "",
	Description:      "",
}

func main() {
	if err := guinea.Run(&rootCommand); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
