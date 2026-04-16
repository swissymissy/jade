package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// get one single product based on product's ID
func (apicfg *ApiConfig) HandlerGetOneProduct(w http.ResponseWriter, r *http.Request) {
	// get product id from url
	productIDStr := r.PathValue("id")
	if productIDStr == "" {
		ResponseWithError(w, http.StatusBadRequest, "ID can't be empty")
		return
	}
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		ResponseWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	// get product information from database
	product, err := apicfg.DB.GetProductByID(r.Context(), int64(productID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Printf("Error finding product in database: %s\n", err)
			ResponseWithError(w, http.StatusNotFound, "Product not found")
			return
		}
		fmt.Printf("Error fetching product from database: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "Error fetching product")
		return
	}

	ResponseWithJSON(w, http.StatusOK, Product{
		ID:          product.ID,
		Name:        product.Name,
		Slug:        product.Slug,
		Type:        product.Type,
		Price:       product.Price,
		Quantity:    product.Quantity,
		Description: product.Description,
		IsAvailable: product.IsAvailable,
		VideoUrl:    product.VideoUrl,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	})
}
