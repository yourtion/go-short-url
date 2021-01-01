package models

import (
	"strings"

	jsoniter "github.com/json-iterator/go"

	"github.com/yourtion/go-short-url/internal/base/config"
	"github.com/yourtion/go-short-url/internal/base/define"
	"github.com/yourtion/go-short-url/internal/base/mysql"
)

// 从Config表获取一条信息
func getOneFromConfig(sql string, params ...interface{}) *define.ConfigRow {
	row := new(define.ConfigRow)
	ok := mysql.FindOne(nil, row, sql, params...)
	if ok {
		return row
	}
	return nil
}

func GetConfigDataByName(name string) *define.ConfigRow {
	return getOneFromConfig("SELECT data FROM "+define.TableConfig+" WHERE name=?", name)
}

func LoadBaseConfig(dynamicInfo *config.DynamicInfo) {
	rows := make([]*define.ConfigRow, 0)
	keys := []string{"white_list", "domains"}
	ok := mysql.FindMany(nil, &rows, "SELECT name, data FROM "+define.TableConfig+" WHERE name in ('"+strings.Join(keys, "','")+"')")
	if !ok {
		return
	}
	info := map[string]string{}
	for _, row := range rows {
		info[row.Name] = row.Data
	}
	for _, key := range keys {
		if info[key] == "" {
			info[key] = "{}"
		}
	}
	jsoniter.Get([]byte(info["white_list"])).ToVal(&dynamicInfo.WhiteList)
	jsoniter.Get([]byte(info["domains"])).ToVal(&dynamicInfo.Domains)
	log.Warnf("%+v [%+v]", dynamicInfo)
}
