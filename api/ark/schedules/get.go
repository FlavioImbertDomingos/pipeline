package schedules

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	common "github.com/banzaicloud/pipeline/api/ark/common"
	"github.com/banzaicloud/pipeline/internal/platform/gin/correlationid"
)

// Get gets an ARK schedule
func Get(c *gin.Context) {
	scheduleName := c.Param("name")

	logger := correlationid.Logger(common.Log, c).WithField("schedule", scheduleName)
	logger.Info("getting schedule")

	schedule, err := common.GetARKService(c.Request).GetSchedulesService().GetByName(scheduleName)
	if err != nil {
		err = errors.Wrap(err, "error getting backup")
		logger.Error(err.Error())
		common.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, schedule)
}
