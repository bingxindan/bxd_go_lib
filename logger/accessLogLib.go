package logger

import (
	"fmt"
	"github.com/bingxindan/bxd_go_lib/logger/builders"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"os"
)

/*

方法：AccessLogLib 接入loglib
功能：接入多样logLib 以支持其他组件中 loggerX 日志打印格式
@param libName 接入logLib库名 默认接入log4go
@param args[0] logrus全局logger实例 使用默认的logrus Builder
@param args[0] zap全局logger实例 使用默认的zap Builder
*/

func AccessLogLib(libName string, args interface{}) {

	if libName == "logrus" {
		instance, ok := args.(*logrus.Logger)
		if !ok {
			fmt.Fprintf(os.Stderr, "SetLogLib Error: Could not get logrus instance!")
			os.Exit(1)
		}
		libBuild := builders.NewLogrusBuilder(instance)
		SetBuilder(libBuild)
	} else if libName == "zap" {

		instance, ok := args.(*zap.Logger)
		if !ok {
			fmt.Fprintf(os.Stderr, "SetLogLib Error: Could not get zap instance!")
			os.Exit(1)
		}
		libBuild := builders.NewZapBuilder(instance)
		SetBuilder(libBuild)
	} else {
		libName = "log4go"
	}
	return
}
