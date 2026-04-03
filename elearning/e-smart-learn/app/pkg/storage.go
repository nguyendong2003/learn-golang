package pkg

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"elearning-api/config"
	"elearning-api/util"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type StorageProvider interface {
	UploadImage(ctx context.Context, fileName string, data []byte) (string, error)
	UploadVideo(ctx context.Context, fileName string, data []byte) (string, error)
	Delete(ctx context.Context, fileURL string) error
	PresignUploadURL(ctx context.Context, filename, filetype string) (string, error)

	IsObjectExist(ctx context.Context, objectName string) (bool, error)
	GetConfig() *config.MinioConfig
}

type storageProvider struct {
	minioClient *minio.Client
	config      *config.MinioConfig
}

func NewStorageProvider(config *config.MinioConfig) StorageProvider {
	minioClient, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: config.UseSSL,
	})
	logger := util.WithLayer(util.LayerInfra)
	if err != nil {
		logger.Error("failed to create minio client", "error", err)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	exists, err := minioClient.BucketExists(ctx, config.BucketName)
	if err != nil {
		logger.Error("failed to check bucket existence", "error", err)
		return nil
	}

	if !exists {
		err = minioClient.MakeBucket(ctx, config.BucketName, minio.MakeBucketOptions{Region: "us-east-1"})
		if err != nil {
			logger.Error("failed to create bucket", "bucket", config.BucketName, "error", err)
			return nil
		}
	}
	policy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Action": ["s3:GetObject"],
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Resource": ["arn:aws:s3:::%s/*"],
				"Sid": "PublicRead"
			}
		]
	}`, config.BucketName)

	err = minioClient.SetBucketPolicy(ctx, config.BucketName, policy)
	if err != nil {
		logger.Warn("failed to set public bucket policy", "bucket", config.BucketName, "error", err)
	}

	return &storageProvider{
		minioClient: minioClient,
		config:      config,
	}
}

func (s *storageProvider) UploadImage(ctx context.Context, fileName string, data []byte) (string, error) {
	return s.upload(ctx, "images", fileName, data)
}

func (s *storageProvider) UploadVideo(ctx context.Context, fileName string, data []byte) (string, error) {
	return s.upload(ctx, "videos", fileName, data)
}

func (s *storageProvider) Delete(ctx context.Context, fileURL string) error {
	if fileURL == "" {
		return nil
	}

	parts := strings.Split(fileURL, s.config.BucketName+"/")
	if len(parts) < 2 {
		return fmt.Errorf("invalid file URL")
	}
	objectName := parts[1]

	return s.minioClient.RemoveObject(ctx, s.config.BucketName, objectName, minio.RemoveObjectOptions{})
}

func (s *storageProvider) upload(ctx context.Context, folder, fileName string, data []byte) (string, error) {
	uniqueName := fmt.Sprintf("%s/%d_%s", folder, time.Now().UnixNano(), fileName)

	contentType := http.DetectContentType(data)

	_, err := s.minioClient.PutObject(ctx, s.config.BucketName, uniqueName,
		bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{
			ContentType: contentType,
		})
	if err != nil {
		logger := util.WithLayer(util.LayerInfra)
		logger.Error("failed to upload file", "folder", folder, "file_name", fileName, "error", err)
		return "", err
	}
	return s.getURL(uniqueName), nil
}

func (s *storageProvider) getURL(objectName string) string {
	schema := "http://"
	if s.config.UseSSL {
		schema = "https://"
	}
	return fmt.Sprintf("%s%s/%s/%s", schema, s.config.ExternalEndpoint, s.config.BucketName, objectName)
}

func (s *storageProvider) PresignUploadURL(ctx context.Context, filename, filetype string) (string, error) {
	objectName := fmt.Sprintf("%s/%d_%s", filetype, time.Now().UnixNano(), filename)
	u, err := s.minioClient.PresignedPutObject(ctx, s.config.BucketName, objectName, time.Minute*15)
	if err != nil {
		return "", err
	}
	urlStr := u.String()
	urlStr = strings.ReplaceAll(urlStr, s.config.Endpoint, s.config.ExternalEndpoint)
	return urlStr, nil
}

// Check
func (s *storageProvider) GetConfig() *config.MinioConfig {
	return s.config
}

func (s *storageProvider) IsObjectExist(ctx context.Context, objectName string) (bool, error) {
	_, err := s.minioClient.StatObject(ctx, s.config.BucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		resp := minio.ToErrorResponse(err)
		if resp.Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
