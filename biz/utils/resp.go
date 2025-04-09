package utils

import (
	"TSMall/biz/errno"
	"context"
	"errors"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

// SendErrResponse  pack error response
func SendErrResponse(ctx context.Context, c *app.RequestContext, code int, err error) {
	// todo edit custom code
	c.String(code, err.Error())
}

// SendSuccessResponse  pack success response
func SendSuccessResponse(ctx context.Context, c *app.RequestContext, code int, data interface{}) {
	// todo edit custom code
	c.JSON(code, data)
}

type BaseResp struct {
	Code int32       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// BuildBaseErr convert error and build BaseResp 接口返回时错误处理
func BuildBaseErr(c *app.RequestContext, err error) *BaseResp {
	if err == nil {
		return baseResp(c, errno.Success)
	}

	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return baseResp(c, e)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return baseResp(c, s)
}

// BuildBaseResp convert data and build BaseResp 接口返回时数据处理
func BuildBaseResp(c *app.RequestContext, data interface{}) *BaseResp {
	return &BaseResp{
		Code: errno.SuccessCode,
		Msg:  "success",
		Data: data,
	}
}

// BaseResp build BaseResp from error
func baseResp(c *app.RequestContext, err errno.ErrNo) *BaseResp {
	// 监控埋点
	c.Response.Header.Set("bizStatusCode", strconv.Itoa(int(err.ErrCode)))
	return &BaseResp{
		Code: err.ErrCode,
		Msg:  err.ErrMsg,
	}
}
