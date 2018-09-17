package buckets

import (
	"github.com/gin-gonic/gin"
)

// AddRoutes adds ARK buckets related API routes
func AddRoutes(group *gin.RouterGroup) {

	group.GET("", List)
	group.POST("", Create)
	item := group.Group("/:name")
	{
		item.GET("", Get)
		item.PUT("", List)
		item.DELETE("", Delete)
	}
}
