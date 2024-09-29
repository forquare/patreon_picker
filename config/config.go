package config

import (
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"runtime"
	"sync"
)

var (
	once   sync.Once
	config Config
)

type Credentials struct {
	Id          string `mapstructure:"id"`
	Secret      string `mapstructure:"secret"`
	RedirectURL string `mapstructure:"redirect_url"`
}

type Connection struct {
	Address string `mapstructure:"address"`
	Port    string `mapstructure:"port"`
}

type Session struct {
	CookieKey string `mapstructure:"cookie"`
	Name      string `mapstructure:"name"`
}

type Config struct {
	Credentials Credentials `mapstructure:"credentials"`
	LogLevel    string      `mapstructure:"log_level"`
	Session     Session     `mapstructure:"session"`
	Connection  Connection  `mapstructure:"connection"`
}

func GetConfig() Config {
	once.Do(func() {
		logger.Debug("Reading config file")
		viper.SetConfigName("picker.yaml")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/config")
		viper.AddConfigPath("/etc/")
		viper.AddConfigPath("/usr/local/etc/")
		err := viper.ReadInConfig()
		if err != nil {
			logger.Fatal("Cannot read config file. File may not exist, or be in the wrong format.")
		}
		err = viper.Unmarshal(&config)
		if err != nil {
			logger.Fatal("Cannot read config file. File may be in the wrong format.")
		}
	})

	if logger.GetLevel() == logger.TraceLevel {
		pc, _, _, ok := runtime.Caller(1)
		details := runtime.FuncForPC(pc)
		if ok && details != nil {
			logger.Tracef("Returning config to %s", details.Name())
		}
	}

	return config
}
