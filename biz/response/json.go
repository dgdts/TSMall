package response

import (
	"TSMall/biz/bizcontext"
	"TSMall/biz/constant"
	"TSMall/biz/response/jsonresult"
	"context"
	"errors"
	"net/http"
	"strconv"

	"TSMall/biz/errno"

	"github.com/cloudwego/hertz/pkg/app"
)

const (
	bizStatusCodeName = "bizStatusCode"
)

type JSONHandler[Req any, Res any] func(c *bizcontext.BizContext, req *Req) (*Res, error)

func JSONErr(c *app.RequestContext, err error) {
	resp := jsonresult.NewJSONErrResult(err)

	c.Set("err", err.Error())
	// bizStatusCode used for prometheus monitor
	c.Response.Header.Set(bizStatusCodeName, strconv.Itoa(int(resp.Code)))
	c.JSON(http.StatusOK, resp)
}

func JSON[Req, Res any](ctx context.Context, c *app.RequestContext, handler JSONHandler[Req, Res]) {
	req := new(Req)
	err := c.BindAndValidate(req)
	if err != nil {
		JSONErr(c, errno.ParameterErr)
		return
	}
	rawBizCtx, ok := c.Get(constant.BizContext)
	if !ok {
		JSONErr(c, errors.New("bizcontext not found"))
		return
	}

	bizCtx, ok := rawBizCtx.(*bizcontext.BizContext)
	if !ok {
		JSONErr(c, errors.New("bizcontext type error"))
		return
	}

	resp, err := handler(bizCtx, req)
	if err != nil {
		JSONErr(c, err)
		return
	}
	c.JSON(http.StatusOK, jsonresult.NewJSONSuccessResult(resp))
}
