package sms

import "github.com/stretchr/testify/mock"

type MockTwilio struct {
	mock.Mock
}

func (m *MockTwilio) SendSMSTwilio(to string, useAsset bool, args ...string) (string, error) {
	argsMock := m.Called(to, useAsset, args)
	return argsMock.String(0), argsMock.Error(1)
}
