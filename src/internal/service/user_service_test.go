package service

import (
	"testing"
	"users-api/src/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(id int64) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	t.Run("valid user", func(t *testing.T) {
		user := &domain.User{
			Name:  "John Doe",
			Email: "john@example.com",
		}

		mockRepo.On("Create", user).Return(nil)

		err := service.CreateUser(user)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid user", func(t *testing.T) {
		user := &domain.User{
			Name:  "", // пустое имя
			Email: "john@example.com",
		}

		err := service.CreateUser(user)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidInput, err)
	})
}

func TestGetUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	t.Run("existing user", func(t *testing.T) {
		expectedUser := &domain.User{
			ID:    1,
			Name:  "John Doe",
			Email: "john@example.com",
		}

		mockRepo.On("GetByID", int64(1)).Return(expectedUser, nil)

		user, err := service.GetUser(1)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("non-existing user", func(t *testing.T) {
		mockRepo.On("GetByID", int64(999)).Return(nil, nil)

		user, err := service.GetUser(999)
		assert.Error(t, err)
		assert.Equal(t, ErrUserNotFound, err)
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	t.Run("valid update", func(t *testing.T) {
		user := &domain.User{
			ID:    1,
			Name:  "John Doe Updated",
			Email: "john.updated@example.com",
		}

		mockRepo.On("Update", user).Return(nil)

		err := service.UpdateUser(user)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid user data", func(t *testing.T) {
		user := &domain.User{
			ID:    1,
			Name:  "", // пустое имя
			Email: "john@example.com",
		}

		err := service.UpdateUser(user)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidInput, err)
	})
}
