package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
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
	if product.IsAvailable != 1 {
		ResponseWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	// fetch product images
	images, err := apicfg.DB.GetImagesByProductID(r.Context(), int64(productID))
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

	videoURL := ""
	if product.VideoUrl.Valid {
		videoURL = apicfg.publicAssetURL(product.VideoUrl.String)
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
		VideoUrl:    videoURL,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
		Images:      imageList,
	})
}
