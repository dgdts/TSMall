package global_init

import (
	"TSMall/biz/r_conf"
	"fmt"

	"TSMall/biz/constant"
	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/cloudwego/hertz/pkg/app/server/render"
	"github.com/hertz-contrib/cors"
	"github.com/hertz-contrib/gzip"
	"github.com/hertz-contrib/logger/accesslog"
	"github.com/hertz-contrib/pprof"
	jsoniter "github.com/json-iterator/go"

	"ssg/kitex-common/common/biz_init"
	"ssg/kitex-common/common/local_conf"
	"ssg/kitex-common/suite/server_suite"

	"github.com/cloudwego/hertz/pkg/app/server"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitServer(service local_conf.Service) (h *server.Hertz) {
	// log初始化
	biz_init.InitHLog(&lumberjack.Logger{
		Filename:   local_conf.GetGConf().Log.LogFileName,
		MaxSize:    local_conf.GetGConf().Log.LogMaxSize,
		MaxBackups: local_conf.GetGConf().Log.LogMaxBackups,
		MaxAge:     local_conf.GetGConf().Log.LogMaxAge,
		Compress:   local_conf.GetGConf().Log.LogCompress,
	}, local_conf.HLogLevel(), local_conf.GetGConf().Log.LogMode)

	// 初始化HertzServer, 这里与注册中心强绑定, 如果不需要注册中心, 可以将此行注释掉自行实现
	h = server.New(server_suite.HertzCommonServerSuite{
		Address: local_conf.GetGConf().Hertz.Service[0].Address,
		CurrentServiceName: fmt.Sprintf("%s.%s.%s.%s",
			constant.FRAME_NAME,
			local_conf.GetGConf().Hertz.App,
			local_conf.GetGConf().Hertz.Server,
			local_conf.GetGConf().Hertz.Service[0].Name),
		RegistryAddr: local_conf.GetGConf().Registry.RegistryAddress,
		RegistryType: local_conf.GetGConf().Registry.Name,
		NamespaceId:  local_conf.GetGConf().Config.Namespace,
		Username:     local_conf.GetGConf().Registry.Username,
		Password:     local_conf.GetGConf().Registry.Password,
	}.Options()...)

	// 注册通用中间件
	registerMiddleware(h)

	// 业务参数初始化，默认远程配置
	r_conf.Biz_conf_init()

	// mysql初始化,远程/本地
	// mysql.RegisterConnByGroup("group", "data_id")
	// mysql.RegisterConnByLocal("group", "data_id")

	// reids初始化,远程/本地
	// redis.RegisterConnByGroup("group", "data_id")
	// redis.RegisterConnByLocal("group", "data_id")

	return h
}

func registerMiddleware(h *server.Hertz) {
	// pprof
	if local_conf.GetGConf().Hertz.EnablePprof {
		pprof.Register(h)
	}
	// gzip
	if local_conf.GetGConf().Hertz.EnableGzip {
		h.Use(gzip.Gzip(gzip.DefaultCompression))
	}

	// access log
	if local_conf.GetGConf().Hertz.EnableAccessLog {
		h.Use(accesslog.New())
	}

	// recovery
	h.Use(recovery.Recovery())

	// cores
	h.Use(cors.Default())

	// 自定义json序列化忽略omitempty
	render.ResetJSONMarshal(jsoniter.Marshal)
}
