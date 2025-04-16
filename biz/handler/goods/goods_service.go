package goods

import (
	"context"

	"TSMall/biz/response"
	logic "TSMall/biz/service/goods"
	_ "TSMall/hertz_gen/goods"

	"github.com/cloudwego/hertz/pkg/app"
)

// GoodsList .
// @router /api/v1/goods/list [POST]
func GoodsList(ctx context.Context, c *app.RequestContext) {
	response.JSON(ctx, c, logic.NewGoodsListService(ctx, c).Run)
}
