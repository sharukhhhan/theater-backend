package service

import (
	"theater/internal/models"
	"theater/internal/repository"
	"theater/pkg/hasher"
	"time"
)

type Auth interface {
	Login(username, password string) (string, error)
	GenerateToken() (string, error)
	ParseToken(tokenString string) error
}

type Event interface {
	CreateEvent(event *models.Event) (uint, error)
	GetAllEvents() (*[]models.Event, error)
	GetEventByID(id uint) (*models.Event, error)
	UpdateEvent(event *models.Event) error
	DeleteEvent(id uint) error
}

type Club interface {
	CreateClub(club *models.Club) (int, error)
	GetAllClubs() (*[]models.Club, error)
	GetClubByID(id int) (*models.Club, error)
	UpdateClub(club *models.Club) error
	DeleteClub(id int) error
}

type Dependency struct {
	Repo     *repository.Repository
	SignKey  string
	TokenTTL time.Duration

	Hasher        hasher.PasswordHasher
	AdminUsername string
	AdminPassword string
}

type Service struct {
	Auth
	Event
	Club
}

func NewService(dependency Dependency) *Service {
	return &Service{
		Auth:  NewAuthService(dependency.Hasher, dependency.SignKey, dependency.TokenTTL, dependency.AdminUsername, dependency.AdminPassword),
		Event: NewEventService(dependency.Repo.Event),
		Club:  NewClubService(dependency.Repo.Club),
	}
}
