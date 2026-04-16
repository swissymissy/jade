package handlers

import (
	"database/sql"
)

type Product struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name"`
	Slug        string         `json:"slug"`
	Type        string         `json:"type"`
	Price       float64        `json:"price"`
	Quantity    int64          `json:"quantity"`
	Description sql.NullString `json:"description"`
	IsAvailable int64          `json:"is_available"`
	VideoUrl    sql.NullString `json:"video"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
}
