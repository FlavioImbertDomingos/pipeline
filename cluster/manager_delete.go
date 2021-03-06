// Copyright © 2018 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cluster

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/banzaicloud/pipeline/dns"
	"github.com/banzaicloud/pipeline/helm"
	pkgCluster "github.com/banzaicloud/pipeline/pkg/cluster"
	"github.com/goph/emperror"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const retry = 3

// DeleteCluster deletes a cluster.
func (m *Manager) DeleteCluster(ctx context.Context, cluster CommonCluster, force bool, kubeProxyCache *sync.Map) error {
	errorHandler := emperror.HandlerWith(
		m.getErrorHandler(ctx),
		"organization", cluster.GetOrganizationId(),
		"cluster", cluster.GetID(),
		"force", force,
	)

	go func() {
		err := m.deleteCluster(ctx, cluster, force, kubeProxyCache)
		if err != nil {
			errorHandler.Handle(err)
		}
	}()

	return nil
}

func deleteAllResource(kubeConfig []byte, logger *logrus.Entry) error {
	client, err := helm.GetK8sConnection(kubeConfig)
	if err != nil {
		return err
	}

	type resourceDeleter interface {
		Delete(name string, options *metav1.DeleteOptions) error
	}

	services := []resourceDeleter{
		client.CoreV1().Services(metav1.NamespaceAll),
		client.AppsV1().Deployments(metav1.NamespaceAll),
		client.AppsV1().DaemonSets(metav1.NamespaceAll),
		client.AppsV1().StatefulSets(metav1.NamespaceAll),
		client.AppsV1().ReplicaSets(metav1.NamespaceAll),
	}

	options := metav1.NewDeleteOptions(0)

	for _, service := range services {
		for i := 0; i < retry; i++ {
			err := service.Delete("", options)
			if err != nil {
				logger.Debugf("deleting resources %T attempt %d/%d failed", service, i, retry)
				time.Sleep(1)
			}
		}
	}
	return nil
}

// deleteDnsRecordsOwnedByCluster deletes DNS records owned by the cluster. These are the DNS records
// created for the public endpoints of the services hosted by the cluster.
func deleteDnsRecordsOwnedByCluster(cluster CommonCluster) error {
	dnsSvc, err := dns.GetExternalDnsServiceClient()
	if err != nil {
		return emperror.Wrap(err, "getting external dns service client failed")
	}

	if dnsSvc == nil {
		return nil
	}

	err = dnsSvc.DeleteDnsRecordsOwnedBy(cluster.GetUID(), cluster.GetOrganizationId())
	if err != nil {
		return emperror.Wrapf(err, "deleting DNS records owned by cluster failed")
	}

	return nil
}

func (m *Manager) deleteCluster(ctx context.Context, cluster CommonCluster, force bool, kubeProxyCache *sync.Map) error {
	logger := m.getLogger(ctx).WithFields(logrus.Fields{
		"organization": cluster.GetOrganizationId(),
		"cluster":      cluster.GetID(),
		"force":        force,
	})

	logger.Info("deleting cluster")

	err := cluster.UpdateStatus(pkgCluster.Deleting, pkgCluster.DeletingMessage)
	if err != nil {
		return emperror.With(
			emperror.Wrap(err, "cluster status update failed"),
			"cluster", cluster.GetID(),
		)
	}

	// get kubeconfig
	c, err := cluster.GetK8sConfig()
	if err != nil {
		if !force {
			cluster.UpdateStatus(pkgCluster.Error, err.Error())

			return emperror.Wrap(err, "error getting kubeconfig")
		}

		logger.Errorf("error during getting kubeconfig: %s", err.Error())
	}

	if !(force && c == nil) {
		// delete deployments
		for i := 0; i < retry; i++ {
			err = helm.DeleteAllDeployment(c)
			// TODO we could check to the Authorization IAM error explicit
			if err != nil {
				logger.Errorf("deleting deployments attempt %d/%d failed: %s", i, retry, err.Error())
				time.Sleep(1)
			} else {
				break
			}
		}
		err = deleteAllResource(c, logger)
		if err != nil {
			if force {
				logger.Errorf("deleting resources failed: %s", err.Error())
			} else {
				return emperror.Wrap(err, "deleting resources failed")
			}
		}

	} else {
		logger.Info("skipping deployment deletion without kubeconfig")
	}

	// clean up dns registrations
	err = deleteDnsRecordsOwnedByCluster(cluster)

	// delete cluster
	err = cluster.DeleteCluster()
	if err != nil {
		if !force {
			cluster.UpdateStatus(pkgCluster.Error, err.Error())

			return emperror.Wrap(err, "error deleting cluster")
		}

		logger.Errorf("error during deleting cluster: %s", err.Error())
	}

	// delete from proxy from kubeProxyCache if any
	// TODO: this should be handled somewhere else
	kubeProxyCache.Delete(fmt.Sprint(cluster.GetOrganizationId(), "-", cluster.GetID()))

	// delete cluster from database
	deleteName := cluster.GetName()
	err = cluster.DeleteFromDatabase()
	if err != nil {
		if !force {
			cluster.UpdateStatus(pkgCluster.Error, err.Error())

			return emperror.Wrap(err, "error deleting cluster from the database")
		}

		logger.Errorf("error during deleting cluster from the database: %s", err.Error())
	}

	// Asyncron update prometheus
	go func() {
		err := UpdatePrometheusConfig()
		if err != nil {
			logger.Warnf("could not update prometheus configmap: %v", err)
		}
	}()

	// clean statestore
	logger.Info("cleaning cluster's statestore folder")
	if err := CleanStateStore(deleteName); err != nil {
		return emperror.Wrap(err, "cleaning cluster statestore failed")
	}

	logger.Info("cluster's statestore folder cleaned")

	logger.Info("cluster deleted successfully")

	return nil
}
