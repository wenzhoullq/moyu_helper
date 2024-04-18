package lib

import (
	"github.com/gin-gonic/gin"
)

func RequestUnmarshal(c *gin.Context, req interface{}) (err error) {
	if err = c.ShouldBind(&req); err != nil {
		return err
	}
	return nil
}
