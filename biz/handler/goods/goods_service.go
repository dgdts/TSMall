package goods

import (
	"context"

	"TSMall/biz/common/response"
	logic "TSMall/biz/service/goods"
	_ "TSMall/hertz_gen/goods"
	"github.com/cloudwego/hertz/pkg/app"
)

// ApiHello .
// @router /api/v1/goods/list [POST]
func ApiHello(ctx context.Context, c *app.RequestContext) {
	response.JSON(ctx, c, logic.NewApiHelloService(ctx, c).Run)
}
