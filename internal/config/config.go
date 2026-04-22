package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Bilibili BilibiliConfig `mapstructure:"bilibili"`
	Email   EmailConfig     `mapstructure:"email"`
	Schedule string         `mapstructure:"schedule"`
}

type BilibiliConfig struct {
	SESSDATA string `mapstructure:"sessdata"`
	BiliJct  string `mapstructure:"bili_jct"`
	Buvid3   string `mapstructure:"buvid3"`
}

type EmailConfig struct {
	SMTPHost  string `mapstructure:"smtp_host"`
	SMTPPort  int    `mapstructure:"smtp_port"`
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
	From      string `mapstructure:"from"`
	To        string `mapstructure:"to"`
}

var cfg *Config

func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg = &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func Get() *Config {
	return cfg
}
