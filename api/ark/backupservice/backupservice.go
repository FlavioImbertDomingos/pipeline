package backupservice

import (
	"github.com/gin-gonic/gin"

	"github.com/banzaicloud/pipeline/api/ark/common"
	"github.com/banzaicloud/pipeline/config"
)

// AddRoutes adds ARK backups related API routes
func AddRoutes(group *gin.RouterGroup) {

	group.Use(common.ARKMiddleware(config.DB(), common.Log))
	group.GET("/status", Status)
	group.POST("/enable", Enable)
	group.POST("/disable", Disable)
}
