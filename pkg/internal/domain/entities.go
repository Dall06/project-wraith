package domain

import "time"

type User struct {
	ID        string
	Username  string
	Email     string
	Name      string
	Phone     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	Status    string
	Meta      map[string]interface{}
}
