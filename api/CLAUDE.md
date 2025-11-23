# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Self-hosted firmware registry API for ESP32 OTA (over-the-air) updates. Built with Go 1.22, SQLite with WAL mode, and filesystem-based binary storage.

## Build & Development Commands

```bash
# Build the application
go build ./cmd/firmware-registry

# Run the application (requires config)
./firmware-registry

# Tidy dependencies
go mod tidy

# Run with custom config
FW_CONFIG_FILE=/path/to/config.yaml ./firmware-registry
```

## Architecture Overview

The codebase follows a **strict 4-layer architecture** with clear separation of concerns:

```
HTTP Handler Layer → Service Layer → Repository Layer → Storage/Database Layer
```

### Layer Responsibilities

**1. Handler Layer** (`internal/api/handlers/`)
- HTTP request/response translation
- Path parsing and routing logic
- Multipart form parsing
- No business logic

**2. Service Layer** (`internal/firmware/service.go`, `internal/webhook/service.go`)
- Contains ALL business logic
- Orchestrates repositories and storage
- Cross-cutting concerns (SHA256 hashing, atomic writes)
- Completely HTTP-agnostic

**3. Repository Layer** (`internal/*/repository_sqlite.go`)
- Abstract data persistence via interfaces
- SQLite implementations provided
- JSON marshaling for complex types (webhook events)
- Interface-based design enables easy database swapping

**4. Storage/Database Layer**
- **Storage**: Filesystem operations with atomic writes (`.tmp` + rename)
- **Database**: SQLite with WAL mode for read/write concurrency
- Storage layout: `{BaseDir}/{type}/{version}/firmware.bin`

### Dependency Flow

```
                    Config
                      |
        +-------------+-------------+
        |             |             |
     SQLite        Storage        Auth
        |             |             |
   Repositories   ---------->  Services
        |             |             |
        +-------------+-------------+
                      |
                  Handlers
                      |
                   Router
                      |
                 HTTP Server
```

**Critical Rule**: Dependencies flow downward only. No circular dependencies allowed.

## Key Architectural Patterns

### Authentication Pattern

Two-tier API-key based authentication (admin vs device):
- Admin: `X-Admin-Key` header (upload, delete, webhook management)
- Device: `X-Device-Key` header (download, list, latest)

Authentication is **NOT global middleware**. Applied per-endpoint using function composition:

```go
h.Auth.RequireAdmin(func(w http.ResponseWriter, r *http.Request) {
    h.upload(w, r, t, v)
})(w, r)
```

OIDC extension point exists (`internal/auth/auth.go`) but returns 501 Not Implemented.

### Repository Interface Pattern

All repositories use interfaces for testability and swappability:

```go
type Repository interface {
    Upsert(Firmware) error
    Get(typeName, version string) (Firmware, error)
    List(typeName string) ([]Firmware, error)
    Delete(typeName, version string) error
}
```

### Domain Model vs DTO Separation

```go
type Firmware struct { ... }          // Internal domain model
type FirmwareDTO struct { ... }       // API representation
func (f Firmware) ToDTO(url string) FirmwareDTO
```

DTOs control JSON structure. Download URLs computed at presentation time.

### Atomic File Writes

Storage always writes to `.tmp` file then renames (atomic on POSIX):

```go
tmp := dest + ".tmp"
os.WriteFile(tmp, data, 0o644)
os.Rename(tmp, dest)  // Atomic
```

### Webhook Event System

Fire-and-forget async notifications with retry logic:
- Events: `firmware.uploaded`, `firmware.deleted`
- Each webhook delivered in separate goroutine
- Exponential backoff: `(attempt + 1) * 500ms`
- HMAC-SHA256 signatures if `FW_WEBHOOK_SECRET` set
- **Critical**: Webhook failures don't block main operations

## Database Design

### Schema

```sql
firmwares:
- PRIMARY KEY (type, version)  -- Composite key
- Metadata only (filename, size_bytes, sha256, created_at)
- created_at as RFC3339 TEXT

webhooks:
- Auto-increment ID
- events as JSON TEXT array
- Boolean enabled flag
```

### Migrations

