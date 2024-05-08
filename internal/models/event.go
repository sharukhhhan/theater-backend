package models

type Event struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Date        string  `json:"date" validate:"required,datetime=2006-01-02 15:04"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"twoDecimalPlaces"`
}
