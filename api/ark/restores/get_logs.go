package restores

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/banzaicloud/pipeline/api/ark/common"
	"github.com/banzaicloud/pipeline/internal/platform/gin/correlationid"
)

// GetLogs get logs for an ARK restore
func GetLogs(c *gin.Context) {
	restoreName := c.Param("name")

	logger := correlationid.Logger(common.Log, c).WithField("restore", restoreName)
	logger.Info("getting restore logs")

	svc := common.GetARKService(c.Request)

	restore, err := svc.GetRestoresService().GetByName(restoreName)
	if err != nil {
		err = errors.Wrap(err, "error getting restore")
		logger.Error(err)
		common.ErrorResponse(c, err)
		return
	}

	err = svc.GetBucketsService().StreamRestoreLogsFromObjectStore(
		restore.Bucket,
		restore.BackupName,
		restore.Name,
		c.Writer,
	)
	if err != nil {
		err = errors.Wrap(err, "error streaming logs")
		logger.Error(err)
		common.ErrorResponse(c, err)
		return
	}
}
