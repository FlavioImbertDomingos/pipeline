package backups

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goph/emperror"

	"github.com/banzaicloud/pipeline/api/ark/common"
	arkAPI "github.com/banzaicloud/pipeline/internal/ark/api"
	"github.com/banzaicloud/pipeline/internal/platform/gin/correlationid"
)

// Create creates an ARK backup
func Create(c *gin.Context) {
	logger := correlationid.Logger(common.Log, c)
	logger.Info("creating backup")

	var req arkAPI.CreateBackupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		err = emperror.Wrap(err, "error parsing request")
		logger.Error(err)
		common.ErrorResponse(c, err)
		return
	}

	err := common.GetARKService(c.Request).GetClusterBackupsService().Create(req)
	if err != nil {
		logger.Error(err)
		common.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, &arkAPI.CreateBackupResponse{
		Name:   req.Name,
		Status: http.StatusOK,
	})
}
