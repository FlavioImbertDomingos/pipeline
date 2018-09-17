package restores

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goph/emperror"

	"github.com/banzaicloud/pipeline/api/ark/common"
	"github.com/banzaicloud/pipeline/internal/platform/gin/correlationid"
)

// Get gets an ARK restore
func Get(c *gin.Context) {
	restoreName := c.Param("name")

	logger := correlationid.Logger(common.Log, c).WithField("restore", restoreName)
	logger.Info("getting restore")

	restore, err := common.GetARKService(c.Request).GetRestoresService().GetByName(restoreName)
	if err != nil {
		err = emperror.Wrap(err, "error getting restore")
		logger.Error(err)
		common.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, restore)
}
