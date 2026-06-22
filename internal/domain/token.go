package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	Id           string    `json:"id" db:"id"`
	RefreshToken string    `json:"refreshToken" db:"refresh_token"`
	UserId       string    `json:"userId" db:"user_id"`
	ClientId     string    `json:"clientId" db:"client_id"`
	ExpiresAt    time.Time `json:"expiresAt" db:"expires_at"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}

type TokensPair struct {
	RefreshToken string `json:"refreshToken"`
	AccessToken  string `json:"accessToken"`
}

// type UserClaims struct {
// 	UserId      string   `json:"userId"`
// 	Username    string   `json:"username"`
// 	Email       string   `json:"email"`
// 	Permissions []string `json:"permissions"`
// }

type AccessClaims struct {
	User UserDto `json:"user"`
	jwt.RegisteredClaims
}
