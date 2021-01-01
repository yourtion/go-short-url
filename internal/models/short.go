package models

import (
	"github.com/yourtion/go-short-url/internal/base/define"
	"github.com/yourtion/go-short-url/internal/base/mysql"
	"github.com/yourtion/go-short-url/internal/utils"
)

// 从Short表中查一条数据
func getOneFromShort(sql string, params ...interface{}) *define.ShortRow {
	row := new(define.ShortRow)
	ok := mysql.FindOne(nil, row, sql, params...)
	if ok {
		return row
	}
	return nil
}

// 通过hash查询短链接是否存在
func getShortInfoByHash(hash string) *define.ShortRow {
	return getOneFromShort("SELECT `id`,`short`,`is_statistics`,`is_access_log` FROM `"+define.TableShort+"` WHERE hash=?", hash)
}

// 通过短链获取链接信息
func GetUrlInfoByShort(short string) *define.ShortRow {
	return getOneFromShort("SELECT `id`,`origin`,`is_active`,`is_statistics`,`is_access_log` FROM `"+define.TableShort+"` WHERE short=?", short)
}

// 创建一条短链接
func AddUrlToShort(origin string) *define.ShortRow {
	row := new(define.ShortRow)

	// 先判断链接hash是否已经存在
	hash := utils.MD5(origin)
	shortInfo := getShortInfoByHash(hash)
	if shortInfo != nil {
		return shortInfo
	}

	// 生成短链接并插入数据库
	short := utils.GenerateShort()
	id := mysql.InsertOne(nil, "INSERT INTO "+define.TableShort+"(short, origin, hash) VALUES (?, ?, ?)", short, origin, hash)
	if id < 0 {
		return nil
	}
	row.Id = id
	row.ActivityId = 0
	row.IsAccessLog = true
	row.IsStatistics = false
	row.Hash = hash
	row.Short = short
	row.Origin = origin
	return row
}
