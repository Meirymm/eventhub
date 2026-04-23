package events

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"

	qrcode "github.com/skip2/go-qrcode"
	"eventhub/internal/models"
)

type TicketRepository struct {
	db *sql.DB
}

func NewTicketRepository(db *sql.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

func generateQR(ticketID, eventID, userID int) string {
	content := fmt.Sprintf(
		"EVENTHUB:ticket=%d:event=%d:user=%d:ts=%d",
		ticketID, eventID, userID, time.Now().Unix(),
	)
	png, err := qrcode.Encode(content, qrcode.Medium, 256)
	if err != nil {
		return ""
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)
}

func (r *TicketRepository) Create(req models.CreateTicketRequest) (*models.Ticket, error) {
	query := `INSERT INTO tickets (event_id, user_id, created_at)
	          VALUES ($1, $2, NOW())
	          RETURNING id, event_id, user_id, created_at`
	var t models.Ticket
	err := r.db.QueryRow(query, req.EventID, req.UserID).
		Scan(&t.ID, &t.EventID, &t.UserID, &t.CreatedAt)
	if err != nil {
		return nil, err
	}

	t.QRCode = generateQR(t.ID, t.EventID, t.UserID)
	_, err = r.db.Exec(`UPDATE tickets SET qr_code = $1 WHERE id = $2`, t.QRCode, t.ID)
	if err != nil {
		t.QRCode = ""
	}

	return &t, nil
}

func (r *TicketRepository) GetByUserID(userID int) ([]models.Ticket, error) {
	query := `
		SELECT t.id, t.event_id, t.user_id, t.created_at, COALESCE(t.qr_code, ''),
		       COALESCE(e.title, 'Unknown Event')
		FROM tickets t
		LEFT JOIN events e ON e.id = t.event_id
		WHERE t.user_id = $1
		ORDER BY t.created_at DESC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tickets []models.Ticket
	for rows.Next() {
		var t models.Ticket
		err := rows.Scan(&t.ID, &t.EventID, &t.UserID, &t.CreatedAt, &t.QRCode, &t.EventTitle)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}
	return tickets, nil
}

func (r *TicketRepository) GetByID(id int) (*models.Ticket, error) {
	query := `SELECT id, event_id, user_id, created_at, COALESCE(qr_code, '') FROM tickets WHERE id = $1`
	var t models.Ticket
	err := r.db.QueryRow(query, id).Scan(&t.ID, &t.EventID, &t.UserID, &t.CreatedAt, &t.QRCode)
	if err != nil {
		return nil, err
	}
	return &t, nil
}