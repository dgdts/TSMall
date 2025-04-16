package jsonresult

import (
	"TSMall/biz/errno"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type statusResult struct {
	Code  int32  `json:"code"`
	Msg   string `json:"msg"`
	BzMsg string `json:"bzMsg"`
}
type JSONResult struct {
	statusResult
	Data any `json:"data"`
}

func NewJSONSuccessResult(data any) *JSONResult {
	return &JSONResult{
		Data: data,
		statusResult: statusResult{
			Code: errno.Success.ErrCode,
			Msg:  errno.Success.ErrMsg,
		},
	}
}

func NewJSONErrResult(err error) *JSONResult {
	r := &JSONResult{}
	if err == nil {
		return r
	}
	if e, ok := err.(errno.ErrNo); ok {
		r.Code = e.ErrCode
		r.Msg = e.ErrMsg
		return r
	}
	hlog.Errorf("unexpected err: %+v", err)
	r.Msg = errno.UnknownErr.ErrMsg
	r.Code = errno.UnknownErr.ErrCode
	return r
}
