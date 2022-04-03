package bootstrap

import (
	"github.com/bingxindan/bxd_go_lib/logger"
	"github.com/bingxindan/bxd_go_lib/logger/builders"
	"github.com/bingxindan/bxd_go_lib/tools/confutil"
)

func InitBxdLogger(department, version string) {
	section := "log"
	//通过配置文件转为map[string]string
	logMap := confutil.GetConfStringMap(section)
	if len(logMap) > 0 {
		config := logger.NewLogConfig()
		config.SetConfigMap(logMap)
		logger.InitLogWithConfig(config)
	} else {
		logger.InitLogger("")
	}
	builder := new(builders.TraceBuilder)
	builder.SetTraceDepartment(department)
	builder.SetTraceVersion(version)
	logger.SetBuilder(builder)
	return
}
