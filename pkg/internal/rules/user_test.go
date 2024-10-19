package rules_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"project-wraith/pkg/internal/domain"
	"project-wraith/pkg/internal/rules"
	"project-wraith/pkg/modules/tools"
	"testing"
)

func TestUserRule(test *testing.T) {
	test.Parallel()

	testCases := []struct {
		name           string
		input          rules.User
		repoReturn     *domain.User
		repoErr        error
		encryptData    bool
		expectedResult *rules.User
		expectedError  error
		method         string
	}{
		{
			name: "Login Success",
			input: rules.User{
				ID:       "123",
				Password: "password",
			},
			repoReturn: &domain.User{
				ID:       "123",
				Password: tools.Sha512("secret", "password"),
			},
			repoErr:     nil,
			encryptData: false,
			expectedResult: &rules.User{
				ID:       "123",
				Password: tools.Sha512("secret", "password"),
			},
			expectedError: nil,
			method:        "Login",
		},

		{
			name: "Register Success",
			input: rules.User{
				Username: "newuser",
				Password: "password",
			},
			repoReturn:  nil,
			repoErr:     nil,
			encryptData: true,
			expectedResult: &rules.User{
				Username: "newuser",
				Password: tools.Sha512("secret", "password"),
			},
			expectedError: nil,
			method:        "Register",
		},

		{
			name: "Edit Success",
			input: rules.User{
				ID:       "123",
				Username: "updateduser",
			},
			repoReturn:     nil,
			repoErr:        nil,
			encryptData:    true,
			expectedResult: nil,
			expectedError:  nil,
			method:         "Edit",
		},

		{
			name: "Get Success",
			input: rules.User{
				ID: "123",
			},
			repoReturn: &domain.User{
				ID: "123",
			},
			repoErr:     nil,
			encryptData: false,
			expectedResult: &rules.User{
				ID: "123",
			},
			expectedError: nil,
			method:        "Get",
		},

		{
			name: "Disable Success",
			input: rules.User{
				ID:       "123",
				Password: "password",
			},
			repoReturn: &domain.User{
				ID:       "123",
				Password: tools.Sha512("secret", "password"),
			},
			repoErr:        nil,
			encryptData:    false,
			expectedResult: nil,
			expectedError:  nil,
			method:         "Disable",
		},
	}

	for _, tc := range testCases {
		test.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := new(domain.MockUserRepository)
			rule := rules.NewUserRule(mockRepo, tc.encryptData, "", "secret")

			// Set up mock behavior
			switch tc.method {
			case "Login":
				mockRepo.On("Update", mock.Anything).Return(tc.repoErr)
				mockRepo.On("Get", mock.Anything).Return(tc.repoReturn, tc.repoErr)
			case "Register":
				mockRepo.On("Duplicated", mock.Anything).Return([]domain.User{}, tc.repoErr)
				mockRepo.On("Create", mock.Anything).Return(tc.repoErr)
			case "Edit":
				mockRepo.On("Update", mock.Anything).Return(tc.repoErr)
			case "Get":
				mockRepo.On("Get", mock.Anything).Return(tc.repoReturn, tc.repoErr)
			case "Disable":
				mockRepo.On("Get", mock.Anything).Return(tc.repoReturn, tc.repoErr)
				mockRepo.On("Update", mock.Anything).Return(tc.repoErr)
			}

			// Run the method under test
			var result *rules.User
			var err error

			switch tc.method {
			case "Login":
				result, err = rule.Login(tc.input)
				// Move assertions outside the switch block
				assert.Equal(t, tc.expectedResult, result)
				assert.Equal(t, tc.expectedError, err)
			case "Register":
				result, err = rule.Register(tc.input)
				// Move assertions outside the switch block
				assert.Equal(t, tc.expectedResult.ID, result.ID)
				assert.Equal(t, tc.expectedError, err)
			case "Edit":
				err = rule.Edit(tc.input)
				// Move assertions outside the switch block
				assert.Equal(t, tc.expectedResult, result)
				assert.Equal(t, tc.expectedError, err)
			case "Get":
				result, err = rule.Get(tc.input)
				// Move assertions outside the switch block
				assert.Equal(t, tc.expectedResult, result)
				assert.Equal(t, tc.expectedError, err)
			case "Disable":
				err = rule.Disable(tc.input)
				// Move assertions outside the switch block
				assert.Equal(t, tc.expectedResult, result)
				assert.Equal(t, tc.expectedError, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
