package rules

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"project-wraith/pkg/internal/domain"
	"project-wraith/pkg/modules/token"
	"project-wraith/pkg/modules/tools"
	"time"
)

type ResetRule interface {
	Start(reset Reset) (*Reset, error)
	Validate(reset Reset) (*Reset, error)
}

type resetRule struct {
	repo      domain.UserRepository
	jwtSecret string
}

func NewResetRule(repository domain.UserRepository, jwtSecret string) ResetRule {
	return &resetRule{
		repo:      repository,
		jwtSecret: jwtSecret,
	}
}

func (rr resetRule) Start(reset Reset) (*Reset, error) {
	entity := domain.User{
		Email: reset.Email,
		Phone: reset.Phone,
	}

	response, err := rr.repo.Get(entity)
	if err != nil {
		return nil, err
	}

	if response == nil {
		return nil, errors.New("user not found")
	}

	claims := map[string]interface{}{
		"id":    response.ID,
		"email": response.Email,
		"phone": response.Phone,
		"reset": true,
	}

	tkn, err := token.CreateJwtToken(rr.jwtSecret, 10*time.Minute, claims)
	if err != nil {
		return nil, err
	}

	result := &Reset{
		ID:       response.ID,
		Username: response.Username,
		Email:    response.Email,
		Phone:    response.Phone,
		Token:    tkn,
	}

	return result, nil
}

func (rr resetRule) Validate(reset Reset) (*Reset, error) {
	res := &Reset{}

	extraValidation := func(claims jwt.MapClaims) error {
		data := claims["data"].(map[string]interface{})

		if !data["reset"].(bool) {
			return errors.New("invalid token")
		}

		entity := domain.User{
			Email: data["email"].(string),
			Phone: data["phone"].(string),
		}

		response, err := rr.repo.Get(entity)
		if err != nil {
			return err
		}

		if response == nil {
			return errors.New("user not found")
		}

		if tools.Sha512(rr.jwtSecret, reset.NewPassword) == response.Password {
			return errors.New("new password cannot be the same as old password")
		}

		res = &Reset{
			ID:    response.ID,
			Token: reset.Token,
		}

		return nil
	}

	valid, err := token.ValidateJwtToken(reset.Token, rr.jwtSecret, extraValidation)
	if err != nil {
		return res, err
	}

	if !valid {
		return res, errors.New("invalid token")
	}

	return res, nil
}
