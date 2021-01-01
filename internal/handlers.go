package internal

import (
	"time"

	"github.com/gorilla/context"

	"github.com/yourtion/go-short-url/internal/services"
	"github.com/yourtion/go-short-url/internal/utils"
)

var today = utils.GetTodayDayString()

// 开启循环
func loop() {
	defer utils.Try()
	go taskOnSeconds()
	go taskOnDayChange()
}

// 每秒任务
func taskOnSeconds() {
	defer utils.Try()
	for {
		time.Sleep(time.Second * 10)
		services.SyncToRedis(false)
		services.FlushAccessLog()
		context.Purge(60)
	}
	//goland:noinspection GoUnreachableCode
	log.Errorln("taskOnSeconds exit")
}

// 日期变化任务
func taskOnDayChange() {
	defer utils.Try()
	for {
		time.Sleep(time.Millisecond * 100)
		day := utils.GetTodayDayString()
		if day == today {
			continue
		}
		log.Infof("taskOnDayChange: %s -> %s [%s]", today, day, time.Now())
		yesterday := today
		today = day
		services.FlushAccessLog()
		services.SyncToRedis(true)
		services.SyncToDB(yesterday)
	}
	//goland:noinspection GoUnreachableCode
	log.Errorln("taskOnDayChange exit")
}
