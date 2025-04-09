package errno

import (
	"errors"
	"fmt"
)

// 基础错误码定义 10000开头, 0 ：表示正常 ====================================
const (
	SuccessCode    = 0
	ServiceErrCode = iota + 10000
	ParamErrCode
	RecordNotFoundErrCode
)

// 基础错误内容
const (
	SuccessMsg           = "Success"
	ServerErrMsg         = "Service is buzy now, Please try again later"
	ParamErrMsg          = "Wrong Parameter has been given"
	RecordNotFoundErrMsg = "Record not found"
)

// 基础自定义错误
var (
	Success           = NewErrNo(SuccessCode, SuccessMsg)
	ServiceErr        = NewErrNo(ServiceErrCode, ServerErrMsg)
	ParamErr          = NewErrNo(ParamErrCode, ParamErrMsg)
	RecordNotFoundErr = NewErrNo(RecordNotFoundErrCode, RecordNotFoundErrMsg)
)

type ErrNo struct {
	ErrCode int32  `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}

func (e ErrNo) Error() string {
	return fmt.Sprintf("err_code=%d, err_msg=%s", e.ErrCode, e.ErrMsg)
}

func NewErrNo(code int32, msg string) ErrNo {
	return ErrNo{code, msg}
}

func NewErrCode(code int32) ErrNo {
	return ErrNo{
		ErrCode: code,
		ErrMsg:  "",
	}
}

func (e ErrNo) WithMessage(msg string) ErrNo {
	e.ErrMsg = msg
	return e
}

// ConvertErr convert error to Errno
func ConvertErr(err error) ErrNo {
	Err := ErrNo{}
	if errors.As(err, &Err) {
		return Err
	}

	s := ServiceErr
	s.ErrMsg = err.Error()
	return s
}
