package handlers 

import (
	"net/http"
	"log"
	"strconv"
	"errors"
	"database/sql"
)

// handler to let admin upload images
func (apicfg *ApiConfig) HandlerUploadImages(w http.ResponseWriter, r *http.Request) {

	// get product ID from URL
	itemIDStr := r.PathValue("id")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		log.Printf("Error converting ID string to int64: %s\n", err)
		ResponseWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	// validate id
	item, err := apicfg.DB.GetProductByID(r.Context(), itemID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ResponseWithError(w, http.StatusNotFound, "Item not found")
			return
		}
		log.Printf("Error getting product by ID from database: %s\n", err)
		ResponseWithError(w, http.StatusNotFound, "Item not found")
		return
	}

	
}