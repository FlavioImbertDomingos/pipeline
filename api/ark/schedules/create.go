package schedules

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goph/emperror"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/banzaicloud/pipeline/api/ark/common"
	"github.com/banzaicloud/pipeline/internal/ark/api"
	"github.com/banzaicloud/pipeline/internal/platform/gin/correlationid"
)

// Create creates an ARK schedule
func Create(c *gin.Context) {
	logger := correlationid.Logger(common.Log, c)
	logger.Info("creating schedule")

	var request api.CreateScheduleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		err = errors.Wrap(err, "error parsing request")
		logger.Error(err.Error())
		common.ErrorResponse(c, err)
		return
	}

	spec := &api.CreateBackupRequest{
		Name:    request.Name,
		Labels:  request.Labels,
		TTL:     request.TTL,
		Options: request.Options,
	}

	svc := common.GetARKService(c.Request)

	if spec.Labels == nil {
		spec.Labels = make(labels.Set, 0)
	}
	spec.Labels["pipeline-distribution"] = svc.GetDeploymentsService().GetCluster().GetDistribution()
	spec.Labels["pipeline-cloud"] = svc.GetDeploymentsService().GetCluster().GetCloud()

	err := svc.GetSchedulesService().Create(spec, request.Schedule)
	if err != nil {
		err = emperror.Wrap(err, "error creating schedule")
		logger.Error(err.Error())
		common.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, &api.CreateScheduleResponse{
		Name:   spec.Name,
		Status: http.StatusOK,
	})
}
