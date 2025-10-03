package config

import (
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// InitMinIO initializes MinIO client
func InitMinIO() *minio.Client {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("MINIO_SECRET_ACCESS_KEY")
	useSSL := os.Getenv("MINIO_USE_SSL") == "true"

	if endpoint == "" {
		endpoint = "localhost:9000"
	}

	if accessKeyID == "" || secretAccessKey == "" {
		log.Fatal("MinIO credentials not set in environment variables")
	}

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		log.Fatalf("Failed to initialize MinIO client: %v", err)
	}

	return minioClient
}