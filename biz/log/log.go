package log

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/server"
	hertzzap "github.com/hertz-contrib/logger/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LogModeConsole = "console"
	LogModeFile    = "file"
)

func InitHLog(ioWriter io.Writer, level hlog.Level, mode string, extKeys ...string) {
	switch mode {
	case LogModeConsole, "":
		setHConsoleMode(level, extKeys...)
	case LogModeFile:
		setHLogFile(ioWriter, level, extKeys...)
	default:
		hlog.Fatalf("unsupport log mode:[%+v]", mode)
	}
}

func setHConsoleMode(level hlog.Level, extKeys ...string) {
	var opts []hertzzap.Option
	hzExterKeys := make([]hertzzap.ExtraKey, 0)
	for _, val := range extKeys {
		hzExterKeys = append(hzExterKeys, hertzzap.ExtraKey(val))
	}
	opts = append(opts,
		hertzzap.WithCoreEnc(
			zapcore.NewJSONEncoder(
				zap.NewDevelopmentEncoderConfig())),
		hertzzap.WithZapOptions(
			zap.AddCaller(),
			zap.AddCallerSkip(3)),
		hertzzap.WithExtraKeys(hzExterKeys),
		hertzzap.WithExtraKeyAsStr())

	logger := hertzzap.NewLogger(opts...)
	fileWriter := io.MultiWriter(os.Stdout)
	logger.SetOutput(fileWriter)
	logger.SetLevel(level)
	hlog.SetLogger(logger)
}

func setHLogFile(ioWriter io.Writer, level hlog.Level, extKeys ...string) {
	var opts []hertzzap.Option
	var output zapcore.WriteSyncer
	// non-prod environment will use sync mode for file output
	if !strings.Contains(strings.ToLower(os.Getenv("GO_ENV")), "prod") {
		opts = append(opts, hertzzap.WithCoreEnc(zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())))
		output = zapcore.AddSync(ioWriter)
	} else {
		opts = append(opts, hertzzap.WithCoreEnc(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())))
		// async log
		output = &zapcore.BufferedWriteSyncer{
			WS:            zapcore.AddSync(ioWriter),
			FlushInterval: time.Minute,
		}
	}
	server.RegisterShutdownHook(func() {
		output.Sync() //nolint:errcheck
	})

	hzExterKeys := make([]hertzzap.ExtraKey, 0)
	for _, val := range extKeys {
		hzExterKeys = append(hzExterKeys, hertzzap.ExtraKey(val))
	}

	opts = append(opts, hertzzap.WithZapOptions(
		zap.AddCaller(),
		zap.AddCallerSkip(3)),
		hertzzap.WithExtraKeys(hzExterKeys),
		hertzzap.WithExtraKeyAsStr())
	logger := hertzzap.NewLogger(opts...)
	logger.SetOutput(output)
	logger.SetLevel(level)
	hlog.SetLogger(logger)
}

func HLogLevel(level string) hlog.Level {
	switch level {
	case "trace":
		return hlog.LevelTrace
	case "debug":
		return hlog.LevelDebug
	case "info":
		return hlog.LevelInfo
	case "notice":
		return hlog.LevelNotice
	case "warn":
		return hlog.LevelWarn
	case "error":
		return hlog.LevelError
	case "fatal":
		return hlog.LevelFatal
	default:
		return hlog.LevelInfo
	}
}
