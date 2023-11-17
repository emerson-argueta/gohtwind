package cmds

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strings"
	"time"
)

type Database struct {
	Adapter  string `yaml:"adapter"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Schema   string `yaml:"schema"`
	Socket   string `yaml:"socket"`
	MaxConn  int    `yaml:"max_conn"`
	MaxIdle  int    `yaml:"max_idle"`
}

type DBsConfig struct {
	Databases map[string]Database `yaml:"databases"`
}

func NewDBsConfig() *DBsConfig {
	// Read YAML file
	data, err := os.ReadFile("./config/database.yml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Replace placeholders with env variables
	configSt, err := parseConfigTemplate(string(data))
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	//Unmarshal the file contents
	var cfg DBsConfig
	err = yaml.Unmarshal([]byte(configSt), &cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return &cfg
}

func CheckDatabaseConnection(db *sql.DB) {
	maxRetries := 5
	retryInterval := 5 * time.Second
	var err error
	for i := 0; i < maxRetries; i++ {
		err = db.Ping()
		if err == nil {
			break
		}
		time.Sleep(retryInterval)
	}
	if err != nil {
		log.Println(err)
	}
}

func parseConfigTemplate(template string) (string, error) {
	// Find placeholders in the format {{VAR_NAME}} and replace them with the actual environment variable values
	placeholders := findPlaceholders(template)

	for _, ph := range placeholders {
		envValue, exists := os.LookupEnv(ph)
		if !exists {
			return "", fmt.Errorf("environment variable %s not set", ph)
		}
		template = strings.ReplaceAll(template, fmt.Sprintf("{{%s}}", ph), envValue)
	}

	return template, nil
}

func findPlaceholders(template string) []string {
	var placeholders []string
	startIdx := 0
	// algorithm: find the first occurrence of {{, then find the first occurrence of }} after that
	// then extract the string between {{ and }} and add it to the placeholders slice
	// then repeat the process until there are no more {{ }} pairs
	for {
		startIdx = strings.Index(template, "{{")
		if startIdx == -1 {
			break
		}
		endIdx := strings.Index(template, "}}")
		if endIdx == -1 {
			break
		}
		placeholders = append(placeholders, template[startIdx+2:endIdx])
		template = template[endIdx+2:]
	}
	return placeholders
}

var dbSetupFuncMap = map[string]func(Database) (*sql.DB, error){
	"mysql": func(d Database) (*sql.DB, error) {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			d.Username, d.Password, d.Host, d.Port, d.Database,
		)
		return setUpDB(d, "mysql", dsn)
	},
	"psql": func(d Database) (*sql.DB, error) {
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable search_path=%s",
			d.Host, d.Port, d.Username, d.Password, d.Database, d.Schema,
		)
		return setUpDB(d, "postgres", dsn)
	},
}

func (cfg *DBsConfig) ConnectAll() (map[string]*sql.DB, error) {
	dbs := make(map[string]*sql.DB)
	for name, d := range cfg.Databases {
		var db *sql.DB
		var err error
		dbSetupFunc, ok := dbSetupFuncMap[d.Adapter]
		if !ok {
			return nil, fmt.Errorf("Adapter %s not supported", d.Adapter)
		}
		db, err = dbSetupFunc(d)
		if err != nil {
			return nil, err
		}
		dbs[name] = db
	}
	return dbs, nil
}

func (cfg *DBsConfig) Connect(dbName string) (*sql.DB, error) {
	d := cfg.Databases[dbName]
	dbSetupFunc, ok := dbSetupFuncMap[d.Adapter]
	if !ok {
		return nil, fmt.Errorf("Adapter %s not supported", d.Adapter)
	}
	return dbSetupFunc(d)
}

func setUpDB(d Database, adapter string, dsn string) (*sql.DB, error) {
	db, err := sql.Open(adapter, dsn)
	if d.MaxConn != 0 {
		db.SetMaxOpenConns(d.MaxConn)
	}
	if d.MaxIdle != 0 {
		db.SetMaxIdleConns(d.MaxIdle)
	}
	return db, err
}
