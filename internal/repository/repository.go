package repository

import (
	"github.com/jmoiron/sqlx"
	"theater/internal/models"
	"theater/internal/repository/postgres"
)

type Event interface {
	CreateEvent(e *models.Event) (uint, error)
	GetAllEvents() (*[]models.Event, error)
	GetEventByID(id uint) (*models.Event, error)
	UpdateEvent(e *models.Event) error
	DeleteEvent(id uint) error
}

type Club interface {
	CreateClub(club *models.Club) (int, error)
	GetAllClubs() (*[]models.Club, error)
	GetClubByID(id int) (*models.Club, error)
	UpdateClub(club *models.Club) error
	DeleteClub(id int) error
}

type Repository struct {
	Event
	Club
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Event: postgres.NewEventPostgres(db),
		Club:  postgres.NewClubPostgres(db),
	}
}
