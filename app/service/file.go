package service

import (
	"cloud.google.com/go/storage"
	"errors"
	"github.com/spf13/viper"
	"go-google-cloud-storage/app/constant"
	"go-google-cloud-storage/app/helper"
	"go-google-cloud-storage/app/integration"
	"go-google-cloud-storage/app/resource/response"
)

type FileService struct {
	Viper *viper.Viper
}

func NewFileService(viper *viper.Viper) *FileService {
	return &FileService{
		viper,
	}
}

func (f *FileService) UploadFile(apiCallID, folder, filename string, file []byte) (*response.UploadFileResponse, constant.ResponseMap) {
	gcs := integration.GCS{
		BucketName:         f.Viper.GetString("GCS_BUCKET_NAME"),
		CredentialFilePath: f.Viper.GetString("GCS_CREDENTIAL_FILE_PATH"),
	}
	uploadedPath, err := gcs.Upload(apiCallID, folder, filename, file, f.Viper.GetBool("GCS_CONFIG_SA"))
	if err != nil {
		helper.LogError(apiCallID, "Error upload file : "+err.Error())
		return nil, constant.Res422SomethingWentWrong
	}

	return &response.UploadFileResponse{Path: uploadedPath}, constant.Res200Save

}

func (f *FileService) DownloadFile(apiCallID, filePath string) (*storage.Reader, string, constant.ResponseMap) {
	gcs := integration.GCS{
		BucketName:         f.Viper.GetString("GCS_BUCKET_NAME"),
		CredentialFilePath: f.Viper.GetString("GCS_CREDENTIAL_FILE_PATH"),
	}

	downloadedFile, contentType, err := gcs.Download(apiCallID, filePath, f.Viper.GetBool("GCS_CONFIG_SA"))
	if err != nil {
		helper.LogError(apiCallID, "Error download file : "+err.Error())
		return nil, "", constant.Res422SomethingWentWrong
	}

	return downloadedFile, contentType, constant.Res200Get
}

func (f *FileService) GetAllFile(apiCallID string) (*[]response.GetFileResponse, constant.ResponseMap) {
	gcs := integration.GCS{
		BucketName:         f.Viper.GetString("GCS_BUCKET_NAME"),
		CredentialFilePath: f.Viper.GetString("GCS_CREDENTIAL_FILE_PATH"),
	}

	listFile, err := gcs.List(apiCallID, "", f.Viper.GetBool("GCS_CONFIG_SA"))
	if err != nil {
		if errors.Is(err, storage.ErrObjectNotExist) {
			return nil, constant.Res400FailedDataNotFound
		}
		helper.LogError(apiCallID, "Error list file : "+err.Error())
		return nil, constant.Res422SomethingWentWrong
	}

	formatted := response.GetFileResponseFormatter(listFile)

	return &formatted, constant.Res200Get
}

func (f *FileService) GetSpecificFile(apiCallID, folder string) (*[]response.GetFileResponse, constant.ResponseMap) {
	gcs := integration.GCS{
		BucketName:         f.Viper.GetString("GCS_BUCKET_NAME"),
		CredentialFilePath: f.Viper.GetString("GCS_CREDENTIAL_FILE_PATH"),
	}

	listFile, err := gcs.List(apiCallID, folder, f.Viper.GetBool("GCS_CONFIG_SA"))
	if err != nil {
		if errors.Is(err, storage.ErrObjectNotExist) {
			return nil, constant.Res400FailedDataNotFound
		}
		helper.LogError(apiCallID, "Error list file : "+err.Error())
		return nil, constant.Res422SomethingWentWrong
	}

	formatted := response.GetFileResponseFormatter(listFile)

	return &formatted, constant.Res200Get
}
