package services

import (
	"strconv"
	"strings"
	"time"

	"github.com/yourtion/go-short-url/internal/base/define"
	"github.com/yourtion/go-short-url/internal/base/redis"
	"github.com/yourtion/go-short-url/internal/models"
	"github.com/yourtion/go-short-url/internal/utils"
)

// Redis同步到数据库后自动设置过期
const ExpireTime = 60 * 24 * time.Hour

func getPvKey(day string) string {
	return redis.GetKey(strings.Join([]string{day, "pv"}, ":"))
}

func getUvKey(day string) string {
	return redis.GetKey(strings.Join([]string{day, "uv"}, ":"))
}

func SyncToRedis(clean bool) {
	if redis.Client == nil {
		return
	}
	syncLock.Lock()
	defer syncLock.Unlock()
	day := utils.GetTodayDayString()
	log.Tracef("SyncToRedis: %s", day)
	pvKey := getPvKey(day)
	uvKey := getUvKey(day)
	pvI, uvI := 0, 0
	pipe := redis.Client.Pipeline()
	// 使用 sync.Map 的 Range 方法遍历
	statics.Range(func(id interface{}, data interface{}) bool {
		pv := data.(*sInfo).pv.Get()
		if pv > 0 {
			data.(*sInfo).pv.Reset()
			pipe.HIncrBy(pvKey, strconv.Itoa(id.(int)), pv)
			pvI++
		}
		uv := data.(*sInfo).uv.Get()
		if uv > 0 {
			data.(*sInfo).uv.Reset()
			pipe.HIncrBy(uvKey, strconv.Itoa(id.(int)), pv)
			uvI++
		}
		// 删除空白的统计
		if clean && pv == 0 && uv == 0 {
			statics.Delete(id)
		}
		return true
	})
	_, err := pipe.Exec()
	if err != nil {
		log.Warnf("SyncToRedisError: %s", err)
	}
	log.Tracef("SyncToRedisDone %s:%d - %s:%d", pvKey, pvI, uvKey, uvI)
}

func FlushAccessLog() {
	day := utils.GetTodayDayString()
	log.Tracef("FlushAccessLog: %s", day)
	err := accessLog.Rotate(utils.GetTodayDayString())
	if err != nil {
		log.Warnf("accessLog.Rotate Error: %s", err)
	}
}

type infoPvUv struct {
	pv int
	uv int
}

func parseMapInfo(key string, val string) (int, int, error) {
	id, err := strconv.Atoi(key)
	if err != nil {
		return 0, 0, err
	}
	v, err := strconv.Atoi(val)
	if err != nil {
		return 0, 0, err
	}
	return id, v, nil
}

func SyncToDB(yesterday string) {
	dbInfo := make(map[int]*infoPvUv)

	// 获取PV信息
	pvs := redis.Client.HGetAll(getPvKey(yesterday)).Val()
	for key, val := range pvs {
		id, pv, err := parseMapInfo(key, val)
		if err != nil {
			continue
		}
		if dbInfo[id] == nil {
			dbInfo[id] = &infoPvUv{}
		}
		dbInfo[id].pv = pv
	}
	// 获取UV信息
	uvs := redis.Client.HGetAll(getUvKey(yesterday)).Val()
	for key, val := range uvs {
		id, uv, err := parseMapInfo(key, val)
		if err != nil {
			continue
		}
		if dbInfo[id] == nil {
			dbInfo[id] = &infoPvUv{}
		}
		dbInfo[id].uv = uv
	}

	// 将PV、UV信息写入数据库
	day, _ := strconv.Atoi(yesterday)
	allOk := true
	for id, val := range dbInfo {
		info := define.StatisticRow{Day: day, SId: id, PV: val.pv, UV: val.uv}
		ok := models.CreateOrUpdateStatisticInfo(&info)
		if !ok {
			log.Warnf("SyncToDB Error: %+v", info)
			allOk = false
		}
	}
	// 如果写入成功，设置RedisKey过期
	if allOk {
		redis.Client.Expire(getPvKey(yesterday), ExpireTime)
		redis.Client.Expire(getUvKey(yesterday), ExpireTime)
	}
	log.Infof("SyncToDB Done: %s", yesterday)
}
