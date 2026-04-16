package handlers

import (
	"fmt"
	"net/http"
	"strconv"
)

// function to handle get all product request
func (apicfg *ApiConfig) HandlerGetAllProducts(w http.ResponseWriter, r *http.Request) {

	// pagination
	limitStr := r.URL.Query().Get("limit")
	limit := 20
	if limitStr != "" {
		parsed, err := strconv.Atoi(limitStr)
		if err != nil {
			fmt.Printf("Error converting string to int: %s\n", err)
			ResponseWithError(w, http.StatusBadRequest, "Invalid limit parameter")
			return
		}
		limit = parsed
	}
	productList := make([]Product, 0, limit)

	// get products from database
	list, err := apicfg.DB.GetAllProducts(r.Context(), int64(limit))
	if err != nil {
		fmt.Printf("Error getting products: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "Failed to fetch products")
		return
	}

	// write to response format
	for _, p := range list {
		productList = append(productList, Product{
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

	ResponseWithJSON(w, http.StatusOK, productList)

}
