package handlers 

import (
	"net/http"
	"fmt"
	"log"
	"strconv"
	"errors"
	"database/sql"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/swissymissy/jade/internal/database"
	"github.com/swissymissy/jade/internal/storage"
)

// handler to let admin upload images
func (apicfg *ApiConfig) HandlerUploadImages(w http.ResponseWriter, r *http.Request) {

	// get product ID from URL
	itemIDStr := r.PathValue("id")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		log.Printf("Error converting ID string to int64: %s\n", err)
		ResponseWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	// check if item exists
	_, err = apicfg.DB.GetProductByID(r.Context(), itemID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ResponseWithError(w, http.StatusNotFound, "Item not found")
			return
		}
		log.Printf("Error getting product by ID from database: %s\n", err)
		ResponseWithError(w, http.StatusNotFound, "Item not found")
		return
	}

	// check images count 
	count, err := apicfg.DB.CountImagesByProductID(r.Context(), itemID)
	if err != nil {
		log.Printf("Error counting images: %s", err)
		ResponseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	// check amount of imgs
	if count >= 5 {
		ResponseWithError(w, http.StatusBadRequest, "Maximum 5 images per product")
		return
	}

	// read file from form
	const maxSize = 10 << 20 //10MB
	r.Body = http.MaxBytesReader(w, r.Body, maxSize)

	// parse file multiparts
	err = r.ParseMultipartForm(maxSize)
	if err != nil {
		log.Printf("Error parsing form: %s", err)
		ResponseWithError(w, http.StatusBadRequest, "File too large or invalid form")
		return
	}

	// retrieve the uploaded image from request
	file, header, err := r.FormFile("image")
	if err != nil {
		log.Printf("Error reading file: %s", err)
		ResponseWithError(w, http.StatusBadRequest, "No image file provided")
		return
	}
	defer file.Close() // close file to prevent data leak

	// validate file type
	extension := strings.ToLower(filepath.Ext(header.Filename))
	allowedTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".webp": "image/webp",
	}

	contentType, ok := allowedTypes[extension]
	if !ok {
		ResponseWithError(w, http.StatusBadRequest, "Only jpg, png, and webp images are allowed")
		return
	}

	// create unique filename
	filename, err := storage.RandomFilename(extension)
	if err != nil {
		log.Printf("Error generating filename: %s", err)
		ResponseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// s3 key
	s3Key := fmt.Sprintf("products/%d/%s", itemID, filename)

	// upload to s3
	
}