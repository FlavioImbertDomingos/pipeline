package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/banzaicloud/pipeline/config"
	pkgCommon "github.com/banzaicloud/pipeline/pkg/common"
)

// Log is a logrus.FieldLogger
var Log logrus.FieldLogger

// init initializes the fieldlogger
func init() {
	Log = config.Logger()
}

// ErrorResponseWithStatus aborts the http request with a JSON error response with the given status code and error
func ErrorResponseWithStatus(c *gin.Context, status int, err error) {

	if c.Writer.Status() != http.StatusOK {
		return
	}

	c.AbortWithStatusJSON(status, pkgCommon.ErrorResponse{
		Code:    status,
		Message: err.Error(),
		Error:   errors.Cause(err).Error(),
	})
}

// ErrorResponse aborts the http request with a JSON error response with a status code and error
func ErrorResponse(c *gin.Context, err error) {

	status := http.StatusBadRequest

	if errors.Cause(err) == gorm.ErrRecordNotFound {
		status = http.StatusNotFound
	}

	ErrorResponseWithStatus(c, status, err)
}
