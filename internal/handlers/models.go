package handlers

import (
	"database/sql"
	"github.com/google/uuid"
)

// =====Admin===================
type AdminLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password'`
}

type AdminCreateRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AdminResponse struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	RecoveryCode string `json:"recovery_code"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type LoginAdmin struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Token     string    `json:"token"`
}

// ====Product================================
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

type ProductCreateRequest struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Price       float64 `json:"price"`
	Quantity    int64   `json:"quantity"`
	Description string  `json:"description"`
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

//====Password Reset=====================
type PasswordResetRequest struct {
	RecoveryCode string `json:"recovery_code"`
	NewPassword  string `json:"new_password"`
}

type PasswordResetResponse struct {
	Message      string `json:"message"`
	RecoveryCode string `json:"recovery_code"`
}