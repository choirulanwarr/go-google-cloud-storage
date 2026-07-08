# GCS Integration — Google Cloud Storage

REST API untuk integrasi Google Cloud Storage menggunakan Golang, Gin Gonic, GORM, dan PostgreSQL.

## Fitur

- Upload file ke GCS bucket
- Download file dari GCS bucket
- List semua file dalam bucket
- List file berdasarkan folder/prefix
- **Presigned URL** — generate temporary download URL tanpa perlu autentikasi GCP
- Create bucket baru
- Delete file dengan generation-match (safety)

## Teknologi

| Library | Versi | Kegunaan |
|---|---|---|
| `cloud.google.com/go/storage` | v1.53.0 | Google Cloud Storage client |
| `github.com/gin-gonic/gin` | v1.10.0 | HTTP framework |
| `gorm.io/gorm` | v1.26.0 | ORM |
| `gorm.io/driver/postgres` | v1.5.11 | PostgreSQL driver |
| `github.com/spf13/viper` | v1.20.1 | Config management |
| `github.com/go-playground/validator/v10` | v10.26.0 | Request validation |
| `github.com/sirupsen/logrus` | v1.9.3 | Structured logging |

## Quick Start

### Prasyarat

- Go 1.23+
- PostgreSQL 16+
- Google Cloud Storage bucket + service account key (JSON)

### Setup Lokal

```bash
# 1. Clone repository
git clone git@github-personal:choirulanwarr/go-google-cloud-storage.git
cd go-google-cloud-storage

# 2. Setup environment
make env
# Edit file .env dengan kredensial GCS & database kamu

# 3. Install dependencies
make setup

# 4. Jalankan aplikasi
make run
```

### Setup dengan Docker

```bash
# 1. Setup .env
make env
# Edit .env — DB_HOST akan otomatis di-override ke "postgres" di dalam container

# 2. Jalankan dengan Docker Compose
make docker-up

# 3. Cek logs
make docker-logs

# 4. Stop
make docker-down
```

## API Endpoints

Semua endpoint berada di prefix `/api/v1`.

### Upload File

```bash
POST /api/v1/upload
Content-Type: multipart/form-data

# Form fields:
#   folder  (required) — target folder di bucket
#   file    (required) — file yang di-upload

curl -X POST http://localhost:4000/api/v1/upload \
  -F "folder=images" \
  -F "file=@photo.jpg"
```

**Allowed MIME types:** `image/jpeg`, `image/png`, `image/jpg`, `application/pdf`, `application/zip`, `application/octet-stream`

**Response:**
```json
{
  "api_id": "API_CALL_...",
  "status": "SUCCESS",
  "message": "Success save data",
  "data": {
    "path": "images/20260708150405_aB3xZ.jpg"
  }
}
```

---

### Download File

```bash
GET /api/v1/download?path=<object-path>

curl -X GET "http://localhost:4000/api/v1/download?path=images/photo.jpg" \
  -o downloaded.jpg
```

Response berupa file binary dengan header `Content-Disposition: attachment`.

---

### List Semua File

```bash
GET /api/v1/list

curl http://localhost:4000/api/v1/list
```

**Response:**
```json
{
  "api_id": "API_CALL_...",
  "status": "SUCCESS",
  "message": "Success get data",
  "data": [
    {
      "name": "photo.jpg",
      "url": "https://storage.googleapis.com/voucher_techno_test/images/photo.jpg",
      "size": "1.2 MB",
      "type": "image/jpeg",
      "createdAt": "2026-07-08T07:04:05Z",
      "updatedAt": "2026-07-08T07:04:05Z"
    }
  ]
}
```

---

### List File per Folder

```bash
GET /api/v1/list/:folder

curl http://localhost:4000/api/v1/list/images
```

---

### Presigned URL

Generate temporary URL untuk mengakses file tanpa kredensial GCP. Berguna untuk verifikasi upload atau share file sementara.

```bash
GET /api/v1/presigned-url?path=<object-path>&expires=<minutes>

# Default 15 menit
curl "http://localhost:4000/api/v1/presigned-url?path=images/photo.jpg"

# Custom expiration — 60 menit
curl "http://localhost:4000/api/v1/presigned-url?path=images/photo.jpg&expires=60"
```

| Parameter | Tipe | Required | Default | Keterangan |
|---|---|---|---|---|
| `path` | string | Ya | — | Path object di GCS bucket |
| `expires` | int | Tidak | 15 | Masa berlaku URL dalam menit (min 1, max 10080 / 7 hari) |

