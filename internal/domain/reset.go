package domain

import "time"

type Reset struct {
	Id        string     `db:"id"`
	UserId    string     `db:"user_id"`
	UsedAt    *time.Time `db:"used_at"`
	CreatedAt time.Time  `db:"created_at"`
}
