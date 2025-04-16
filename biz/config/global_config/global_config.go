package global_config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Global struct {
	Namespace     string `yaml:"namespace"`
	EnvName       string `yaml:"env_name"`
	LocalIP       string `yaml:"local_ip"`
	ContainerName string `yaml:"container_name"`
}

type Hertz struct {
	App             string  `yaml:"app"`
	Server          string  `yaml:"server"`
	BinPath         string  `yaml:"bin_path"`
	ConfPath        string  `yaml:"conf_path"`
	DataPath        string  `yaml:"data_path"`
	EnablePprof     bool    `yaml:"enable_pprof"`
	EnableGzip      bool    `yaml:"enable_gzip"`
	EnableAccessLog bool    `yaml:"enable_access_log"`
	Service         Service `yaml:"service"`
}

type Service struct {
	Name    string `yaml:"name"`
	Address string `yaml:"addr"`
}

type Log struct {
	LogMode       string   `yaml:"log_mode"`
	LogLevel      string   `yaml:"log_level"`
	LogFileName   string   `yaml:"log_file_name"`
	LogMaxSize    int      `yaml:"log_max_size"`
	LogMaxBackups int      `yaml:"log_max_backups"`
	LogMaxAge     int      `yaml:"log_max_age"`
	LogCompress   bool     `yaml:"log_compress"`
	ExtKeys       []string `yaml:"ext_keys"`
}

type Prometheus struct {
	Enable bool   `yaml:"enable"`
	Addr   string `yaml:"addr"`
	Path   string `yaml:"path"`
}

type Registry struct {
	Name            string   `yaml:"name"`
	RegistryAddress []string `yaml:"registry_address"`
	Namespace       string   `yaml:"namespace"`
	Username        string   `yaml:"username"`
	Password        string   `yaml:"password"`
}

type Selector struct {
	Name       string   `yaml:"name"`
	ServerAddr []string `yaml:"server_addr"`
	Namespace  string   `yaml:"namespace"`
	Username   string   `yaml:"username"`
	Password   string   `yaml:"password"`
}

type RemoteConfig struct {
	Name       string   `yaml:"name"`
	ServerAddr []string `yaml:"server_addr"`
	Namespace  string   `yaml:"namespace"`
	Username   string   `yaml:"username"`
	Password   string   `yaml:"password"`
}

type ModeType string

const (
	ModeTypeLocal ModeType = "local"
	ModeTypeNacos ModeType = "nacos"
)

// Config is for biz config
type Config struct {
	Mode ModeType `yaml:"mode"`

	// used when mode is local, path is the config file path
	LocalConfigPath string `yaml:"local_config_path"`

	// used when mode is remote, below is the nacos server config
	RemoteConfig *RemoteConfig `yaml:"remote_config"`
	ConfigGroup  string        `yaml:"config_group"`
	ConfigDataID string        `yaml:"config_data_id"`
}

type GlobalConfig struct {
	Env    string
	Global *Global `yaml:"global"`
	Hertz  *Hertz  `yaml:"hertz"`
	Log    *Log    `yaml:"log"`

	Config *Config `yaml:"config"`

	Mode       ModeType    `yaml:"mode"`
	Prometheus *Prometheus `yaml:"prometheus"`
	Registry   *Registry   `yaml:"registry"`
	Selector   *Selector   `yaml:"selector"`
}

var (
	GConfig *GlobalConfig
)

func InitGlobalConfig(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var config GlobalConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	GConfig = &config
}
