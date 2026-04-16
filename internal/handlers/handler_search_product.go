package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/swissymissy/jade/internal/database"
)

// handler to let user search for products based on the url params
func (apicfg *ApiConfig) HandlerSearchProduct(w http.ResponseWriter, r *http.Request) {
	// param values from url
	query := r.URL.Query().Get("q")
	if query == "" {
		ResponseWithError(w, http.StatusBadRequest, "Search can't be empty")
		return
	}

	limit := 20
	// search for items in database
	rows, err := apicfg.DB.SearchProduct(r.Context(), database.SearchProductParams{
		Column1: query,
		Column2: query,
		Column3: query,
		Column4: query,
		Limit:   int64(limit),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Printf("Error fetching items from db: %s\n", err)
			ResponseWithError(w, http.StatusNotFound, "Items not found")
			return
		}
		fmt.Printf("Error fetching items from db: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "Error finding items")
		return
	}

	// writing each item to response format
	itemList := make([]Product, 0, limit)
	for _, p := range rows {
		itemList = append(itemList, Product{
			ID:          p.ID,
			Name:        p.Name,
			Slug:        p.Slug,
			Type:        p.Type,
			Price:       p.Price,
			Quantity:    p.Quantity,
			Description: p.Description,
			IsAvailable: p.IsAvailable,
			VideoUrl:    p.VideoUrl,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		})
	}

	// repsonse to client
	ResponseWithJSON(w, http.StatusOK, itemList)
}
