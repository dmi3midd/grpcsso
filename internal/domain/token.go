package domain

type Token struct {
	ID           string `json:"id" db:"id"`
	UserID       string `json:"userId" db:"user_id"`
	RefreshToken string `json:"refreshToken" db:"refresh_token"`
	UserAgent    string `json:"userAgent" db:"user_agent"`
	IpAddress    string `json:"ipAddress" db:"ip_address"`
	ExpiresAt    string `json:"expiresAt" db:"expires_at"`
	CreatedAt    string `json:"createdAt" db:"created_at"`
}
