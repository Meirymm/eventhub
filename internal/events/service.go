package events

import (
	"eventhub/internal/models"
)

type EventService struct {
	repo *EventRepository
}

func NewEventService(repo *EventRepository) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) GetAllEvents() ([]models.Event, error) {
	// Здесь можно добавить логику (например, фильтрацию)
	return s.repo.GetAll()
}
