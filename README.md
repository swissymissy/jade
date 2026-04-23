# Jadeite USA

Full-stack product showcase and admin dashboard for a small jade jewelry business, built with Go, SQLite, and AWS S3.

**Live:** [jadeiteusa.com](https://www.jadeiteusa.com)

---

## Architecture

```
                        ┌──────────────────────────────────────────┐
                        │              Go Server                   │
                        │                                          │
   Client Request ─────▶│  Middleware (JWT Auth)                   │
                        │       │                                  │
                        │       ▼                                  │
                        │  Router ──▶ Handler                      │
                        │               │                          │
                        │        ┌──────┴──────┐                   │
                        │        ▼             ▼                   │
                        │    SQLite        AWS S3                  │
                        │  (product data)  (images & videos)       │
                        │                                          │
                        └──────────────────────────────────────────┘

   Static files (HTML/CSS/JS) are served directly by the Go server.
```

The server handles everything: routing, authentication, file validation, database operations, S3 uploads, and serving the frontend. There's no separate frontend framework — the admin dashboard and public-facing pages are plain HTML, CSS, and JavaScript served by Go's built-in HTTP server.

---

## Technical Decisions

**Why server-side uploads instead of presigned URLs?**
Files go through the Go server before reaching S3 so that validation happens in a controlled environment. The server checks file size (10MB for images, 100MB for video), validates the file extension (JPEG, PNG, WebP for images; MP4, QuickTime, WebM for video), generates a unique filename, and constructs the S3 key before uploading. This prevents malicious or oversized files from ever reaching the bucket — the client can't bypass these checks.

**Why SQLite instead of PostgreSQL?**
This is a small business with a single admin user and low concurrent traffic. SQLite keeps the deployment simple — no separate database server to manage, no connection strings to configure, no extra cost. The entire database ships as a single file alongside the binary.

---

## File Upload Pipeline

When an admin uploads a product image or video, the request goes through the following flow:

1. **Auth check** — JWT middleware verifies the token before the request reaches the handler
2. **File size validation** — the server enforces a maximum upload size to reject oversized files early
3. **Extension validation** — the server checks the file extension against an allowlist:
   - Images: `.jpeg`, `.jpg`, `.png`, `.webp`
   - Videos: `.mp4`, `.mov` (QuickTime), `.webm`
4. **Filename generation** — a unique filename is generated to avoid collisions and prevent path traversal
5. **S3 upload** — the file is uploaded to the S3 bucket with the generated key
6. **Database update** — the S3 key / URL is stored in SQLite alongside the product record

Files that fail any validation step are rejected with an appropriate error response before touching S3.

---

## API Overview

The server exposes 16 RESTful endpoints, split between public and admin-protected routes.

### Public Endpoints

| Method | Route | Description |
|--------|-------|-------------|
| `GET` | `/api/products` | List all products |
| `GET` | `/api/products/{id}` | Get a single product with its media |
|`GET`| `/api/products/slug/{slug}`| get a single product with its slug |
| `GET` | `/api/products/search` | Search for specific products |
| `GET` | `/api/products/filter` | Filter products by price range |

### Admin Endpoints (JWT Protected)

| Method | Route | Description |
|--------|-------|-------------|
| `POST` | `/api/admin/products` | Create a new product |
| `GET` | `/api/admin/products` | Get list of all products in admin dashboard |
| `PUT` | `/api/admin/products/{id}` | Update product information |
| `DELETE` | `/api/admin/products/{id}` | Delete a product and its associated media |
| `POST` | `/api/admin/products/{id}/images` | Upload image(s) for a product |
| `POST` | `/api/admin/products/{id}/video` | Upload video for a product |
| `DELETE` | `/api/admin/images/{id}` | Delete a specific image |

> **Note:** The routes above are representative of the API structure. Refer to the source code for the full list of all 16 endpoints.

---

## Project Structure

```
jade/
├── frontend/            # HTML, CSS, JS served to the client
├── internal/
│   ├── auth/            # Authorization check functions
│   ├── handlers/        # HTTP handler functions for each route
│   ├── middleware/       # JWT authentication middleware
│   ├── storage/         # S3 upload/delete operations
│   └── database/        # SQLite connection and queries
├── sql/
│   ├── queries/         # SQL query definitions (used by sqlc)
│   └── schema/          # Database schema and migrations
├── main.go              # Entry point, router setup
├── sqlc.yaml            # sqlc configuration for type-safe query generation
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

---

## Running Locally

### Prerequisites

- Go 1.25+
- AWS account with an S3 bucket and IAM user credentials

### Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/swissymissy/jade.git
   cd jade
   ```

2. Create a `.env` file and fill in your environmant variables:
   ```
   PORT="8080"
   PLATFORM="dev"
   DB_PATH="./jade.db"
   BASE_URL="http://localhost"
   JWT_SECRET=your-jwt-secret
   
   AWS_ACCESS_KEY_ID=your-access-key
   AWS_SECRET_ACCESS_KEY=your-secret-key
   S3_REGION=your-region
   S3_BUCKET_NAME=your-bucket-name
   S3_BASE_URL=your-s3-base-url
   ```

4. Run the server:
   ```bash
   go run .
   ```

   Or build and run:
   ```bash
   go build -o jade && ./jade
   ```

The server starts on `http://localhost:8080` (or whichever port you've configured).

### With Docker

```bash
docker build -t jade .
docker run --env-file .env -p 8080:8080 jade
```

---

## What I'd Improve

These are things I'd tackle next if I continued developing this project:

- **Presigned URLs** — move to direct client-to-S3 uploads with presigned URLs to reduce server load and avoid proxying large video files through the backend
- **Image compression** — resize and compress images server-side before uploading to S3 to reduce storage costs and improve page load times
- **Rate limiting** — add rate limiting to the upload and auth endpoints to prevent abuse
