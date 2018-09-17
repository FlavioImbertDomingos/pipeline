package buckets

import (
	"net/http"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"
	"github.com/goph/emperror"
	"github.com/pkg/errors"

	"github.com/banzaicloud/pipeline/api/ark/common"
	"github.com/banzaicloud/pipeline/auth"
	"github.com/banzaicloud/pipeline/config"
	"github.com/banzaicloud/pipeline/internal/ark"
	"github.com/banzaicloud/pipeline/internal/ark/api"
	"github.com/banzaicloud/pipeline/internal/platform/gin/correlationid"
)

// Create creates an ARK backup bucket
func Create(c *gin.Context) {
	logger := correlationid.Logger(common.Log, c)
	logger.Info("creating bucket")

	var request api.CreateBucketRequest
	if err := c.BindJSON(&request); err != nil {
		err = errors.Wrap(err, "error parsing request")
		logger.Error(err.Error())
		common.ErrorResponse(c, err)
		return
	}

	org := auth.GetCurrentOrganization(c.Request)
	bs := ark.BucketsServiceFactory(org, config.DB(), logger)

	_, err := bs.GetByRequest(api.FindBucketRequest{
		Cloud:      request.Cloud,
		BucketName: request.BucketName,
		Location:   request.Location,
	})
	if err == nil {
		err = errors.New("bucket already exists")
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		err = emperror.Wrap(err, "error creating bucket")
		logger.Error(err.Error())
		common.ErrorResponse(c, err)
		return
	}

	bucket, err := bs.FindOrCreateBucket(&api.CreateBucketRequest{
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

	c.JSON(http.StatusOK, bucket.ConvertModelToEntity())
}
