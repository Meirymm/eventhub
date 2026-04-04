package events

import (
	"database/sql"
	"eventhub/internal/models"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

// Твоя первая функция: получить все ивенты из базы
func (r *EventRepository) GetAll() ([]models.Event, error) {
	query := `SELECT id, title, description, venue_id, category_id, start_time, created_at FROM events`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var e models.Event
		err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.VenueID, &e.CategoryID, &e.StartTime, &e.CreatedAt)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}
