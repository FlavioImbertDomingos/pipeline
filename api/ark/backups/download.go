package backups

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/banzaicloud/pipeline/api/ark/common"
	"github.com/banzaicloud/pipeline/internal/platform/gin/correlationid"
)

// Download downloads an ARK backup contents from object store
func Download(c *gin.Context) {
	backupName := c.Param("name")

	logger := correlationid.Logger(common.Log, c).WithField("backup", backupName)
	logger.Info("downloading backup contents")

	svc := common.GetARKService(c.Request)
	backup, err := svc.GetBackupsService().GetByName(backupName)
	if err != nil {
		err = errors.Wrap(err, "error getting backup")
		logger.Error(err)
		common.ErrorResponse(c, err)
		return
	}

	c.Writer.Header().Set("Content-Type", "application/x-gzip")
	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+backupName+".tgz")

	err = svc.GetBucketsService().StreamBackupContentsFromObjectStore(backup.Bucket, backupName, c.Writer)
	if err != nil {
		err = errors.Wrap(err, "error streaming contents")
		logger.Error(err)
		common.ErrorResponse(c, err)
		return
	}
}
