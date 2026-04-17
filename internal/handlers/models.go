package handlers

import (
	"database/sql"
	"github.com/google/uuid"
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

type AdminLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password'`
}

type LoginAdmin struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Token     string    `json:"token"`
}

type ProductCreateRequest struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Price       float64 `json:"price"`
	Quantity    int64   `json:"quantity"`
	Description string  `json:"description"`
}

type ProductResponse struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name"`
	Slug        string         `json:"slug"`
	Type        string         `json:"type"`
	Price       float64        `json:"price"`
	Quantity    int64          `json:"quantity"`
	Description sql.NullString `json:"description"`
	VideoUrl    sql.NullString `json:"video_url"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
}

type UploadedImage struct {
	ID        int64  `json:"id"`
	ProductID int64  `json:"product_id"`
	S3Key     string `json:"s3_key"`
	CreatedAt string `json:"created_at"`
}

type ProductUpdateRequest struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Price       float64 `json:"price"`
	Quantity    int64   `json:"quantity"`
	Description string  `json:"description"`
	IsAvailable int64   `json:"is_available"`
}