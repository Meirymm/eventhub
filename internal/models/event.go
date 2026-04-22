package models

import "time"

type Event struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	VenueID     int       `json:"venue_id"`
	CategoryID  *int       `json:"category_id"`
	StartTime   time.Time `json:"start_time"`
	CreatedAt   time.Time `json:"created_at"`
}


type CreateEventRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	VenueID     int       `json:"venue_id" binding:"required"`
	StartTime   time.Time `json:"start_time" binding:"required"`
}
