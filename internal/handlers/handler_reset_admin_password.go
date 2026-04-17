package handlers 

import (
	"net/http"
	"log"

	"github.com/swissymissy/jade/internal/auth"
	"github.com/swissymissy/jade/internal/database"
)

func (apicfg *ApiConfig) HandlerResetPassword(w http.ResponseWriter, r *http.Request) {
	// decode request 
	var req PasswordResetRequest
	err := DecodeRequest(r, &req)
	if err != nil {
		log.Printf("Error decoding request: %s", err)
		ResponseWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if req.RecoveryCode == "" || req.NewPassword == "" {
		ResponseWithError(w, http.StatusBadRequest, "Recovery code and new password are required")
		return
	}

	// get admin info
	admin, err := apicfg.DB.GetAdmin(r.Context())
	if err != nil {
		log.Printf("Error getting admin: %s", err)
		ResponseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// validate the recovery code
	match, err := auth.CheckPasswordHash(req.RecoveryCode, admin.RecoveryHash)
	if err != nil || !match {
		ResponseWithError(w, http.StatusUnauthorized, "Invalid recovery code")
		return
	}

	// hash new password
	newPasswordHash, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		log.Printf("Error hashing new password: %s", err)
		ResponseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// generate new recovery code after the old one used
	newRecoveryCode, err := auth.GenerateRecoveryCode()
	if err != nil {
		log.Printf("Error generating new recovery code: %s", err)
		ResponseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// hash new recovery code to store in db
	newRecoveryHash, err := auth.HashPassword(newRecoveryCode)
	if err != nil {
		log.Printf("Error hashing new recovery code: %s", err)
		ResponseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// update password in database
	err = apicfg.DB.UpdateAdminPassword(r.Context(), database.UpdateAdminPasswordParams{
		PasswordHash: newPasswordHash,
		ID:           admin.ID,
	})
	if err != nil {
		log.Printf("Error updating password: %s", err)
		ResponseWithError(w, http.StatusInternalServerError, "Failed to update password")
		return
	}

	// update recovery code hash
	err = apicfg.DB.UpdateAdminRecoveryHash(r.Context(), database.UpdateAdminRecoveryHashParams{
		RecoveryHash: newRecoveryHash,
		ID:           admin.ID,
	})
	if err != nil {
		log.Printf("Error updating recovery hash: %s", err)
		ResponseWithError(w, http.StatusInternalServerError, "Failed to update recovery code")
		return
	}

	// log and response
	log.Printf("Admin %s reset their password", admin.Email)
	ResponseWithJSON(w, http.StatusOK, PasswordResetResponse{
		Message:      "Password reset successfully",
		RecoveryCode: newRecoveryCode,
	})
}