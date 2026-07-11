package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dmi3midd/grpcsso/internal/config"
	"github.com/dmi3midd/grpcsso/internal/domain"
	"github.com/dmi3midd/grpcsso/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrInvalidRefreshToken     = errors.New("invalid refresh token")
	ErrInvalidAccessToken      = errors.New("invalid access token")
	ErrSubjectAndIDNotFound    = errors.New("subject and id not found")
	ErrTokenNotFound           = errors.New("token not found")
)

type TokensPair struct {
	RefreshToken string `json:"refreshToken"`
	AccessToken  string `json:"accessToken"`
}

type AccessClaims struct {
	User domain.UserDto `json:"user"`
	jwt.RegisteredClaims
}

type TokenManager interface {
	// GenerateTokens generates pair with access and refresh tokens and token id (TokensPair, tokenId, error).
	GenerateTokens(user *domain.UserDto) (*TokensPair, string, error)
	// ValidateRefreshToken validates refresh token and returns token id and user id (tokenId, userId, error).
	// It returns ("", "", error) if validation go wrong.
	// It returns [ErrUnexpectedSigningMethod] if the token uses an unexpected signing method.
	// It returns [ErrInvalidRefreshToken] if the token is invalid.
	// It returns [ErrSubjectAndIDNotFound] if subject or token ID are not found in claims.
	ValidateRefreshToken(refreshToken string) (string, string, error)
	// ValidateAccessToken validates access token and returns userDto and token id (userDto, tokenId, error).
	// It returns (nil, "", error) if validation go wrong.
	// It returns [ErrUnexpectedSigningMethod] if the token uses an unexpected signing method.
	// It returns [ErrInvalidAccessToken] if the token is invalid.
	// It returns [ErrSubjectAndIDNotFound] if subject or token ID are not found in claims.
	ValidateAccessToken(accessToken string) (*domain.UserDto, string, error)
	// FindToken finds and returns a Token entity by its id string.
	// It returns [ErrTokenNotFound] if no token are found.
	FindToken(ctx context.Context, id string) (*domain.Token, error)
	// SaveToken creates refresh token for the user.
	SaveToken(ctx context.Context, token *domain.Token) (string, error)
	// RemoveToken removes refresh token.
	RemoveToken(ctx context.Context, id string) error
}

type tokenManager struct {
	tokenRepo repository.TokenRepository
	jwtCfg    *config.JWTConfig
	keys      config.KeysPair
}

func NewTokenManager(tokenRepo repository.TokenRepository, keys *config.KeysPair, jwtCfg *config.JWTConfig) TokenManager {
	return &tokenManager{
		tokenRepo: tokenRepo,
		jwtCfg:    jwtCfg,
		keys:      *keys,
	}
}

func (s *tokenManager) GenerateTokens(user *domain.UserDto) (*TokensPair, string, error) {
	op := "TokenManager.GenerateTokens"
	accessExpiry := s.jwtCfg.AccessTokenTTL
	refreshExpiry := s.jwtCfg.RefreshTokenTTL
	now := time.Now()
	v7uuid, _ := uuid.NewV7()
	id := v7uuid.String()

	// Access token
	accessToken, err := jwt.NewWithClaims(
		jwt.SigningMethodRS256,
		&AccessClaims{
			User: *user,
			RegisteredClaims: jwt.RegisteredClaims{
				ID:        id,
				Issuer:    "grpcsso",
				Subject:   user.ID,
				Audience:  jwt.ClaimStrings{s.jwtCfg.Audience},
				ExpiresAt: jwt.NewNumericDate(now.Add(accessExpiry)),
				IssuedAt:  jwt.NewNumericDate(now),
				NotBefore: jwt.NewNumericDate(now),
			},
		},
	).SignedString(s.keys.PrivateKey)
	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", op, err)
	}

	// Refresh token
	refreshClaims := jwt.RegisteredClaims{
		ID:        id,
		Issuer:    "grpcsso",
		Subject:   user.ID,
		Audience:  jwt.ClaimStrings{s.jwtCfg.Audience},
		ExpiresAt: jwt.NewNumericDate(now.Add(refreshExpiry)),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodRS256, &refreshClaims).SignedString(s.keys.PrivateKey)
	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", op, err)
	}

	return &TokensPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, id, nil
}

func (s *tokenManager) ValidateRefreshToken(refreshToken string) (string, string, error) {
	op := "TokenManager.ValidateRefreshToken"
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("%s: %w %v", op, ErrUnexpectedSigningMethod, token.Header["alg"])
		}
		return s.keys.PublicKey, nil
	})

	if err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	if !token.Valid {
		return "", "", fmt.Errorf("%s: %w", op, ErrInvalidRefreshToken)
	}

	userId := claims.Subject
	tokenId := claims.ID

	if userId == "" || tokenId == "" {
		return "", "", fmt.Errorf("%s: %w", op, ErrSubjectAndIDNotFound)
	}

	return tokenId, userId, nil
}

func (s *tokenManager) ValidateAccessToken(accessToken string) (*domain.UserDto, string, error) {
	op := "TokenManager.ValidateAccessToken"
	claims := &AccessClaims{}
	token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("%s: %w %v", op, ErrUnexpectedSigningMethod, token.Header["alg"])
		}
		return s.keys.PublicKey, nil
	})

	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", op, err)
	}

	if !token.Valid {
		return nil, "", fmt.Errorf("%s: %w", op, ErrInvalidAccessToken)
	}

	userId := claims.Subject
	tokenId := claims.ID

	if userId == "" || tokenId == "" {
		return nil, "", fmt.Errorf("%s: %w", op, ErrSubjectAndIDNotFound)
	}

	return &domain.UserDto{
		ID:       userId,
		Username: claims.User.Username,
		Email:    claims.User.Email,
	}, tokenId, nil
}

func (s *tokenManager) FindToken(ctx context.Context, id string) (*domain.Token, error) {
	op := "TokenManager.FindToken"
	token, err := s.tokenRepo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrTokenNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrTokenNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return token, nil
}

// TODO: Review method
func (s *tokenManager) SaveToken(ctx context.Context, token *domain.Token) (string, error) {
	op := "TokenManager.SaveToken"

	claims := &jwt.RegisteredClaims{}
	_, _, err := jwt.NewParser().ParseUnverified(token.RefreshToken, claims)
	if err != nil {
		return "", fmt.Errorf("%s: failed to parse refresh token: %w", op, err)
	}

	token.ExpiresAt = time.Now().Add(s.jwtCfg.RefreshTokenTTL)
	token.CreatedAt = time.Now()

	id, err := s.tokenRepo.Create(ctx, token)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *tokenManager) RemoveToken(ctx context.Context, id string) error {
	op := "TokenManager.RemoveToken"
	if err := s.tokenRepo.DeleteById(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
