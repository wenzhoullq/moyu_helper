package source

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"weixin_LLM/dto/request"
	"weixin_LLM/lib"
	"weixin_LLM/service/source"
)

func UpdateSource(c *gin.Context) {
	var req request.UpdateSourceRequest
	if err := lib.RequestUnmarshal(c, &req); err != nil {
		lib.SetContextErrorResponse(c, err)
		return
	}
	ss := source.NewSourceService()
	resp, err := ss.UpdateSource(c, &req)
	if err != nil {
		//打印日志
		ss.Logln(logrus.ErrorLevel, err.Error())
	}
	lib.SetContextResponse(c, resp)
}
