package config

import (
	"errors"

	"github.com/spf13/viper"
)

const config_path = "config"

var ErrParsingConfigFileFail = errors.New("設定ファイルの読み込み中にエラーが発生しました")

type Config struct {
	config *viper.Viper
}

func (c *Config) GetInt(config string) int {
	return c.config.GetInt(config)
}

func (c *Config) GetString(config string) string {
	return c.config.GetString(config)
}

func (c *Config) GetStringSlice(config string) []string {
	return c.config.GetStringSlice(config)
}

func (c *Config) Get(config string) any {
	return c.config.Get(config)
}

func NewConfig(env string) (*Config, error) {
	conf := viper.New()
	conf.SetConfigType("toml")

	conf.SetConfigName("general")
	conf.AddConfigPath(config_path)
	err := conf.ReadInConfig()
	if err != nil {
		return nil, ErrParsingConfigFileFail
	}

	conf.SetConfigName(env)
	conf.AddConfigPath(config_path)
	err = conf.MergeInConfig()

	if err != nil {
		return nil, ErrParsingConfigFileFail
	}

	return &Config{
		config: conf,
	}, nil
}
