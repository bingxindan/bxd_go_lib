package logger

import (
	"context"
	"fmt"
	"github.com/bingxindan/bxd_go_lib/logger/builders"
	"github.com/bingxindan/bxd_go_lib/logger/logtrace"
	"testing"
	"time"
)

func TestLogger(m *testing.T) {
	config := NewLogConfig()
	fmt.Printf("c: %+v\n", config)
	config.LogPath = "/Users/lauren/work/logs/jz_api.log"
	config.SuffixEnv = "DEV"
	InitLogWithConfig(config)
	//或使用xml配置 logger.InitLogger("conf/log.xml")
	defer Close()
	builder := new(builders.TraceBuilder)
	builder.SetTraceDepartment("HS-Golang")
	builder.SetTraceVersion("0.1")
	SetBuilder(builder)

	//初始化trace信息 一次完整调用前执行
	ctx := context.WithValue(context.Background(), "start", time.Now())
	ctx = context.WithValue(ctx, logtrace.GetMetadataKey(), logtrace.GenLogTraceMetadata())

	Ix(ctx, "tag", "Iaa: %+v", 2222)

	time.Sleep(5 * time.Second)

	for i := 0; i < 1000; i++ {
		Ix(ctx, "tag", "Iaa: %+v", 3333)
	}

	time.Sleep(10 * time.Second)
}
