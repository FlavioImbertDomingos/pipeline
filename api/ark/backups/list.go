package backups

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goph/emperror"

	"github.com/banzaicloud/pipeline/api/ark/common"
	"github.com/banzaicloud/pipeline/internal/platform/gin/correlationid"
)

// List lists ARK backups
func List(c *gin.Context) {
	logger := correlationid.Logger(common.Log, c)
	logger.Info("getting backups")

	backups, err := common.GetARKService(c.Request).GetClusterBackupsService().List()
	if err != nil {
		err = emperror.Wrap(err, "error getting backups")
		logger.Error(err)
		common.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, backups)
}
