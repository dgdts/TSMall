package middelware

import (
	"context"

	"TSMall/biz/bizcontext"
	"TSMall/biz/constant"

	"github.com/cloudwego/hertz/pkg/app"
)

var _ Middelware = (*bizContextMiddelware)(nil)

type bizContextMiddelware struct{}

func (a *bizContextMiddelware) Init() {}

func (a *bizContextMiddelware) GetOrder() int {
	return 0
}

func (a *bizContextMiddelware) Name() string {
	return "bizcontext"
}

func (a *bizContextMiddelware) Do(ctx context.Context, c *app.RequestContext) {
	bizContext := &bizcontext.BizContext{
		Context: ctx,
	}
	c.Set(constant.BizContext, bizContext)
}
