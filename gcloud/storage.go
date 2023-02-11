package gcloud

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"cloud.google.com/go/storage"
)

func (c *Client) getStorageService(project string) (*storage.Client, error) {
	var err error
	svc := c.services.storageService

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

	c.services.storageService = svc

	return svc, nil
}

// StorageBucketCreate creates a storage buck in Cloud Storage
func (c *Client) StorageBucketCreate(project, bucket string) error {
	svc, err := c.getStorageService(project)
	if err != nil {
		return err
	}

	if err := svc.Bucket(bucket).Create(context.Background(), project, &storage.BucketAttrs{}); err != nil {
		return fmt.Errorf("could not create bucket (%s): %s", bucket, err)
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

	w := obj.NewWriter(context.Background())
	defer w.Close()

	if _, err := io.Copy(w, file); err != nil {
		return "", fmt.Errorf("could not write file (%s) to bucket (%s): %s", path, bucket, err)
	}

	result := fmt.Sprintf("gs://%s/%s", obj.BucketName(), obj.ObjectName())

	return result, nil
}
