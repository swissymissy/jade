package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/swissymissy/jade/database"
)

type ApiConfig struct {
	port      string
	platform  string
	db        *database.Queries
	jwtSecret string
	s3Bucket  string
	s3Region  string
	s3Client  *s3.Client
}

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
	log.Println("Connected to database")

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	s3Bucket := os.Getenv("S3_BUCKET")
	if s3Bucket == "" {
		log.Fatal("S3_BUCKET environment variable is not set")
	}

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
	apicfg := ApiConfig{
		port: port,
		platform: platform,
		db: db,
		jwtSecret: jwtSecret,
		s3Bucket: s3Bucket,
		s3Region: s3Region,
		s3Client: s3NewClient,
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

	// TODO: register handlers here

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
