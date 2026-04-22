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

func (r *EventRepository) GetAll() ([]models.Event, error) {
	query := `SELECT id, title, description, venue_id, category_id, start_time, created_at
	          FROM events ORDER BY start_time ASC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var e models.Event
		err := rows.Scan(
			&e.ID, &e.Title, &e.Description,
			&e.VenueID, &e.CategoryID, &e.StartTime, &e.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}

func (r *EventRepository) GetByID(id int) (*models.Event, error) {
	query := `SELECT id, title, description, venue_id, category_id, start_time, created_at
	          FROM events WHERE id = $1`
	var e models.Event
	err := r.db.QueryRow(query, id).Scan(
		&e.ID, &e.Title, &e.Description,
		&e.VenueID, &e.CategoryID, &e.StartTime, &e.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *EventRepository) Create(req models.CreateEventRequest) (*models.Event, error) {
	query := `INSERT INTO events (title, description, venue_id, start_time, created_at)
	          VALUES ($1, $2, $3, $4, NOW())
	          RETURNING id, title, description, venue_id, start_time, created_at`
	var e models.Event
	err := r.db.QueryRow(query, req.Title, req.Description, req.VenueID, req.StartTime).
		Scan(&e.ID, &e.Title, &e.Description, &e.VenueID, &e.StartTime, &e.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *EventRepository) Update(id int, req models.CreateEventRequest) (*models.Event, error) {
	query := `UPDATE events
	          SET title=$1, description=$2, venue_id=$3, start_time=$4
	          WHERE id=$5
	          RETURNING id, title, description, venue_id, start_time, created_at`
	var e models.Event
	err := r.db.QueryRow(query, req.Title, req.Description, req.VenueID, req.StartTime, id).
		Scan(&e.ID, &e.Title, &e.Description, &e.VenueID, &e.StartTime, &e.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *EventRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM events WHERE id = $1`, id)
	return err
}