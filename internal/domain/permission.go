package domain

type Permission struct {
	Id   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}
