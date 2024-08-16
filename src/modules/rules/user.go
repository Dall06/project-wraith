package rules

import (
	"errors"
	"project-wraith/src/modules/domain"
	"project-wraith/src/pkg/tools"
	"time"
)

type UserRule interface {
	Login(model User) (*User, error)
	Register(model User) (*User, error)
	Edit(model User) error
	Get(model User) (*User, error)
	Remove(model User) error
}

type userRule struct {
	repo      *domain.UserRepository
	shaSecret string
}

func NewRule(repo *domain.UserRepository, shaSecret string) UserRule {
	return &userRule{
		repo:      repo,
		shaSecret: shaSecret,
	}
}

func (r userRule) Login(model User) (*User, error) {
	entity := domain.User{
		ID:       model.ID,
		Username: model.Username,
		Email:    model.Email,
		Name:     model.Name,
		Phone:    model.Phone,
		Password: model.Password,
	}

	response, err := r.repo.Get(entity)
	if err != nil {
		return nil, err
	}

	if response == nil {
		return nil, errors.New("user not found")
	}

	if response.Password != tools.Sha512(r.shaSecret, entity.Password) {
		return nil, errors.New("password incorrect")
	}

	result := &User{
		ID:       response.ID,
		Username: response.Username,
		Email:    response.Email,
		Name:     response.Name,
		Phone:    response.Phone,
		Password: response.Password,
	}

	return result, nil
}

func (r userRule) Register(model User) (*User, error) {
	entity := domain.User{
		ID:        model.ID,
		Username:  model.Username,
		Email:     model.Email,
		Name:      model.Name,
		Phone:     model.Phone,
		Password:  tools.Sha512(r.shaSecret, model.Password),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := r.repo.Create(entity)
	if err != nil {
		return nil, err
	}

	return &model, nil
}

func (r userRule) Edit(model User) error {
	entity := domain.User{
		ID:        model.ID,
		Username:  model.Username,
		Email:     model.Email,
		Name:      model.Name,
		Phone:     model.Phone,
		Password:  tools.Sha512(r.shaSecret, model.Password),
		UpdatedAt: time.Now(),
	}

	err := r.repo.Update(entity)
	if err != nil {
		return err
	}

	return nil
}

func (r userRule) Get(model User) (*User, error) {
	entity := domain.User{
		ID:       model.ID,
		Username: model.Username,
		Email:    model.Email,
		Phone:    model.Phone,
	}

	response, err := r.repo.Get(entity)
	if err != nil {
		return nil, err
	}

	if response == nil {
		return nil, errors.New("user not found")
	}

	result := &User{
		ID:       response.ID,
		Username: response.Username,
		Email:    response.Email,
		Name:     response.Name,
		Phone:    response.Phone,
		Password: response.Password,
	}

	return result, nil
}

func (r userRule) Remove(model User) error {
	entity := domain.User{
		ID:       model.ID,
		Username: model.Username,
		Email:    model.Email,
		Phone:    model.Phone,
		Password: model.Password,
	}

	response, err := r.repo.Get(entity)
	if err != nil {
		return nil
	}

	if response == nil {
		return errors.New("user not found")
	}

	if response.Password != tools.Sha512(r.shaSecret, entity.Password) {
		return errors.New("password incorrect")
	}

	err = r.repo.Delete(entity.ID)
	if err != nil {
		return err
	}

	return nil
}
