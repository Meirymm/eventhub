package models

type Ticket struct {
	ID      int     `json:"id"`
	EventID int     `json:"event_id"`
	UserID  int     `json:"user_id"`
	Price   float64 `json:"price"`
}
