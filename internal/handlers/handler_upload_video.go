package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/swissymissy/jade/internal/database"
	"github.com/swissymissy/jade/internal/storage"
)

func (apicfg *ApiConfig) HandlerUploadVideo(w http.ResponseWriter, r *http.Request) {
	// get product ID from url
	itemIDStr := r.PathValue("id")
	itemID, err := strconv.ParseInt(itemIDStr, 10, 64)
	if err != nil {
		ResponseWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	// check if product exist
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

	// if video already exists, replace old one from S3
	if product.VideoUrl.Valid && product.VideoUrl.String != "" {
		err = storage.DeleteFromS3(r.Context(), apicfg.S3Client, apicfg.S3Bucket, product.VideoUrl.String)
		if err != nil {
			log.Printf("Error deleting old video from S3: %s", err)
			// no return — continue with upload even if delete fails
		}
	}

	// read file from form
	const maxVideoSize = 100 << 20 // 100MB
	r.Body = http.MaxBytesReader(w, r.Body, maxVideoSize)

	err = r.ParseMultipartForm(maxVideoSize)
	if err != nil {
		log.Printf("Error parsing form: %s", err)
		ResponseWithError(w, http.StatusBadRequest, "File too large or invalid form")
		return
	}

	file, header, err := r.FormFile("video")
	if err != nil {
		log.Printf("Error reading file: %s", err)
		ResponseWithError(w, http.StatusBadRequest, "No video file provided")
		return
	}
	defer file.Close() // close file to prevent data leak

	// validate file type
	extension := strings.ToLower(filepath.Ext(header.Filename))
	allowedTypes := map[string]string{
		".mp4":  "video/mp4",
		".mov":  "video/quicktime",
		".webm": "video/webm",
	}
	contentType, ok := allowedTypes[extension]
	if !ok {
		ResponseWithError(w, http.StatusBadRequest, "Only mp4, mov, and webm videos are allowed")
		return
	}

	// generate unique filename and s3Key
	filename, err := storage.RandomFilename(extension)
	if err != nil {
		log.Printf("Error generating filename: %s", err)
		ResponseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	s3Key := fmt.Sprintf("products/%d/video/%s", itemID, filename)

	// upload to S3
	err = storage.UploadToS3(r.Context(), apicfg.S3Client, apicfg.S3Bucket, s3Key, contentType, file)
	if err != nil {
		log.Printf("Error uploading video to S3: %s", err)
		ResponseWithError(w, http.StatusInternalServerError, "Failed to upload video")
		return
	}

	// update product to database
	err = apicfg.DB.UpdateProductVideoURL(r.Context(), database.UpdateProductVideoURLParams{
		VideoUrl: sql.NullString{String: s3Key, Valid: true},
		ID:       itemID,
	})
	if err != nil {
		log.Printf("Error updating video URL in database: %s", err)
		ResponseWithError(w, http.StatusInternalServerError, "Failed to save video")
		return
	}

	log.Printf("Video uploaded for product %d: %s", itemID, s3Key)
	ResponseWithJSON(w, http.StatusOK, map[string]string{
		"message":   "Video uploaded successfully",
		"video_key": s3Key,
	})

}
