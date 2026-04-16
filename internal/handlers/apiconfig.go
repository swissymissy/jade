package handlers

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/swissymissy/jade/internal/database"
)

type ApiConfig struct {
	Port      string
	Platform  string
	DB        *database.Queries
	JWTSecret string
	S3Bucket  string
	S3Region  string
	S3Client  *s3.Client
}
