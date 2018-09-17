package restores

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/banzaicloud/pipeline/api/ark/common"
	"github.com/banzaicloud/pipeline/internal/platform/gin/correlationid"
)

// GetResults get results for a completed ARK restore
func GetResults(c *gin.Context) {
	restoreName := c.Param("name")

	logger := correlationid.Logger(common.Log, c).WithField("restore", restoreName)
	logger.Info("getting restore results")

	restore, err := common.GetARKService(c.Request).GetRestoresService().GetByName(restoreName)
	if err != nil {
		err = errors.Wrap(err, "error getting restore")
		logger.Error(err)
		common.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, restore.Results)
}
