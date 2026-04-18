package handlers

import (
	"net/http"
	"strconv"
	"fmt"
)

func (apicfg *ApiConfig) HandlerGetAllProductsAdmin(w http.ResponseWriter, r *http.Request) {
	// get limit from url and sanitize input
	limitStr := SanitizeString(r.URL.Query().Get("limit"))
	limit := 100
	if limitStr != "" {
		parsed, err := strconv.Atoi(limitStr)
		if err != nil {
			fmt.Printf("Error converting string to int: %s\n", err)
			ResponseWithError(w, http.StatusBadRequest, "Invalid limit parameter")
			return
		}
		limit = parsed
	}

	productList := make([]AdminProductListing, 0, limit)

	list, err := apicfg.DB.GetAllProductsAdmin(r.Context(), int64(limit))
	if err != nil {
		fmt.Printf("Error getting admin products: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "Failed to fetch products")
		return
	}

	for _, p := range list {
		var coverImage *ProductImage
		cover, err := apicfg.DB.GetCoverImageByProductID(r.Context(), p.ID)
		if err == nil {
			coverImage = &ProductImage{
				ID:        cover.ID,
				ProductID: cover.ProductID,
				S3Key:     cover.S3Key,
				ImageURL:  apicfg.publicAssetURL(cover.S3Key),
				Cover:     cover.Cover,
				CreatedAt: cover.CreatedAt,
			}
		}

		images, err := apicfg.DB.GetImagesByProductID(r.Context(), p.ID)
		if err != nil {
			fmt.Printf("Error getting images for product %d: %s\n", p.ID, err)
			ResponseWithError(w, http.StatusInternalServerError, "Failed to fetch product images")
			return
		}
		imageList := make([]ProductImage, 0, len(images))
		for _, img := range images {
			imageList = append(imageList, ProductImage{
				ID:        img.ID,
				ProductID: img.ProductID,
				S3Key:     img.S3Key,
				ImageURL:  apicfg.publicAssetURL(img.S3Key),
				Cover:     img.Cover,
				CreatedAt: img.CreatedAt,
			})
		}

		productList = append(productList, AdminProductListing{
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
			CoverImage:  coverImage,
			Images:      imageList,
		})
	}

	ResponseWithJSON(w, http.StatusOK, productList)
}