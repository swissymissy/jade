package handlers

import (
	"log"
	"net/http"
)

func (apicfg *ApiConfig) HandlerReset(w http.ResponseWriter, r *http.Request) {
	// check if "dev"
	if apicfg.Platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err := apicfg.DB.ResetAdmins(r.Context())
	if err != nil {
		log.Prinf("Error reseting admin table: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "Unable to reset admin table")
		return
	}
	w.WriteHeader(200)
}