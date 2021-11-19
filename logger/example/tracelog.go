package example

import (
	"context"
	"github.com/bingxindan/bxd_go_lib/logger"
	"github.com/bingxindan/bxd_go_lib/logger/builders"
	"github.com/bingxindan/bxd_go_lib/logger/logtrace"
	"strconv"
	"time"
)

func example_tracelog() {
	//初始化logger main函数指定
	config := logger.NewLogConfig()
	logger.InitLogWithConfig(config)
	//或使用xml配置 logger.InitLogger("conf/log.xml")
	defer logger.Close()
	builder := new(builders.TraceBuilder)
	builder.SetTraceDepartment("HS-Golang")
	builder.SetTraceVersion("0.1")
	logger.SetBuilder(builder)

	//初始化trace信息 一次完整调用前执行
	ctx := context.WithValue(context.Background(), "logid", strconv.FormatInt(logger.Id(), 10))
	ctx = context.WithValue(ctx, "start", time.Now())
	ctx = context.WithValue(ctx, logtrace.GetMetadataKey(), logtrace.GenLogTraceMetadata())

	//logger
	logger.Ix(ctx, "Example", "example log time:%v,module:%s", time.Now(), "test1")
	logger.Ex(ctx, "Example", "example log time:%v,module:%s", time.Now(), "test2")
	logger.Wx(ctx, "Example", "example log time:%v,module:%s", time.Now(), "test3")

}
