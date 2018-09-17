package backupservice

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/banzaicloud/pipeline/api/ark/common"
	"github.com/banzaicloud/pipeline/internal/ark/api"
	"github.com/banzaicloud/pipeline/internal/platform/gin/correlationid"
)

// Status gets an ARK backup deployment status by trying to create ARK client
func Status(c *gin.Context) {
	logger := correlationid.Logger(common.Log, c)
	logger.Info("checking ARK deployment status")

	schedulesSvc := common.GetARKService(c.Request).GetSchedulesService()
	_, err := schedulesSvc.List()
	if err != nil {
		err = errors.New("backup service not deployed")
		logger.Error(err.Error())
		common.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, api.BackupServiceStatusResponse{
		Status: http.StatusOK,
	})
}
