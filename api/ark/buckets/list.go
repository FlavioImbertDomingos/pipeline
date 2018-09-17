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

// List lists ARK backup buckets
func List(c *gin.Context) {
	logger := correlationid.Logger(common.Log, c)
	logger.Info("getting buckets")

	org := auth.GetCurrentOrganization(c.Request)
	bs := ark.BucketsServiceFactory(org, config.DB(), logger)
	buckets, err := bs.List()
	if err != nil {
		err = emperror.Wrap(err, "error getting buckets")
		logger.Error(err)
		common.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, buckets)
}
