package mysql

import (
	"database/sql"
	"fmt"
	"net/url"
	"sync/atomic"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/yourtion/go-short-url/internal/base/config"
	"github.com/yourtion/go-short-url/internal/base/define"
	"github.com/yourtion/go-short-url/internal/base/logger"
)

var DB *sqlx.DB
var log *logger.Entry
var queryCounter int64

func incrQueueCounter() {
	atomic.AddInt64(&queryCounter, 1)
}

// 打开数据库连接
func Open(opts *config.MySQLConfig) {
	log = logger.NewModuleLogger("mysql")

	// 检查当前服务必须在 +8:00 时区运行
	location, offset := time.Now().Zone()
	if offset != 28800 {
		log.Fatalf("invalid time zone: %s, this program must be run in Asia/Shanghai (+8:00) timezone", location)
	}

	// 连接到数据库
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local&time_zone=%s",
		opts.User,
		opts.Password,
		opts.Host,
		opts.Port,
		opts.Database,
		opts.Charset,
		url.QueryEscape("'+8:00'"),
	)
	log.Debugf("connect to %s, SetMaxOpenConns=%d", connStr, config.Config.MySQL.MaxConnections)
	DB = sqlx.MustConnect("mysql", connStr)
	DB.SetMaxOpenConns(config.Config.MySQL.MaxConnections)

	// 生成表名称
	define.TableConfig = opts.Prefix + "config"
	define.TableShort = opts.Prefix + "short"
	define.TableStatistic = opts.Prefix + "statistic"
}

// 查询一条数据
func FindOne(tx *sqlx.Tx, dest interface{}, query string, args ...interface{}) (success bool) {
	incrQueueCounter()
	var err error
	log.Debugf("#%d FindOne: %s %+v", queryCounter, query, args)
	if tx != nil {
		err = tx.Get(dest, query, args...)
	} else {
		err = DB.Get(dest, query, args...)
	}
	if err != nil {
		if err != sql.ErrNoRows {
			log.Warningf("#%d FindOne failed: %s => %s %+v", queryCounter, err, query, args)
		}
		log.Debugf("#%d FindMany: success=false", queryCounter)
		return false
	}
	log.Debugf("#%d FindMany: success=true", queryCounter)
	return true
}

// 查询多条数据
func FindMany(tx *sqlx.Tx, dest interface{}, query string, args ...interface{}) (success bool) {
	incrQueueCounter()
	var err error
	log.Debugf("%#d FindMany: %s %+v", queryCounter, query, args)
	if tx != nil {
		err = tx.Select(dest, query, args...)
	} else {
		err = DB.Select(dest, query, args...)
	}
	if err != nil {
		if err != sql.ErrNoRows {
			log.Warningf("#%d FindMany failed: %s => %s %+v", queryCounter, err, query, args)
		}
		log.Debugf("#%d FindMany: success=false", queryCounter)
		return false
	}
	log.Debugf("#%d FindMany: success=true", queryCounter)
	return true
}

// 插入一条数据
func InsertOne(tx *sqlx.Tx, query string, args ...interface{}) (insertId int) {
	incrQueueCounter()
	var err error
	var res sql.Result
	log.Debugf("#%d InsertOne: %s %+v", queryCounter, query, args)
	if tx != nil {
		res, err = tx.Exec(query, args...)
	} else {
		res, err = DB.Exec(query, args...)
	}
	if err != nil {
		log.Warningf("#%dInsertOne failed: %s => %s %+v", queryCounter, err, query, args)
		return 0
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Warningf("#%d InsertOne failed: %s => %s %+v", queryCounter, err, query, args)
	}
	insertId = int(id)
	log.Debugf("#%d InsertOne: insertId=%d", queryCounter, insertId)
	return insertId
}

// 更新多条数据
func UpdateMany(tx *sqlx.Tx, query string, args ...interface{}) (rowsAffected int) {
	incrQueueCounter()
	var err error
	var res sql.Result
	log.Debugf("#%d UpdateMany: %s %+v", queryCounter, query, args)
	if tx != nil {
		res, err = tx.Exec(query, args...)
	} else {
		res, err = DB.Exec(query, args...)
	}
	if err != nil {
		log.Warningf("UpdateMany failed: %s => %s %+v", err, query, args)
		return 0
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Warningf("UpdateMany failed: %s => %s %+v", err, query, args)
	}
	rowsAffected = int(rows)
	log.Debugf("#%d UpdateMany: rowsAffected=%d", queryCounter, rowsAffected)
	return rowsAffected
}

// 更新一条数据
func UpdateOne(tx *sqlx.Tx, query string, args ...interface{}) (rowsAffected int) {
	incrQueueCounter()
	rowsAffected = UpdateMany(tx, query+" LIMIT 1", args...)
	return rowsAffected
}

type queryCountRow struct {
	Count int `db:"count"`
}

// 查询记录数量，需要 SELECT count(*) AS count FROM ... 这样的格式
func FindCount(tx *sqlx.Tx, query string, args ...interface{}) (count int, success bool) {
	row := new(queryCountRow)
	ok := FindOne(tx, row, query, args...)
	if ok {
		return row.Count, true
	}
	return 0, false
}
