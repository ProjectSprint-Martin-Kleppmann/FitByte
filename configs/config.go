package configs

import (
	"fmt"

	"github.com/spf13/viper"

	"FitByte/pkg/log"
)

type option struct {
	ConfigFolder []string
	ConfigType   string
	ConfigFile   string
}

type Option func(*option)

func LoadConfig(options ...Option) Config {
	opt := &option{
		ConfigFolder: getDefaultConfigFolder(),
		ConfigType:   getDefaultConfigType(),
		ConfigFile:   getDefaultConfigFile(),
	}

	for _, optFunc := range options {
		optFunc(opt)
	}

	for _, folder := range opt.ConfigFolder {
		viper.AddConfigPath(folder)
	}

	viper.SetConfigType(opt.ConfigType)
	viper.SetConfigType(opt.ConfigType)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	log.Logger.Info().Msg("Configuration loaded successfully")
	log.Logger.Info().Interface("config", cfg).Msg("Loaded configuration details")
	return cfg
}

func WithConfigFolder(folder []string) Option {
	return func(opt *option) {
		opt.ConfigFolder = folder
	}
}

func WithConfigType(configType string) Option {
	return func(opt *option) {
		opt.ConfigType = configType
	}
}

func WithConfigFile(configFile string) Option {
	return func(opt *option) {
		opt.ConfigFile = configFile
	}
}

func getDefaultConfigFolder() []string {
	return []string{"./files/config"}
}

func getDefaultConfigFile() string {
	return "config"
}

func getDefaultConfigType() string {
	return "yaml"
}
