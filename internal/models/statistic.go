package models

import (
	"github.com/yourtion/go-short-url/internal/base/define"
	"github.com/yourtion/go-short-url/internal/base/mysql"
)

// 创建或更新 UV、PV 记录
func CreateOrUpdateStatisticInfo(info *define.StatisticRow) bool {
	sql := "INSERT INTO " + define.TableStatistic + " (sid,day,pv,uv) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE pv=pv+?, uv=uv+?;"
	id := mysql.InsertOne(nil, sql, info.SId, info.Day, info.PV, info.UV, info.PV, info.UV)
	return id > 0
}
