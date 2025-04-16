package global_init

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"TSMall/biz/config/biz_config"
	"TSMall/biz/config/global_config"
	"TSMall/biz/constant"
	"TSMall/biz/hertz"
	"TSMall/biz/log"
	"TSMall/biz/middelware"

	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/cloudwego/hertz/pkg/app/server/render"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/cors"
	"github.com/hertz-contrib/gzip"
	"github.com/hertz-contrib/logger/accesslog"
	"github.com/hertz-contrib/pprof"
	jsoniter "github.com/json-iterator/go"

	"github.com/cloudwego/hertz/pkg/app/server"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitServer(service global_config.Service) (h *server.Hertz) {
	// init hlog
	log.InitHLog(&lumberjack.Logger{
		Filename:   global_config.GConfig.Log.LogFileName,
		MaxSize:    global_config.GConfig.Log.LogMaxSize,
		MaxBackups: global_config.GConfig.Log.LogMaxBackups,
		MaxAge:     global_config.GConfig.Log.LogMaxAge,
		Compress:   global_config.GConfig.Log.LogCompress,
	}, log.HLogLevel(global_config.GConfig.Log.LogLevel), global_config.GConfig.Log.LogMode)

	if global_config.GConfig.Mode == global_config.ModeTypeNacos {
		// register to nacos
		h = server.New(hertz.HertzCommonServerSuite{
			Address: global_config.GConfig.Hertz.Service.Address,
			CurrentServiceName: fmt.Sprintf("%s.%s.%s.%s",
				constant.FrameName,
				global_config.GConfig.Hertz.App,
				global_config.GConfig.Hertz.Server,
				global_config.GConfig.Hertz.Service.Name),
			RegistryAddr: global_config.GConfig.Registry.RegistryAddress,
			RegistryType: global_config.GConfig.Registry.Name,
			NamespaceId:  global_config.GConfig.Registry.Namespace,
			Username:     global_config.GConfig.Registry.Username,
			Password:     global_config.GConfig.Registry.Password,
		}.Options()...)
	} else {
		h = server.Default(server.WithHostPorts(global_config.GConfig.Hertz.Service.Address))
	}

	// register middleware
	registerMiddleware(h)

	// init biz config
	biz_config.BizConfigInit(global_config.GConfig.Config)

	h.SetCustomSignalWaiter(func(err chan error) error {
		signalToNotify := []os.Signal{syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM}
		if signal.Ignored(syscall.SIGHUP) {
			signalToNotify = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
		}
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, signalToNotify...)
		select {
		case sig := <-signals:
			switch sig {
			// case syscall.SIGTERM:
			//     // force exit
			//     return errors.NewPublic(sig.String()) // nolint
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM:
				hlog.SystemLogger().Infof("Received signal: %s\n", sig)
				// graceful shutdown
				return nil
			}
		case err := <-err:
			// error occurs, exit immediately
			return err
		}
		return nil
	})

	// init redis
	// tsredis.InitRedis(nil)

	// init mongodb
	// tsmongodb.InitMongoDB(nil)

	middelware.InitMiddeleware(h)

	return h
}

func registerMiddleware(h *server.Hertz) {
	// pprof
	if global_config.GConfig.Hertz.EnablePprof {
		pprof.Register(h)
	}
	// gzip
	if global_config.GConfig.Hertz.EnableGzip {
		h.Use(gzip.Gzip(gzip.DefaultCompression))
	}

	// access log
	if global_config.GConfig.Hertz.EnableAccessLog {
		h.Use(accesslog.New())
	}

	// recovery
	h.Use(recovery.Recovery())

	// cores
	h.Use(cors.Default())

	// customize json marshal ingore omitempty
	render.ResetJSONMarshal(jsoniter.Marshal)
}
