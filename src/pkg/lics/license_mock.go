package lics

import "github.com/stretchr/testify/mock"

type MockLicenseRepository struct {
	mock.Mock
}

func (m *MockLicenseRepository) Get(lic License) (*License, error) {
	args := m.Called(lic)
	return args.Get(0).(*License), args.Error(1)
}
