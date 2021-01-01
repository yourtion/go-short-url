package redirect

import (
	"net/http"

	"github.com/yourtion/go-short-url/internal/base/define"
	"github.com/yourtion/go-short-url/internal/controller/helper"
	"github.com/yourtion/go-short-url/internal/services"
	"github.com/yourtion/go-short-url/internal/utils"
)

var sidKey = define.ServiceName + "-id"

// 短链接Handler
func shortHandler(ctx *helper.Context) {
	short := ctx.GetParamsString("short")
	shortInfo, ok := services.GetOriginUrlFromShort(short)
	if !ok || shortInfo == nil {
		ctx.Res.WriteHeader(http.StatusNotFound)
		ctx.ResponseText("Short not Found: " + short)
	} else {
		name := services.BuildCookieKey(shortInfo.Short)
		addUv := false
		if shortInfo.IsStatistics && ctx.GetCookie(name) == nil {
			// 使用 Path 策略保证上报的 cookie 数量少
			ctx.SetCookie(&http.Cookie{Name: name, Value: "1", Path: ctx.Req.URL.Path, MaxAge: 3600 * 24 * 30})
			addUv = true
		}
		sidCookie := ctx.GetCookie(sidKey)
		uid := ""
		if sidCookie != nil {
			uid = sidCookie.Value
		} else {
			uid = utils.GenerateUid()
			ctx.SetCookie(&http.Cookie{Name: sidKey, Value: uid, Path: "/"})
		}
		go services.AddStatisticsInfo(ctx, shortInfo, uid, addUv)
		// ctx.ResponseOk(shortInfo.Origin)
		ctx.Redirect(shortInfo.Origin)
	}
}
