package image

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

// Service interface defines the image service operations
type Service interface {
	UploadImage(file *multipart.FileHeader) (string, error)
}

type service struct {
	minioClient *minio.Client
	bucketName  string
}

// NewService creates a new image service
func NewService(minioClient *minio.Client) Service {
	bucketName := os.Getenv("MINIO_BUCKET_NAME")
	if bucketName == "" {
		bucketName = "images"
	}

	s := &service{
		minioClient: minioClient,
		bucketName:  bucketName,
	}

	// Ensure bucket exists
	s.ensureBucketExists()

	return s
}

// ensureBucketExists creates the bucket if it doesn't exist
func (s *service) ensureBucketExists() {
	ctx := context.Background()
	exists, err := s.minioClient.BucketExists(ctx, s.bucketName)
	if err != nil {
		return
	}

	if !exists {
		err = s.minioClient.MakeBucket(ctx, s.bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return
		}

		// Set bucket policy to public read
		policy := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {"AWS": ["*"]},
					"Action": ["s3:GetObject"],
					"Resource": ["arn:aws:s3:::%s/*"]
				}
			]
		}`, s.bucketName)

		s.minioClient.SetBucketPolicy(ctx, s.bucketName, policy)
	}
}

// UploadImage uploads an image to MinIO and returns the URL
func (s *service) UploadImage(fileHeader *multipart.FileHeader) (string, error) {
	// Validate file extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != ".jpg" && ext != ".jpeg" {
		return "", fmt.Errorf("invalid file format: must be jpg or jpeg")
	}

	// Validate file size (min 10KB, max 2MB)
	const minSize = 10 * 1024      // 10KB
	const maxSize = 2 * 1024 * 1024 // 2MB

	if fileHeader.Size < minSize {
		return "", fmt.Errorf("file size too small: minimum 10KB required")
	}

	if fileHeader.Size > maxSize {
		return "", fmt.Errorf("file size too large: maximum 2MB allowed")
	}

	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Generate unique filename using UUID
	objectName := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// Upload to MinIO
	ctx := context.Background()
	_, err = s.minioClient.PutObject(
		ctx,
		s.bucketName,
		objectName,
		file,
		fileHeader.Size,
		minio.PutObjectOptions{
			ContentType: fileHeader.Header.Get("Content-Type"),
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	// Construct the file URL
	// Use public endpoint if set, otherwise use MINIO_ENDPOINT
	publicEndpoint := os.Getenv("MINIO_PUBLIC_ENDPOINT")
	if publicEndpoint == "" {
		publicEndpoint = os.Getenv("MINIO_ENDPOINT")
	}

	useSSL := os.Getenv("MINIO_USE_SSL") == "true"
	protocol := "http"
	if useSSL {
		protocol = "https"
	}

	imageURL := fmt.Sprintf("%s://%s/%s/%s", protocol, publicEndpoint, s.bucketName, objectName)

	return imageURL, nil
}

// Helper function to read file content (if needed for validation)
func readFileContent(file multipart.File) ([]byte, error) {
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return content, nil
}