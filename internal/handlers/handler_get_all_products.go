package handlers

import (
	"fmt"
	"net/http"
	"strconv"
)

// function to handle get all product request
func (apicfg *ApiConfig) HandlerGetAllProducts(w http.ResponseWriter, r *http.Request) {

	// pagination
	limitStr := SanitizeString(r.URL.Query().Get("limit"))
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
	productList := make([]ProductListing, 0, limit)

	// get products from database
	list, err := apicfg.DB.GetAllProducts(r.Context(), int64(limit))
	if err != nil {
		fmt.Printf("Error getting products: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "Failed to fetch products")
		return
	}

	// write to response format
	for _, p := range list {

		// if a product has no image, it serializes as "cover_image: null"
		// in JSON instead of an empty object. Clean for frontend to handle
		var coverImage *ProductImage
		cover, err := apicfg.DB.GetCoverImageByProductID(r.Context(), p.ID)
		if err == nil {
			coverImage = &ProductImage{
				ID:        cover.ID,
				ProductID: cover.ProductID,
				S3Key:     cover.S3Key,
				ImageURL: apicfg.publicAssetURL(cover.S3Key),
				Cover:     cover.Cover,
				CreatedAt: cover.CreatedAt,
			}
		}

		productList = append(productList, ProductListing{
			ID:          p.ID,
			Name:        p.Name,
			Slug:        p.Slug,
			Type:        p.Type,
			Price:       p.Price,
			Quantity:    p.Quantity,
			IsAvailable: p.IsAvailable,
			CoverImage:  coverImage,
		})
	}

	ResponseWithJSON(w, http.StatusOK, productList)

}
