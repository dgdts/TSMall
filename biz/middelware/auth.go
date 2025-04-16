package middelware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

var _ Middelware = (*auth)(nil)

type auth struct{}

func (a *auth) Init() {}

func (a *auth) GetOrder() int {
	return 1
}

func (a *auth) Name() string {
	return "auth"
}

var excludePath = map[string]struct{}{
	"/ping": {},
}

func (a *auth) Do(ctx context.Context, c *app.RequestContext) {
	// check exclude path
	if _, ok := excludePath[string(c.URI().Path())]; ok {
		return
	}
	// do auth
}
