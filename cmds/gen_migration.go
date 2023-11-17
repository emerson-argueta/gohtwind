package cmds

import (
	"flag"
	"fmt"
	"os"
	"time"
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
		-schema-name string (required for postgres adapter)
			Name of the schema to apply the migration to (use the name of the schema in the config/database.yml file)
	Info:
		* Generates a migration file for the specified model
		* The migration file is generated in the db/migrations/<database_name> folder for mysql and sqlite3
		* The migration file is generated in the db/migrations/<database_name>/<schema_name> folder for postgres
		* The migration file is named <timestamp>_<feature_name>_<table_name>.sql
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
	schemaName   *string
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
	schemaName := genMigrationFlags.String("schema-name", "", "Name of the schema to apply the migration to (use the name of the schema in the config/database.yml file)")
	return &migration{
		flagSet:      genMigrationFlags,
		featName:     featName,
		tableName:    tableName,
		action:       action,
		adapter:      adapter,
		databaseName: databaseName,
		schemaName:   schemaName,
		projectPath:  pjp,
	}
}

func GenMigration() {
	m := newMigration()
	if len(os.Args) < 2 {
		fmt.Println(genMigrationUsageString())
		os.Exit(1)
	}
	m.flagSet.Parse(os.Args[2:])
	if *m.tableName == "" || *m.action == "" || *m.adapter == "" || *m.databaseName == "" {
		fmt.Println(genMigrationUsageString())
		os.Exit(1)
	}
	if *m.adapter == "postgres" && *m.schemaName == "" {
		fmt.Println(genMigrationUsageString())
		os.Exit(1)
	}
	m.genMigration()
}

func (m *migration) genMigration() {
	// TODO: Implement
	/* Generates a migration file for the specified model
	* The migration file is generated in the db/migrations/<database_name> folder
	* The migration file is named <timestamp>_<feature_name>_<table_name>.sql
	* The migration file contains a create table statement, an alter table statement, or a drop table statement
	 */
	t := time.Now().Format("20060102150405")
	fn := fmt.Sprintf("%s_%s.sql", t, *m.tableName)
	if *m.featName != "" {
		fn = fmt.Sprintf("%s_%s_%s.sql", t, *m.featName, *m.tableName)
	}
	fp := m.creteFilePath(fn)
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

func (m *migration) creteFilePath(fn string) string {
	switch *m.adapter {
	case "postgres":
		fp := fmt.Sprintf("%s/db/migrations/%s/%s", m.projectPath, *m.databaseName, *m.schemaName)
		err := os.MkdirAll(fp, os.ModePerm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return fmt.Sprintf("%s/%s", fp, fn)
	case "mysql":
		fp := fmt.Sprintf("%s/db/migrations/%s", m.projectPath, *m.databaseName)
		err := os.MkdirAll(fp, os.ModePerm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return fmt.Sprintf("%s/%s", fp, fn)
	case "sqlite3":
		fp := fmt.Sprintf("%s/db/migrations/%s", m.projectPath, *m.databaseName)
		err := os.MkdirAll(fp, os.ModePerm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return fmt.Sprintf("%s/%s", fp, fn)
	default:
		fmt.Fprintf(os.Stderr, "Error: %v\n", "Invalid adapter")
	}
	return ""
}

func (m *migration) CreateTableStatement() []byte {
	switch *m.adapter {
	case "postgres":
		// make sure it has database name prefix
		return []byte(fmt.Sprintf("CREATE TABLE %s (\n\tid SERIAL PRIMARY KEY\n);", *m.tableName))
	case "mysql":
		return []byte(fmt.Sprintf("CREATE TABLE %s (\n\tid INT NOT NULL AUTO_INCREMENT,\n\tPRIMARY KEY (id)\n);", *m.tableName))
	case "sqlite3":
		return []byte(fmt.Sprintf("CREATE TABLE %s (\n\tid INTEGER PRIMARY KEY AUTOINCREMENT\n);", *m.tableName))
	default:
		fmt.Fprintf(os.Stderr, "Error: %v\n", "Invalid adapter")

	}
	return []byte("")
}

func (m *migration) AlterTableStatement() []byte {
	switch *m.adapter {
	case "postgres":
		return []byte(fmt.Sprintf("ALTER TABLE %s ADD COLUMN column_name data_type;", *m.tableName))
	case "mysql":
		return []byte(fmt.Sprintf("ALTER TABLE %s ADD COLUMN column_name data_type;", *m.tableName))
	case "sqlite3":
		return []byte(fmt.Sprintf("ALTER TABLE %s ADD COLUMN column_name data_type;", *m.tableName))
	default:
		fmt.Fprintf(os.Stderr, "Error: %v\n", "Invalid adapter")
	}
	return []byte("")
}

func (m *migration) DropTableStatement() []byte {
	switch *m.adapter {
	case "postgres":
		return []byte(fmt.Sprintf("DROP TABLE %s;", *m.tableName))
	case "mysql":
		return []byte(fmt.Sprintf("DROP TABLE %s;", *m.tableName))
	case "sqlite3":
		return []byte(fmt.Sprintf("DROP TABLE %s;", *m.tableName))
	default:
		fmt.Fprintf(os.Stderr, "Error: %v\n", "Invalid adapter")
	}
	return []byte("")
}
