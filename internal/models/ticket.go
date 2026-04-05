package models

import "time"

type Ticket struct {
	ID        int       `json:"id"`
	EventID   int       `json:"event_id"`
	UserID    int       `json:"user_id"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateTicketRequest struct {
	EventID int `json:"event_id" binding:"required"`
	UserID  int `json:"user_id" binding:"required"`
}
