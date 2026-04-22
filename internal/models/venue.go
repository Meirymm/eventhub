package models

type Venue struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Capacity int    `json:"capacity"`
}
