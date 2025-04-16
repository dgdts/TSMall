package errno

import (
	"errors"
	"fmt"
)

const (
	SuccessCode = 0

	GlobalErrorCode = 1000
)

// Global error
var (
	Success = NewErrNo(SuccessCode, "success")

	UnknownErr   = NewErrNo(GlobalErrorCode+1, "unknown error")
	ParameterErr = NewErrNo(GlobalErrorCode+2, "parameter error")
)

type ErrNo struct {
	ErrCode  int32  `json:"err_code"`
	ErrMsg   string `json:"err_msg"`
	ErrCause error  `json:"-"`
}

func (e ErrNo) Error() string {
	return fmt.Sprintf("err_code=%d, err_msg=%s", e.ErrCode, e.ErrMsg)
}

func NewErrNo(code int32, msg string) ErrNo {
	return ErrNo{code, msg, nil}
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

func (e ErrNo) WithCause(cause error) ErrNo {
	e.ErrCause = cause
	return e
}

func (e ErrNo) Is(target error) bool {
	targetErr, ok := target.(ErrNo)
	if !ok {
		return false
	}
	return e.ErrCode == targetErr.ErrCode
}

// ConvertErr convert error to Errno
func ConvertErr(err error) ErrNo {
	Err := ErrNo{}
	if errors.As(err, &Err) {
		return Err
	}

	s := UnknownErr
	s.ErrMsg = err.Error()
	return s
}
