package integration

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"google.golang.org/api/iterator"
	"io"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
	"go-google-cloud-storage/app/helper"
	"google.golang.org/api/option"
)

type GCS struct {
	ProjectID          string
	BucketName         string
	StorageClassBucket string
	LocationBucket     string
	CredentialFilePath string
	ConfigSA           bool
}

type gcsConfigChecklist struct {
	Key   string
	Valid bool
}

func GCSInstance(v *viper.Viper) (*GCS, error) {
	gcs := &GCS{
		ProjectID:          v.GetString("GCS_PROJECT_ID"),
		BucketName:         v.GetString("GCS_BUCKET_NAME"),
		CredentialFilePath: v.GetString("GCS_CREDENTIAL_FILE_PATH"),
		StorageClassBucket: v.GetString("GCS_STORAGE_CLASS_BUCKET"),
		LocationBucket:     v.GetString("GCS_STORAGE_LOCATION_BUCKET"),
		ConfigSA:           v.GetBool("GCS_CONFIG_SA"),
	}

	checks := []gcsConfigChecklist{
		{"GCS_PROJECT_ID", gcs.ProjectID != ""},
		{"GCS_BUCKET_NAME", gcs.BucketName != ""},
		{"GCS_CREDENTIAL_FILE_PATH", gcs.CredentialFilePath != ""},
		{"GCS_STORAGE_CLASS_BUCKET", gcs.StorageClassBucket != ""},
		{"GCS_STORAGE_LOCATION_BUCKET", gcs.LocationBucket != ""},
	}

	for _, check := range checks {
		if !check.Valid {
			return nil, fmt.Errorf("missing or invalid required GCS configuration: %s", check.Key)
		}
	}

	return gcs, nil
}

func (g *GCS) initClient(ctx context.Context, apiCallID string, useLocalCredential bool) (*storage.Client, error) {
	var opts []option.ClientOption
	if useLocalCredential {
		opts = append(opts, option.WithCredentialsFile(g.CredentialFilePath))
		helper.LogInfo(apiCallID, "Config GCS with local credential")
	}
	return storage.NewClient(ctx, opts...)
}

func (g *GCS) Upload(apiCallID, folder, filename string, fileData []byte) (string, error) {
	path := filepath.Join(folder, helper.GenerateUniqueFilename()+filepath.Ext(filename))
	helper.LogInfo(apiCallID, "Uploading file: "+path)

	ctx := context.Background()
	client, err := g.initClient(ctx, apiCallID, g.ConfigSA)
	if err != nil {
		return "", fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	o := client.Bucket(g.BucketName).Object(path)
	o = o.If(storage.Conditions{DoesNotExist: true})
	wc := o.NewWriter(ctx)

	if _, err := io.Copy(wc, bytes.NewReader(fileData)); err != nil {
		return "", fmt.Errorf("io.Copy: %w", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("Writer.Close: %w", err)
	}
	helper.LogInfo(apiCallID, "Uploaded file successfully: "+path)

	return path, nil
}

func (g *GCS) Download(apiCallID, objectPath string) (*storage.Reader, string, error) {
	ctx := context.Background()
	client, err := g.initClient(ctx, apiCallID, g.ConfigSA)
	if err != nil {
		return nil, "", fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*100)
	defer cancel()

	rc := client.Bucket(g.BucketName).Object(objectPath)
	nr, err := rc.NewReader(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("Object(%q).NewReader: %w", objectPath, err)
	}

	attr, err := rc.Attrs(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("io.ReadAll: %w", err)
	}

	helper.LogInfo(apiCallID, "Download file successfully: "+objectPath)

	return nr, attr.ContentType, nil
}

func (g *GCS) List(apiCallID string, folder string) (*[]storage.ObjectAttrs, error) {
	ctx := context.Background()
	client, err := g.initClient(ctx, apiCallID, g.ConfigSA)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*100)
	defer cancel()

	var results []storage.ObjectAttrs
	query := &storage.Query{}
	if folder != "" {
		query.Prefix = folder
	}

	it := client.Bucket(g.BucketName).Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Bucket(%q).Objects: %w", g.BucketName, err)
		}
		results = append(results, *attrs)
	}

	if len(results) == 0 {
		return nil, storage.ErrObjectNotExist
	}

	return &results, nil
}

func (g *GCS) CreateBucket(apiCallID, bucketName string) error {
	ctx := context.Background()
	client, err := g.initClient(ctx, apiCallID, g.ConfigSA)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	storageClassAndLocation := &storage.BucketAttrs{
		StorageClass: g.StorageClassBucket,
		Location:     g.LocationBucket,
	}
	bucket := client.Bucket(bucketName)
	if err := bucket.Create(ctx, g.ProjectID, storageClassAndLocation); err != nil {
		return fmt.Errorf("Bucket(%q).Create: %w", bucketName, err)
	}
	return nil
}

func (g *GCS) DeleteFile(apiCallID, path string) error {
	ctx := context.Background()
	client, err := g.initClient(ctx, apiCallID, g.ConfigSA)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	o := client.Bucket(g.BucketName).Object(path)

	attrs, err := o.Attrs(ctx)
	if err != nil {
		return fmt.Errorf("object.Attrs: %w", err)
	}
	o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

	if err := o.Delete(ctx); err != nil {
		return fmt.Errorf("Object(%q).Delete: %w", path, err)
	}
	return nil
}
