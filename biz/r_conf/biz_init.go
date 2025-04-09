package r_conf

import (
	"ssg/kitex-common/common/remote_conf"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gopkg.in/yaml.v3"
)

// 初始化业务参数配置
var conf_group = "Dev"
var conf_data_id = "biztest"
var BizCongfig *BizConfig

type BizConfig struct {
	//test_biz_config
	TestBizConfig string `yaml:"test_biz_config"`
}

func Biz_conf_init() {
	tmpConfig := &BizConfig{}
	//需要自行指定group  和  dataid
	remote_conf.GetConfig(conf_group, conf_data_id, tmpConfig)
	hlog.Infof("remote config : [%+v]", tmpConfig.TestBizConfig)

	BizCongfig = tmpConfig
	go watch(conf_group, conf_data_id)
}

func watch(group string, key string) {
	c, err := remote_conf.Watch(group, key)
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

		hlog.Infof("remote config Change tmp: \n %v", tmpConfig.TestBizConfig)
		BizCongfig = tmpConfig
	}
	hlog.Infof("watch config end")
}
