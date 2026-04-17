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

// handler to let admin upload images
func (apicfg *ApiConfig) HandlerUploadImages(w http.ResponseWriter, r *http.Request) {

	// get product ID from URL
	itemIDStr := r.PathValue("id")
	itemID, err := strconv.ParseInt(itemIDStr, 10, 64)
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

	// handle uploading image
	// 1. read file from form
	const maxSize = 10 << 20 //10MB
	r.Body = http.MaxBytesReader(w, r.Body, maxSize)

	// parse file multiparts
	err = r.ParseMultipartForm(maxSize)
	if err != nil {
		log.Printf("Error parsing form: %s", err)
		ResponseWithError(w, http.StatusBadRequest, "File too large or invalid form")
		return
	}

	// 2. retrieve the uploaded image from request
	file, header, err := r.FormFile("image")
	if err != nil {
		log.Printf("Error reading file: %s", err)
		ResponseWithError(w, http.StatusBadRequest, "No image file provided")
		return
	}
	defer file.Close() // close file to prevent data leak

	// 3. validate file type
	extension := strings.ToLower(filepath.Ext(header.Filename))
	allowedTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".webp": "image/webp",
	}
	// 4. check if the extension is valid type
	contentType, ok := allowedTypes[extension]
	if !ok {
		ResponseWithError(w, http.StatusBadRequest, "Only jpg, png, and webp images are allowed")
		return
	}

	// 5. create unique filename
	filename, err := storage.RandomFilename(extension)
	if err != nil {
		log.Printf("Error generating filename: %s", err)
		ResponseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// 6. create s3 key
	s3Key := fmt.Sprintf("products/%d/%s", itemID, filename)

	// 7. upload to s3
	err = storage.UploadToS3(r.Context(), apicfg.S3Client, apicfg.S3Bucket, s3Key, contentType, file)
	if err != nil {
		log.Printf("Error uploading to S3: %s", err)
		ResponseWithError(w, http.StatusInternalServerError, "Failed to upload image")
		return
	}

	// set cover photo
	cover := int64(0)
	if count == 0 {
		cover = 1 // first image is automatically the cover
	}

	// save to database
	image, err := apicfg.DB.CreateProductImage(r.Context(), database.CreateProductImageParams{
		ProductID: itemID,
		S3Key:     s3Key,
		Cover:     cover,
	})
	if err != nil {
		log.Printf("Error saving image to database: %s", err)
		ResponseWithError(w, http.StatusInternalServerError, "Failed to save image record")
		return
	}

	log.Printf("Image uploaded for product %d: %s", itemID, s3Key)
	ResponseWithJSON(w, http.StatusCreated, UploadedImage{
		ID:        image.ID,
		ProductID: image.ProductID,
		S3Key:     image.S3Key,
		CreatedAt: image.CreatedAt,
	})

}
