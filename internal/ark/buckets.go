package ark

import (
	"github.com/pkg/errors"

	"github.com/banzaicloud/pipeline/auth"
	"github.com/banzaicloud/pipeline/internal/ark/api"
)

// ValidateCreateBucketRequest validates a CreateBucketRequest
func ValidateCreateBucketRequest(req *api.CreateBucketRequest, org *auth.Organization) error {

	err := IsProviderSupported(req.Cloud)
	if err != nil {
		return errors.Wrap(err, req.Cloud)
	}

	if req.Cloud == "amazon" {
		if req.Location == "" {
			return errors.Wrap(errors.New("location must not be empty"), "error validating create bucket request")
		}
	} else {
		req.Location = ""
	}

	secret, err := GetSecretWithValidation(req.SecretID, org.ID, req.Cloud)
	if err != nil {
		return errors.Wrap(err, "error validating create bucket request")
	}

	os, err := NewObjectStore(req.Cloud)
	if err != nil {
		return errors.Wrap(err, "error validating create bucket request")
	}
	os.Initialize(secret)

	_, err = os.ListCommonPrefixes(req.BucketName, "/")
	if err != nil {
		return errors.Wrap(err, "error validating create bucket request")
	}

	return nil
}
