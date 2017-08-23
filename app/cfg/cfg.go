package cfg

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

func IsProduction() bool {
	return GetEnv() == "production"
}

func GetConfigDir() string {
	dir := "config"
	if GetEnv() == "test" {
		dir = os.Getenv("GOPATH") + "/src/github.com/jirfag/beepcar-telegram-bot/" + dir
	}

	return dir
}

func GetEnv() string {
	env := os.Getenv("GO_ENV")
	if env != "" {
		return env
	}

	if flag.Lookup("test.v") != nil {
		return "test"
	}

	return "development"
}

func GetYamlConfig(name string, config interface{}) error {
	configPath := fmt.Sprintf("%s/%s.yaml", GetConfigDir(), name)
	configContent, err := ioutil.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("can't read config %q: %s", configPath, err)
	}

	if err = yaml.Unmarshal(configContent, config); err != nil {
		return fmt.Errorf("invalid yaml in config %q: %s", configPath, err)
	}

	return nil
}

type App struct {
	Telegram struct {
		APIKey string `yaml:"api_key"`
	}
}

var app App

func GetApp() App {
	return app
}

func init() {
	if err := GetYamlConfig("app", &app); err != nil {
		log.Fatalf("can't read app config: %s", err)
	}
}
