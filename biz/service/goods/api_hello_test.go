package goods

import (
	"context"
	"testing"

	goods "TSMall/hertz_gen/goods"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/test/assert"
	"ssg/trade-admin/biz/bizcontext"
)

func TestApiHelloService_Run(t *testing.T) {
	ctx := context.Background()
	c := app.NewContext(1)
	s := NewApiHelloService(ctx, c)
	// init req and assert value
	req := &goods.Empty{}
	bizctx := &bizcontext.BizContext{}
	resp, err := s.Run(bizctx, req)
	assert.DeepEqual(t, nil, resp)
	assert.DeepEqual(t, nil, err)
	// todo edit your unit test.
}
