package network

import (
	"fmt"
	"github.com/bingxindan/bxd_go_lib/logger"
	"github.com/bingxindan/bxd_go_lib/tools/confutil"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	CN = "cn"
	En = "en"
)

var Lang = CN

func SetLang(lang string) {
	Lang = lang
}

func GetLang() string {
	return Lang
}

func configInit() {
	lang := confutil.GetConfDefault("Lang", "lang", "cn")
	if lang != CN && lang != En {
		panic(fmt.Sprintf("can not support lang %v", lang))
	}
	SetLang(lang)
}

var EnErrorMap = map[int]string{
	10001: "Parameter check missing",
	10002: "Parameter check error",
	10101: "User name missing",
	10102: "User name error",
	10200: "Incorrect cell phone number",
	10300: "Incorrect phone number",
	10400: "Incorrect email address",
	20000: "Not logged in",
	20100: "Session timeout",
	20200: "You have been kicked out",
	20300: "Password had been changed",
	20400: "The login name has been changed",
	20500: "The phone number has been changed",
	30100: "Version not supported",
	30200: "Version not supported",
	30300: "Version not supported",
	40100: "No permission to view",
	40200: "No permission to modify",
	40300: "No permission to add",
	40400: "No permission to delete",
	50000: "system error",
	50100: "system not supported",
	50201: "Abnormal system connection",
	50202: "Abnormal system connection",
	50203: "Abnormal system connection",
	50401: "System connection timed out",
	50402: "System connection timed out",
	50403: "System connection timed out",
}

// Response the unified json structure
type response struct {
	Code    int         `json:"code"`
	Stat    int         `json:"stat"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func Success(v interface{}) interface{} {
	ret := response{Stat: 1, Code: 0, Message: "ok", Data: v}
	return ret
}

// JSON respond unified JSON structure with 200 http status code
func JSON(ctx *gin.Context, xe logger.BxdError, data interface{}) {
	Respond(ctx, http.StatusOK, xe, data)
}

// Respond encapsulates ctx.JSON
func Respond(ctx *gin.Context, status int, xe logger.BxdError, data interface{}) {
	respStat := 0
	if xe.Code == 0 {
		respStat = 1
	}
	if data == nil {
		data = gin.H{}
	}
	resp := response{
		Stat:    respStat,
		Code:    xe.Code,
		Message: xe.Msg,
		Data:    data,
	}
	ctx.JSON(status, resp)
}
