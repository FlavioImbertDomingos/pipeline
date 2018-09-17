package buckets

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goph/emperror"

	"github.com/banzaicloud/pipeline/api/ark/common"
	"github.com/banzaicloud/pipeline/auth"
	"github.com/banzaicloud/pipeline/config"
	"github.com/banzaicloud/pipeline/internal/ark"
	arkAPI "github.com/banzaicloud/pipeline/internal/ark/api"
	"github.com/banzaicloud/pipeline/internal/platform/gin/correlationid"
)

// Delete deletes an ARK backup bucket
func Delete(c *gin.Context) {
	bucketName := c.Param("name")

	logger := correlationid.Logger(common.Log, c).WithField("bucket", bucketName)
	logger.Info("deleting bucket")

	org := auth.GetCurrentOrganization(c.Request)
	bs := ark.BucketsServiceFactory(org, config.DB(), logger)
	err := bs.DeleteByName(bucketName)
	if err != nil {
		err = emperror.Wrap(err, "error deleting bucket")
		logger.Error(err)
		common.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, &arkAPI.DeleteBucketResponse{
		Name:   bucketName,
		Status: http.StatusOK,
	})
}
