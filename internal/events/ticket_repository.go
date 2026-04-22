package events

import (
	"database/sql"
	"eventhub/internal/models"
)

type TicketRepository struct {
	db *sql.DB
}

func NewTicketRepository(db *sql.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

func (r *TicketRepository) Create(req models.CreateTicketRequest) (*models.Ticket, error) {
	query := `INSERT INTO tickets (event_id, user_id, created_at) 
	          VALUES ($1, $2, NOW()) RETURNING id, event_id, user_id, created_at`
	var t models.Ticket
	err := r.db.QueryRow(query, req.EventID, req.UserID).
		Scan(&t.ID, &t.EventID, &t.UserID, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TicketRepository) GetByUserID(userID int) ([]models.Ticket, error) {
	query := `SELECT id, event_id, user_id, created_at FROM tickets WHERE user_id = $1`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tickets []models.Ticket
	for rows.Next() {
		var t models.Ticket
		err := rows.Scan(&t.ID, &t.EventID, &t.UserID, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}
	return tickets, nil
}