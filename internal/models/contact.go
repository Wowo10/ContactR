package models

import "time"

type Contact struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Contact     string    `json:"contact"`
	DateCreated time.Time `json:"dateCreated"`
	DateUpdated time.Time `json:"dateUpdated"`
	Tags        []string  `json:"is_valid"`
	LinkedIn    string    `json:"linkedInUrl"`
	Credly      string    `json:"credlyInUrl"`
}
