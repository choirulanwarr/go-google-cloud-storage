package integration

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
	"go-google-cloud-storage/app/helper"
	"google.golang.org/api/option"
)

type GCS struct {
	BucketName         string
	CredentialFilePath string
}

var allowedMIMETypes = map[string]bool{
	"image/jpeg":               true,
	"image/png":                true,
	"image/jpg":                true,
	"application/pdf":          true,
	"application/zip":          true,
	"application/octet-stream": true,
}

func IsAllowedFileType(apiCallID string, file multipart.File) bool {
	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		return false
	}
	_, _ = file.Seek(0, io.SeekStart)

	mimeType := http.DetectContentType(buffer)
	helper.LogInfo(apiCallID, "Content-Type: "+mimeType)
	return allowedMIMETypes[mimeType]
}

func DefaultMIME(mime string) string {
	if mime == "" {
		return "application/octet-stream"
	}
	return mime
}

func generateUniqueFilename() string {
	now := time.Now()
	timestamp := fmt.Sprintf("%d%02d%02d%02d", now.Unix(), now.Hour(), now.Minute(), now.Second())
	randomSuffix := rand.Intn(10_000_000) + 1
	return fmt.Sprintf("%s%d", timestamp, randomSuffix)
}

func (g *GCS) initClient(ctx context.Context, apiCallID string, useLocalCredential bool) (*storage.Client, error) {
	var opts []option.ClientOption
	if useLocalCredential {
		opts = append(opts, option.WithCredentialsFile(g.CredentialFilePath))
		helper.LogInfo(apiCallID, "Config GCS with local credential")
	}
	return storage.NewClient(ctx, opts...)
}

func (g *GCS) Upload(apiCallID, folder, filename string, fileData []byte, useLocalCredential bool) (string, error) {
	path := filepath.Join(folder, generateUniqueFilename()+filepath.Ext(filename))
	helper.LogInfo(apiCallID, "Uploading file: "+path)

	ctx := context.Background()
	client, err := g.initClient(ctx, apiCallID, useLocalCredential)
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

func (g *GCS) Download(apiCallID, objectPath string, useLocalCredential bool) (*storage.Reader, string, error) {
	ctx := context.Background()
	client, err := g.initClient(ctx, apiCallID, useLocalCredential)
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
