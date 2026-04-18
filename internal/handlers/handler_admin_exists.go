package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
)

func (apicfg *ApiConfig) HandlerAdminExists(w http.ResponseWriter, r *http.Request) {
	// check if admin already created
	_, err := apicfg.DB.GetAdmin(r.Context())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ResponseWithJSON(w, 200, struct {
				Exists bool `json:"exists"`
			}{
				Exists: false,
			})
			return
		}
		log.Printf("Error fetching admin table: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	ResponseWithJSON(w, 200, struct {
		Exists bool `json:"exists"`
	}{
		Exists: true,
	})
}
