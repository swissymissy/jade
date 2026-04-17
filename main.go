package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"

	"github.com/swissymissy/jade/internal/database"
	"github.com/swissymissy/jade/internal/handlers"
	"github.com/swissymissy/jade/internal/middleware"
)

func main() {

	godotenv.Load()

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM environment variable is not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DB_PATH")
	if dbURL == "" {
		log.Fatal("DB_URL should be set")
	}
	// connect to database
	db, err := sql.Open("sqlite", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database", err)
	}
	defer db.Close()
	// ping to verify db connection
	if err := db.Ping(); err != nil {
		log.Fatal("Cannot connect to database")
	}
	dbQuery := database.New(db)
	log.Println("Connected to database")

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	// get s3 bucket name
	s3Bucket := os.Getenv("S3_BUCKET")
	if s3Bucket == "" {
		log.Fatal("S3_BUCKET environment variable is not set")
	}

	// get s3 region
	s3Region := os.Getenv("S3_REGION")
	if s3Region == "" {
		log.Fatal("S3_REGION environment variable is not set")
	}

	awsCfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(s3Region))
	if err != nil {
		log.Fatal("Failed to load AWS config: %v", err)
	}

	s3NewClient := s3.NewFromConfig(awsCfg)

	// initialize apiconfig
	apicfg := handlers.ApiConfig{
		Port:      port,
		Platform:  platform,
		DB:        dbQuery,
		JWTSecret: jwtSecret,
		S3Bucket:  s3Bucket,
		S3Region:  s3Region,
		S3Client:  s3NewClient,
	}

	// serve mux
	mux := http.NewServeMux()

	// create new http server
	address := fmt.Sprintf(":%s", port)
	jadeServer := http.Server{
		Addr:    address,
		Handler: mux,
	}

	// create handler
	fileServer := http.FileServer(http.Dir("./frontend/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// register handlers
	// public routes
	mux.HandleFunc("GET /api/products", apicfg.HandlerGetAllProducts)
	mux.HandleFunc("GET /api/products/{id}", apicfg.HandlerGetOneProduct)
	mux.HandleFunc("GET /api/products/search", apicfg.HandlerSearchProduct)
	mux.HandleFunc("GET /api/products/filter", apicfg.HandlerFilterByPrice)

	// admind routes - protected
	mux.HandleFunc("POST /api/admin/products", middleware.AuthRequired(apicfg.HandlerCreateProduct, apicfg.JWTSecret))
	mux.HandleFunc("POST /api/admin/products/{id}/images", middleware.AuthRequired(apicfg.HandlerUploadImages, apicfg.JWTSecret))
	mux.HandleFunc("POST /api/admin/products/{id}/video", middleware.AuthRequired(apicfg.HandlerUploadVideo, apicfg.JWTSecret))
	mux.HandleFunc("DELETE /api/admin/products/{id}", middleware.AuthRequired(apicfg.HandlerDeleteProduct, apicfg.JWTSecret))
	mux.HandleFunc("PUT /api/admin/products/{id}", middleware.AuthRequired(apicfg.HandlerUpdateProduct, apicfg.JWTSecret))
	mux.HandleFunc("DELETE /api/admin/images/{id}", middleware.AuthRequired(apicfg.HandlerDeleteImage, apicfg.JWTSecret))

	// Auth
	mux.HandleFunc("POST /api/admin/register", apicfg.HandlerCreateAdmin)
	mux.HandleFunc("POST /api/admin/login", apicfg.AdminLogin)
	mux.HandleFunc("POST /api/admin/reset-password", apicfg.HandlerResetPassword)

	// run server in background
	go func() {
		fmt.Printf("Serving on: http://localhost:%s/static/\n", port)
		if err := jadeServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %s\n", err)
		}
	}()

	// blocks until OS sends SIGTERM or SIGINT
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")

	// give in-flight requests up to 10s to finish
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := jadeServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error. Forced shutdown: %s\n", err)
	}
	log.Println("Graceful shutdown complete")
}
