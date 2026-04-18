package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
)

// get product by slug
func (apicfg *ApiConfig) HandlerGetProductBySlug(w http.ResponseWriter, r *http.Request) {
	// get product slug from url
	slugStr := r.PathValue("slug")
	if slugStr == "" {
		ResponseWithError(w, http.StatusBadRequest, "slug can't be empty")
		return
	}

	// get product infor from database
	product, err := apicfg.DB.GetProductBySlug(r.Context(), slugStr)
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

	// fetch product images
	images, err := apicfg.DB.GetImagesByProductID(r.Context(), product.ID)
	if err != nil {
		log.Printf("Error fetching images: %s", err)
		// no return — product still works without images
	}

	imageList := make([]ProductImage, 0, len(images))
	for _, img := range images {
		imageList = append(imageList, ProductImage{
			ID:        img.ID,
			ProductID: img.ProductID,
			S3Key:     img.S3Key,
			ImageURL: apicfg.publicAssetURL(img.S3Key),
			Cover:     img.Cover,
			CreatedAt: img.CreatedAt,
		})
	}

	ResponseWithJSON(w, http.StatusOK, ProductDetail{
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
		Images:      imageList,
	})
}
