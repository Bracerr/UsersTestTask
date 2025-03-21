package postgres

import (
	"database/sql"
	"testing"
	"time"
	"users-api/src/internal/domain"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	t.Run("successful creation", func(t *testing.T) {
		user := &domain.User{
			Name:  "John Doe",
			Email: "john@example.com",
		}

		mock.ExpectQuery("INSERT INTO users").
			WithArgs(user.Name, user.Email, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		err := repo.Create(user)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), user.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	t.Run("user exists", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "email", "created_at", "updated_at"}).
			AddRow(1, "John Doe", "john@example.com", time.Now(), time.Now())

		mock.ExpectQuery("SELECT (.+) FROM users").
			WithArgs(1).
			WillReturnRows(rows)

		user, err := repo.GetByID(1)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "John Doe", user.Name)
		assert.Equal(t, "john@example.com", user.Email)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("user not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM users").
			WithArgs(999).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetByID(999)
		assert.NoError(t, err)
		assert.Nil(t, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	t.Run("successful update", func(t *testing.T) {
		user := &domain.User{
			ID:    1,
			Name:  "John Doe Updated",
			Email: "john.updated@example.com",
		}

		mock.ExpectExec("UPDATE users").
			WithArgs(user.Name, user.Email, sqlmock.AnyArg(), user.ID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Update(user)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("user not found", func(t *testing.T) {
		user := &domain.User{
			ID:    999,
			Name:  "John Doe",
			Email: "john@example.com",
		}

		mock.ExpectExec("UPDATE users").
			WithArgs(user.Name, user.Email, sqlmock.AnyArg(), user.ID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Update(user)
		assert.Equal(t, sql.ErrNoRows, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDeleteUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	t.Run("successful deletion", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM users").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete(1)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("user not found", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM users").
			WithArgs(999).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Delete(999)
		assert.Equal(t, sql.ErrNoRows, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
