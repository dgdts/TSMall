package biz_config

import (
	"TSMall/biz/utils"
	"os"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gopkg.in/yaml.v3"

	"TSMall/biz/config/global_config"
)

var bizConfig *BizConfig

type BizConfig struct {
	//test_biz_config
	TestBizConfig string `yaml:"test_biz_config"`
}

func BizConfigInit(config *global_config.Config) {
	switch config.Mode {
	case global_config.ModeTypeNacos:
		bizConfigNacosInit(config)
	default:
		bizConfigLocalInit(config)
	}
}

func bizConfigLocalInit(config *global_config.Config) {
	configBytes, err := os.ReadFile(config.LocalConfigPath)
	if err != nil {
		panic(err)
	}

	var tmpConfig BizConfig
	err = yaml.Unmarshal(configBytes, &tmpConfig)
	if err != nil {
		panic(err)
	}
	setBizConfig(&tmpConfig)
	hlog.Infof("local config: %v", GetBizConfig())
}

func bizConfigNacosInit(config *global_config.Config) {
	utils.InitRemoteConfig(*config.RemoteConfig)

	var tmpConfig BizConfig
	err := utils.GetRemoteConfig(config.ConfigDataID, config.ConfigGroup, &tmpConfig)
	if err != nil {
		panic(err)
	}

	setBizConfig(&tmpConfig)
	hlog.Infof("remote config: %v", GetBizConfig())
	// could add watch for hot update, but we should consider the race condition when update config
	// go watch(config.ConfigGroup, config.ConfigDataID)
}

func setBizConfig(config *BizConfig) {
	bizConfig = config
}

func GetBizConfig() *BizConfig {
	return bizConfig
}

func watch(group string, key string) {
	c, err := utils.WatchRemoteConfig(group, key)
	if err != nil {
		hlog.Errorf("watch failed. error: %s", err.Error())
		panic(err)
	}

	for resp := range c {
		tmpConfig := &BizConfig{}
		err := yaml.Unmarshal([]byte(resp), &tmpConfig)
		if err != nil {
			panic(err)
		}
		hlog.Infof("remote config Change tmp: %v", tmpConfig)
		setBizConfig(tmpConfig)
	}
	hlog.Infof("watch config end")
}
