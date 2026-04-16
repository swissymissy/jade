package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/swissymissy/jade/internal/database"
)

// create new product in database
func (apicfg *ApiConfig) HandlerCreateProduct(w http.ResponseWriter, r *http.Request) {

	// decode request
	var newItem ProductCreateRequest
	err := DecodeRequest(r, &newItem)
	if err != nil {
		log.Printf("Error decoding request: %s\n", err)
		ResponseWithError(w, http.StatusBadRequest, "Can't create new item")
		return
	}

	// validate fields
	if newItem.Name == "" || newItem.Type == "" {
		ResponseWithError(w, http.StatusBadRequest, "Name or Type should not be empty")
		return
	}
	if newItem.Price <= 0.0 {
		ResponseWithError(w, http.StatusBadRequest, "Price should be greater than 0")
		return
	}
	if newItem.Quantity < 0 {
		ResponseWithError(w, http.StatusBadRequest, "Quantity should not be less than 0")
		return
	}
	if newItem.Description == "" {
		newItem.Description = "No description"
	}

	// generate slug from give name
	slug := strings.Join(strings.Fields(strings.ToLower(strings.TrimSpace(newItem.Name))), "-")

	// insert into database
	item, err := apicfg.DB.CreateProduct(r.Context(), database.CreateProductParams{
		Name:        newItem.Name,
		Slug:        slug,
		Type:        newItem.Type,
		Price:       newItem.Price,
		Quantity:    newItem.Quantity,
		Description: ToNullString(newItem.Description),
		VideoUrl:    sql.NullString{String: "", Valid: false},
	})
	if err != nil {
		log.Printf("Error creating product: %s", err)
		ResponseWithError(w, http.StatusInternalServerError, "Failed to create product")
		return
	}

	ResponseWithJSON(w, http.StatusCreated, ProductResponse{
		ID:          item.ID,
		Name:        item.Name,
		Slug:        item.Slug,
		Type:        item.Type,
		Price:       item.Price,
		Quantity:    item.Quantity,
		Description: item.Description,
		VideoUrl:    item.VideoUrl,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	})
}
