package service

import (
	"context"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"strings"

	"elearning-api/apperror"
	"elearning-api/pkg"
)

type UploadService interface {
	UploadImage(ctx context.Context, file *multipart.FileHeader) (string, error)
	PresignUploadURL(ctx context.Context, filename, filetype string) (string, error)
	DeleteImage(ctx context.Context, fileURL string) error

	ValidateVideoURL(ctx context.Context, fileURL string) (bool, error)
	ValidateDocumentURL(ctx context.Context, fileURL string) (bool, error)
}

type uploadService struct {
	storageProvider pkg.StorageProvider
}

func NewUploadService(storageProvider pkg.StorageProvider) UploadService {
	return &uploadService{
		storageProvider: storageProvider,
	}
}

func (u *uploadService) UploadImage(ctx context.Context, file *multipart.FileHeader) (string, error) {
	if file.Size > 5*1024*1024 {
		return "", apperror.NewBadRequestError("Image size must be less than 5MB")
	}

	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[ext] {
		return "", apperror.NewBadRequestError("Invalid image format. Allowed: jpg, jpeg, png, gif, webp")
	}

	openFile, err := file.Open()
	if err != nil {
		return "", apperror.NewBadRequestError("Failed to open file")
	}
	defer openFile.Close()

	fileContent := make([]byte, file.Size)
	if _, err := openFile.Read(fileContent); err != nil {
		return "", apperror.NewInternalServerError("Failed to read file")
	}

	url, err := u.storageProvider.UploadImage(ctx, file.Filename, fileContent)
	if err != nil {
		return "", apperror.NewInternalServerError("Failed to upload image")
	}
	return url, nil
}

func (u *uploadService) PresignUploadURL(ctx context.Context, filename, filetype string) (string, error) {
	if filetype != "images" && filetype != "videos" && filetype != "documents" {
		return "", apperror.NewBadRequestError("Invalid file type. Allowed: images, videos, documents")
	}

	url, err := u.storageProvider.PresignUploadURL(ctx, filename, filetype)
	if err != nil {
		return "", apperror.NewInternalServerError("Failed to generate presigned URL")
	}
	return url, nil
}

func (u *uploadService) DeleteImage(ctx context.Context, fileURL string) error {
	if fileURL == "" {
		return nil
	}

	if err := u.storageProvider.Delete(ctx, fileURL); err != nil {
		return apperror.NewInternalServerError("Failed to delete old image")
	}
	return nil
}

func (u *uploadService) ValidateVideoURL(ctx context.Context, fileURL string) (bool, error) {
	return u.validateStorageURL(ctx, fileURL, "videos")
}

func (u *uploadService) ValidateDocumentURL(ctx context.Context, fileURL string) (bool, error) {
	return u.validateStorageURL(ctx, fileURL, "documents")
}

func (u *uploadService) validateStorageURL(ctx context.Context, fileURL string, folder string) (bool, error) {
	if fileURL == "" {
		return false, nil
	}

	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return false, apperror.NewBadRequestError("Invalid URL")
	}

	cfg := u.storageProvider.GetConfig()

	// normalize host
	expectedHost := strings.TrimPrefix(cfg.ExternalEndpoint, "http://")
	expectedHost = strings.TrimPrefix(expectedHost, "https://")

	// check host
	if parsedURL.Host != expectedHost {
		return false, apperror.NewBadRequestError("Invalid host")
	}

	// extract objectName
	objectName, err := extractObjectName(fileURL, cfg.BucketName)
	if err != nil {
		return false, err
	}

	// check folder
	if !strings.HasPrefix(objectName, folder+"/") {
		return false, apperror.NewBadRequestError("Invalid file type")
	}

	// check file exist in minio
	exists, err := u.storageProvider.IsObjectExist(ctx, objectName)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func extractObjectName(fileURL string, bucketName string) (string, error) {
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return "", apperror.NewBadRequestError("Invalid URL")
	}

	// /elearning/videos/xxx.mp4
	path := strings.TrimPrefix(parsedURL.Path, "/")

	parts := strings.SplitN(path, "/", 2)
	if len(parts) < 2 {
		return "", apperror.NewBadRequestError("Invalid path")
	}

	bucket := parts[0]
	objectName := parts[1]

	if bucket != bucketName {
		return "", apperror.NewBadRequestError("Invalid bucket")
	}

	return objectName, nil
}
