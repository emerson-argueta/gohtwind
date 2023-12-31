package cmds

import (
	"database/sql"
	"flag"
	"fmt"
	"gohtwind/utils"
	"log"
	"os"
	"time"
)

func genSchemaUsageString() string {
	return `
Usage: gohtwind gen-schema [options]
	Options:
		-database-name string (required)
			Name of the database to generate a schema for
		-schema-name string (required for postgres adapter)
			Name of the schema to generate a schema for
		-adapter string (required)
			Name of the database to apply the migration to (use the name of the database in the config/database.yml file)
		-env string
			Environment to use (default "development")
	Info:
		* Generates a schema for the specified database 
		* The schema is generated in the db/schemas folder	
		* The schema is named <timestamp>_<database_name>.sql or <timestamp>_<database_name>_<schema_name>.sql (for postgres)
		* The schema contains ddl statements for creating the database and its tables 
`
}

type schema struct {
	flagSet      *flag.FlagSet
	databaseName *string
	adapter      *string
	schemaName   *string
	env          *string
	projectPath  string
	db           *sql.DB
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
	sc := genSchemaFlags.String("schema-name", "", "Name of the schema to generate a schema for (postgres only)")
	env := genSchemaFlags.String("env", "development", "Environment (production or development)")
	if len(os.Args) < 2 {
		fmt.Println(genSchemaUsageString())
		os.Exit(1)
	}
	genSchemaFlags.Parse(os.Args[2:])
	utils.SetUpEnv(*env)
	dc := NewDBsConfig()
	db, err := dc.Connect(*databaseName)
	return &schema{
		flagSet:      genSchemaFlags,
		databaseName: databaseName,
		adapter:      adapter,
		schemaName:   sc,
		projectPath:  pjp,
		db:           db,
	}
}

func GenSchema() {
	s := newSchema()
	if *s.databaseName == "" || *s.adapter == "" {
		fmt.Println(genSchemaUsageString())
		os.Exit(1)
	}
	if *s.adapter == "postgres" && *s.schemaName == "" {
		fmt.Println(genSchemaUsageString())
		os.Exit(1)
	}
	f := s.createSchemaFile()
	defer f.Close()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(s.db)
	CheckDatabaseConnection(s.db)
	s.writeDDLs(f)
}

func (s *schema) createSchemaFile() *os.File {
	t := time.Now().Format("20060102150405")
	fn := fmt.Sprintf("%s_%s.sql", t, *s.databaseName)
	if *s.adapter == "postgres" {
		fn = fmt.Sprintf("%s_%s_%s.sql", t, *s.databaseName, *s.schemaName)
	}
	fp := s.creteFilePath(fn)
	f, err := os.Create(fp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	return f
}

func (s *schema) creteFilePath(fn string) string {
	fp := fmt.Sprintf("%s/db/schemas", s.projectPath)
	err := os.MkdirAll(fp, os.ModePerm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	return fmt.Sprintf("%s/%s", fp, fn)
}

func (s *schema) writeDDLs(f *os.File) {
	qq := s.ddlQueries()
	for _, q := range qq {
		// execute the query and get the ddl statement
		var tableName, ddl string
		err := s.db.QueryRow(q).Scan(&tableName, &ddl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		_, err = f.WriteString(fmt.Sprintf("%s;\n\n", ddl))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}
}

func (s *schema) ddlQueries() []string {
	var tq, dq, qp string
	switch *s.adapter {
	case "mysql":
		tq = "SELECT TABLE_NAME FROM information_schema.tables WHERE table_schema = ?"
		dq = "SHOW CREATE TABLE %s.%s"
		qp = *s.databaseName
	case "postgres":
		tq = "SELECT table_name, pg_get_ddl('TABLE ' || table_name) FROM information_schema.tables WHERE table_schema = ?"
		dq = "SELECT '%s'::text AS table_name, pg_get_ddl('TABLE %s.%s')::text AS ddl"
		qp = *s.schemaName
	case "sqlite3":
		tq = "SELECT name FROM sqlite_master WHERE type='table'"
		dq = "SELECT '%s' AS table_name, sql AS ddl FROM sqlite_master WHERE type='table' AND name='%s'"
		qp = *s.databaseName
	}
	stmt, err := s.db.Prepare(tq)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	qq := []string{}
	rows, err := stmt.Query(qp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		var q string
		switch *s.adapter {
		case "mysql":
			q = fmt.Sprintf(dq, *s.databaseName, tableName)
		case "postgres":
			q = fmt.Sprintf(dq, tableName, *s.databaseName, tableName)
		case "sqlite3":
			q = fmt.Sprintf(dq, tableName, tableName)
		default:
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		qq = append(qq, q)
	}
	return qq
}
