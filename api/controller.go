package api

import (
	"github.com/bingxindan/bxd_go_lib/gokit/network"
	"github.com/bingxindan/bxd_go_lib/logger"
	"github.com/bingxindan/bxd_go_lib/tools/bxd_errors"
	"github.com/gin-gonic/gin"
)

type BaseController struct {
}

// Success 成功的格式化输出
func (ctl *BaseController) Success(ctx *gin.Context, data interface{}) {
	ctl.output(ctx, bxd_errors.Success, data)
}

// Error 自定义错误码，格式化输出
func (ctl *BaseController) Error(ctx *gin.Context, code int, message string, data interface{}) {
	xe := bxd_errors.NewError(code, message)
	ctl.output(ctx, xe, data)
}

// 处理data
func (ctl *BaseController) output(ctx *gin.Context, xe logger.BxdError, data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	network.JSON(ctx, xe, data)
}
