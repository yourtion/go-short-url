package utils

import (
	"runtime/debug"

	"github.com/yourtion/go-short-url/internal/base/logger"
)

var log *logger.Entry

func init() {
	log = logger.NewModuleLogger("try_recover")
}

func Try() {
	err := recover()
	if err != nil {
		log.Errorf("try recover: %+v\n%s", err, string(debug.Stack()))
	}
}
