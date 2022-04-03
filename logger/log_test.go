package logger

import (
	"context"
	"github.com/bingxindan/bxd_go_lib/logger/builders"
	"testing"
	"time"
)

func TestLogger(m *testing.T) {
	config := NewLogConfig()
	config.LogPath = "/Users/lauren/work/logs/jz_api.log"
	InitLogWithConfig(config)
	//或使用xml配置 logger.InitLogger("conf/log.xml")
	defer Close()
	builder := new(builders.TraceBuilder)
	builder.SetTraceDepartment("HS-Golang")
	builder.SetTraceVersion("0.1")
	SetBuilder(builder)

	Ix(context.Background(), "tag", "Iaa: %+v", 2222)

	time.Sleep(5 * time.Second)

	for i := 0; i < 1000; i++ {
		Ix(context.Background(), "tag", "Iaa: %+v", 3333)
	}

	time.Sleep(10 * time.Second)
}
