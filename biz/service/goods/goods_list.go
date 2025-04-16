package goods

import (
	"context"

	"TSMall/biz/bizcontext"
	goods "TSMall/hertz_gen/goods"

	"github.com/cloudwego/hertz/pkg/app"
)

type GoodsListService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGoodsListService(Context context.Context, RequestContext *app.RequestContext) *GoodsListService {
	return &GoodsListService{RequestContext: RequestContext, Context: Context}
}

func (h *GoodsListService) Run(ctx *bizcontext.BizContext, req *goods.Empty) (resp *goods.Empty, err error) {
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
