package models

type Contact struct {
	Id       int      `json:"id"`
	Name     string   `json:"name"`
	Contact  string   `json:"contact"`
	Tags     []string `json:"tags"`
	LinkedIn string   `json:"linkedInUrl"`
	Credly   string   `json:"credlyInUrl"`
}
