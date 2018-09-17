package backups

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/banzaicloud/pipeline/api/ark/common"
	"github.com/banzaicloud/pipeline/internal/ark/api"
	"github.com/banzaicloud/pipeline/internal/platform/gin/correlationid"
)

// Delete deletes an ARK backup
func Delete(c *gin.Context) {
	backupName := c.Param("name")

	logger := correlationid.Logger(common.Log, c).WithField("backup", backupName)
	logger.Info("deleting backup")

	err := common.GetARKService(c.Request).GetClusterBackupsService().DeleteByName(backupName)
	if err != nil {
		err = errors.Wrap(err, "error deleting backup")
		logger.Error(err.Error())
		common.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, &api.DeleteBackupResponse{
		Name:   backupName,
		Status: http.StatusOK,
	})
}
