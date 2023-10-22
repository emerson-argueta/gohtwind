package infra

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type Database struct {
	Adapter  string `yaml:"adapter"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
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

	//Unmarshal the file contents
	var cfg DBsConfig
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return &cfg
}

var (
	dbSetupFuncMap = map[string]func(Database) (*sql.DB, error){
		"mysql": setUpMysqlDB,
		"psql":  setUpPsqlDB,
	}
)

func (cfg *DBsConfig) ConnectAll() (map[string]*sql.DB, error) {
	dbs := make(map[string]*sql.DB)
	for name, db := range cfg.Databases {
		var d *sql.DB
		var err error
		if dbSetupFunc, ok := dbSetupFuncMap[db.Adapter]; ok {
			d, err = dbSetupFunc(db)
		} else {
			return nil, fmt.Errorf("Adapter %s not supported", db.Adapter)
		}
		if err != nil {
			return nil, err
		}
		dbs[name] = d
	}
	return dbs, nil
}

func setUpMysqlDB(d Database) (*sql.DB, error) {
	connectString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", d.Username, d.Password, d.Host, d.Port, d.Database)
	db, err := sql.Open("mysql", connectString)
	if d.MaxConn != 0 {
		db.SetMaxOpenConns(d.MaxConn)
	}
	if d.MaxIdle != 0 {
		db.SetMaxIdleConns(d.MaxIdle)
	}
	return db, err
}

func setUpPsqlDB(d Database) (*sql.DB, error) {
	var connectString = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", d.Host, d.Port, d.Username, d.Password, d.Database)
	db, err := sql.Open("postgres", connectString)
	if d.MaxConn != 0 {
		db.SetMaxOpenConns(d.MaxConn)
	}
	if d.MaxIdle != 0 {
		db.SetMaxIdleConns(d.MaxIdle)
	}
	return db, err
}
