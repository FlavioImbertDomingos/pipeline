package restores

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goph/emperror"

	"github.com/banzaicloud/pipeline/api/ark/common"
	"github.com/banzaicloud/pipeline/internal/ark/api"
	"github.com/banzaicloud/pipeline/internal/platform/gin/correlationid"
)

// Delete deletes an ARK restore
func Delete(c *gin.Context) {
	restoreName := c.Param("name")

	logger := correlationid.Logger(common.Log, c).WithField("restore", restoreName)
	logger.Info("getting restore")

	err := common.GetARKService(c.Request).GetRestoresService().DeleteByName(restoreName)
	if err != nil {
		err = emperror.Wrap(err, "error deleting restore")
		logger.Error(err)
		common.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, &api.DeleteRestoreResponse{
		Name:   restoreName,
		Status: http.StatusOK,
	})
}
