package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var c Config

type Config struct {
	Server       Server       `yaml:"server"`
	WxConfig     WxConfig     `yaml:"wxConfig"`
	OpenAIConfig OpenAIConfig `yaml:"openAIConfig"`
}

type Server struct {
	Host string `yaml:"host"`
	Post int    `yaml:"post"`
}

type WxConfig struct {
	Token          string `yaml:"token"`
	EncodingAeskey string `yaml:"encoding_aeskey"`
	CorpId         string `yaml:"corp_id"`
	CorpSecret     string `yaml:"corp_secret"`
}

type OpenAIConfig struct {
	Model string `yaml:"model"`
	Token string `yaml:"token"`
}

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("/etc/appname/")
	viper.AddConfigPath("$HOME/.appname")
	viper.AddConfigPath("./")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	viper.Unmarshal(&c)
}

func NewConfig() Config {
	return c
}
