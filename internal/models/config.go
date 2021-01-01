package models

import (
	"strings"

	jsoniter "github.com/json-iterator/go"

	"github.com/yourtion/go-short-url/internal/base/config"
	"github.com/yourtion/go-short-url/internal/base/define"
	"github.com/yourtion/go-short-url/internal/base/mysql"
)

// 从 Config 表获取一条信息
func getOneFromConfig(sql string, params ...interface{}) *define.ConfigRow {
	row := new(define.ConfigRow)
	ok := mysql.FindOne(nil, row, sql, params...)
	if ok {
		return row
	}
	return nil
}

// 通过配置名从 Config 表加载数据
func GetConfigDataByName(name string) *define.ConfigRow {
	return getOneFromConfig("SELECT data FROM "+define.TableConfig+" WHERE name=?", name)
}

// 加载基础动态配置
func LoadBaseConfig(dynamicInfo *config.DynamicInfo) {
	rows := make([]*define.ConfigRow, 0)
	// 配置 key 名字列表
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
	// 写入配置
	jsoniter.Get([]byte(info["white_list"])).ToVal(&dynamicInfo.WhiteList)
	jsoniter.Get([]byte(info["domains"])).ToVal(&dynamicInfo.Domains)
	log.Warnf("dynamicInfo [%+v]", dynamicInfo)
}
