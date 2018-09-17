package restores

import (
	"github.com/banzaicloud/pipeline/api/ark/common"
	"github.com/banzaicloud/pipeline/config"
	"github.com/gin-gonic/gin"
)

// AddRoutes adds ARK restores related API routes
func AddRoutes(group *gin.RouterGroup) {

	group.Use(common.ARKMiddleware(config.DB(), common.Log))
	group.GET("", List)
	group.POST("", Create)
	item := group.Group("/:name")
	{
		item.GET("", Get)
		item.DELETE("", Delete)
		item.GET("/logs", GetLogs)
		item.GET("/results", GetResults)
	}
}
