package request

import "mime/multipart"

type UploadFileRequest struct {
	Folder string                `form:"folder" validate:"required,not_only_space"`
	File   *multipart.FileHeader `form:"file" validate:"required,not_only_space"`
}

type DownloadFileRequest struct {
	Path string `validate:"required,not_only_space"`
}

type CreateBucketRequest struct {
	BucketName string `validate:"required,not_only_space"`
}
