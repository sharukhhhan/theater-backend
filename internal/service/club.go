package service

import (
	"theater/internal/models"
	"theater/internal/repository"
)

type ClubService struct {
	repo repository.Club
}

func NewClubService(repo repository.Club) *ClubService {
	return &ClubService{
		repo: repo,
	}
}

func (c *ClubService) CreateClub(club *models.Club) (int, error) {
	id, err := c.repo.CreateClub(club)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (c *ClubService) GetAllClubs() (*[]models.Club, error) {
	return c.repo.GetAllClubs()
}

func (c *ClubService) GetClubByID(id int) (*models.Club, error) {
	return c.repo.GetClubByID(id)
}

func (c *ClubService) UpdateClub(club *models.Club) error {
	return c.repo.UpdateClub(club)
}

func (c *ClubService) DeleteClub(id int) error {
	return c.repo.DeleteClub(id)
}
