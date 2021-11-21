package bxd_errors

import (
	"errors"
	"fmt"
	"github.com/bingxindan/bxd_go_lib/logger"
	"github.com/spf13/cast"
	"strings"
)

var (
	Success         logger.BxdError = logger.BxdError{Code: BxdSuccessCode, Msg: "请求成功"}
	DefaultError    logger.BxdError = logger.BxdError{Code: BxdDefaultErrorCode, Msg: "请求失败"}
	BxdUnknownError logger.BxdError = logger.BxdError{Code: BxdUnknownErrorCode, Msg: "未知错误"}
)

func NewError(code int, msg string) logger.BxdError {
	return logger.BxdError{Code: code, Msg: msg}
}

func NewDiyError(code int, err error) error {
	msg := err.Error()
	if errCode, ok := Transfer(err); ok {
		msg = errCode.Msg
	}
	return errors.New(fmt.Sprintf("%d|%s", code, msg))
}

func Transfer(err error) (xCode logger.BxdError, ok bool) {
	segments := strings.SplitN(err.Error(), "|", 2)
	code := cast.ToInt(segments[0])
	if len(segments) > 1 && code > 0 {
		xCode.Code = code
		xCode.Msg = segments[1]
		ok = true
	}
	return
}
