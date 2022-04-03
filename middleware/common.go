package middleware

import (
	"bytes"
	"context"
	"fmt"
	"github.com/bingxindan/bxd_go_lib/bxdgin"
	"github.com/bingxindan/bxd_go_lib/logger"
	"github.com/bingxindan/bxd_go_lib/logger/logtrace"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http/httputil"
	"os"
	"runtime"
	"strconv"
	"time"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

type CheckError func([]byte) bool

var errCheck CheckError

func CheckStat1(in []byte) bool {
	return (bytes.Index(in, []byte(`"stat":1`)) < 0) && (bytes.Index(in, []byte(`"stat": 1`)) < 0)
}

func Logger(fns ...CheckError) gin.HandlerFunc {
	if len(fns) > 0 {
		errCheck = fns[0]
	}
	return func(c *gin.Context) {
		blw := &bodyLogWriter{body: bytes.NewBuffer([]byte{}), ResponseWriter: c.Writer}
		c.Writer = blw
		// start time
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		var body []byte
		if c.Request.Body != nil {
			body, _ = ioutil.ReadAll(c.Request.Body)
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		// process request
		c.Next()

		_, skip := c.Get("SKIPLOG")
		if skip {
			return
		}

		// stop timer
		end := time.Now()
		latency := end.Sub(start)

		statusCode := c.Writer.Status()
		if raw != "" {
			path = path + "?" + raw
		}

		buf := blw.body.Bytes()
		ctx := bxdgin.TransferToContext(c)
		meta := ctx.Value(logtrace.GetMetadataKey())
		if metadata, ok := meta.(*logtrace.TraceNode); ok {
			xrf := "\"%s\""
			if len(buf) > 0 && bytes.HasPrefix(buf, []byte("{")) {
				xrf = "%s"
			}
			metadata.Set("x_response", fmt.Sprintf(xrf, buf))
			metadata.Set("x_status", fmt.Sprintf("\"%s\"", strconv.Itoa(statusCode)))
			metadata.Set("x_request_time", fmt.Sprintf("\"%s\"", fmt.Sprintf("%v", latency)))
			ctx = context.WithValue(ctx, logtrace.GetMetadataKey(), metadata)
		}
		if errCheck != nil && errCheck(buf) {
			logger.Ex(ctx, "[GIN]", "request failed, response: %s", buf)
		} else {
			logger.Ix(ctx, "[GIN]", "request success")
		}
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	if _, err := w.body.Write(b); err != nil {
		fmt.Printf("bodyLogWriter err: %v", err)
	}
	return w.ResponseWriter.Write(b)
}

// 创建logid
func HandleCtxInfo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("logid", strconv.FormatInt(logger.Id(), 10))
		ctx.Set("start", time.Now())
	}
}

func SkipLogInfo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("SKIPLOG", "1")
	}
}

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := stack(3)
				httprequest, _ := httputil.DumpRequest(c.Request, false)
				ctx := bxdgin.TransferToContext(c)
				logger.Ex(ctx, "[Recovery]", "%s panic recovered:\n%s\n%s\n%s", time.Now().Format("2006/01/02 - 15:04:05"), string(httprequest), err, string(stack))
				_, _ = fmt.Fprintf(os.Stdout, "[Recovery] %s panic recovered:\n%s\n%s\n%s", time.Now().Format("2006/01/02 - 15:04:05"), string(httprequest), err, string(stack))
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}

func stack(skip int) []byte {
	var (
		// the returned data
		buf = new(bytes.Buffer)

		// As we loop, we open files and read them. These variables record the currently
		// loaded file.
		lines    [][]byte
		lastFile string
	)
	// Skip the expected number of frames
	for i := skip; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//  runtime/debug.*T·ptrmethod
	// and want
	//  *T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

func source(lines [][]byte, n int) []byte {
	// in stack trace, lines are 1-indexed but our array is 0-indexed
	n--
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}
