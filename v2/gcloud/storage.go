// Copyright 2023 Google LLC
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

package gcloud

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"cloud.google.com/go/storage"
)

func (c *Client) getStorageService(project string) (*storage.Client, error) {
	var err error
	svc := c.services.storage

	if svc != nil {
		return svc, nil
	}

	if err := c.ServiceEnable(project, "storage.googleapis.com"); err != nil {
		return nil, fmt.Errorf("error activating service for polling: %s", err)
	}

	svc, err = storage.NewClient(c.ctx, c.opts)
	if err != nil {
		return nil, err
	}

	c.services.storage = svc

	return svc, nil
}

// StorageBucketCreate creates a storage bucket in Cloud Storage
func (c *Client) StorageBucketCreate(project, bucket string) error {
	svc, err := c.getStorageService(project)
	if err != nil {
		return err
	}

	if err := svc.Bucket(bucket).Create(c.ctx, project, &storage.BucketAttrs{}); err != nil {
		return fmt.Errorf("could not create bucket (%s): %s", bucket, err)
	}

	return nil
}

// StorageBucketDelete deletes a storage bucket in Cloud Storage
func (c *Client) StorageBucketDelete(project, bucket string) error {
	svc, err := c.getStorageService(project)
	if err != nil {
		return err
	}

	if err := svc.Bucket(bucket).Delete(c.ctx); err != nil {
		return fmt.Errorf("could not delete bucket (%s): %s", bucket, err)
	}

	return nil
}

// StorageObjectCreate creates an object in a particular bucket in Cloud Storage
func (c *Client) StorageObjectCreate(project, bucket, path string) (string, error) {
	svc, err := c.getStorageService(project)
	if err != nil {
		return "", err
	}
	name := filepath.Base(path)

	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	obj := svc.Bucket(bucket).Object(name)

	w := obj.NewWriter(c.ctx)
	defer w.Close()

	if _, err := io.Copy(w, file); err != nil {
		return "", fmt.Errorf("could not write file (%s) to bucket (%s): %s", path, bucket, err)
	}

	result := fmt.Sprintf("gs://%s/%s", obj.BucketName(), obj.ObjectName())

	return result, nil
}

// StorageObjectDelete deletes an object in a particular bucket in Cloud Storage
func (c *Client) StorageObjectDelete(project, bucket, gspath string) error {
	svc, err := c.getStorageService(project)
	if err != nil {
		return err
	}
	name := filepath.Base(gspath)

	obj := svc.Bucket(bucket).Object(name)

	obj.Delete(c.ctx)

	return nil
}
