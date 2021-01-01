package internal

import (
	"time"

	"github.com/gorilla/context"

	"github.com/yourtion/go-short-url/internal/services"
	"github.com/yourtion/go-short-url/internal/utils"
)

var today = utils.GetTodayDayString()

func loop() {
	defer utils.Try()
	go taskOnSeconds()
	go taskOnDayChange()
}

func taskOnSeconds() {
	defer utils.Try()
	for {
		time.Sleep(time.Second * 10)
		services.SyncToRedis(false)
		services.FlushAccessLog()
		context.Purge(60)
	}
	log.Errorln("taskOnSeconds exit")
}

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
	log.Errorln("taskOnDayChange exit")
}
