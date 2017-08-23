package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/jirfag/beepcar-telegram-bot/app/cfg"
)

var db *gorm.DB

type envConfig struct {
	Adapter  string
	Username string
	Password string
	Host     string
	Port     int
	Database string
}

func (c envConfig) ConnString() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s",
		c.Host, c.Port, c.Username, c.Database, c.Password)
}

func getEnvConfig() (*envConfig, error) {
	dbConfigPath := cfg.GetConfigDir() + "/db.yaml"
	configContent, err := ioutil.ReadFile(dbConfigPath)
	if err != nil {
		return nil, fmt.Errorf("can't read db config %q: %s", dbConfigPath, err)
	}

	var configs map[string]envConfig
	if err = yaml.Unmarshal(configContent, &configs); err != nil {
		return nil, fmt.Errorf("invalid yaml in db config %q: %s", dbConfigPath, err)
	}

	c, ok := configs[cfg.GetEnv()]
	if !ok {
		return nil, fmt.Errorf("no current env %q db config in %q", cfg.GetEnv(), dbConfigPath)
	}

	return &c, nil
}

func init() {
	if err := initDB(); err != nil {
		log.Fatalf("can't init DB: %s", err)
	}
}

func initDB() error {
	cfg, err := getEnvConfig()
	if err != nil {
		return fmt.Errorf("can't get db env config: %s", err)
	}

	db, err = gorm.Open(cfg.Adapter, cfg.ConnString())
	if err != nil {
		return fmt.Errorf("can't open gorm connection for cfg %+v: %s", cfg, err)
	}

	return nil
}

func Get() *gorm.DB {
	return db
}

func Int64FK(v int64) sql.NullInt64 {
	return sql.NullInt64{
		Int64: v,
		Valid: v != 0,
	}
}
