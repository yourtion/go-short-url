package models

import "github.com/yourtion/go-short-url/internal/base/logger"

var log *logger.Entry

func init() {
	log = logger.NewModuleLogger("models")
}
