package domain

import "time"

type User struct {
	Id           string    `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"passwordHash" db:"password_hash"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}

type UserDto struct {
	Id       string `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
}

func (u *User) ToUserDto() *UserDto {
	return &UserDto{
		Id:       u.Id,
		Username: u.Username,
		Email:    u.Email,
	}
}
