package postgres

import (
	"database/sql"
	"time"

	"users-api/src/internal/domain"

	"github.com/Masterminds/squirrel"
)

type UserRepository struct {
	db      *sql.DB
	builder squirrel.StatementBuilderType
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db:      db,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *UserRepository) Create(user *domain.User) error {
	now := time.Now().Format(time.RFC3339)
	user.CreatedAt = now
	user.UpdatedAt = now

	query := r.builder.
		Insert("users").
		Columns("name", "email", "created_at", "updated_at").
		Values(user.Name, user.Email, user.CreatedAt, user.UpdatedAt).
		Suffix("RETURNING id")

	err := query.RunWith(r.db).QueryRow().Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetByID(id int64) (*domain.User, error) {
	user := &domain.User{}

	query := r.builder.
		Select("id", "name", "email", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{"id": id})

	err := query.RunWith(r.db).QueryRow().Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Update(user *domain.User) error {
	user.UpdatedAt = time.Now().Format(time.RFC3339)

	query := r.builder.
		Update("users").
		Set("name", user.Name).
		Set("email", user.Email).
		Set("updated_at", user.UpdatedAt).
		Where(squirrel.Eq{"id": user.ID})

	result, err := query.RunWith(r.db).Exec()
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *UserRepository) Delete(id int64) error {
	query := r.builder.
		Delete("users").
		Where(squirrel.Eq{"id": id})

	result, err := query.RunWith(r.db).Exec()
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
