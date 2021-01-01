package define

import (
	"time"
)

// 配置表
var TableConfig string

// 短链映射表
var TableShort string

// 统计信息表
var TableStatistic string

type ShortRow struct {
	Id           int       `db:"id" json:"id"`
	Short        string    `db:"short" json:"short"`
	Origin       string    `db:"origin" json:"origin"`
	Hash         string    `db:"hash" json:"hash"`
	IsActive     bool      `db:"is_active" json:"is_active"`
	IsStatistics bool      `db:"is_statistics" json:"is_statistics"`
	IsAccessLog  bool      `db:"is_access_log" json:"is_access_log"`
	ActivityId   int       `db:"activity_id" json:"activity_id"`
	CreateTime   time.Time `db:"created_at" json:"created_at"`
	UpdateTime   time.Time `db:"updated_at" json:"updated_at"`
}

type StatisticRow struct {
	Id  int `db:"id" json:"id"`
	SId int `db:"sid" json:"sid"`
	Day int `db:"day" json:"day"`
	PV  int `db:"pv" json:"pv"`
	UV  int `db:"uv" json:"uv"`
}

type ConfigRow struct {
	Name       string    `db:"name" json:"name"`
	Note       string    `db:"note" json:"note"`
	Data       string    `db:"data" json:"data"`
	Schema     string    `db:"schema" json:"schema"`
	CreateTime time.Time `db:"created_at" json:"created_at"`
	UpdateTime time.Time `db:"updated_at" json:"updated_at"`
}
