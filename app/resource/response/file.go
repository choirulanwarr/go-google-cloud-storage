package response

import (
	"cloud.google.com/go/storage"
	"go-google-cloud-storage/app/helper"
	"path/filepath"
	"time"
)

type UploadFileResponse struct {
	Path string `json:"path"`
}

type GetFileResponse struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	Size      string `json:"size"`
	Type      string `json:"type"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func GetFileResponseFormatter(listFile *[]storage.ObjectAttrs) []GetFileResponse {
	var result []GetFileResponse

	for _, file := range *listFile {
		result = append(result, GetFileResponse{
			Name:      filepath.Base(file.Name),
			URL:       helper.GeneratePublicURL(file.Bucket, file.Name),
			Size:      helper.FormatFileSize(file.Size),
			Type:      file.ContentType,
			CreatedAt: file.Created.Format(time.RFC3339),
			UpdatedAt: file.Updated.Format(time.RFC3339),
		})
	}

	return result
}
