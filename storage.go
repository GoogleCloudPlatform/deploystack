package deploystack

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"cloud.google.com/go/storage"
)

var storageService *storage.Client

func getStorageService() (*storage.Client, error) {
	if storageService != nil {
		return storageService, nil
	}

	ctx := context.Background()
	svc, err := storage.NewClient(ctx, opts)
	if err != nil {
		return nil, err
	}

	storageService = svc

	return svc, nil
}

func CreateStorageBucket(project, bucket string) error {
	svc, err := getStorageService()
	if err != nil {
		return err
	}

	if err := svc.Bucket(bucket).Create(context.Background(), project, &storage.BucketAttrs{}); err != nil {
		return fmt.Errorf("could not create bucket (%s): %s", bucket, err)
	}

	return nil
}

func CreateStorageObject(project, bucket, path string) (string, error) {
	svc, err := getStorageService()
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
