package events

import (
	"fmt"
	"testing"
	"time"

	"eventhub/internal/models"
)

// ══════════════════════════════════════════════════
//  MOCK — фейковый репозиторий событий
// ══════════════════════════════════════════════════

type mockEventRepo struct {
	events []models.Event
	nextID int
}

func (m *mockEventRepo) GetAll() ([]models.Event, error) {
	return m.events, nil
}

func (m *mockEventRepo) GetByID(id int) (*models.Event, error) {
	for i := range m.events {
		if m.events[i].ID == id {
			return &m.events[i], nil
		}
	}
	return nil, fmt.Errorf("event not found")
}

func (m *mockEventRepo) Create(req models.CreateEventRequest) (*models.Event, error) {
	if req.Title == "" {
		return nil, fmt.Errorf("title is required")
	}
	m.nextID++
	e := models.Event{
		ID:          m.nextID,
		Title:       req.Title,
		Description: req.Description,
		VenueID:     req.VenueID,
		StartTime:   req.StartTime,
		CreatedAt:   time.Now(),
	}
	m.events = append(m.events, e)
	return &m.events[len(m.events)-1], nil
}

func (m *mockEventRepo) Update(id int, req models.CreateEventRequest) (*models.Event, error) {
	for i := range m.events {
		if m.events[i].ID == id {
			m.events[i].Title       = req.Title
			m.events[i].Description = req.Description
			m.events[i].VenueID     = req.VenueID
			m.events[i].StartTime   = req.StartTime
			return &m.events[i], nil
		}
	}
	return nil, fmt.Errorf("event not found")
}

func (m *mockEventRepo) Delete(id int) error {
	for i, e := range m.events {
		if e.ID == id {
			m.events = append(m.events[:i], m.events[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("event not found")
}

// ══════════════════════════════════════════════════
//  TABLE-DRIVEN ТЕСТ: CreateEvent
//  Запуск: go test ./internal/events/... -v -run TestCreateEvent
// ══════════════════════════════════════════════════

func TestCreateEvent(t *testing.T) {
	tests := []struct {
		name    string
		req     models.CreateEventRequest
		wantErr bool
	}{
		{
			name: "успешное создание",
			req:  models.CreateEventRequest{
				Title:       "Tech Conference 2026",
				Description: "Большая IT конференция",
				VenueID:     1,
				StartTime:   time.Now().Add(24 * time.Hour),
			},
			wantErr: false,
		},
		{
			name: "пустое название",
			req:  models.CreateEventRequest{
				Title:     "",
				VenueID:   1,
				StartTime: time.Now().Add(24 * time.Hour),
			},
			wantErr: true,
		},
		{
			name: "второе событие",
			req:  models.CreateEventRequest{
				Title:       "Летний фестиваль",
				Description: "Музыка под открытым небом",
				VenueID:     2,
				StartTime:   time.Now().Add(48 * time.Hour),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockEventRepo{}
			svc  := NewEventService(repo)

			event, err := svc.CreateEvent(tt.req)

			if tt.wantErr {
				if err == nil {
					t.Error("FAIL: ожидали ошибку — её нет")
				}
				return
			}
			if err != nil {
				t.Errorf("FAIL: не ожидали ошибку: %v", err)
				return
			}
			if event.ID == 0 {
				t.Error("FAIL: ID события равен 0")
			}
			if event.Title != tt.req.Title {
				t.Errorf("FAIL: ожидали title '%s', получили '%s'", tt.req.Title, event.Title)
			}
		})
	}
}

// ══════════════════════════════════════════════════
//  TABLE-DRIVEN ТЕСТ: GetEventByID
//  Запуск: go test ./internal/events/... -v -run TestGetEventByID
// ══════════════════════════════════════════════════

func TestGetEventByID(t *testing.T) {
	repo := &mockEventRepo{}
	svc  := NewEventService(repo)

	// Создаём событие для тестов
	created, _ := svc.CreateEvent(models.CreateEventRequest{
		Title:     "Тестовое событие",
		VenueID:   1,
		StartTime: time.Now().Add(24 * time.Hour),
	})

	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{
			name:    "найти существующее",
			id:      created.ID,
			wantErr: false,
		},
		{
			name:    "несуществующий ID",
			id:      999,
			wantErr: true,
		},
		{
			name:    "отрицательный ID",
			id:      -1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event, err := svc.GetEventByID(tt.id)

			if tt.wantErr {
				if err == nil {
					t.Error("FAIL: ожидали ошибку — её нет")
				}
				return
			}
			if err != nil {
				t.Errorf("FAIL: не ожидали ошибку: %v", err)
				return
			}
			if event.Title != "Тестовое событие" {
				t.Errorf("FAIL: неверный title: %s", event.Title)
			}
		})
	}
}

// ══════════════════════════════════════════════════
//  TABLE-DRIVEN ТЕСТ: DeleteEvent
//  Запуск: go test ./internal/events/... -v -run TestDeleteEvent
// ══════════════════════════════════════════════════

func TestDeleteEvent(t *testing.T) {
	repo := &mockEventRepo{}
	svc  := NewEventService(repo)

	e1, _ := svc.CreateEvent(models.CreateEventRequest{Title: "Удалить меня", VenueID: 1, StartTime: time.Now()})
	e2, _ := svc.CreateEvent(models.CreateEventRequest{Title: "Оставить меня", VenueID: 1, StartTime: time.Now()})

	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{
			name:    "удалить существующее",
			id:      e1.ID,
			wantErr: false,
		},
		{
			name:    "удалить несуществующее",
			id:      999,
			wantErr: true,
		},
		{
			name:    "удалить повторно",
			id:      e1.ID, // уже удалено выше
			wantErr: true,
		},
	}

	_ = e2 // e2 не удаляем — проверяем что не задели чужое

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.DeleteEvent(tt.id)

			if tt.wantErr {
				if err == nil {
					t.Error("FAIL: ожидали ошибку — её нет")
				}
				return
			}
			if err != nil {
				t.Errorf("FAIL: не ожидали ошибку: %v", err)
			}
		})
	}

	// Проверяем что e2 на месте
	remaining, _ := svc.GetAllEvents()
	if len(remaining) != 1 {
		t.Errorf("FAIL: должно остаться 1 событие, осталось %d", len(remaining))
	}
}