package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go-google-cloud-storage/app/constant"
	"go-google-cloud-storage/app/helper"
	"go-google-cloud-storage/app/resource/request"
	"go-google-cloud-storage/app/service"
	"io"
	"net/http"
	"path/filepath"
)

type FileHandler struct {
	Service   *service.FileService
	Validator *validator.Validate
}

func NewFileHandler(service *service.FileService, validator *validator.Validate) *FileHandler {
	return &FileHandler{
		service,
		validator,
	}
}

func (f *FileHandler) UploadFile(ctx *gin.Context) {
	apiCallID := ctx.GetString(constant.RequestIDKey)

	var req request.UploadFileRequest
	if err := ctx.ShouldBind(&req); err != nil {
		helper.LogError(apiCallID, "Failed to bind request: "+err.Error())
		helper.ResponseAPI(ctx, constant.Res400InvalidPayload)
		return
	}

	if err := f.Validator.Struct(req); err != nil {
		helper.LogError(apiCallID, "Payload validation failed: "+err.Error())
		formattedErrors := helper.ErrorValidationFormatter(err.(validator.ValidationErrors))
		helper.ResponseAPI(ctx, constant.Res400InvalidPayload, formattedErrors)
		return
	}

	formFile, err := ctx.FormFile("file")
	if err != nil {
		helper.LogError(apiCallID, "Failed to retrieve uploaded file: "+err.Error())
		helper.ResponseAPI(ctx, constant.Res400InvalidPayload)
		return
	}

	uploadedFile, err := formFile.Open()
	if err != nil {
		helper.LogError(apiCallID, "Failed to open uploaded file: "+err.Error())
		helper.ResponseAPI(ctx, constant.Res400InvalidPayload)
		return
	}
	defer uploadedFile.Close()

	if !helper.IsAllowedFileType(apiCallID, uploadedFile) {
		helper.LogError(apiCallID, "Rejected file: unsupported content type")
		helper.ResponseAPI(ctx, constant.Res400InvalidPayload)
		return
	}

	fileBytes, err := io.ReadAll(uploadedFile)
	if err != nil {
		helper.LogError(apiCallID, "Failed to read uploaded file: "+err.Error())
		helper.ResponseAPI(ctx, constant.Res400InvalidPayload)
		return
	}

	result, response := f.Service.UploadFile(apiCallID, req.Folder, formFile.Filename, fileBytes)
	helper.ResponseAPI(ctx, response, result)
}

func (f *FileHandler) DownloadFile(ctx *gin.Context) {
	apiCallID := ctx.GetString(constant.RequestIDKey)

	filePath := ctx.Query("path")
	req := request.DownloadFileRequest{
		Path: filePath,
	}

	if err := f.Validator.Struct(req); err != nil {
		helper.LogError(apiCallID, "Payload validation failed: "+err.Error())
		formattedErrors := helper.ErrorValidationFormatter(err.(validator.ValidationErrors))
		helper.ResponseAPI(ctx, constant.Res400InvalidPayload, formattedErrors)
		return
	}

	fileStream, contentType, result := f.Service.DownloadFile(apiCallID, filePath)
	if result.Code != http.StatusOK {
		helper.LogError(apiCallID, fmt.Sprintf("Failed to download file: %s", result.Message))
		helper.ResponseAPI(ctx, result)
		return
	}

	if fileStream == nil {
		helper.LogError(apiCallID, "Download failed: file stream is nil")
		helper.ResponseAPI(ctx, constant.Res422SomethingWentWrong)
		return
	}

	fileName := filepath.Base(filePath)

	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	ctx.Header("Content-Type", helper.DefaultMIME(contentType))

	helper.LogInfo(apiCallID, "Serving file: "+fileName+" with Content-Type: "+contentType)

	if _, err := io.Copy(ctx.Writer, fileStream); err != nil {
		helper.LogError(apiCallID, "Failed to write file to response: "+err.Error())
		helper.ResponseAPI(ctx, constant.Res422SomethingWentWrong)
		return
	}

	if err := fileStream.Close(); err != nil {
		helper.LogError(apiCallID, "Failed to close file: "+err.Error())
		helper.ResponseAPI(ctx, constant.Res422SomethingWentWrong)
		return
	}
}

func (f *FileHandler) GetAllFile(ctx *gin.Context) {
	apiCallID := ctx.GetString(constant.RequestIDKey)
	result, response := f.Service.GetAllFile(apiCallID)
	helper.ResponseAPI(ctx, response, result)
}

func (f *FileHandler) GetSpecificFile(ctx *gin.Context) {
	apiCallID := ctx.GetString(constant.RequestIDKey)
	folder := ctx.Param("folder")

	result, response := f.Service.GetSpecificFile(apiCallID, folder)
	helper.ResponseAPI(ctx, response, result)
}

func (f *FileHandler) CreateBucket(ctx *gin.Context) {
	apiCallID := ctx.GetString(constant.RequestIDKey)

	var req request.CreateBucketRequest
	if err := ctx.ShouldBind(&req); err != nil {
		helper.LogError(apiCallID, "Failed to bind request: "+err.Error())
		helper.ResponseAPI(ctx, constant.Res400InvalidPayload)
		return
	}

	if err := f.Validator.Struct(req); err != nil {
		helper.LogError(apiCallID, "Payload validation failed: "+err.Error())
		formattedErrors := helper.ErrorValidationFormatter(err.(validator.ValidationErrors))
		helper.ResponseAPI(ctx, constant.Res400InvalidPayload, formattedErrors)
		return
	}

	response := f.Service.CreateBucket(apiCallID, req.BucketName)
	helper.ResponseAPI(ctx, response)
}
