package admin

import (
	"strings"

	"github.com/yourtion/go-short-url/internal/base/config"
	"github.com/yourtion/go-short-url/internal/controller/helper"
	"github.com/yourtion/go-short-url/internal/models"
	"github.com/yourtion/go-short-url/internal/services"
)

func verifyUrl(url string) bool {
	if !strings.HasPrefix(url, "http") {
		return false
	}
	return true
}

// 短链接Handler
func createShortHandler(ctx *helper.Context) {
	body, err := ctx.ParseJsonBody()
	if err != nil {
		ctx.ResponseError(err.Error())
		return
	}
	url := body.Get("url").ToString()
	if url == "" || !verifyUrl(url) {
		ctx.ResponseJson("url not verify")
		return
	}
	info := models.AddUrlToShort(url)
	// 新创建的url放到缓存里面
	if info.Id > 0 && info.Origin != "" {
		services.UpdateCache(info.Short, info)
	}
	ctx.ResponseOk(map[string]interface{}{
		"short": config.Config.Server.URL + "/s/" + info.Short,
	})
}
