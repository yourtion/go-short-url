package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// 默认起始时间是 20190101
var START = time.Unix(1546272000, 0)

// 生成短链接URL（四位随机数+十二位距离START毫秒数）=> HEX
func GenerateShort() string {
	now := time.Now()
	rand.Seed(now.UnixNano())

	t := now.Sub(START).Nanoseconds() / 1000000
	i1 := rand.Intn(9999)
	s := fmt.Sprintf("%04d%012d", i1, t)
	v, _ := strconv.ParseInt(s, 10, 64)

	return strconv.FormatInt(v, 36)
}

// 生成短链接URL（八位随机数+十二位距离START毫秒数）=> HEX
func GenerateUid() string {
	now := time.Now()
	rand.Seed(now.UnixNano())

	t := now.Sub(START).Nanoseconds() / 1000000
	i := rand.Intn(99999999)
	s := fmt.Sprintf("%012d%04d", t, i)
	v, _ := strconv.ParseInt(s, 10, 64)

	return strconv.FormatInt(v, 36)
}

// 获取日期字符串
func GetDayString(t time.Time) string {
	return t.Format("20060102")
}

// 获取当前日期字符串
func GetTodayDayString() string {
	return GetDayString(time.Now())
}
