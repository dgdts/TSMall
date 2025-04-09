package goods

import (
	"context"

	goods "TSMall/hertz_gen/goods"
	"github.com/cloudwego/hertz/pkg/app"
	"ssg/trade-admin/biz/bizcontext"
)

type ApiHelloService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewApiHelloService(Context context.Context, RequestContext *app.RequestContext) *ApiHelloService {
	return &ApiHelloService{RequestContext: RequestContext, Context: Context}
}

func (h *ApiHelloService) Run(ctx *bizcontext.BizContext, req *goods.Empty) (resp *goods.Empty, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code
	// define your error in errno
	// if err != nil {
	// 	return nil, err
	//  建议提前定义好错误码，这样可以直接返回错误码
	//  return nil, errno.ServiceErr
	// }
	return
}
