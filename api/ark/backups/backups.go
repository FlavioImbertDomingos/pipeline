package backups

import (
	"github.com/gin-gonic/gin"

	"github.com/banzaicloud/pipeline/api/ark/common"
	"github.com/banzaicloud/pipeline/config"
)

// AddOrgRoutes adds routes for managing ARK backups within an organization
func AddOrgRoutes(group *gin.RouterGroup) {
	group.GET("", ListAll)
}

// AddRoutes adds ARK backups related API routes
func AddRoutes(group *gin.RouterGroup) {

	group.Use(common.ARKMiddleware(config.DB(), common.Log))
	group.GET("", List)
	group.POST("", Create)
	item := group.Group("/:name")
	{
		item.GET("", Get)
		item.DELETE("", Delete)
		item.GET("/download", Download)
		item.GET("/logs", GetLogs)
	}
}
