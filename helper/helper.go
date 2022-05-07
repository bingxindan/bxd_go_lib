package helper

import (
	"context"
	"github.com/bingxindan/bxd_go_lib/logger"
	"github.com/pkg/errors"
	"math/rand"
	"net"
	"strconv"
	"time"
)

// @Desc 当前时间
// @Date 2021-05-12
// @Param
// @Return
// @Author zhangming16
func GetCurTime() string {
	loc, _ := time.LoadLocation("Asia/Shanghai") //上海
	curTime := time.Now().In(loc).Format("2006-01-02 15:04:05")

	return curTime
}

// @Desc 获取当前时间戳
// @Date 2021-02-14
// @Param
// @Return
// @Author zhangming16
func GetCurTimestampStr() string {
	loc, _ := time.LoadLocation("Asia/Shanghai") //上海
	time := time.Now().In(loc).Unix()
	return strconv.FormatInt(time, 10)
}

// @Desc 获取指定格式时间
// @Date 2021-06-03
// @Param
// @Return
// @Author zhangming16
func GetTimeFormat(format string) (curTime string) {
	loc, _ := time.LoadLocation("Asia/Shanghai") //上海
	obj := time.Now().In(loc)

	switch format {
	case "day":
		curTime = obj.Format("20060102")
	}

	return curTime
}

// @Desc 获取当前时间戳
// @Date 2021-06-03
// @Param
// @Return
// @Author zhangming16
func GetTimestamp() int64 {
	loc, _ := time.LoadLocation("Asia/Shanghai") //上海
	time := time.Now().In(loc).Unix()
	return time
}

// @Desc 时间戳时字符串改成整型
// @Date 2021-06-06
// @Param
// @Return
// @Author zhangming16
func ConvertStrToTimestamp(timeStr string) int64 {
	timeLayout := "2006-01-02 15:04:05"                          //转化所需模板
	loc, _ := time.LoadLocation("Local")                         //重要：获取时区
	theTime, _ := time.ParseInLocation(timeLayout, timeStr, loc) //使用模板在对应时区转化为time.time类型
	toTime := theTime.Unix()                                     //转化为时间戳 类型是int64
	return toTime
}

// @Desc 生成随机字符串
// @Date 2022-01-05
// @Param
// @Return
// @Author zhangming
func GetRandomString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// desc 获取服务器IP
// @author zhangming
// @date 2022/5/7-18:15
func GetSrvIP(ctx context.Context) (ip string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		logger.Ex(ctx, "helper.GetSrvIP", "net.InterfaceAddrs fail, err: %+v", err)
		return "", err
	}

	for _, addr := range addrs {
		// 检查ip地址判断是否回环地址
		if inet, ok := addr.(*net.IPNet); ok && !inet.IP.IsLoopback() {
			if inet.IP.To4() != nil {
				return inet.IP.String(), nil
			}
		}
	}

	return "", errors.New("服务器IP获取失败！")
}
