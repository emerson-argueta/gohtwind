package cmds

import (
	"flag"
	"fmt"
	"os"
)

func genSchemaUsageString() string {
	return `
Usage: gohtwind gen-schema [options]
	Options:
		-database-name string (required)
			Name of the database to generate a schema for
		-adapter string (required)
			Name of the database to apply the migration to (use the name of the database in the config/database.yml file)
	Info:
		* Generates a schema for the specified database 
		* The schema is generated in the db/schema folder	
		* The schema is named schema_<database_name>.sql
		* The schema contains ddl statements for creating the database and its tables 
`
}

type schema struct {
	flagSet      *flag.FlagSet
	databaseName *string
	adapter      *string
	projectPath  string
}

func newSchema() *schema {
	genSchemaFlags := flag.NewFlagSet("gohtwind gen-schema", flag.ExitOnError)
	pjp, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	databaseName := genSchemaFlags.String("database-name", "", "Name of the database to generate a schema for")
	adapter := genSchemaFlags.String("adapter", "", "Database adapter to use")
	return &schema{
		flagSet:      genSchemaFlags,
		databaseName: databaseName,
		adapter:      adapter,
		projectPath:  pjp,
	}
}

func GenSchema(args []string) {
	s := newSchema()
	if len(os.Args) < 2 {
		fmt.Println(genSchemaUsageString())
		os.Exit(1)
	}
	s.flagSet.Parse(os.Args[2:])
	if *s.databaseName == "" || *s.adapter == "" {
		fmt.Println(genSchemaUsageString())
		os.Exit(1)
	}
	s.genSchema()
}

func (s *schema) genSchema() {
	// TODO: Implement
}
