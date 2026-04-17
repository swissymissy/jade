package handlers

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/swissymissy/jade/internal/auth"
	"github.com/swissymissy/jade/internal/database"
)

// register new admin.
// there can be only one admin in total.
func (apicfg *ApiConfig) HandlerCreateAdmin(w http.ResponseWriter, r *http.Request) {
	// decode request
	var req AdminCreateRequest
	err := DecodeRequest(r, &req)
	if err != nil {
		log.Printf("Error decoding request: %s", err)
		ResponseWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if req.Email == "" || req.Password == "" {
		ResponseWithError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	if len(req.Password) < 8 {
		ResponseWithError(w, http.StatusBadRequest, "Password must be at least 8 characters")
		return
	}

	// hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		ResponseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	//generate a recovery code
	recoveryCode, err := auth.GenerateRecoveryCode()
	if err != nil {
		log.Printf("Error generating recovery code: %s", err)
		ResponseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	// hash recovery code
	recoveryHash, err := auth.HashPassword(recoveryCode)
	if err != nil {
		log.Printf("Error hashing recovery code: %s", err)
		ResponseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// create new id
	adminID := uuid.New().String()

	// save admin to database
	admin, err := apicfg.DB.CreateAdmin(r.Context(), database.CreateAdminParams{
		ID:           adminID,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		RecoveryHash: recoveryHash,
	})
	if err != nil {
		log.Printf("Error creating admin: %s", err)
		ResponseWithError(w, http.StatusConflict, "Admin already exists")
		return
	}

	// log and response
	log.Printf("Admin created: %s", admin.Email)
	ResponseWithJSON(w, http.StatusCreated, AdminResponse{
		ID:           admin.ID,
		Email:        admin.Email,
		RecoveryCode: recoveryCode,
		CreatedAt:    admin.CreatedAt,
		UpdatedAt:    admin.UpdatedAt,
	})
}
