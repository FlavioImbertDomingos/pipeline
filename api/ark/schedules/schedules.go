package schedules

import (
	"github.com/gin-gonic/gin"

	common "github.com/banzaicloud/pipeline/api/ark/common"
	"github.com/banzaicloud/pipeline/config"
)

// AddRoutes adds ARK schedules related API routes
func AddRoutes(group *gin.RouterGroup) {

	group.Use(common.ARKMiddleware(config.DB(), common.Log))
	group.GET("", List)
	group.POST("", Create)
	item := group.Group("/:name")
	{
		item.GET("", Get)
		item.DELETE("", Delete)
	}
}
