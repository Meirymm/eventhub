package models

import "time"

type Ticket struct {
	ID         int       `json:"id"`
	EventID    int       `json:"event_id"`
	UserID     int       `json:"user_id"`
	Price      float64   `json:"price"`
	QRCode     string    `json:"qr_code,omitempty"`
	EventTitle string    `json:"event_title,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type CreateTicketRequest struct {
	EventID int `json:"event_id" binding:"required"`
	UserID  int `json:"-"` // берётся из JWT, не из запроса
}