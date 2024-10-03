package mail

import (
	"github.com/stretchr/testify/mock"
)

type MockMail struct {
	mock.Mock
}

func (m *MockMail) Send(tmpl string, content interface{}, subject string, to []string) error {
	args := m.Called(tmpl, content, subject, to)
	return args.Error(0)
}
