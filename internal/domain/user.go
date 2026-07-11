package domain

type User struct {
	ID           string `json:"id" db:"id"`
	Username     string `json:"username" db:"username"`
	Email        string `json:"email" db:"email"`
	PasswordHash string `json:"passwordHash" db:"password_hash"`
	CreatedAt    string `json:"createdAt" db:"created_at"`
	UpdatedAt    string `json:"updatedAt" db:"updated_at"`
}

type UserDto struct {
	ID       string `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
}

func (u *User) ToUserDto() *UserDto {
	return &UserDto{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
	}
}
