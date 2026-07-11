package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dmi3midd/grpcsso/internal/domain"
	"github.com/dmi3midd/grpcsso/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExist = errors.New("user already exist")
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidPassword  = errors.New("invalid password")
)

type AuthDto struct {
	User         domain.UserDto
	AccessToken  string
	RefreshToken string
}

type UserService interface {
	// Registration performs user registration.
	// It returns ErrUserAlreadyExist if the user exist.
	Registration(ctx context.Context, username, email, password string) error
	// Login performs user login and returns LoginResult struct.
	// It returns [ErrUserNotFound] if no user are found.
	// It returns [ErrInvalidPassword] if the password is invalid.
	Login(ctx context.Context, email, password, userAgent, ipAddress string) (*AuthDto, error)
	// Logout performs logout user.
	// Look at TokenService.ValidateRefreshToken for errors.
	Logout(ctx context.Context, refreshToken string) error
	// Refresh performs refreshing access and refresh tokens.
	// It returns [ErrUserNotFound] if no user are found.
	// Look at TokenService.ValidateRefreshToken for other errors.
	Refresh(ctx context.Context, refreshToken, ipAddress, userAgent string) (*AuthDto, error)
}

type userService struct {
	userRepo     repository.UserRepository
	tokenManager TokenManager
}

func NewUserService(userRepo repository.UserRepository, tokenManager TokenManager) UserService {
	return &userService{
		userRepo:     userRepo,
		tokenManager: tokenManager,
	}
}

func (s *userService) Registration(ctx context.Context, username, email, password string) error {
	op := "UserService.Registration"

	candidate, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		return fmt.Errorf("%s: %w", op, err)
	}

	if candidate != nil {
		return fmt.Errorf("%s: %w", op, ErrUserAlreadyExist)
	}

	v7uuid, _ := uuid.NewV7()
	id := v7uuid.String()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	user := &domain.User{
		ID:           id,
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if _, err := s.userRepo.Create(ctx, user); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *userService) Login(ctx context.Context, email, password, userAgent, ipAddress string) (*AuthDto, error) {
	op := "UserService.Login"

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidPassword)
	}

	userDto := user.ToUserDto()
	tokens, tokenId, err := s.tokenManager.GenerateTokens(userDto)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	_, err = s.tokenManager.SaveToken(ctx, &domain.Token{
		ID:           tokenId,
		UserID:       userDto.ID,
		RefreshToken: tokens.RefreshToken,
		UserAgent:    userAgent,
		IpAddress:    ipAddress,
		IsRevoked:    false,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &AuthDto{
		User:         *userDto,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (s *userService) Logout(ctx context.Context, refreshToken string) error {
	op := "UserService.Logout"
	tokenId, _, err := s.tokenManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if err := s.tokenManager.RemoveToken(ctx, tokenId); err != nil {
		if errors.Is(err, ErrTokenNotFound) {
			return nil
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *userService) Refresh(ctx context.Context, refreshToken, ipAddress, userAgent string) (*AuthDto, error) {
	op := "UserService.Refresh"
	tokenId, userId, err := s.tokenManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// if token is not found user is unauthorized and need to login
	if err := s.tokenManager.RevokeToken(ctx, tokenId); err != nil {
		// if errors.Is(err, ErrTokenNotFound) {
		// 	return nil, fmt.Errorf("%s: %w", op, ErrTokenNotFound)
		// }
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	user, err := s.userRepo.GetById(ctx, userId)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	userDto := user.ToUserDto()

	tokens, newTokenId, err := s.tokenManager.GenerateTokens(userDto)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if _, err := s.tokenManager.SaveToken(ctx, &domain.Token{
		ID:           newTokenId,
		UserID:       userId,
		RefreshToken: tokens.RefreshToken,
		UserAgent:    userAgent,
		IpAddress:    ipAddress,
		IsRevoked:    false,
	}); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &AuthDto{
		User:         *userDto,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
