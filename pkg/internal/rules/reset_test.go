package rules_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"project-wraith/pkg/internal/domain"
	"project-wraith/pkg/internal/rules"
	"project-wraith/pkg/modules/tools"
	"testing"
)

func TestResetRule(test *testing.T) {
	test.Parallel()

	testCases := []struct {
		name           string
		input          rules.Reset
		repoReturn     *domain.User
		repoErr        error
		encryptData    bool
		expectedResult *rules.Reset
		method         string
	}{
		{
			name: "Start Success",
			input: rules.Reset{
				ID:    "123",
				Email: "ZwUeh@example.com",
			},
			repoReturn: &domain.User{
				ID:       "123",
				Password: tools.Sha512("secret", "password"),
			},
			repoErr:     nil,
			encryptData: false,
			expectedResult: &rules.Reset{
				ID:       "123",
				Email:    "ZwUeh@example.com",
				Phone:    "",
				Username: "",
				Token:    "",
			},
			method: "Start",
		},
		{
			name: "Validate Success",
			input: rules.Reset{
				ID:          "123",
				Email:       "ZwUeh@example.com",
				NewPassword: "secret_password",
			},
			repoReturn: &domain.User{
				ID:       "123",
				Password: tools.Sha512("secret", "password"),
			},
			repoErr:     nil,
			encryptData: true,
			expectedResult: &rules.Reset{
				ID:       "123",
				Email:    "ZwUeh@example.com",
				Phone:    "",
				Username: "",
				Token:    "",
			},
			method: "Validate",
		},
	}

	for _, tc := range testCases {
		test.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockRepo := new(domain.MockUserRepository)
			rule := rules.NewResetRule(mockRepo, "secret")

			// Set up mock behavior for repository based on method
			switch tc.method {
			case "Start":
				mockRepo.On("Get", mock.Anything).Return(tc.repoReturn, tc.repoErr)
			case "Validate":
				mockRepo.On("Get", mock.Anything).Return(tc.repoReturn, tc.repoErr)
			}

			// Execute the method under test
			var result *rules.Reset
			var err error

			switch tc.method {
			case "Start":
				result, err = rule.Start(tc.input)
				// Check that the result is as expected
				assert.NotNil(t, result)
				assert.Equal(t, tc.expectedResult.ID, result.ID)
				assert.NoError(t, err)

			case "Validate":
				// Call Start to get a token, then validate
				result, err = rule.Start(tc.input)
				assert.NoError(t, err)

				tc.input.Token = result.Token
				result, err = rule.Validate(tc.input)
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.Token)
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
