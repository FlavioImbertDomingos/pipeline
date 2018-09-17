package common

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goph/emperror"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	"github.com/banzaicloud/pipeline/auth"
	"github.com/banzaicloud/pipeline/cluster"
	"github.com/banzaicloud/pipeline/internal/ark"
	intCluster "github.com/banzaicloud/pipeline/internal/cluster"
	"github.com/banzaicloud/pipeline/internal/platform/gin/utils"
	"github.com/banzaicloud/pipeline/model"
)

const (
	currentClusterModelName = "clusterModel"
	arkServiceName          = "arkService"
)

// ARKMiddleware is a middleware for initializing a CommonCluster and an ARKService
// from the request parameters for later use
func ARKMiddleware(db *gorm.DB, logger logrus.FieldLogger) gin.HandlerFunc {

	return func(c *gin.Context) {
		logger = logger.WithField("middleware", "ARK")

		clusters := intCluster.NewClusters(db)
		clusterID, ok := ginutils.UintParam(c, "id")
		if !ok {
			logger.Error("invalid ID parameter")
			c.Abort()
			return
		}
		org := auth.GetCurrentOrganization(c.Request)
		if org == nil {
			err := errors.New("invalid organization")
			ErrorResponse(c, err)
			logger.Error(err)
			c.Abort()
			return
		}
		cl, err := clusters.FindOneByID(org.ID, clusterID)
		if err != nil {
			err = emperror.Wrap(err, "error getting cluster model")
			ErrorResponse(c, err)
			logger.Error(err)
			c.Abort()
			return
		}
		c.Request = setVariableToContext(c.Request, currentClusterModelName, cl)

		cluster, err := cluster.GetCommonClusterFromModel(cl)
		if err != nil {
			err = emperror.Wrap(err, "error getting cluster")
			ErrorResponse(c, err)
			logger.Error(err)
			c.Abort()
			return
		}

		svc := ark.NewARKService(org, cluster, db, logger)
		c.Request = setVariableToContext(c.Request, arkServiceName, svc)

		c.Next()
	}
}

// GetCurrentClusterModel return the current cluster model
func GetCurrentClusterModel(req *http.Request) *model.ClusterModel {
	if cluster := req.Context().Value(currentClusterModelName); cluster != nil {
		return cluster.(*model.ClusterModel)
	}
	return nil
}

// GetARKService return the current ark.Service model
func GetARKService(req *http.Request) *ark.Service {
	if svc := req.Context().Value(arkServiceName); svc != nil {
		return svc.(*ark.Service)
	}
	return nil
}

func setVariableToContext(req *http.Request, key interface{}, val interface{}) *http.Request {

	return req.WithContext(context.WithValue(req.Context(), key, val))
}
