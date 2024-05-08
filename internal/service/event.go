package service

import (
	"theater/internal/models"
	"theater/internal/repository"
)

type EventService struct {
	repo repository.Event
}

func NewEventService(repo repository.Event) *EventService {
	return &EventService{
		repo: repo,
	}
}

func (e *EventService) CreateEvent(event *models.Event) (uint, error) {
	id, err := e.repo.CreateEvent(event)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (e *EventService) GetAllEvents() (*[]models.Event, error) {
	return e.repo.GetAllEvents()
}

func (e *EventService) GetEventByID(id uint) (*models.Event, error) {
	return e.repo.GetEventByID(id)
}

func (e *EventService) UpdateEvent(event *models.Event) error {
	return e.repo.UpdateEvent(event)
}

func (e *EventService) DeleteEvent(id uint) error {
	return e.repo.DeleteEvent(id)
}
