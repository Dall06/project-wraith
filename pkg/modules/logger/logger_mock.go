package logger

import "github.com/stretchr/testify/mock"

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Initialize() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockLogger) Info(msg string, args ...interface{}) {
	m.Called(msg)
}

func (m *MockLogger) Warn(msg string, args ...interface{}) {
	m.Called(msg)
}

func (m *MockLogger) Error(msg string, args ...interface{}) {
	m.Called(msg)
}
