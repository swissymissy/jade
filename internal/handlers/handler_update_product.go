package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/swissymissy/jade/internal/database"
)

func (apicfg *ApiConfig) HandlerUpdateProduct(w http.ResponseWriter, r *http.Request) {
	// get item ID from URL
	itemIDStr := r.PathValue("id")
	itemID, err := strconv.ParseInt(itemIDStr, 10, 64)
	if err != nil {
		ResponseWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	// check if product exists
	item, err := apicfg.DB.GetProductByID(r.Context(), itemID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ResponseWithError(w, http.StatusNotFound, "Product not found")
			return
		}
		log.Printf("Error getting product: %s", err)
		ResponseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// decode request body
	var req ProductUpdateRequest
	err = DecodeRequest(r, &req)
	if err != nil {
		log.Printf("Error decoding request: %s", err)
		ResponseWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	// validate fields
	if req.Name == "" || req.Type == "" {
		ResponseWithError(w, http.StatusBadRequest, "Name and type are required")
		return
	}
	if req.Price <= 0.0 {
		ResponseWithError(w, http.StatusBadRequest, "Price must be greater than 0")
		return
	}
	if req.Quantity < 0 {
		ResponseWithError(w, http.StatusBadRequest, "Quantity cannot be negative")
		return
	}

	// generate new slug if name changed
	slug := strings.Join(strings.Fields(strings.ToLower(strings.TrimSpace(req.Name))), "-")

	// update product in database
	updated, err := apicfg.DB.UpdateProduct(r.Context(), database.UpdateProductParams{
		Name:        req.Name,
		Slug:        slug,
		Type:        req.Type,
		Price:       req.Price,
		Quantity:    req.Quantity,
		Description: ToNullString(req.Description),
		VideoUrl:    item.VideoUrl,
		IsAvailable: req.IsAvailable,
		ID:          itemID,
	})
	if err != nil {
		log.Printf("Error updating product: %s", err)
		ResponseWithError(w, http.StatusInternalServerError, "Failed to update product")
		return
	}

	// response
	ResponseWithJSON(w, http.StatusOK, Product{
		ID:          updated.ID,
		Name:        updated.Name,
		Slug:        updated.Slug,
		Type:        updated.Type,
		Price:       updated.Price,
		Quantity:    updated.Quantity,
		Description: updated.Description,
		IsAvailable: updated.IsAvailable,
		VideoUrl:    updated.VideoUrl,
		CreatedAt:   updated.CreatedAt,
		UpdatedAt:   updated.UpdatedAt,
	})
}
