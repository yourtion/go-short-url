package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()
	// 使用 JSON 格式记录
	// Logger.SetFormatter(&logrus.JSONFormatter{})
	// 输出到 stdout
	Logger.SetOutput(os.Stdout)

	// 初始化 syslog
	// hook, err := lSyslog.NewSyslogHook("", "", syslog.LOG_WARNING, define.ServiceName)
	// if err != nil {
	// 	Logger.Errorf("init syslog hook failed: %s", err)
	// } else {
	// 	Logger.Hooks.Add(hook)
	// }
}

type Entry = logrus.Entry

func WithFields(key string, value interface{}) *Entry {
	return Logger.WithField(key, value)
}

func NewModuleLogger(name string) *Entry {
	return Logger.WithField("module", name)
}
