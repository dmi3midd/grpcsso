package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/dmi3midd/grpcsso/internal/domain"

	"github.com/jmoiron/sqlx"
)

var (
	ErrTokenNotFound error = errors.New("token not found")
)

type TokenRepository interface {
	// Get retrieves a Token entity by its id.
	// It returns [ErrTokenNotFound] if no token are found.
	GetById(ctx context.Context, id string) (*domain.Token, error)
	// Get retrieves a Token entity by its refresh token.
	// It returns [ErrTokenNotFound] if no token are found.
	GetByToken(ctx context.Context, refreshToken string) (*domain.Token, error)
	// Create creates a Token entity.
	Create(ctx context.Context, token *domain.Token) (string, error)
	// Update updates refresh token in the Token entity.
	Update(ctx context.Context, id string, token *domain.Token) error
	// DeleteById removes the Token entity by its id.
	DeleteById(ctx context.Context, id string) error
	// DeleteByToken removes the Token entity by its refresh token.
	DeleteByToken(ctx context.Context, refreshToken string) error
	// DeleteExpired removes expired tokens.
	DeleteExpired(ctx context.Context) error
}

type tokenRepository struct {
	db *sqlx.DB
}

func NewTokenRepo(db *sqlx.DB) TokenRepository {
	return &tokenRepository{db: db}
}

func (r *tokenRepository) GetById(ctx context.Context, id string) (*domain.Token, error) {
	op := "TokenRepository.GetById"
	query := `
	SELECT id, user_id, refresh_token, user_agent, ip_address, is_revoked, expires_at, created_at, updated_at
	FROM refresh_tokens WHERE id = $1
	`
	executor := ExtractTx(ctx, r.db)
	var token domain.Token
	err := sqlx.GetContext(ctx, executor, &token, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrTokenNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &token, nil
}

func (r *tokenRepository) GetByToken(ctx context.Context, refreshToken string) (*domain.Token, error) {
	op := "TokenRepository.GetByToken"
	query := `
	SELECT id, user_id, refresh_token, user_agent, ip_address, is_revoked, expires_at, updated_at, created_at
	FROM refresh_tokens WHERE refresh_token = $1
	`
	executor := ExtractTx(ctx, r.db)
	var token domain.Token
	err := sqlx.GetContext(ctx, executor, &token, query, refreshToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrTokenNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &token, nil
}

func (r *tokenRepository) Create(ctx context.Context, token *domain.Token) (string, error) {
	op := "TokenRepository.Create"
	query := `
	INSERT INTO refresh_tokens (id, user_id, refresh_token, user_agent, ip_address, is_revoked, expires_at, updated_at, created_at)
	VALUES (:id, :user_id, :refresh_token, :user_agent, :ip_address, :is_revoked, :expires_at, :updated_at, :created_at)
	`
	executor := ExtractTx(ctx, r.db)
	if _, err := sqlx.NamedExecContext(ctx, executor, query, token); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return token.Id, nil
}

func (r *tokenRepository) Update(ctx context.Context, id string, token *domain.Token) error {
	op := "TokenRepository.Update"
	query := `
	UPDATE refresh_tokens
	SET refresh_token = :refresh_token, user_agent = :user_agent, ip_address = :ip_address, is_revoked = :is_revoked,
	expires_at = :expires_at, updated_at = :updated_at, created_at = :created_at
	WHERE id = :id
	`
	executor := ExtractTx(ctx, r.db)
	if _, err := sqlx.NamedExecContext(ctx, executor, query, token); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *tokenRepository) DeleteById(ctx context.Context, id string) error {
	op := "TokenRepository.DeleteById"
	query := `
	DELETE FROM refresh_tokens
	WHERE id = $1
	`
	executor := ExtractTx(ctx, r.db)
	if _, err := executor.ExecContext(ctx, query, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *tokenRepository) DeleteByToken(ctx context.Context, refreshToken string) error {
	op := "TokenRepository.DeleteByToken"
	query := `
	DELETE FROM refresh_tokens
	WHERE refresh_token = $1
	`
	executor := ExtractTx(ctx, r.db)
	if _, err := executor.ExecContext(ctx, query, refreshToken); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *tokenRepository) DeleteExpired(ctx context.Context) error {
	op := "TokenRepository.DeleteExpired"
	query := "DELETE FROM refresh_tokens WHERE expires_at < $1"
	_, err := r.db.ExecContext(ctx, query, time.Now())
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
