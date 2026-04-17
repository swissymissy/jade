package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/swissymissy/jade/internal/storage"
)

func (apicfg *ApiConfig) HandlerDeleteProduct(w http.ResponseWriter, r *http.Request) {
	// get item ID from URL
	itemIDStr := r.PathValue("id")
	itemID, err := strconv.ParseInt(itemIDStr, 10, 64)
	if err != nil {
		ResponseWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	// check product exists
	product, err := apicfg.DB.GetProductByID(r.Context(), itemID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ResponseWithError(w, http.StatusNotFound, "Product not found")
			return
		}
		log.Printf("Error getting product: %s", err)
		ResponseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// delete all images from S3
	// 1. get all images from databse
	images, err := apicfg.DB.GetImagesByProductID(r.Context(), itemID)
	if err != nil {
		log.Printf("Error fetching images: %s", err)
		// no return — still try to delete the product
	}
	// 2. loop deleting all images in S3
	for _, img := range images {
		err = storage.DeleteFromS3(r.Context(), apicfg.S3Client, apicfg.S3Bucket, img.S3Key)
		if err != nil {
			log.Printf("Error deleting image %s from S3: %s", img.S3Key, err)
			// continue deleting the rest
		}
	}

	// delete video from S3 if exists
	if product.VideoUrl.Valid && product.VideoUrl.String != "" {
		err = storage.DeleteFromS3(r.Context(), apicfg.S3Client, apicfg.S3Bucket, product.VideoUrl.String)
		if err != nil {
			log.Printf("Error deleting video from S3: %s", err)
		}
	}

	// finally, delete the product from database
	err = apicfg.DB.DeleteProduct(r.Context(), itemID)
	if err != nil {
		log.Printf("Error deleting product: %s", err)
		ResponseWithError(w, http.StatusInternalServerError, "Failed to delete product")
		return
	}

	log.Printf("Product %d deleted with all associated files", itemID)
	ResponseWithJSON(w, http.StatusOK, map[string]string{
		"message": "Product deleted successfully",
	})

}