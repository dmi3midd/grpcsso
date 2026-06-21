package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dmi3midd/grpcsso/internal/domain"
	"github.com/jmoiron/sqlx"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository interface {
	GetById(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	CreateUser(ctx context.Context, user domain.User) (*domain.User, error)
	UpdateUser(ctx context.Context, user domain.User) (*domain.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetById(ctx context.Context, id string) (*domain.User, error) {
	op := "UserRepository.GetById"
	query := `
	SELECT id, username, email, hashed_password, created_at, updated_at 
	FROM users 
	WHERE id = $1
	`
	var user domain.User
	if err := r.db.GetContext(ctx, &user, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	op := "UserRepository.GetByEmail"
	query := `
	SELECT id, username, email, hashed_password, created_at, updated_at 
	FROM users 
	WHERE email = $1
	`
	var user domain.User
	if err := r.db.GetContext(ctx, &user, query, email); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user domain.User) (*domain.User, error) {
	op := "UserRepository.CreateUser"
	query := `
	INSERT INTO users (id, username, email, hashed_password, created_at, updated_at)
	VALUES (:id, :username, :email, :hashed_password, :created_at, :updated_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user domain.User) (*domain.User, error) {
	op := "UserRepository.UpdateUser"
	query := `
	UPDATE users 
	SET username = :username, email = :email, hashed_password = :hashed_password, updated_at = :updated_at
	WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &user, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id string) error {
	op := "UserRepository.DeleteUser"
	query := `
	DELETE FROM users 
	WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
