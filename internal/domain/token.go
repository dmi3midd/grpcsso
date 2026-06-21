package domain

import "time"

type Token struct {
	Id           string    `json:"id" db:"id"`
	RefreshToken string    `json:"refreshToken" db:"refresh_token"`
	UserId       string    `json:"userId" db:"user_id"`
	ClientId     string    `json:"clientId" db:"client_id"`
	ExpiresAt    time.Time `json:"expiresAt" db:"expires_at"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}
