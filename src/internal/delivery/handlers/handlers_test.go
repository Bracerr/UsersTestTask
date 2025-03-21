package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"users-api/src/internal/domain"
	"users-api/src/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) GetUser(id int64) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) DeleteUser(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateUser(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	t.Run("valid user", func(t *testing.T) {
		user := &domain.User{
			Name:  "John Doe",
			Email: "john@example.com",
		}

		mockService.On("CreateUser", user).Return(nil)

		body, _ := json.Marshal(user)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.CreateUser(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid request body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString("invalid json"))
		w := httptest.NewRecorder()

		handler.CreateUser(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetUser(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	t.Run("existing user", func(t *testing.T) {
		user := &domain.User{
			ID:    1,
			Name:  "John Doe",
			Email: "john@example.com",
		}

		mockService.On("GetUser", int64(1)).Return(user, nil)

		req := httptest.NewRequest(http.MethodGet, "/users?id=1", nil)
		w := httptest.NewRecorder()

		handler.GetUser(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)

		var response domain.User
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, user.ID, response.ID)
		assert.Equal(t, user.Name, response.Name)
		assert.Equal(t, user.Email, response.Email)
	})

	t.Run("invalid user id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users?id=invalid", nil)
		w := httptest.NewRecorder()

		handler.GetUser(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("user not found", func(t *testing.T) {
		mockService.On("GetUser", int64(999)).Return(nil, service.ErrUserNotFound)

		req := httptest.NewRequest(http.MethodGet, "/users?id=999", nil)
		w := httptest.NewRecorder()

		handler.GetUser(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUpdateUser(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	t.Run("valid update", func(t *testing.T) {
		user := &domain.User{
			ID:    1,
			Name:  "John Doe Updated",
			Email: "john.updated@example.com",
		}

		mockService.On("UpdateUser", user).Return(nil)

		body, _ := json.Marshal(user)
		req := httptest.NewRequest(http.MethodPut, "/users", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.UpdateUser(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid request body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/users", bytes.NewBufferString("invalid json"))
		w := httptest.NewRecorder()

		handler.UpdateUser(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("user not found", func(t *testing.T) {
		user := &domain.User{
			ID:    999,
			Name:  "John Doe",
			Email: "john@example.com",
		}

		mockService.On("UpdateUser", user).Return(service.ErrUserNotFound)

		body, _ := json.Marshal(user)
		req := httptest.NewRequest(http.MethodPut, "/users", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.UpdateUser(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}
