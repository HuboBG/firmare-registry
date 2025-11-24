# Firmware Registry Deploy (Docker Compose)

This folder runs:
- api (Go)
- ui (Vue nginx prod server)
- nginx reverse proxy (TLS optional)

## Quick start
1. Copy `.env.example` to `.env` and configure it as needed
2. Copy `api/.env.example` to `api/.env` and configure it as needed
3. Copy `ui/.env.example` to `api/.env` and configure it as needed
4. `docker compose up -d --build`
5. Open http://localhost:8080 (UI) or http://localhost:8080/api/health

## Notes
- API stores firmware binaries in volume `fw_data`
- SQLite DB stored in volume `fw_db`
- Nginx protects UI with basic auth (optional) and proxies /api to backend.
