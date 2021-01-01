// 配置中心
package config

import (
	"github.com/BurntSushi/toml"

	"github.com/yourtion/go-short-url/internal/base/logger"
)

var log *logger.Entry

// 核心配置（通过配置文件加载）
var Config MainConfig

// 动态配置（通过数据库 config 表加载）
var Dynamic DynamicInfo

func init() {
	log = logger.NewModuleLogger("config")
}

// 加载配置文件
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