**Response:**
```json
{
  "api_id": "API_CALL_...",
  "status": "SUCCESS",
  "message": "Success get data",
  "data": {
    "url": "https://storage.googleapis.com/voucher_techno_test/images/photo.jpg?X-Goog-Algorithm=...",
    "path": "images/photo.jpg",
    "expires_at": "2026-07-08T07:19:05Z"
  }
}
```

> **Catatan:** Untuk generate presigned URL di environment lokal (non-GCP), set `GCS_CONFIG_SA=true` di `.env` agar menggunakan service account JSON key untuk signing.

---

### Create Bucket

```bash
POST /api/v1/bucket/create
Content-Type: application/json

curl -X POST http://localhost:4000/api/v1/bucket/create \
  -H "Content-Type: application/json" \
  -d '{"bucketName": "my-new-bucket"}'
```

Bucket dibuat dengan Storage Class dan Location sesuai konfigurasi di `.env`.

---

### Delete File

```bash
DELETE /api/v1/delete
Content-Type: application/json

curl -X DELETE http://localhost:4000/api/v1/delete \
  -H "Content-Type: application/json" \
  -d '{"path": "images/photo.jpg"}'
```

Delete menggunakan **generation-match** untuk mencegah race condition — file hanya dihapus jika generation-nya cocok.

---

## Konfigurasi (.env)

```ini
# App
APP_STATUS=development
APP_PORT=4000
AUTO_MIGRATION_SWITCH=1

# Database
DB_USERNAME=postgres
DB_PASSWORD=root
DB_NAME=gcs
DB_HOST=localhost
DB_PORT=5432
DB_SSLMODE=disable
DB_CONN_MAX_IDLE_TIME=30m
DB_CONN_MAX_LIFE_TIME=30m
DB_MAX_OPEN_CONN=50
DB_MAX_IDLE_CONN=10

# Logger
LOGGER_STDOUT=true
LOGGER_FILE_LOCATION=app.log
LOGGER_LEVEL=info

# Google Cloud Storage
GCS_CONFIG_SA=false                     # true = pakai local credential JSON file
GCS_PROJECT_ID=my-gcp-project-id
GCS_BUCKET_NAME=my-bucket-name
GCS_CREDENTIAL_FILE_PATH=./my-key.json  # path ke service account JSON key
GCS_STORAGE_CLASS_BUCKET=COLDLINE
GCS_STORAGE_LOCATION_BUCKET=asia
```

## Makefile Commands

```bash
make help           # Lihat semua perintah yang tersedia

# Docker
make docker-up      # Start semua service (PostgreSQL + App)
make docker-down    # Stop semua service
make docker-build   # Rebuild image tanpa cache
make docker-logs    # Tail logs container app
make docker-shell   # Shell ke dalam container app
make docker-db-shell# PostgreSQL shell

# Development
make run            # Jalankan app secara lokal
make build          # Build binary
make watch          # Auto-reload dengan air

# Testing
make test           # Run unit tests
make test-verbose   # Run tests dengan output verbose
make test-cover     # Run tests + coverage report
make test-html      # Coverage report dalam HTML (buka di browser)

# Code Quality
make fmt            # Format semua file Go
make vet            # Run go vet
make lint           # Run golangci-lint
make tidy           # Tidy Go modules

# Setup & Cleanup
make setup          # Install dependencies
make env            # Copy .env.example → .env
make clean          # Hapus build artifacts, coverage, logs
make clean-all      # clean + hapus Docker volumes
```

## Struktur Project

```
go-google-cloud-storage/
  main.go                      # Entry point
  Dockerfile                   # Multi-stage Docker build
  docker-compose.yml           # PostgreSQL + App services
  Makefile                     # Development automation
  .env                         # Environment config (gitignored)
  .env.example                 # Config template
  app/
    config/                    # App initialization, DB, server, validator, viper
    constant/                  # Response status/message constants
    handler/                   # HTTP handlers (FileHandler)
    helper/                    # Logging, response builder, file utils, validator
    integration/               # GCS client wrapper (Upload, Download, List, CreateBucket, Delete, PresignedURL)
    middleware/                # Request ID middleware
    model/                     # GORM models (Config)
    resource/
      request/                 # Request DTOs
      response/                # Response DTOs & formatters
    router/                    # Route definitions
    service/                   # Business logic layer
  test/
    handler/                   # Handler tests
  logs/                        # Log output (gitignored)
```

## Lisensi

MIT
