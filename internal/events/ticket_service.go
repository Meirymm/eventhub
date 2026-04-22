package events

import "eventhub/internal/models"

type TicketService struct {
	repo *TicketRepository
}

func NewTicketService(repo *TicketRepository) *TicketService {
	return &TicketService{repo: repo}
}

func (s *TicketService) BookTicket(req models.CreateTicketRequest) (*models.Ticket, error) {
	return s.repo.Create(req)
}

func (s *TicketService) GetTicketsByUser(userID int) ([]models.Ticket, error) {
	return s.repo.GetByUserID(userID)
}
