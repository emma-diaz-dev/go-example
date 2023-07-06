package cfg

import (
	"github.com/joeshaw/envdecode"
)

var (
	config *Config
)

type (
	Config struct {
		ClickPath string `env:"CLICK_PATH,default=/tmp/clicks.pipe"`
		ImpPath   string `env:"IMP_PATH,default=/tmp/impressions.pipe"`
	}
)

func initCfg() {
	if config != nil {
		return
	}
	config = &Config{}
	if err := envdecode.Decode(config); err != nil {
		panic(err)
	}
}

func GetConfig() *Config {
	initCfg()
	return config
}
