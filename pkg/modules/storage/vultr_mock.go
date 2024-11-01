package storage

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/mock"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) UploadObject(bucket, directory, filename, permission, localPath string) (*s3.PutObjectOutput, error) {
	args := m.Called(bucket, directory, filename, permission, localPath)
	return args.Get(0).(*s3.PutObjectOutput), args.Error(1)
}
