package router

import (
	"github.com/gin-gonic/gin"
	"zbx-api/api"
	"zbx-api/logger"
)

func Setup() *gin.Engine {

	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	v1 := r.Group("/v1")
	{

		//v1.POST("/hostFilter/", api.HistoryFiterApi)
		v1.GET("/getzabbix/:id/",api.GetZabbixKey)

	}

	return r
}