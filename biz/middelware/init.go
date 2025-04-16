package middelware

import (
	"context"
	"sort"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type Middelware interface {
	Init()
	GetOrder() int
	Do(ctx context.Context, c *app.RequestContext)
	Name() string
}

func getRegisterMiddleware() []Middelware {
	ret := []Middelware{
		&bizContextMiddelware{},
		&auth{},
	}

	sort.Slice(ret, func(i, j int) bool {
		return ret[i].GetOrder() < ret[j].GetOrder()
	})

	for _, m := range ret {
		m.Init()
	}

	return ret
}

func InitMiddeleware(s *server.Hertz) {
	middleware := getRegisterMiddleware()

	for _, m := range middleware {
		s.Use(m.Do)
		hlog.Infof("load middleware: %s", m.Name())
	}
}
