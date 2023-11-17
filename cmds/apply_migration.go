package cmds

import (
	"flag"
	"fmt"
	"os"
)

func applyMigrationUsageString() string {
	return `
Usage: gohtwind apply-migration [options]
	Options:
		-file-name string (required)
			Name of the migration file to apply
		-adapter string (required)
			Database adapter to use
		-database-name string (required)
			Name of the database to apply the migration to (use the name of the database in the config/database.yml file)
	Info:
		* Applies the specified migration file
		* The migration file is expected to be in the db/migrations/<database_name> folder
		* The migration file is expected to be named <feature_name>_<table_name>.sql
		* The migration file is expected to contain a create table statement or an alter table statement
`
}

type applyMigration struct {
	flagSet      *flag.FlagSet
	fileName     *string
	adapter      *string
	databaseName *string
	projectPath  string
}

func newApplyMigration() *applyMigration {
	applyMigrationFlags := flag.NewFlagSet("gohtwind apply-migration", flag.ExitOnError)
	pjp, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fileName := applyMigrationFlags.String("file-name", "", "Name of the migration file to apply")
	adapter := applyMigrationFlags.String("adapter", "", "Database adapter to use")
	databaseName := applyMigrationFlags.String("database-name", "", "Name of the database to apply the migration to (use the name of the database in the config/database.yml file)")
	return &applyMigration{
		flagSet:      applyMigrationFlags,
		fileName:     fileName,
		adapter:      adapter,
		databaseName: databaseName,
		projectPath:  pjp,
	}
}

func ApplyMigration(args []string) {
	m := newApplyMigration()
	if len(os.Args) < 2 {
		fmt.Println(applyMigrationUsageString())
		os.Exit(1)
	}
	m.flagSet.Parse(os.Args[2:])
	if *m.fileName == "" || *m.adapter == "" || *m.databaseName == "" {
		fmt.Println(applyMigrationUsageString())
		os.Exit(1)
	}
	m.applyMigration()
}

func (m *applyMigration) applyMigration() {
	// TODO: Implement
}
