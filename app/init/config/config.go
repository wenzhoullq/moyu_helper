package config

import (
	"github.com/BurntSushi/toml"
	"weixin_LLM/dto/config"
)

var Config *config.Configuration

func ConfigInit(filename string) error {
	if _, err := toml.DecodeFile(filename, &Config); err != nil {
		return err
	}
	return nil
}
