package cfg

import (
	"flag"
	"os"
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
	env := os.Getenv("GO_MODE")
	if env != "" {
		return env
	}

	if flag.Lookup("test.v") != nil {
		return "test"
	}

	return "development"
}
