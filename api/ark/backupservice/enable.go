package backupservice

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goph/emperror"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/banzaicloud/pipeline/api/ark/common"
	"github.com/banzaicloud/pipeline/internal/ark/api"
	"github.com/banzaicloud/pipeline/internal/platform/gin/correlationid"
)

// Enable create an ARK service deployment and adding a base scheduled full backup
func Enable(c *gin.Context) {
	logger := correlationid.Logger(common.Log, c)
	logger.Debug("deploying backup service to cluster")

	svc := common.GetARKService(c.Request)

	var request api.EnableBackupServiceRequest
	if err := c.BindJSON(&request); err != nil {
		err = errors.Wrap(err, "error parsing request")
		logger.Error(err.Error())
		common.ErrorResponse(c, err)
		return
	}

	_, err := svc.GetDeploymentsService().GetActiveDeployment()
	if err == nil {
		err = errors.New("backup service already deployed")
		logger.Error(err.Error())
		common.ErrorResponse(c, err)
		return
	}

	bucketService := svc.GetBucketsService()
	bucket, err := bucketService.FindOrCreateBucket(&api.CreateBucketRequest{
		Cloud:      request.Cloud,
		BucketName: request.BucketName,
		Location:   request.Location,
		SecretID:   request.SecretID,
	})
	if err != nil {
		err = emperror.Wrap(err, "error persisting bucket")
		logger.Error(err.Error())
		common.ErrorResponse(c, err)
		return
	}

	err = bucketService.IsBucketInUse(bucket)
	if err != nil {
		logger.Error(err.Error())
		common.ErrorResponse(c, err)
		return
	}

	err = svc.GetDeploymentsService().Deploy(bucket, false)
	if err != nil {
		err = emperror.Wrap(err, "error deploying backup service")
		logger.Error(err.Error())
		common.ErrorResponse(c, err)
		return
	}

	duration, _ := time.ParseDuration("5m")
	spec := &api.CreateBackupRequest{
		Name:   "pipeline-full-backup",
		Labels: request.Labels,
		TTL: metav1.Duration{
			Duration: duration,
		},
	}

	if spec.Labels == nil {
		spec.Labels = make(labels.Set, 0)
	}
	spec.Labels["pipeline-distribution"] = svc.GetDeploymentsService().GetCluster().GetDistribution()
	spec.Labels["pipeline-cloud"] = svc.GetDeploymentsService().GetCluster().GetCloud()

	err = svc.GetSchedulesService().Create(spec, request.Schedule)
	if err != nil {
		err = emperror.Wrap(err, "error creating schedule")
		logger.Error(err.Error())
		common.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, api.EnableBackupServiceResponse{
		Status: http.StatusOK,
	})
}
