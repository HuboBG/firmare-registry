# Firmware Registry API (Go)

Self-hosted firmware registry for ESP32 OTA.

## Features
- Multiple firmware types, each with multiple semantic versions
- Binary storage on local filesystem
- SQLite metadata with automatic migrations
- API-key auth (admin vs device)
- Webhook notifications (HMAC signed, retry w/backoff)
- Clean separation: handlers, services, repositories, storage, config

## Storage layout
`{FW_STORAGE_DIR}/{type}/{version}/firmware.bin`

## Auth
- Admin endpoints require header: `X-Admin-Key: <FW_ADMIN_KEY>`
- Device endpoints require header: `X-Device-Key: <FW_DEVICE_KEY>`

## Endpoints
- GET  `/api/health`
- POST `/api/firmware/{type}/{version}` (admin, multipart field `file`)
- GET  `/api/firmware/{type}/{version}` (device, streams binary)
- DELETE `/api/firmware/{type}/{version}` (admin)
- GET  `/api/firmware/{type}` (device, list)
- GET  `/api/firmware/{type}/latest` (device, semantic latest)
- GET/POST `/api/webhooks` (admin)
- PUT/DELETE `/api/webhooks/{id}` (admin)

Webhook events:
- `firmware.uploaded`
- `firmware.deleted`

Signature:
If `FW_WEBHOOK_SECRET` is set, a header is added:
`X-Webhook-Signature = hex(HMAC-SHA256(secret, raw_body))`

## Config
Use env vars or YAML (set `FW_CONFIG_FILE=/path/config.yaml`).
See internal/config for all keys.

## Migrations
Runs on boot from `./migrations` using golang-migrate.
