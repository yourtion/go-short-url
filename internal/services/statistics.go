package services

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/yourtion/go-short-url/internal/utils"

	"github.com/yourtion/go-short-url/internal/base/define"
	"github.com/yourtion/go-short-url/internal/controller/helper"
)

type sInfo struct {
	pv utils.AtomicInt
	uv utils.AtomicInt
}

var statics sync.Map
var syncLock sync.Mutex

func init() {
	syncLock = sync.Mutex{}
}

func BuildCookieKey(id int) string {
	return define.ServiceName + "-" + strconv.Itoa(id)
}

func AddStatisticsInfo(ctx *helper.Context, shortInfo *define.ShortRow, uid string, uv bool) {
	if shortInfo.IsStatistics {
		addPvUv(shortInfo.Id, uv)
	}
	if shortInfo.IsAccessLog {
		addAccessLog(ctx, shortInfo, uid)
	}
}

func addAccessLog(ctx *helper.Context, shortInfo *define.ShortRow, uid string) {
	log.Tracef("addAccessLog:%d", shortInfo.Id)
	ua := ctx.GetUserAgent()
	ip := ctx.GetIp()
	ref := ctx.GetReferer()
	t := time.Now().Unix()
	info := fmt.Sprintf("%d,%d,\"%s\",\"%s\",\"%s\",\"%s\"\n", shortInfo.Id, t, uid, ip, ua, ref)
	_, err := accessLog.WriteString(info)
	if err != nil {
		log.Warnf("accessLog.WriteString Error: %s", err)
		log.Warnln(info)
	}
}

func addPvUv(id int, uv bool) {
	log.Tracef("addPvUv:%d, %v", id, uv)
	if val, ok := statics.LoadOrStore(id, new(sInfo)); ok {
		val.(*sInfo).pv.Incr(1)
		if uv {
			val.(*sInfo).uv.Incr(1)
		}
	}
}
