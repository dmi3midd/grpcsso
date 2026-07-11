package domain

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
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
