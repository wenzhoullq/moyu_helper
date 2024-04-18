package source

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"weixin_LLM/lib"
	"weixin_LLM/service/source"
)

func GetSource(c *gin.Context) {
	ss := source.NewSourceService()
	resp, err := ss.GetSource(c)
	if err != nil {
		//打印日志
		ss.Logln(logrus.ErrorLevel, err.Error())
	}
	lib.SetContextResponse(c, resp)
}
