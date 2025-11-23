# Firmware Registry Deploy (Docker Compose)

This folder runs:
- firmware-registry-api (Go)
- firmware-registry-ui (Vue dev server)
- nginx reverse proxy (TLS optional)

## Quick start
1. Copy `.env.example` to `.env` and edit keys/paths
2. `docker compose up -d --build`
3. Open http://localhost (UI) or http://localhost/api/health

## Notes
- API stores firmware binaries in volume `fw_data`
- SQLite DB stored in volume `fw_db`
- Nginx protects UI with basic auth (optional) and proxies /api to backend.
