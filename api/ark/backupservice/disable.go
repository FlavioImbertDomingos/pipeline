package backupservice

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goph/emperror"

	"github.com/banzaicloud/pipeline/api/ark/common"
	"github.com/banzaicloud/pipeline/internal/ark/api"
	"github.com/banzaicloud/pipeline/internal/platform/gin/correlationid"
)

// Disable removes ARK deployment from the cluster
func Disable(c *gin.Context) {
	logger := correlationid.Logger(common.Log, c)
	logger.Info("removing backup service from cluster")

	err := common.GetARKService(c.Request).GetDeploymentsService().Remove()
	if err != nil {
		err = emperror.Wrap(err, "error removing backup service")
		logger.Error(err)
		common.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, api.DisableBackupServiceResponse{
		Status: http.StatusOK,
	})
}
