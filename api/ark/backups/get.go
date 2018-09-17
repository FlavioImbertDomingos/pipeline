package backups

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/banzaicloud/pipeline/api/ark/common"
	"github.com/banzaicloud/pipeline/internal/platform/gin/correlationid"
)

// Get gets an ARK backup
func Get(c *gin.Context) {
	backupName := c.Param("name")

	logger := correlationid.Logger(common.Log, c).WithField("backup", backupName)
	logger.Info("getting backup")

	backup, err := common.GetARKService(c.Request).GetClusterBackupsService().GetByName(backupName)
	if err != nil {
		err = errors.Wrap(err, "error getting backup")
		logger.Error(err.Error())
		common.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, backup)
}
