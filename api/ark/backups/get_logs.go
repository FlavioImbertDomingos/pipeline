package backups

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/banzaicloud/pipeline/api/ark/common"
	"github.com/banzaicloud/pipeline/internal/platform/gin/correlationid"
)

// GetLogs gets ARK backup logs from object store
func GetLogs(c *gin.Context) {
	backupName := c.Param("name")

	logger := correlationid.Logger(common.Log, c).WithField("backup", backupName)
	logger.Info("getting backup logs")

	svc := common.GetARKService(c.Request)

	backup, err := svc.GetBackupsService().GetByName(backupName)
	if err != nil {
		err = errors.Wrap(err, "error getting backup")
		logger.Error(err)
		common.ErrorResponse(c, err)
		return
	}

	err = svc.GetBucketsService().StreamBackupLogsFromObjectStore(backup.Bucket, backupName, c.Writer)
	if err != nil {
		err = errors.Wrap(err, "error streaming logs")
		logger.Error(err)
		common.ErrorResponse(c, err)
		return
	}
}
