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
		gen-form
			Generate a form for a model
		gen-migration
			Generate a migration file
		apply-migration
			Apply a migration file
		gen-schema
			Generate a schema file
		gen-encryption-key
			Generates a 32 byte encryption key for use with AES-256
`
}

var cmdFuncs = map[string]func(){
	"new":                cmds.GenProject,
	"gen-feature":        cmds.GenFeature,
	"gen-models":         cmds.GenModels,
	"gen-repository":     cmds.GenRepository,
	"gen-form":           cmds.GenForm,
	"gen-migration":      cmds.GenMigration,
	"apply-migration":    cmds.ApplyMigration,
	"gen-schema":         cmds.GenSchema,
	"gen-encryption-key": cmds.GenEncryptionKey,
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
