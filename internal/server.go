package internal

import (
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/yourtion/go-short-url/internal/base/config"
	"github.com/yourtion/go-short-url/internal/base/define"
	"github.com/yourtion/go-short-url/internal/base/logger"
	"github.com/yourtion/go-short-url/internal/base/mysql"
	"github.com/yourtion/go-short-url/internal/base/redis"
	"github.com/yourtion/go-short-url/internal/controller/admin"
	"github.com/yourtion/go-short-url/internal/controller/helper"
	"github.com/yourtion/go-short-url/internal/controller/redirect"
	"github.com/yourtion/go-short-url/internal/models"
	"github.com/yourtion/go-short-url/internal/services"
)

var log *logger.Entry
var server *http.Server

func init() {
	log = logger.NewModuleLogger("main").WithField("version", define.Version)
}

func Server() *mux.Router {
	log.Infof("pid: %d, gid: %d, uid: %d", os.Getpid(), os.Getgid(), os.Getuid())

	// 根据运行目录获取配置文件名
	configFile := "config.toml"
	workingDir := "./"

	// 载入配置
	config.LoadConfig(workingDir, configFile)
	log.Infof("server name is %s", config.Config.Server.Name)
	log.Infof("config: %+v", config.Config)

	// 切换到指定的工作目录
	if err := os.Chdir(config.Config.CWD); err != nil {
		log.Errorf("change working directory to %s failed: %s", config.Config.CWD, err)
	} else {
		log.Infof("current working directory is %s", config.Config.CWD)
	}

	// 初始化日志记录器
	if level, err := logrus.ParseLevel(config.Config.Log.Level); err != nil {
		log.Errorf("invalid log level: %s", config.Config.Log.Level)
	} else {
		logger.Logger.SetLevel(level)
		if level >= logrus.DebugLevel {
			// 调试模式下输出带颜色的日志，方便阅读
			logger.Logger.SetFormatter(&logrus.TextFormatter{
				ForceColors: true,
			})
		}
		log.Infof("log level is %s", config.Config.Log.Level)
	}

	// 初始化MySQL连接
	mysql.Open(&config.Config.MySQL)
	// 初始化Redis连接
	redis.Open(&config.Config.Redis)

	// 加载动态配置
	models.LoadBaseConfig(&config.Dynamic)
	// 初始化缓存
	services.CacheInit()
	// 初始化路由
	r := mux.NewRouter()
	log.Debugf("server prefix: %s", "/")
	helper.CorsOptions(r, "/api/")
	redirect.Register(r)
	admin.Register(r)

	return r
}

func StartServer() {
	r := Server()

	// 判断是否需要启动 pprof
	if len(config.Config.Server.PProf) > 0 {
		log.Warnf("start pprof web interface on %s", config.Config.Server.PProf)
		log.Infof("Open on http://%s/debug/pprof/", config.Config.Server.PProf)
		go func() {
			if err := http.ListenAndServe(config.Config.Server.PProf, nil); err != nil {
				log.Warnln(err)
			}
		}()
	}

	// 启动服务器
	log.Printf("server listen on %s", config.Config.Server.Listen)
	go loop()

	server = new(http.Server)
	server.ReadTimeout = 5 * time.Second
	server.WriteTimeout = 5 * time.Second
	server.Addr = config.Config.Server.Listen
	server.Handler = context.ClearHandler(r)
	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Warnln(err)
		} else {
			log.Fatalf("listen http failed: %s", err)
		}
	}
}
