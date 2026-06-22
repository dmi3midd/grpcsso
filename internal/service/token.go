package service

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/dmi3midd/grpcsso/internal/config"
	"github.com/dmi3midd/grpcsso/internal/domain"
	"github.com/dmi3midd/grpcsso/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/xid"
)

var (
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrInvalidRefreshToken     = errors.New("invalid refresh token")
	ErrInvalidAccessToken      = errors.New("invalid access token")
	ErrSubjectAndIDNotFound    = errors.New("subject and id not found")
	ErrTokenNotFound           = errors.New("token not found")
)

type TokenService interface {
	// GenerateTokens generates pair with access and refresh tokens and token id (TokensPair, tokenId, error).
	GenerateTokens(user domain.UserDto, clientId string) (*domain.TokensPair, string, error)
	// ValidateRefreshToken validates refresh token and returns token and user id (tokenId, userId, error).
	// It returns ("", "", error) if validation go wrong.
	// It returns ErrUnexpectedSigningMethod if the token uses an unexpected signing method.
	// It returns ErrInvalidRefreshToken if the token is invalid.
	// It returns ErrSubjectAndIDNotFound if subject or token ID are not found in claims.
	ValidateRefreshToken(refreshToken string) (string, string, error)
	// ValidateAccessToken validates access token and returns userDto and token id (userDto, tokenId, error).
	// It returns (nil, "", error) if validation go wrong.
	// It returns ErrUnexpectedSigningMethod if the token uses an unexpected signing method.
	// It returns ErrInvalidAccessToken if the token is invalid.
	// It returns ErrSubjectAndIDNotFound if subject or token ID are not found in claims.
	ValidateAccessToken(accessToken string) (*domain.UserDto, string, error)
	// SaveToken creates refresh token for the user.
	SaveToken(ctx context.Context, refreshToken, userId, clientId, tokenId string) (string, error)
	// RemoveToken removes refresh token.
	// It returns ErrTokenNotFound if no token are found.
	RemoveToken(ctx context.Context, id string) error
	// FindToken finds and returns a Token entity by its refresh token string.
	// It returns ErrTokenNotFound if no token are found.
	FindToken(ctx context.Context, id string) (*domain.Token, error)
	// GetPublicKey returns public rsa keys
	GetPublicKey() rsa.PublicKey
}

type tokenService struct {
	tokenStore repository.TokenRepository
	keys       config.KeysPair
}

func NewTokenService(tokenStore repository.TokenRepository, keys *config.KeysPair) TokenService {
	return &tokenService{
		tokenStore: tokenStore,
		keys:       *keys,
	}
}

func (s *tokenService) GenerateTokens(user domain.UserDto, clientId string) (*domain.TokensPair, string, error) {
	op := "TokenService.GenerateTokens"
	accessExpiry, _ := time.ParseDuration("30m")
	refreshExpiry, _ := time.ParseDuration("336h")
	now := time.Now()
	id := xid.New().String()

	// Access token
	accessToken, err := jwt.NewWithClaims(
		jwt.SigningMethodRS256,
		&domain.AccessClaims{
			User: user,
			RegisteredClaims: jwt.RegisteredClaims{
				ID:        id,
				Issuer:    "grpcsso",
				Subject:   user.Id,
				Audience:  jwt.ClaimStrings{clientId},
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
		Subject:   user.Id,
		Audience:  jwt.ClaimStrings{clientId},
		ExpiresAt: jwt.NewNumericDate(now.Add(refreshExpiry)),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodRS256, &refreshClaims).SignedString(s.keys.PrivateKey)
	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", op, err)
	}

	return &domain.TokensPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, id, nil
}

func (s *tokenService) ValidateRefreshToken(refreshToken string) (string, string, error) {
	op := "TokenService.ValidateRefreshToken"
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

func (s *tokenService) ValidateAccessToken(accessToken string) (*domain.UserDto, string, error) {
	op := "TokenService.ValidateAccessToken"
	claims := &domain.AccessClaims{}
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
		Id:       userId,
		Username: claims.User.Username,
		Email:    claims.User.Email,
	}, tokenId, nil
}

// FindToken implements [TokenService].
func (t *tokenService) FindToken(ctx context.Context, id string) (*domain.Token, error) {
	panic("unimplemented")
}

// GetPublicKey implements [TokenService].
func (t *tokenService) GetPublicKey() rsa.PublicKey {
	panic("unimplemented")
}

// RemoveToken implements [TokenService].
func (t *tokenService) RemoveToken(ctx context.Context, id string) error {
	panic("unimplemented")
}

// SaveToken implements [TokenService].
func (t *tokenService) SaveToken(ctx context.Context, refreshToken string, userId string, clientId string, tokenId string) (string, error) {
	panic("unimplemented")
}
