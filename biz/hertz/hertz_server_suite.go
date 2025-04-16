package hertz

import (
	"os"
	"strconv"
	"strings"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"

	hertz_utils "github.com/cloudwego/hertz/pkg/common/utils"
	consulapi "github.com/hashicorp/consul/api"
	registryconsul "github.com/hertz-contrib/registry/consul"
	registrynacos "github.com/hertz-contrib/registry/nacos"
)

type HertzCommonServerSuite struct {
	CurrentServiceName string   `yaml:"current_service_name"`
	RegistryAddr       []string `yaml:"registry_addr"`
	RegistryType       string   `yaml:"registry_type"`
	NamespaceId        string   `yaml:"namespace_id"`
	Username           string   `yaml:"username"`
	Password           string   `yaml:"password"`
	Address            string   `yaml:"address"`

	EnablePrometheusTracer bool   `yaml:"enable_prometheus_tracer"`
	PrometheusTracerAddr   string `yaml:"prometheus_tracer_addr"`
	PrometheusTracerPath   string `yaml:"prometheus_tracer_path"`

	LogLevel string `yaml:"log_level"`
}

func (h HertzCommonServerSuite) Options() []config.Option {
	opts := []config.Option{
		server.WithHostPorts(h.Address),
	}

	var r registry.Registry
	switch h.RegistryType {
	case RegTypeConsul:
		r = initHertzConsulReg(&h)
	case RegTypeNacos:
		r = initHertzNacosReg(&h)
	case RegTypeETCD:
		panic("etcd unsupported registry type now")
	default:
		panic("unsupported registry type")
	}

	opts = append(opts, server.WithRegistry(r, &registry.Info{
		ServiceName: h.CurrentServiceName,
		Addr:        hertz_utils.NewNetAddr("tcp", h.Address),
		Weight:      10,
		Tags:        nil,
	}))

	if h.EnablePrometheusTracer {
		opts = append(opts, server.WithTracer(NewPrometheusTracer(
			h.PrometheusTracerAddr,
			h.PrometheusTracerPath,
		)))
	}

	return opts
}

// initHertzConsulReg use consul as registry center
func initHertzConsulReg(h *HertzCommonServerSuite) registry.Registry {
	consulConfig := &consulapi.Config{
		Address: h.RegistryAddr[0],
		Token:   h.Password,
	}
	consulClient, err := consulapi.NewClient(consulConfig)
	if err != nil {
		panic(err)
	}

	r := registryconsul.NewConsulRegister(consulClient, registryconsul.WithCheck(&consulapi.AgentServiceCheck{
		Interval:                       "7s",
		Timeout:                        "5s",
		DeregisterCriticalServiceAfter: "15s",
	}))
	return r
}

// initHertzNacosReg 使用nacos作为注册中心
func initHertzNacosReg(h *HertzCommonServerSuite) registry.Registry {
	sc := make([]constant.ServerConfig, 0)
	for _, addr := range h.RegistryAddr {
		ipAndPort := strings.Split(addr, ":")
		uintPort, err := strconv.Atoi(ipAndPort[1])
		if err != nil {
			panic(err)
		}
		sc = append(sc, constant.ServerConfig{
			IpAddr: ipAndPort[0],
			Port:   uint64(uintPort),
		})
	}
	if h.Username == "" {
		h.Username = os.Getenv("NACOS_USERNAME")
	}
	if h.Password == "" {
		h.Password = os.Getenv("NACOS_PASSWORD")
	}
	cc := constant.ClientConfig{
		NamespaceId:         h.NamespaceId,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            h.LogLevel,
		Username:            h.Username,
		Password:            h.Password,
	}

	cli, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err)
	}

	return registrynacos.NewNacosRegistry(cli)
}
