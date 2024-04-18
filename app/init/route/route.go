package route

import (
	"github.com/gin-gonic/gin"
	source2 "weixin_LLM/api/source"
)

func RouteInit() *gin.Engine {
	r := gin.New()
	r.Use()
	imaotai := r.Group("moyuHelper")
	{
		logIn := imaotai.Group("source")
		{
			logIn.POST("/updateSource", source2.UpdateSource)
			logIn.GET("/getSource", source2.GetSource)
			logIn.POST("/createSource", source2.CreateSource)
		}
	}
	return r
}
