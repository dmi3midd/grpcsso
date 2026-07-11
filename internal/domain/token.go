package domain

import "time"

type Token struct {
	ID           string    `json:"id" db:"id"`
	UserID       string    `json:"userId" db:"user_id"`
	RefreshToken string    `json:"refreshToken" db:"refresh_token"`
	UserAgent    string    `json:"userAgent" db:"user_agent"`
	IpAddress    string    `json:"ipAddress" db:"ip_address"`
	ExpiresAt    time.Time `json:"expiresAt" db:"expires_at"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
}
