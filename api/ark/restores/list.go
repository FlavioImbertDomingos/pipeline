package restores

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goph/emperror"

	"github.com/banzaicloud/pipeline/api/ark/common"
	"github.com/banzaicloud/pipeline/internal/platform/gin/correlationid"
)

// List lists ARK restores
func List(c *gin.Context) {
	logger := correlationid.Logger(common.Log, c)
	logger.Info("getting restores")

	restores, err := common.GetARKService(c.Request).GetRestoresService().List()
	if err != nil {
		err = emperror.Wrap(err, "error getting backups")
		logger.Error(err)
		common.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, restores)
}
