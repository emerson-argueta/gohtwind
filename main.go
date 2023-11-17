package main

import (
	"fmt"
	"gohtwind/cmds"
	"os"
)

func usageString() string {
	return `
Usage: gohtwind <command> [options]
	Commands:
		new
			Generate a new project
		gen-feature
			Generate a new feature
		gen-models
			Generate models from a database
		gen-repository
			Generate a repository for a feature
`
}

var cmdFuncs = map[string]func(){
	"new":            cmds.GenProject,
	"gen-feature":    cmds.GenFeature,
	"gen-models":     cmds.GenModels,
	"gen-repository": cmds.GenRepository,
	"gen-form":       cmds.GenForm,
	"gen-migration":  cmds.GenMigration,
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(usageString())
		os.Exit(1)
	}
	cmd := os.Args[1]
	if f, ok := cmdFuncs[cmd]; !ok {
		fmt.Println(usageString())
		os.Exit(1)
	} else {
		f()
	}
}
