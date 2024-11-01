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

type Storage interface {
	UploadObject(bucket, directory, filename, permission, localPath string) (*s3.PutObjectOutput, error)
}

// ObjectStorage struct for interacting with S3 or compatible storage
type storageClient struct {
	client *s3.S3
}

func NewObjectStorage(accessKey, secretKey string) Storage {
	str := s3.New(session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("ewr"),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:    aws.String("https://ewr1.vultrobjects.com/"),
	})))

	return &storageClient{client: str}
}

func (s *storageClient) UploadObject(bucket, directory, filename, permission, localPath string) (*s3.PutObjectOutput, error) {
	file, err := os.Open(localPath)
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

	output, err := s.client.PutObject(&s3.PutObjectInput{
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
