package schedules

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/banzaicloud/pipeline/api/ark/common"
	arkAPI "github.com/banzaicloud/pipeline/internal/ark/api"
	"github.com/banzaicloud/pipeline/internal/platform/gin/correlationid"
)

// Delete deletes an ARK schedule
func Delete(c *gin.Context) {
	scheduleName := c.Param("name")

	logger := correlationid.Logger(common.Log, c).WithField("schedule", scheduleName)
	logger.Info("deleting schedule")

	err := common.GetARKService(c.Request).GetSchedulesService().DeleteByName(scheduleName)
	if err != nil {
		err = errors.Wrap(err, "error deleting schedule")
		logger.Error(err.Error())
		common.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, &arkAPI.DeleteScheduleResponse{
		Name:   scheduleName,
		Status: http.StatusOK,
	})
}
