package hertz

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/tracer"
	"github.com/cloudwego/hertz/pkg/common/tracer/stats"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type serverTracer struct {
	serverHandledCounter   *prom.CounterVec
	serverHandledHistogram *prom.HistogramVec
}

func (s *serverTracer) Start(ctx context.Context, c *app.RequestContext) context.Context {
	return ctx
}

func (s *serverTracer) Finish(ctx context.Context, c *app.RequestContext) {
	if c.GetTraceInfo().Stats().Level() == stats.LevelDisabled {
		return
	}

	httpStart := c.GetTraceInfo().Stats().GetEvent(stats.HTTPStart)
	httpFinish := c.GetTraceInfo().Stats().GetEvent(stats.HTTPFinish)
	if httpFinish == nil || httpStart == nil {
		return
	}
	cost := httpFinish.Time().Sub(httpStart.Time())
	_ = counterAdd(s.serverHandledCounter, 1, genLabels(c))
	_ = histogramObserve(s.serverHandledHistogram, cost, genLabels(c))
}

func counterAdd(counterVec *prom.CounterVec, value int, labels prom.Labels) error {
	counter, err := counterVec.GetMetricWith(labels)
	if err != nil {
		return err
	}
	counter.Add(float64(value))
	return nil
}

func histogramObserve(histogramVec *prom.HistogramVec, value time.Duration, labels prom.Labels) error {
	histogram, err := histogramVec.GetMetricWith(labels)
	if err != nil {
		return err
	}
	histogram.Observe(float64(value.Microseconds()))
	return nil
}

func genLabels(ctx *app.RequestContext) prom.Labels {
	labels := make(prom.Labels)
	labels[labelMethod] = defaultValIfEmpty(string(ctx.Request.Method()), unknownLabelValue)
	labels[labelStatusCode] = defaultValIfEmpty(strconv.Itoa(ctx.Response.Header.StatusCode()), unknownLabelValue)
	labels[labelPath] = defaultValIfEmpty(ctx.FullPath(), unknownLabelValue)
	labels[labelBizCode] = defaultValIfEmpty(ctx.Response.Header.Get(labelBizCode), succBizStatus)

	return labels
}

func defaultValIfEmpty(val, def string) string {
	if val == "" {
		return def
	}
	return val
}

type prometheusTracerConfig struct {
	buckets            []float64
	enableGoCollector  bool
	registry           *prom.Registry
	runtimeMetricRules []collectors.GoRuntimeMetricsRule
	disableServer      bool
}

func defaultConfig() *prometheusTracerConfig {
	return &prometheusTracerConfig{
		buckets:           defaultBuckets,
		enableGoCollector: false,
		registry:          prom.NewRegistry(),
		disableServer:     false,
	}
}

func NewPrometheusTracer(addr, path string, options ...PrometheusTracerConfigOption) tracer.Tracer {
	cfg := defaultConfig()

	if len(options) > 0 {
		for _, option := range options {
			option.apply(cfg)
		}
	}

	if !cfg.disableServer {
		http.Handle(path, promhttp.HandlerFor(cfg.registry, promhttp.HandlerOpts{
			ErrorHandling: promhttp.ContinueOnError,
		}))
		go func() {
			if err := http.ListenAndServe(addr, nil); err != nil {
				hlog.Fatal("HERTZ: Unable to start a promhttp server, err: " + err.Error())
			}
		}()
	}

	serverHandledCounter := prom.NewCounterVec(
		prom.CounterOpts{
			Namespace: "hertz_server_throughput",
			Help:      "Total number of HTTPs completed by the server, regardless of success or failure.",
		},
		[]string{
			labelMethod,
			labelStatusCode,
			labelPath,
			labelBizCode,
		},
	)
	cfg.registry.MustRegister(serverHandledCounter)

	serverHandledHistogram := prom.NewHistogramVec(
		prom.HistogramOpts{
			Namespace: "hertz_server_latency_us",
			Help:      "Latency (microseconds) of HTTP that had been application-level handled by the server.",
			Buckets:   cfg.buckets,
		},
		[]string{
			labelMethod,
			labelStatusCode,
			labelPath,
			labelBizCode,
		},
	)
	cfg.registry.MustRegister(serverHandledHistogram)

	if cfg.enableGoCollector {
		cfg.registry.MustRegister(collectors.NewGoCollector(
			collectors.WithGoCollectorRuntimeMetrics(cfg.runtimeMetricRules...),
		))
	}

	return &serverTracer{
		serverHandledCounter:   serverHandledCounter,
		serverHandledHistogram: serverHandledHistogram,
	}
}

type PrometheusTracerConfigOption interface {
	apply(cfg *prometheusTracerConfig)
}

type prometheusTracerConfigOption func(*prometheusTracerConfig)

func (fn prometheusTracerConfigOption) apply(cfg *prometheusTracerConfig) {
	fn(cfg)
}

func WithEnableGoCollector(enable bool) PrometheusTracerConfigOption {
	return prometheusTracerConfigOption(func(cfg *prometheusTracerConfig) {
		cfg.enableGoCollector = enable
	})
}

func WithGoCollectorRule(rules ...collectors.GoRuntimeMetricsRule) PrometheusTracerConfigOption {
	return prometheusTracerConfigOption(func(cfg *prometheusTracerConfig) {
		cfg.runtimeMetricRules = rules
	})
}

func WithDisableServer(disable bool) PrometheusTracerConfigOption {
	return prometheusTracerConfigOption(func(cfg *prometheusTracerConfig) {
		cfg.disableServer = disable
	})
}

func WithHistogramBuckets(buckets []float64) PrometheusTracerConfigOption {
	return prometheusTracerConfigOption(func(cfg *prometheusTracerConfig) {
		cfg.buckets = buckets
	})
}

func WithRegistry(reg *prom.Registry) PrometheusTracerConfigOption {
	return prometheusTracerConfigOption(func(cfg *prometheusTracerConfig) {
		cfg.registry = reg
	})
}
