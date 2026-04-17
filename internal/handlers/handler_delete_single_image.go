package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/swissymissy/jade/internal/storage"
)

func (apicfg *ApiConfig) HandlerDeleteImage(w http.ResponseWriter, r *http.Request) {
	// get image ID from URL
	imageIDStr := r.PathValue("id")
	imageID, err := strconv.ParseInt(imageIDStr, 10, 64)
	if err != nil {
		ResponseWithError(w, http.StatusBadRequest, "Invalid image ID")
		return
	}

	// get image record to get s3 key
	image, err := apicfg.DB.GetImageByID(r.Context(), imageID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ResponseWithError(w, http.StatusNotFound, "Image not found")
			return
		}
		log.Printf("Error getting image: %s", err)
		ResponseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// delete image from s3
	err = storage.DeleteFromS3(r.Context(), apicfg.S3Client, apicfg.S3Bucket, image.S3Key)
	if err != nil {
		log.Printf("Error deleting image from S3: %s", err)
		// no return, continue — still remove from database
	}

	// delete image from databse
	err = apicfg.DB.DeleteImage(r.Context(), imageID)
	if err != nil {
		log.Printf("Error deleting image from database: %s", err)
		ResponseWithError(w, http.StatusInternalServerError, "Failed to delete image")
		return
	}

	log.Printf("Image %d deleted: %s", imageID, image.S3Key)
	ResponseWithJSON(w, http.StatusOK, map[string]string{
		"message": "Image deleted successfully",
	})
}
