package config

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"weixin_LLM/lib"
	"weixin_LLM/service/config"
)

func LoadConfig(c *gin.Context) {
	cs := config.NewConfigService()
	resp, err := cs.LoadConfig(c)
	if err != nil {
		//打印日志
		cs.Logln(logrus.ErrorLevel, err.Error())
	}
	lib.SetContextResponse(c, resp)
}
