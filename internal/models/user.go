package models

import "time"

type User struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	ValidUntil time.Time `json:"valid_until"`
	IsValid    bool      `json:"is_valid"`
	IsAdmin    bool      `json:"is_admin"`
}
