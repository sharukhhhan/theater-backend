package postgres

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
	"theater/internal/models"
)

type ClubPostgres struct {
	db *sqlx.DB
}

func NewClubPostgres(db *sqlx.DB) *ClubPostgres {
	return &ClubPostgres{db: db}
}

func (r *ClubPostgres) CreateClub(club *models.Club) (int, error) {
	query := `INSERT INTO clubs (name, description, contacts, price, spots_number) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id int
	err := r.db.QueryRow(query, club.Name, club.Description, club.Contact, club.Price, club.SpotsNumber).Scan(&id)
	if err != nil {
		if ok, _ := IsDuplicateKeyError(err); ok {
			return 0, errors.New(ErrClubNameAlreadyExists)
		}
		return 0, err
	}

	return id, nil
}

func (r *ClubPostgres) GetAllClubs() (*[]models.Club, error) {
	query := `SELECT * FROM clubs`
	var clubs []models.Club
	err := r.db.Select(&clubs, query)
	if err != nil {
		return nil, err
	}

	return &clubs, nil
}

func (r *ClubPostgres) GetClubByID(id int) (*models.Club, error) {
	query := `SELECT * FROM clubs WHERE id = $1`
	var club models.Club
	err := r.db.Get(&club, query, id)
	if err != nil {
		return nil, err
	}

	return &club, nil
}

func (r *ClubPostgres) UpdateClub(club *models.Club) error {
	// Start with the base update statement
	query := `UPDATE clubs SET`
	args := []interface{}{}
	updateFields := []string{}

	// Check each field for non-zero values and add them to the query
	if club.Name != "" {
		updateFields = append(updateFields, " name = $"+strconv.Itoa(len(args)+1))
		args = append(args, club.Name)
	}
	if club.Description != "" {
		updateFields = append(updateFields, " description = $"+strconv.Itoa(len(args)+1))
		args = append(args, club.Description)
	}
	if club.Price != 0 {
		updateFields = append(updateFields, " price = $"+strconv.Itoa(len(args)+1))
		args = append(args, club.Price)
	}
	if club.SpotsNumber != 0 {
		updateFields = append(updateFields, " spots_number = $"+strconv.Itoa(len(args)+1))
		args = append(args, club.SpotsNumber)
	}
	if club.IsActive != "" {
		updateFields = append(updateFields, " is_active = $"+strconv.Itoa(len(args)+1))
		args = append(args, club.IsActive)
	}

	// If no fields to update, return early
	if len(updateFields) == 0 {
		return nil
	}

	// Add the fields to the query
	query += strings.Join(updateFields, ",")
	query += " WHERE id = $" + strconv.Itoa(len(args)+1)
	args = append(args, club.ID)

	// Prepare and execute the query
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	result, err := stmt.Exec(args...)
	if err != nil {
		if ok, _ := IsDuplicateKeyError(err); ok {
			return errors.New(ErrClubNameAlreadyExists)
		}
		return err
	}

	if rowsAffected, err := result.RowsAffected(); rowsAffected == 0 {
		if err != nil {
			return fmt.Errorf("error while counting rows affecting: %v", err)
		}
		return errors.New(ErrClubNotFound)
	}

	return nil
}

func (r *ClubPostgres) DeleteClub(id int) error {
	query := `DELETE FROM clubs WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	if rowsAffected, err := result.RowsAffected(); rowsAffected == 0 {
		if err != nil {
			return fmt.Errorf("error while counting rows affecting: %v", err)
		}
		return errors.New(ErrClubNotFound)
	}

	return nil
}
