package services

import (
	"errors"
	"time"

	"github.com/yourtion/go-utils/cache"
	"github.com/yourtion/go-utils/memo"

	"github.com/yourtion/go-short-url/internal/base/define"
	"github.com/yourtion/go-short-url/internal/base/logger"
	"github.com/yourtion/go-short-url/internal/models"
	"github.com/yourtion/go-short-url/internal/utils"
)

var shortCache *cache.LRU
var shortMemo *memo.Memo
var accessLog *logger.RotateWriter

func CacheInit() {
	// TODO: 配置AccessLog路径
	accessLog = logger.NewRotateWriter("/tmp", utils.GetTodayDayString())
	// TODO: 通过配置获取缓存时间
	shortCache = cache.NewLRU(1024, 60*time.Second)
	shortMemo = memo.NewMemo(getShortInfoByShort)
}

// 获取短链缓存数据
func GetCacheStatus() *cache.Status {
	return shortCache.Status()
}

// 获取短链接信息Memo方法
func getShortInfoByShort(short string) (interface{}, error) {
	info := models.GetUrlInfoByShort(short)
	if info != nil {
		return info, nil
	}
	return nil, errors.New("not found")
}

// 通过短链获取原始链接地址
func GetOriginUrlFromShort(short string) (*define.ShortRow, bool) {
	// 先从缓存取
	row, ok := shortCache.Get(short)
	log.Tracef("shortCache: %v -> %+v", ok, row)
	if row != nil {
		return row.(*define.ShortRow), true
	}

	// 缓存没有信息通过Memo获取（防止并发导致数据库压力过大）
	shortInfo, err := shortMemo.Get(short)
	log.Tracef("shortMemo: %v -> %+v", err, shortInfo)
	if err == nil && shortInfo != nil {
		shortCache.Set(short, shortInfo)
		return shortInfo.(*define.ShortRow), true
	}

	// 缓存空结果（防止缓存穿透导致DDos）
	// TODO: 通过配置获取缓存时间
	shortCache.SetEx(short, &define.ShortRow{}, 10*time.Second)
	return nil, false
}

func UpdateCache(key string, info *define.ShortRow) {
	shortCache.Set(key, info)
}