- Uses `golang-migrate/migrate/v4`
- Runs on EVERY boot (idempotent)
- Location: `./migrations/`
- Naming: `{version}_{name}.{up|down}.sql`
- Fail-fast on errors

### SQLite Configuration

```
Connection: {path}?_busy_timeout=5000&_journal_mode=WAL
- WAL mode: readers don't block writers
- 5-second busy timeout for concurrent writes
- Single connection pool per instance
```

## Configuration Management

Three-layer precedence (highest wins):
1. Hard-coded defaults (in code)
2. YAML file (via `FW_CONFIG_FILE` env var)
3. Environment variable overrides

### Key Environment Variables

```bash
# Server
FW_LISTEN_ADDR=:8080
FW_PUBLIC_BASE_URL=https://example.com

# Storage
FW_STORAGE_DIR=/data/firmware
FW_DB_PATH=/data/db/firmware-registry.db

# Auth
FW_ADMIN_KEY=<required>
FW_DEVICE_KEY=<required>

# Upload
FW_MAX_UPLOAD_MB=50

# Webhooks
FW_WEBHOOK_SECRET=<hmac-signing-key>
FW_WEBHOOK_TIMEOUT_SEC=5
FW_WEBHOOK_RETRIES=3
```

## Routing & HTTP Handling

Uses standard library `net/http.ServeMux` - no framework dependency.

### Manual Path Parsing

Handlers implement `http.Handler` interface with custom path parsing:

```go
path := strings.TrimPrefix(r.URL.Path, "/api/firmware/")
parts := filterEmpty(strings.Split(path, "/"))
```

Special endpoint: `/api/firmware/{type}/latest` uses semantic version sorting.

## Important Conventions

### Error Handling

Services return `(result, error)` tuples. Handlers convert to HTTP status codes:
- Database errors → 500
- Not found → 404
- Validation errors → 400
- No custom error types (standard Go errors)

### Dependency Injection

All dependencies injected via struct fields in `main.go`:

```go
fwHandler := &handlers.FirmwareHandler{
    Auth:     authHandler,
    Service:  fwSvc,
    Webhooks: whSvc,
    MaxBytes: cfg.MaxUploadMB * 1024 * 1024,
}
```

No global state or singletons.

### File Organization

- One concern per file (`models.go`, `service.go`, `repository_sqlite.go`)
- Repository implementations: `repository_{db_type}.go`
- Handlers in `internal/api/handlers/`
- Utilities in `internal/util/`

### Semantic Versioning

Custom `CompareSemver` utility (`internal/util/semver.go`):
- Handles `major.minor.patch` comparison
- `/latest` endpoint uses semantic sorting
- Graceful fallback to lexical comparison

## Critical Implementation Details

### Composite Primary Key Strategy

`(type, version)` is the primary key, enabling:
- Multiple firmware types (e.g., "esp32-main", "esp32-bootloader")
- Multiple semantic versions per type
- Upsert allows re-uploading same version (overwrites)

### Multipart Upload Handling

- `MaxBytesReader` prevents DoS
- Form field name hardcoded as `"file"`
- Original filename preserved in metadata

### Time Storage

Timestamps stored as RFC3339 strings in SQLite, parsed back to `time.Time` in Go.

### WAL Mode Benefits

- Readers don't block writers
- Better concurrency for OTA downloads during uploads
- Trades disk space for performance

### Graceful Webhook Failures

Webhook dispatch errors are swallowed. Main operation success independent of notification success. No database transaction coupling.

## When Making Changes

1. **Adding new endpoints**: Follow handler → service → repository pattern
2. **Database changes**: Create new migration in `migrations/`, never modify existing ones
3. **Authentication changes**: Modify `internal/auth/auth.go`, maintain interface compatibility
4. **Storage changes**: Maintain atomic write pattern (`.tmp` + rename)
5. **New dependencies**: Run `go mod tidy` after adding imports
6. **Avoid circular imports**: Never import `internal/api` from handlers, use `internal/auth` instead

## Security Considerations

- API keys stored in plaintext (no hashing)
- No rate limiting implemented
- HMAC webhook signatures for integrity verification
- Input validation via `MaxBytesReader` and form parsing
