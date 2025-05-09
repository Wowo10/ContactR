package models

import (
	"time"
)

type Contact struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Contact     string    `json:"contact"`
	Tags        []string  `json:"tags"`
	LinkedIn    string    `json:"linkedInUrl"`
	Credly      string    `json:"credlyInUrl"`
	DateCreated time.Time `json:"dateCreated",omitempty`
	DateUpdated time.Time `json:"dateUpdated",omitempty`
}
