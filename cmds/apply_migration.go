package cmds

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log"
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
		-schema-name string (required for postgres adapter)
			Name of the schema to apply the migration to (use the name of the schema in the config/database.yml file)
		-env string (optional)
			Environment to use (development, test, production) (default "development")
	Info:
		* Applies the specified migration file
		* The migration file is expected to be in the db/migrations/<database_name> folder
		* The migration file is expected to be in the db/migrations/<database_name>/<schema_name> folder (for postgres)
		* The migration file is expected to be named <feature_name>_<table_name>.sql
		* The migration file is expected to contain a create table statement or an alter table statement
`
}

type applyMigration struct {
	flagSet      *flag.FlagSet
	fileName     *string
	adapter      *string
	databaseName *string
	schemaName   *string
	env          *string
	projectPath  string
}

func setUpEnv() string {
	// Define the -env flag
	env := flag.String("env", "development", "Environment (production or development)")
	flag.Parse()
	var envFileName string
	if *env == "production" {
		envFileName = ".env.production"
	} else {
		envFileName = ".env"
	}
	return envFileName
}

func newApplyMigration() *applyMigration {
	ef := setUpEnv()
	err := godotenv.Load(ef)
	applyMigrationFlags := flag.NewFlagSet("gohtwind apply-migration", flag.ExitOnError)
	pjp, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fileName := applyMigrationFlags.String("file-name", "", "Name of the migration file to apply")
	adapter := applyMigrationFlags.String("adapter", "", "Database adapter to use")
	databaseName := applyMigrationFlags.String("database-name", "", "Name of the database to apply the migration to (use the name of the database in the config/database.yml file)")
	schemaName := applyMigrationFlags.String("schema-name", "", "Name of the schema to apply the migration to (use the name of the schema in the config/database.yml file)")
	env := applyMigrationFlags.String("env", "development", "Environment to use (development, test, production)")
	return &applyMigration{
		flagSet:      applyMigrationFlags,
		fileName:     fileName,
		adapter:      adapter,
		databaseName: databaseName,
		schemaName:   schemaName,
		env:          env,
		projectPath:  pjp,
	}
}

func ApplyMigration() {
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
	if *m.adapter == "postgres" && *m.schemaName == "" {
		fmt.Println(applyMigrationUsageString())
		os.Exit(1)
	}
	m.applyMigration()
}

func (m *applyMigration) applyMigration() {
	switch *m.adapter {
	case "mysql":
		m.applyMySQLMigration()
	case "postgres":
		m.applyPostgresMigration()
	case "sqlite3":
		m.applySQLite3Migration()
	default:
		fmt.Fprintf(os.Stderr, "Error: %v\n", fmt.Errorf("Unsupported adapter: %s", *m.adapter))
		os.Exit(1)
	}
}

func (m *applyMigration) applyMySQLMigration() {
	fp := fmt.Sprintf("%s/db/migrations/%s/%s", m.projectPath, *m.databaseName, *m.fileName)
	fb, err := os.ReadFile(fp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fs := string(fb)
	dc := NewDBsConfig()
	db, err := dc.Connect(*m.databaseName)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)
	CheckDatabaseConnection(db)
	_, err = db.Exec(fs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func (m *applyMigration) applyPostgresMigration() {
	fp := fmt.Sprintf("%s/db/migrations/%s/%s/%s", m.projectPath, *m.databaseName, *m.schemaName, *m.fileName)
	fb, err := os.ReadFile(fp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fs := string(fb)
	dc := NewDBsConfig()
	db, err := dc.Connect(*m.databaseName)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)
	CheckDatabaseConnection(db)
	_, err = db.Exec(fs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func (m *applyMigration) applySQLite3Migration() {
	fp := fmt.Sprintf("%s/db/migrations/%s/%s", m.projectPath, *m.databaseName, *m.fileName)
	fb, err := os.ReadFile(fp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fs := string(fb)
	dc := NewDBsConfig()
	db, err := dc.Connect(*m.databaseName)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)
	CheckDatabaseConnection(db)
	_, err = db.Exec(fs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
