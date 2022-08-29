package logger

import (
	"bytes"
	"context"
	"fmt"
	"github.com/bingxindan/bxd_go_lib/logger/builders"
	"github.com/bingxindan/bxd_go_lib/logger/logtrace"
	kratoslog "github.com/go-kratos/kratos/v2/log"
	"log"
	"os"
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

	time.Sleep(5 * time.Second)

	for i := 0; i < 1000; i++ {
		Ix(ctx, "tag", "Iaa: %+v", PrintlnTxt())
	}

	time.Sleep(100 * time.Second)
}

func PrintlnTxt() string {
	a := "aaaaa我的中国心我的中国心我的中国心我的中国心我的中国心我的中国心我的中国心我的中国心我的中国心" +
		"我的中国心我的中国心我的中国心我的中国心我的中国心我的中国心我的中国心我的中国心我的中国心我的中国心" +
		"我的中国心我的中国心我的中国心我的中国心我的中国心我的中国心我的中国心我的中国心我的中国心我的中国心" +
		"我的中国心我的中国心我的中国心我的中国心我的中国心我的中国心我的中国心我的中国心我的中国心我的中国心" +
		"我的中国心我的中国心我的中国心我的中国心我的中国心我的中国心我的中国心我的中国心我的中国心我的中国心" +
		"我的中国心我的中国心我的中国心我的中国心我的中国心我的中国心我的中国心我的中国心我的中国心我的中国心"
	return a
}

func TestXesLog(t *testing.T) {
	var (
		str = []byte("bbbb")
	)
	fmt.Println(string(str))

	fileName := "/Users/lauren/work/logs/jz_api.log"

	fd, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664)
	if err != nil {
		log.Fatalln(err)
	}

	buf := bytes.NewBuffer(make([]byte, 0, 4096))
	buf.Write(str)
	buf.WriteTo(fd)
}

func TestKratosLog(t *testing.T) {
	f, err := os.OpenFile("/Users/lauren/work/logs/kratos.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return
	}

	std := kratoslog.NewStdLogger(f)

	h := kratoslog.NewHelper(std)

	h.Infof("aaaaaaa: %+v", "bbbbb")
}
