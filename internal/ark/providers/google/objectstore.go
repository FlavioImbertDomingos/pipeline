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

package google

import (
	"context"
	"encoding/json"
	"io"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	apiStorage "google.golang.org/api/storage/v1"

	"github.com/banzaicloud/pipeline/secret"
)

type objectStore struct {
	client         *storage.Client
	googleAccessID string
	privateKey     []byte
}

// NewObjectStore creates a new objectStore
func NewObjectStore() *objectStore {
	return &objectStore{}
}

func (o *objectStore) Initialize(s *secret.SecretItemResponse) error {
	creds, err := json.Marshal(s.Values)
	if err != nil {
		return err
	}

	jwtConfig, err := google.JWTConfigFromJSON(creds)
	if err != nil {
		return errors.WithStack(err)
	}
	if jwtConfig.Email == "" {
		return errors.Errorf("credentials does not contain an email")
	}
	if len(jwtConfig.PrivateKey) == 0 {
		return errors.Errorf("credentials does not contain a private key")
	}

	o.googleAccessID = jwtConfig.Email
	o.privateKey = jwtConfig.PrivateKey

	ctx := context.Background()

	credentials, err := google.CredentialsFromJSON(ctx, creds, apiStorage.DevstorageFullControlScope)
	if err != nil {
		return err
	}

	client, err := storage.NewClient(ctx, option.WithScopes(storage.ScopeReadWrite), option.WithCredentials(credentials))
	if err != nil {
		return errors.WithStack(err)
	}
	o.client = client

	return nil
}

func (o *objectStore) Init(config map[string]string) error {

	return nil
}

func (o *objectStore) PutObject(bucket string, key string, body io.Reader) error {

	w := o.client.Bucket(bucket).Object(key).NewWriter(context.Background())

	_, copyErr := io.Copy(w, body)
	closeErr := w.Close()
	if copyErr != nil {
		return copyErr
	}

	return closeErr
}

func (o *objectStore) GetObject(bucket string, key string) (io.ReadCloser, error) {
	r, err := o.client.Bucket(bucket).Object(key).NewReader(context.Background())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return r, nil
}

func (o *objectStore) ListCommonPrefixes(bucket string, delimiter string) ([]string, error) {
	q := &storage.Query{
		Delimiter: delimiter,
	}

	var res []string

	iter := o.client.Bucket(bucket).Objects(context.Background(), q)

	for {
		obj, err := iter.Next()
		if err == iterator.Done {
			return res, nil
		}
		if err != nil {
			return nil, errors.WithStack(err)
		}

		if obj.Prefix != "" {
			res = append(res, obj.Prefix[0:strings.LastIndex(obj.Prefix, delimiter)])
		}
	}
}

func (o *objectStore) ListObjects(bucket, prefix string) ([]string, error) {
	q := &storage.Query{
		Prefix: prefix,
	}

	var res []string

	iter := o.client.Bucket(bucket).Objects(context.Background(), q)

	for {
		obj, err := iter.Next()
		if err == iterator.Done {
			return res, nil
		}
		if err != nil {
			return nil, errors.WithStack(err)
		}

		res = append(res, obj.Name)
	}
}

func (o *objectStore) DeleteObject(bucket string, key string) error {
	return errors.Wrapf(o.client.Bucket(bucket).Object(key).Delete(context.Background()), "error deleting object %s", key)
}

func (o *objectStore) CreateSignedURL(bucket, key string, ttl time.Duration) (string, error) {
	return storage.SignedURL(bucket, key, &storage.SignedURLOptions{
		GoogleAccessID: o.googleAccessID,
		PrivateKey:     o.privateKey,
		Method:         "GET",
		Expires:        time.Now().Add(ttl),
	})
}
