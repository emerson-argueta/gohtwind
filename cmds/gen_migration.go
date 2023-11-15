package cmds

import (
	"flag"
	"fmt"
	"os"
)

func genMigrationUsageString() string {
	return `
Usage: gohtwind gen-migration [options]
	Options:
		-feature-name string (optional)
			Prefixes the table name with the feature name 
		-table-name string (required)
			What the table is named 
		-action string (required)
			What the migration does (create, alter, drop)
		-adapter string (required)
			Database adapter to use
		-database-name string (required)
			Name of the database to apply the migration to (use the name of the database in the config/database.yml file)
	Info:
		* Generates a migration file for the specified model
		* The migration file is generated in the db/migrations/<database_name> folder
		* The migration file is named <feature_name>_<table_name>.sql
		* The migration file contains a create table statement, an alter table statement, or a drop table statement
`
}

type migration struct {
	flagSet      *flag.FlagSet
	featName     *string
	tableName    *string
	action       *string
	adapter      *string
	databaseName *string
	projectPath  string
}

func newMigration() *migration {
	genMigrationFlags := flag.NewFlagSet("gohtwind gen-migration", flag.ExitOnError)
	pjp, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	featName := genMigrationFlags.String("feature-name", "", "Prefixes the table name with the feature name")
	tableName := genMigrationFlags.String("table-name", "", "What the table is named")
	action := genMigrationFlags.String("action", "", "What the migration does (create, alter, drop)")
	adapter := genMigrationFlags.String("adapter", "", "Database adapter to use")
	databaseName := genMigrationFlags.String("database-name", "", "Name of the database to apply the migration to (use the name of the database in the config/database.yml file)")
	return &migration{
		flagSet:      genMigrationFlags,
		featName:     featName,
		tableName:    tableName,
		action:       action,
		adapter:      adapter,
		databaseName: databaseName,
		projectPath:  pjp,
	}
}

func GenMigration() {
	migration := newMigration()
	if len(os.Args) < 2 {
		fmt.Println(genMigrationUsageString())
		os.Exit(1)
	}
	migration.flagSet.Parse(os.Args[2:])
	if *migration.tableName == "" || *migration.action == "" || *migration.adapter == "" || *migration.databaseName == "" {
		fmt.Println(genMigrationUsageString())
		os.Exit(1)
	}
	migration.genMigration()
}

func (m *migration) genMigration() {
	// TODO: Implement
	/* Generates a migration file for the specified model
	* The migration file is generated in the db/migrations/<database_name> folder
	* The migration file is named <feature_name>_<table_name>.sql
	* The migration file contains a create table statement, an alter table statement, or a drop table statement
	 */
	fn := fmt.Sprintf("%s_%s.sql", *m.featName, *m.tableName)
	fp := fmt.Sprintf("%s/db/migrations/%s/%s", m.projectPath, *m.databaseName, fn)
	f, err := os.Create(fp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	switch *m.action {
	case "create":
		f.WriteAt(m.CreateTableStatement(), 0)
	case "alter":
		f.WriteAt(m.AlterTableStatement(), 0)
	case "drop":
		f.WriteAt(m.DropTableStatement(), 0)
	default:
		fmt.Fprintf(os.Stderr, "Error: %v\n", "Invalid action")
	}
}

func (m *migration) CreateTableStatement() []byte {
	// TODO: Implement
	return []byte("")
}

func (m *migration) AlterTableStatement() []byte {
	// TODO: Implement
	return []byte("")
}

func (m *migration) DropTableStatement() []byte {
	// TODO: Implement
	return []byte("")
}
