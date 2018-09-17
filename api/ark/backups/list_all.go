package backups

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

// ListAll lists every ARK backup for the organization
func ListAll(c *gin.Context) {
	logger := correlationid.Logger(common.Log, c)
	logger.Info("getting backups")

	org := auth.GetCurrentOrganization(c.Request)
	bs := ark.BackupsServiceFactory(org, config.DB(), logger)

	backups, err := bs.List()
	if err != nil {
		err = emperror.Wrap(err, "error getting backups")
		logger.Error(err)
		common.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, backups)
}
