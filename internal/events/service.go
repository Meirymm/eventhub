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

	return s.repo.GetAll()
}

func (s *EventService) GetEventByID(id int) (*models.Event, error) {
	return s.repo.GetByID(id)
}

func (s *EventService) CreateEvent(req models.CreateEventRequest) (*models.Event, error) {
	return s.repo.Create(req)
}

func (s *EventService) UpdateEvent(id int, req models.CreateEventRequest) (*models.Event, error) {
	return s.repo.Update(id, req)
}

func (s *EventService) DeleteEvent(id int) error {
	return s.repo.Delete(id)
}
