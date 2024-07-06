package route

import (
	"github.com/gin-gonic/gin"
	config2 "weixin_LLM/api/config"
	source2 "weixin_LLM/api/source"
)

func RouteInit() *gin.Engine {
	r := gin.New()
	r.Use()
	imaotai := r.Group("moyuHelper")
	{
		source := imaotai.Group("source")
		{
			source.POST("/updateSource", source2.UpdateSource)
			source.GET("/getSource", source2.GetSource)
			source.POST("/createSource", source2.CreateSource)
		}
		config := imaotai.Group("config")
		{
			config.GET("/loadConfig", config2.LoadConfig)
		}

	}
	return r
}
