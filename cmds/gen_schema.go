package cmds

import (
	"database/sql"
	"flag"
	"fmt"
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
	dc := NewDBsConfig()
	db, err := dc.Connect(*sc)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)
	CheckDatabaseConnection(db)
	return &schema{
		flagSet:      genSchemaFlags,
		databaseName: databaseName,
		adapter:      adapter,
		schemaName:   sc,
		projectPath:  pjp,
		db:           db,
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
	if *s.adapter == "postgres" && *s.schemaName == "" {
		fmt.Println(genSchemaUsageString())
		os.Exit(1)
	}
	s.genSchema()
}

func (s *schema) genSchema() {
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
	defer f.Close()
	switch *s.adapter {
	case "mysql":
		s.writeMysqlDDL(f)
	case "postgres":
		s.writePostgresDDL(f)
	case "sqlite3":
		s.writeSQLite3DDL(f)
	default:
		fmt.Fprintf(os.Stderr, "Error: %v\n", fmt.Errorf("Unsupported adapter: %s", *s.adapter))
		os.Exit(1)
	}
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

func (s *schema) writeMysqlDDL(f *os.File) {
	qq := s.mysqlDdlQueries()
	for _, q := range qq {
		// execute the query and get the ddl statement
		var tableName, ddl string
		err := s.db.QueryRow(q).Scan(&tableName, &ddl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		// write the ddl statement to the file
		_, err = f.WriteString(fmt.Sprintf("%s;\n", ddl))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}
}

func (s *schema) mysqlDdlQueries() []string {
	stmt, err := s.db.Prepare("SELECT TABLE_NAME FROM information_schema.tables WHERE table_schema = ?")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	qq := []string{}
	rows, err := stmt.Query(*s.databaseName)
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
		q := fmt.Sprintf("SHOW CREATE TABLE %s.%s;\n", *s.databaseName, tableName)
		qq = append(qq, q)
	}
	return qq
}

func (s *schema) writePostgresDDL(f *os.File) {
	qq := s.postgresDdlQueries()
	for _, q := range qq {
		// execute the query and get the ddl statement
		var tableName, ddl string
		err := s.db.QueryRow(q).Scan(&tableName, &ddl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		// write the ddl statement to the file
		_, err = f.WriteString(fmt.Sprintf("%s;\n", ddl))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}
}

func (s *schema) postgresDdlQueries() []string {
	stmt, err := s.db.Prepare("SELECT table_name FROM information_schema.tables WHERE table_schema = ?")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	qq := []string{}
	rows, err := stmt.Query(*s.schemaName)
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
		q := fmt.Sprintf("SELECT '%s'::text AS table_name, pg_get_ddl('TABLE %s.%s')::text AS ddl", tableName, *s.schemaName, tableName)
		qq = append(qq, q)
	}
	return qq
}

func (s *schema) writeSQLite3DDL(f *os.File) {
	qq := s.sqlite3DdlQueries()
	for _, q := range qq {
		// execute the query and get the ddl statement
		var tableName, ddl string
		err := s.db.QueryRow(q).Scan(&tableName, &ddl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		// write the ddl statement to the file
		_, err = f.WriteString(fmt.Sprintf("%s;\n", ddl))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}
}

func (s *schema) sqlite3DdlQueries() []string {
	stmt, err := s.db.Prepare("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	qq := []string{}
	rows, err := stmt.Query()
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
		q := fmt.Sprintf("SELECT '%s' AS table_name, sql AS ddl FROM sqlite_master WHERE type='table' AND name='%s'", tableName, tableName)
		qq = append(qq, q)
	}
	return qq
}
