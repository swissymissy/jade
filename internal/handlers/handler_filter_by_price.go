package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/swissymissy/jade/internal/database"
)

// handler to help users filter items by price range min, max
func (apicfg *ApiConfig) HandlerFilterByPrice(w http.ResponseWriter, r *http.Request) {
	// default min and max
	min := 0.0
	max := 999999.99

	// price range from url
	minStr := r.URL.Query().Get("min")
	if minStr != "" {
		parsed, err := strconv.ParseFloat(minStr, 64)
		if err != nil {
			ResponseWithError(w, http.StatusBadRequest, "Invalid min price")
			return
		}
		min = parsed
	}

	maxStr := r.URL.Query().Get("max")
	if maxStr != "" {
		parsed, err := strconv.ParseFloat(maxStr, 64)
		if err != nil {
			ResponseWithError(w, http.StatusBadRequest, "Invalid max price")
			return
		}
		max = parsed
	}

	if min > max {
		ResponseWithError(w, http.StatusBadRequest, "Min price can't be greater than max price")
		return
	}

	// getting items from database
	limit := 20
	rows, err := apicfg.DB.FilterProductByPrice(r.Context(), database.FilterProductByPriceParams{
		Price:   min,
		Price_2: max,
		Limit:   int64(limit),
	})
	if err != nil {
		fmt.Printf("Error filtering items: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "Error fetching products")
		return
	}

	// writing items to response format
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
	ResponseWithJSON(w, http.StatusOK, itemList)
}
