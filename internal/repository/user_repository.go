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
	ErrUserNotFound error = errors.New("user not found")
)

type UserRepository interface {
	// UserExists checks if a user with the given Id exists.
	IsExists(ctx context.Context, ext sqlx.ExtContext, userId string) (bool, error)
	// GetById retrieves a User entity by its id.
	// It returns [ErrUserNotFound] if no user are found.
	GetById(ctx context.Context, ext sqlx.ExtContext, userId string) (*domain.User, error)
	// GetByEmail retrieves a User entity by its email.
	// It returns [ErrUserNotFound] if no user are found.
	GetByEmail(ctx context.Context, ext sqlx.ExtContext, email string) (*domain.User, error)
	// Create creates a User entity and returns it.
	Create(ctx context.Context, ext sqlx.ExtContext, user *domain.User) (string, error)
	// Update updates the User entity.
	Update(ctx context.Context, ext sqlx.ExtContext, user *domain.User) (string, error)
	// Delete removes the User entity.
	Delete(ctx context.Context, ext sqlx.ExtContext, userId string) error
}

type userRepository struct {
}

func NewUserRepo() UserRepository {
	return &userRepository{}
}

func (r *userRepository) IsExists(ctx context.Context, ext sqlx.ExtContext, userId string) (bool, error) {
	op := "UserRepository.IsExists"
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
	var exists bool
	err := sqlx.GetContext(ctx, ext, &exists, query, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return exists, nil
}

func (r *userRepository) GetById(ctx context.Context, ext sqlx.ExtContext, userId string) (*domain.User, error) {
	op := "UserRepository.GetById"
	query := `
	SELECT id, username, email, password_hash, created_at, updated_at
	FROM users WHERE id = $1
	`
	var user domain.User
	err := sqlx.GetContext(ctx, ext, &user, query, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, ext sqlx.ExtContext, email string) (*domain.User, error) {
	op := "UserRepository.GetByEmail"
	query := `
	SELECT id, username, email, password_hash, created_at, updated_at
	FROM users WHERE email = $1
	`
	var user domain.User
	err := sqlx.GetContext(ctx, ext, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, ext sqlx.ExtContext, user *domain.User) (string, error) {
	op := "UserRepository.Create"
	query := `
	INSERT INTO users (id, username, email, password_hash, created_at, updated_at)
	VALUES (:id, :username, :email, :password_hash, :created_at, :updated_at)
	`
	if _, err := sqlx.NamedExecContext(ctx, ext, query, user); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return user.Id, nil
}

func (r *userRepository) Update(ctx context.Context, ext sqlx.ExtContext, user *domain.User) (string, error) {
	op := "UserRepository.Update"
	query := `
	UPDATE users
	SET username = :username, email = :email, password_hash = :password_hash, updated_at = :updated_at
	WHERE id = :id
	`
	_, err := sqlx.NamedExecContext(ctx, ext, query, user)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return user.Id, nil
}

func (r *userRepository) Delete(ctx context.Context, ext sqlx.ExtContext, userId string) error {
	op := "UserRepository.Delete"
	query := `
	DELETE FROM users
	WHERE id = $1
	`
	if _, err := ext.ExecContext(ctx, query, userId); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
