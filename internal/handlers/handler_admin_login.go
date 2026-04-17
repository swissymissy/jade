package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/swissymissy/jade/internal/auth"
)

// admin login handler
func (apicfg *ApiConfig) AdminLogin(w http.ResponseWriter, r *http.Request) {

	// decode request
	var adminLogin AdminLoginRequest
	err := DecodeRequest(r, &adminLogin)
	if err != nil {
		log.Printf("Error decoding request: %s\n", err)
		ResponseWithError(w, 400, "invalid request")
		return
	}

	// find admin by email
	admin, err := apicfg.DB.GetAdminByEmail(r.Context(), adminLogin.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Login attempt with unknown email: %s\n", err)
			ResponseWithError(w, 401, "Incorrect Email or Password")
			return
		}
		log.Printf("Error getting admin from db: %s\n", err)
		ResponseWithError(w, http.StatusUnauthorized, "Incorrect Email or Password ")
		return
	}

	// check password
	match, err := auth.CheckPasswordHash(adminLogin.Password, admin.PasswordHash)
	if err != nil {
		log.Printf("%s\n", err)
		ResponseWithError(w, http.StatusUnauthorized, "Incorrect Email or Password")
		return
	}
	if !match {
		ResponseWithError(w, http.StatusUnauthorized, "Incorrect Email or Password")
		return
	}

	// create new token for new login session
	adminID, err := uuid.Parse(admin.ID)
	if err != nil {
		log.Printf("Error converting ID string to UUID: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	token, err := auth.MakeJWT(adminID, apicfg.JWTSecret)
	if err != nil {
		log.Printf("Error making new access token: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	
	// set cookie 
	http.SetCookie(w, &http.Cookie{
		Name: "jade_session",
		Value: token,
		Path: "/",
		HttpOnly: true,
		Secure: apicfg.Platform == "prod",
		SameSite: http.SameSiteLaxMode,
		MaxAge: 60*60*24, // 24h
	})

	log.Printf("Admin %s has logged in\n", admin.Email)

	ResponseWithJSON(w, http.StatusOK, LoginAdmin{
		ID:        adminID,
		Email:     admin.Email,
		CreatedAt: admin.CreatedAt,
		UpdatedAt: admin.UpdatedAt,
	})

}
