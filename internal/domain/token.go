package domain

type Token struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
	UserAgent    string `json:"user_agent"`
	IpAddress    string `json:"ip_address"`
	ExpiresAt    string `json:"expires_at"`
	CreatedAt    string `json:"created_at"`
}
