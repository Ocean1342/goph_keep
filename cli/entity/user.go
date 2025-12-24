package entity

import "time"

type User struct {
	ID              int       `json:"id" db:"id"`
	Login           string    `json:"login" db:"login"`
	Password        string    `json:"password" db:"password"`
	Phrase          string    `json:"phrase" db:"phrase"`
	Token           string    `json:"token" db:"token"`
	IsAuthenticated bool      `json:"is_authenticated" db:"is_authenticated"`
	AuthDt          time.Time `json:"auth_dt" db:"auth_dt"`
}
