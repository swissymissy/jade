package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// upload to s3 bucket
func UploadToS3(ctx context.Contextclient *s3.Client, bucket, key, contentType string, body io.Reader) error {
	// put object in bucket
	_, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key: aws.String(key),
		Body: body,
		contentType: aws.String(contentType),
	})
	if err != nil {
		return fmt.Errorf("Failed to upload to S3: %w", err)
	}
	return nil
}

// delete from s3 bucket
func DeleteFromS3(ctx context.Context, client *s3.Client, bucket, key string ) error {
	// delete object from bucket
	_, err := client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key: aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete from S3: %w", err)
	}
	return nil
}