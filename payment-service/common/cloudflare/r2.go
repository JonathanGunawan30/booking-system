package cloudflare

import (
	"context"
	"io"
	"payment-service/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type R2Client struct {
	client *s3.Client
	bucket string
}

func NewR2Client(client *s3.Client) *R2Client {
	return &R2Client{
		client: client,
		bucket: config.AppConfig.R2BucketName,
	}
}

func (r *R2Client) Upload(key string, body io.Reader, contentType string) (string, error) {
	_, err := r.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(r.bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
	})

	if err != nil {
		return "", err
	}

	publicURL := config.AppConfig.R2PublicURL + "/" + key
	return publicURL, nil
}

func (r *R2Client) Delete(key string) error {
	_, err := r.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})
	return err
}
