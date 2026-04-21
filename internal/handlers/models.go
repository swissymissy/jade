package handlers

import (
	"database/sql"
	"github.com/google/uuid"
)

// =====Admin===================
type AdminLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
}

type AdminProductListing struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name"`
	Slug        string         `json:"slug"`
	Type        string         `json:"type"`
	Price       float64        `json:"price"`
	Quantity    int64          `json:"quantity"`
	Description sql.NullString `json:"description"`
	About       sql.NullString `json:"about"`
	IsAvailable int64          `json:"is_available"`
	VideoUrl    sql.NullString `json:"video_url"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
	CoverImage  *ProductImage  `json:"cover_image"`
	Images      []ProductImage `json:"images"`
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
	About       sql.NullString `json:"about"`
	IsAvailable int64          `json:"is_available"`
	VideoUrl    sql.NullString `json:"video"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
}

type ProductDetail struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name"`
	Slug        string         `json:"slug"`
	Type        string         `json:"type"`
	Price       float64        `json:"price"`
	Quantity    int64          `json:"quantity"`
	Description sql.NullString `json:"description"`
	About       sql.NullString `json:"about"`
	IsAvailable int64          `json:"is_available"`
	VideoUrl    string         `json:"video_url"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
	Images      []ProductImage `json:"images"`
}

type ProductListing struct {
	ID          int64         `json:"id"`
	Name        string        `json:"name"`
	Slug        string        `json:"slug"`
	Type        string        `json:"type"`
	Price       float64       `json:"price"`
	Quantity    int64         `json:"quantity"`
	IsAvailable int64         `json:"is_available"`
	CoverImage  *ProductImage `json:"cover_image"`
}

type ProductCreateRequest struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Price       float64 `json:"price"`
	Quantity    int64   `json:"quantity"`
	Description string  `json:"description"`
	About       string  `json:"about"`
}

type UploadedImage struct {
	ID        int64  `json:"id"`
	ProductID int64  `json:"product_id"`
	S3Key     string `json:"s3_key"`
	ImageURL  string `json:"image_url"`
	CreatedAt string `json:"created_at"`
}

type ProductImage struct {
	ID        int64  `json:"id"`
	ProductID int64  `json:"product_id"`
	S3Key     string `json:"s3_key"`
	ImageURL  string `json:"image_url"`
	Cover     int64  `json:"cover"`
	CreatedAt string `json:"created_at"`
}

type ProductUpdateRequest struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Price       float64 `json:"price"`
	Quantity    int64   `json:"quantity"`
	Description string  `json:"description"`
	About       string  `json:"about"`
	IsAvailable int64   `json:"is_available"`
}

// ====Password Reset=====================
type PasswordResetRequest struct {
	RecoveryCode string `json:"recovery_code"`
	NewPassword  string `json:"new_password"`
}

type PasswordResetResponse struct {
	Message      string `json:"message"`
	RecoveryCode string `json:"recovery_code"`
}
