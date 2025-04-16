package utils

import (
	"TSMall/biz/config/global_config"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/stretchr/testify/assert/yaml"
)

const (
	NacosUsernameEnvKey = "NACOS_USERNAME"
	NacosPasswordEnvKey = "NACOS_PASSWORD"

	NacosClientLogDir   = "/tmp/log"
	NacosClientCacheDir = "/tmp/cache"

	NacosClientTimeoutMs = 5000
)

type RemoteConfigManager struct {
	client       config_client.IConfigClient
	serverConfig global_config.RemoteConfig
}

var gConfMgr *RemoteConfigManager
var ones sync.Once

func InitRemoteConfig(config global_config.RemoteConfig) {
	ones.Do(func() {
		gConfMgr = &RemoteConfigManager{
			serverConfig: config,
		}

		nacosServerConfig := make([]constant.ServerConfig, 0)
		for _, addr := range gConfMgr.serverConfig.ServerAddr {
			ipAndPort := strings.Split(addr, ":")
			uintPort, err := strconv.ParseUint(ipAndPort[1], 10, 64)
			if err != nil {
				panic(err)
			}
			nacosServerConfig = append(nacosServerConfig, constant.ServerConfig{
				IpAddr: ipAndPort[0],
				Port:   uintPort,
			})
		}

		userName := gConfMgr.serverConfig.Username
		if userName == "" {
			userName = os.Getenv(NacosUsernameEnvKey)
		}
		password := gConfMgr.serverConfig.Password
		if password == "" {
			password = os.Getenv(NacosPasswordEnvKey)
		}

		clientConfig := constant.ClientConfig{
			NamespaceId:         gConfMgr.serverConfig.Namespace,
			TimeoutMs:           NacosClientTimeoutMs,
			NotLoadCacheAtStart: true,
			LogDir:              NacosClientLogDir,
			CacheDir:            NacosClientCacheDir,
			Username:            userName,
			Password:            password,
		}

		client, err := clients.CreateConfigClient(map[string]any{
			constant.KEY_SERVER_CONFIGS: nacosServerConfig,
			constant.KEY_CLIENT_CONFIG:  clientConfig,
		})

		if err != nil {
			panic(err)
		}
		gConfMgr.client = client
	})
}

func GetRemoteConfig(dataId, group string, out any) error {
	content, err := gConfMgr.client.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})

	if err != nil {
		return err
	}

	err = yaml.Unmarshal([]byte(content), out)
	if err != nil {
		return err
	}

	return nil
}

func WatchRemoteConfig(dataId, group string) (chan string, error) {
	configChan := make(chan string)
	err := gConfMgr.client.ListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			configChan <- data
		},
	})
	return configChan, err
}
