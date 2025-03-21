package service

import (
	"users-api/src/internal/domain"
	"users-api/src/internal/errors"
)

var (
	ErrUserNotFound = errors.ErrUserNotFound
	ErrInvalidInput = errors.ErrInvalidInput
)

type UserService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(user *domain.User) error {
	if user.Name == "" || user.Email == "" {
		return errors.ErrInvalidInput
	}
	return s.repo.Create(user)
}

func (s *UserService) GetUser(id int64) (*domain.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) UpdateUser(user *domain.User) error {
	if user.ID == 0 {
		return errors.ErrInvalidInput
	}

	currentUser, err := s.repo.GetByID(user.ID)
	if err != nil {
		return err
	}
	if currentUser == nil {
		return errors.ErrUserNotFound
	}

	if user.Name != "" {
		currentUser.Name = user.Name
	}
	if user.Email != "" {
		currentUser.Email = user.Email
	}

	err = s.repo.Update(currentUser)
	if err != nil {
		return err
	}

	*user = *currentUser
	return nil
}

func (s *UserService) DeleteUser(id int64) error {
	if id == 0 {
		return errors.ErrInvalidInput
	}
	return s.repo.Delete(id)
}
