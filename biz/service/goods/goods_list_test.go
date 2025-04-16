package goods

import (
	"context"
	"testing"

	"TSMall/biz/bizcontext"
	goods "TSMall/hertz_gen/goods"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/test/assert"
)

func TestGoodsListService_Run(t *testing.T) {
	ctx := context.Background()
	c := app.NewContext(1)
	s := NewGoodsListService(ctx, c)
	// init req and assert value
	req := &goods.Empty{}
	bizctx := &bizcontext.BizContext{}
	resp, err := s.Run(bizctx, req)
	assert.DeepEqual(t, nil, resp)
	assert.DeepEqual(t, nil, err)
	// todo edit your unit test.
}
