package models

type Club struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" db:"name" validate:"required"`
	Description string  `json:"description" db:"description"`
	Contact     string  `json:"contact" db:"contacts"`
	Price       float64 `json:"price" db:"price" validate:"twoDecimalPlaces"`
	SpotsNumber int     `json:"spots_number" db:"spots_number" validate:"number,min=1"`
	IsActive    string  `json:"is_active" db:"is_active"`
}
