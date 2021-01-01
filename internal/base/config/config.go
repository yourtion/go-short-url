package config

import (
	"github.com/BurntSushi/toml"

	"github.com/yourtion/go-short-url/internal/base/logger"
)

var log *logger.Entry
var Config MainConfig
var Dynamic DynamicInfo

func init() {
	log = logger.NewModuleLogger("config")
}

func LoadConfig(workingDir string, file string) {
	log.Infof("load config from file: %s", file)
	_, err := toml.DecodeFile(file, &Config)
	if err != nil {
		log.Fatalf("load config failed: %s", err)
	}

	if len(Config.CWD) < 1 {
		Config.CWD = workingDir
	}
}
