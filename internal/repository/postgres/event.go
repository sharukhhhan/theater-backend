package postgres

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
	"theater/internal/models"
)

type EventPostgres struct {
	db *sqlx.DB
}

func NewEventPostgres(db *sqlx.DB) *EventPostgres {
	return &EventPostgres{db: db}
}

func (r *EventPostgres) CreateEvent(e *models.Event) (uint, error) {
	query := `INSERT INTO events (name, description, date, price) VALUES ($1, $2, $3, $4) RETURNING id`
	var id uint
	err := r.db.QueryRow(query, e.Name, e.Description, e.Date, e.Price).Scan(&id)
	if err != nil {
		if ok, _ := IsDuplicateKeyError(err); ok {
			return 0, errors.New(ErrEventNameAlreadyExists)
		}
		return 0, err
	}

	return id, nil
}

func (r *EventPostgres) GetAllEvents() (*[]models.Event, error) {
	query := `SELECT * FROM events`
	var events []models.Event
	err := r.db.Select(&events, query)
	if err != nil {
		return nil, err
	}

	return &events, nil
}

func (r *EventPostgres) GetEventByID(id uint) (*models.Event, error) {
	query := `SELECT * FROM events WHERE id = $1`
	var event models.Event
	err := r.db.Get(&event, query, id)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *EventPostgres) UpdateEvent(e *models.Event) error {
	// Start with the base update statement
	query := `UPDATE events SET`
	args := []interface{}{}
	updateFields := []string{}

	// Check each field for non-zero values and add them to the query
	if e.Name != "" {
		updateFields = append(updateFields, " name = $"+strconv.Itoa(len(args)+1))
		args = append(args, e.Name)
	}
	if e.Date != "" {
		updateFields = append(updateFields, " date = $"+strconv.Itoa(len(args)+1))
		args = append(args, e.Date)
	}
	if e.Description != "" {
		updateFields = append(updateFields, " description = $"+strconv.Itoa(len(args)+1))
		args = append(args, e.Description)
	}
	if e.Price != 0 {
		updateFields = append(updateFields, " price = $"+strconv.Itoa(len(args)+1))
		args = append(args, e.Price)
	}

	// If no fields to update, return early
	if len(updateFields) == 0 {
		return nil
	}

	// Add the fields to the query
	query += strings.Join(updateFields, ",")
	query += " WHERE id = $" + strconv.Itoa(len(args)+1)
	args = append(args, e.ID)

	// Prepare and execute the query
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	result, err := stmt.Exec(args...)
	if err != nil {
		if ok, _ := IsDuplicateKeyError(err); ok {
			return errors.New(ErrEventNameAlreadyExists)
		}
		return err
	}

	if rowsAffected, err := result.RowsAffected(); rowsAffected == 0 {
		if err != nil {
			return fmt.Errorf("error while counting rows affecting: %v", err)
		}
		return errors.New(ErrEventNotFound)
	}

	return nil
}

func (r *EventPostgres) DeleteEvent(id uint) error {
	query := `DELETE FROM events WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	if rowsAffected, err := result.RowsAffected(); rowsAffected == 0 {
		if err != nil {
			return fmt.Errorf("error while counting rows affecting: %v", err)
		}
		return errors.New(ErrEventNotFound)
	}

	return nil
}
