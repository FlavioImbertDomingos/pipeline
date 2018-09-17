package buckets

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goph/emperror"

	"github.com/banzaicloud/pipeline/api/ark/common"
	"github.com/banzaicloud/pipeline/auth"
	"github.com/banzaicloud/pipeline/config"
	"github.com/banzaicloud/pipeline/internal/ark"
	"github.com/banzaicloud/pipeline/internal/platform/gin/correlationid"
)

// Get gets an ARK backup bucket
func Get(c *gin.Context) {
	bucketName := c.Param("name")

	logger := correlationid.Logger(common.Log, c).WithField("bucket", bucketName)
	logger.Info("getting buckets")

	org := auth.GetCurrentOrganization(c.Request)
	bs := ark.BucketsServiceFactory(org, config.DB(), logger)
	bucket, err := bs.GetByName(bucketName)
	if err != nil {
		err = emperror.Wrap(err, "error getting bucket")
		logger.Error(err)
		common.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, bucket)
}
