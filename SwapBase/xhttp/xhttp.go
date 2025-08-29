package xhttp

import (
	"context"

	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/ProjectsTask/SwapBase/errcode"

	"github.com/ProjectsTask/SwapBase/kit/convert"

	"github.com/ProjectsTask/SwapBase/logger/xzap"
)

// Error 错误响应返回
func Error(c *gin.Context, err error) {
	ctx := c.Request.Context()
	e := errcode.ParseErr(err)
	if e == errcode.ErrUnexpected || e == errcode.ErrCustom {
		xzap.WithContext(ctx).Error("request handle err",
			zap.Error(err),
			zap.Uint32("code", e.Code()),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery))
	}

	WriteHeader(c.Writer, e)
	c.JSON(e.HTTPCode(), &Response{
		TraceId: GetTraceId(ctx),
		Code:    e.Code(),
		Msg:     e.Error(),
		Data:    nil,
	})
}

// Response 业务通用响应体
type Response struct {
	TraceId string      `json:"trace_id" example:"a1b2c3d4e5f6g7h8" extensions:"x-order=000"` // 链路追踪id
	Code    uint32      `json:"code" example:"200" extensions:"x-order=001"`                  // 状态码
	Msg     string      `json:"msg" example:"OK" extensions:"x-order=002"`                    // 消息
	Data    interface{} `json:"data" extensions:"x-order=003"`                                // 数据
}

// GetTraceId 获取链路追踪id
func GetTraceId(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		return spanCtx.TraceID().String()
	}

	return ""
}

// WriteHeader 写入自定义响应header
func WriteHeader(w http.ResponseWriter, err ...error) {
	var ee error
	if len(err) > 0 {
		ee = err[0]
	}

	e := errcode.ParseErr(ee)
	w.Header().Set(HeaderGWErrorCode, convert.ToString(e.Code()))
	w.Header().Set(HeaderGWErrorMessage, url.QueryEscape(e.Error()))
}
