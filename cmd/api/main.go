package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"TSMall/biz/config/global_config"
	"TSMall/biz/global_init"
	"TSMall/biz/router"
	_ "TSMall/kit/xjsoniter"

	protocolconfig "github.com/hertz-contrib/http2/config"
	"github.com/hertz-contrib/http2/factory"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func env() {
	env := strings.ToLower(os.Getenv("ENV"))
	if env == "" {
		env = "dev"
	}
	_ = os.Setenv("ENV", env)
}

func main() {
	env()

	workspacePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	confPath := filepath.Join(workspacePath, "configs", os.Getenv("ENV"), "conf.yaml")

	// init global config
	global_config.InitGlobalConfig(confPath)

	// init hertz server
	h := global_init.InitServer(global_config.GConfig.Hertz.Service)

	h.AddProtocol("h2", factory.NewServerFactory(
		protocolconfig.WithReadTimeout(time.Minute),
		protocolconfig.WithDisableKeepAlive(false),
	))

	// add a ping route to test
	h.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, utils.H{
			"ping":   "pong",
			"pod_ip": os.Getenv("POD_IP"),
		})
	})

	// register biz router
	router.GeneratedRegister(h)
	// server start
	h.Spin()
}
