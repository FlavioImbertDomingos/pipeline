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

package ark

import (
	"github.com/spf13/viper"

	"github.com/banzaicloud/pipeline/internal/ark/providers/google"
	pkgErrors "github.com/banzaicloud/pipeline/pkg/errors"
	"github.com/banzaicloud/pipeline/pkg/providers"
	"github.com/banzaicloud/pipeline/secret"
)

// ChartConfig describes an ARK deployment chart config
type ChartConfig struct {
	Namespace      string
	Chart          string
	Name           string
	ValueOverrides []byte
}

// ValueOverrides descibes values to be overridden in a deployment
type ValueOverrides struct {
	Configuration configuration `json:"configuration"`
	Credentials   credentials   `json:"credentials"`
	Image         image         `json:"image"`
}

type image struct {
	Repository string `json:"repository"`
	Tag        string `json:"tag"`
	PullPolicy string `json:"pullPolicy"`
}

type credentials struct {
	SecretContents secretContents `json:"secretContents"`
}

type secretContents struct {
	Cloud string `json:"cloud"`
}

type configuration struct {
	PersistentVolumeProvider persistentVolumeProvider `json:"persistentVolumeProvider"`
	BackupStorageProvider    backupStorageProvider    `json:"backupStorageProvider"`
	RestoreOnlyMode          bool
}

type persistentVolumeProvider struct {
	Name string `json:"name"`
}

type backupStorageProvider struct {
	Name   string `json:"name"`
	Bucket string `json:"bucket"`
	Region string `json:"region,omitempty"`
}

// ConfigRequest describes an ARK config request
type ConfigRequest struct {
	Cloud          string
	BucketName     string
	BucketProvider string
	BucketLocation string
	RestoreMode    bool

	Secret *secret.SecretItemResponse
}

// GetChartConfig get a ChartConfig
func GetChartConfig() ChartConfig {

	return ChartConfig{
		Name:      viper.GetString("ark.name"),
		Namespace: viper.GetString("ark.namespace"),
		Chart:     viper.GetString("ark.chart"),
	}
}

// Get gets helm deployment value overrides
func (req ConfigRequest) Get() (values ValueOverrides, err error) {

	pvp, err := req.getPVPConfig()
	if err != nil {
		return values, err
	}

	bsp, err := req.getBSPConfig()
	if err != nil {
		return values, err
	}

	cred, err := req.getCredentials()
	if err != nil {
		return values, err
	}

	return ValueOverrides{
		Configuration: configuration{
			PersistentVolumeProvider: pvp,
			BackupStorageProvider:    bsp,
			RestoreOnlyMode:          req.RestoreMode,
		},
		Credentials: cred,
		Image: image{
			Repository: viper.GetString("ark.image"),
			Tag:        viper.GetString("ark.imagetag"),
			PullPolicy: viper.GetString("ark.pullpolicy"),
		},
	}, nil
}

func (req ConfigRequest) getPVPConfig() (persistentVolumeProvider, error) {

	var config persistentVolumeProvider
	var pvc string

	switch req.Cloud {
	case providers.Google:
		pvc = google.PersistentVolumeProvider
	default:
		return config, pkgErrors.ErrorNotSupportedCloudType
	}

	return persistentVolumeProvider{
		Name: pvc,
	}, nil
}

func (req ConfigRequest) getBSPConfig() (backupStorageProvider, error) {

	var config backupStorageProvider
	var bsp string

	switch req.BucketProvider {
	case providers.Google:
		bsp = google.PersistentVolumeProvider
	default:
		return config, pkgErrors.ErrorNotSupportedCloudType
	}

	return backupStorageProvider{
		Name:   bsp,
		Bucket: req.BucketName,
		Region: req.BucketLocation,
	}, nil
}

func (req ConfigRequest) getCredentials() (credentials, error) {

	var config credentials
	var sContents string
	var err error

	switch req.Cloud {
	case providers.Google:
		sContents, err = google.GetSecretContents(req.Secret)
		if err != nil {
			return config, err
		}
	default:
		return config, pkgErrors.ErrorNotSupportedCloudType
	}

	return credentials{
		SecretContents: secretContents{
			Cloud: sContents,
		},
	}, err
}
