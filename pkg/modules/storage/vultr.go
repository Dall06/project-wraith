package storage

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gabriel-vasile/mimetype"
	"log"
	"os"
	"path"
)

func NewObjectStorage(accessKey, secretKey string) *s3.S3 {
	storage := s3.New(session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("ewr"),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:    aws.String("https://ewr1.vultrobjects.com/"),
	})))

	return storage
}

func UploadObject(storage *s3.S3, bucket, directory, filename, permission string) (*s3.PutObjectOutput, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	mimeType, err := mimetype.DetectFile(filename)
	if err != nil {
		err := file.Close()
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed to detect content type: %w", err)
	}
	contentType := mimeType.String()

	objectKey := path.Join(directory, path.Base(filename))
	log.Printf("Uploading Object: %s", filename)

	output, err := storage.PutObject(&s3.PutObjectInput{
		Body:        file,
		Bucket:      aws.String(bucket),
		Key:         aws.String(objectKey),
		ACL:         aws.String(permission),
		ContentType: aws.String(contentType),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to upload file %q: %w", filename, err)
	}

	err = file.Close()
	if err != nil {
		return nil, err
	}

	log.Printf("File %q uploaded successfully.", filename)
	return output, nil
}
