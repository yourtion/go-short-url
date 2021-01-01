package admin

import (
	"github.com/gorilla/mux"

	"github.com/yourtion/go-short-url/internal/base/logger"
	"github.com/yourtion/go-short-url/internal/controller/helper"
)

var log *logger.Entry

func init() {
	log = logger.NewModuleLogger("api.admin")
}

func Register(r *mux.Router) {
	var register = func(method string, path string, cors bool, list ...func(ctx *helper.Context)) {
		r.HandleFunc(path, helper.Wrap(helper.CombineHandlers(list...), cors)).Methods(method)
	}
	register("POST", "/api/create", true, createShortHandler)
}
