package helper

import (
	"net/http"
	"regexp"

	"github.com/gorilla/mux"

	"github.com/yourtion/go-short-url/internal/base/config"
	"github.com/yourtion/go-short-url/internal/base/define"
	"github.com/yourtion/go-short-url/internal/base/logger"
	"github.com/yourtion/go-short-url/internal/utils"
)

var log *logger.Entry
var IpRegexp *regexp.Regexp

func init() {
	log = logger.NewModuleLogger("helper")
	IpRegexp = regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`)
}

func CombineHandlers(list ...func(ctx *Context)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(w, r)
		for _, f := range list {
			ctx.goNext = false
			f(ctx)
			if !ctx.goNext {
				break
			}
		}
	}
}

func Wrap(f func(w http.ResponseWriter, r *http.Request), cors bool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer utils.Try()

		w.Header().Set("X-Powered-By", define.Version)
		w.Header().Set("X-Server", config.Config.Server.Name)

		if cors {
			// CORS
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("origin"))
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS, PATCH")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Authorization, Accept, X-Requested-With, Cookie")
		}

		f(w, r)
	}
}

func CorsOptions(r *mux.Router, prefix string) {
	r.PathPrefix(prefix).HandlerFunc(Wrap(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}, true)).Methods("OPTIONS")
}
