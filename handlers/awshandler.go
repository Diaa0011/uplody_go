package handlers

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

func uploadToS3(file []byte, filename string) error {
	err := godotenv.Load(".env")
	if err != nil {
		// If it fails, print the actual directory so we can see why
		currDir, _ := os.Getwd()
		return fmt.Errorf("could not find .env file. Current directory: %s", currDir)
	}

	// 2. Explicitly read the keys from the environment
	appKey := os.Getenv("AWS_ACCESS_KEY_ID")
	appSecret := os.Getenv("AWS_SECRET_ACCESS_KEY")
	appRegion := os.Getenv("AWS_REGION")
	bucketName := os.Getenv("AWS_BUCKET_NAME")

	// 3. Load configuration using the keys we just grabbed
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(appRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(appKey, appSecret, "")),
	)
	if err != nil {
		return fmt.Errorf("unable to load SDK config, %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)

	// 4. Upload the file
	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(file),
		// ACL: aws.String("public-read"), // Remember our talk about public access!
	})

	if err != nil {
		return fmt.Errorf("failed to upload file to S3, %v", err)
	}

	return nil
}
