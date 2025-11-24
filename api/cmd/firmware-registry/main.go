package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"firmware-registry-api/internal/api"
	"firmware-registry-api/internal/api/handlers"
	"firmware-registry-api/internal/auth"
	"firmware-registry-api/internal/config"
	"firmware-registry-api/internal/db"
	"firmware-registry-api/internal/firmware"
	"firmware-registry-api/internal/webhook"
)

func main() {
	cfgPath := os.Getenv("FW_CONFIG_FILE")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatal("config load failed:", err)
	}

	// Ensure directories exist
	_ = os.MkdirAll(cfg.StorageDir, 0o755)
	_ = os.MkdirAll(filepath.Dir(cfg.DBPath), 0o755)

	// DB + migrations
	database := db.OpenSQLite(cfg.DBPath)
	db.RunMigrations(cfg.DBPath, "./migrations")

	// Firmware layer
	fwRepo := &firmware.SQLiteRepo{DB: database}
	fwSvc := &firmware.Service{
		Repo:       fwRepo,
		Storage:    firmware.Storage{BaseDir: cfg.StorageDir},
		PublicBase: cfg.PublicBaseURL,
	}

	// Webhook layer
	whRepo := &webhook.SQLiteRepo{DB: database}
	whSvc := &webhook.Service{
		Repo:       whRepo,
		Secret:     cfg.Webhooks.Secret,
		TimeoutSec: cfg.Webhooks.TimeoutSec,
		Retries:    cfg.Webhooks.Retries,
	}

	// Initialize OIDC verifier if enabled
	var oidcVerifier *auth.OIDCVerifier
	if cfg.OIDC.Enabled {
		ctx := context.Background()
		var err error
		oidcVerifier, err = auth.NewOIDCVerifier(
			ctx,
			cfg.OIDC.IssuerURL,
			cfg.OIDC.ClientID,
			cfg.OIDC.Audience,
			cfg.OIDC.AdminRole,
			cfg.OIDC.DeviceRole,
		)
		if err != nil {
			log.Printf("WARNING: OIDC enabled but failed to initialize: %v", err)
			log.Println("Falling back to API key authentication only")
			cfg.OIDC.Enabled = false
		} else {
			log.Println("OIDC authentication enabled with issuer:", cfg.OIDC.IssuerURL)
		}
	}

	authHandler := auth.Auth{
		AdminKey:     cfg.AdminKey,
		DeviceKey:    cfg.DeviceKey,
		OIDCEnabled:  cfg.OIDC.Enabled,
		OIDCVerifier: oidcVerifier,
	}

	fwHandler := &handlers.FirmwareHandler{
		Auth:     authHandler,
		Service:  fwSvc,
		Webhooks: whSvc,
		MaxBytes: cfg.MaxUploadMB * 1024 * 1024,
	}
	whHandler := &handlers.WebhookHandler{
		Auth: authHandler,
		Repo: whRepo,
	}

	router := api.NewRouter(fwHandler, whHandler)
	handler := api.CORSMiddleware(router)

	log.Println("Firmware Registry API listening on", cfg.ListenAddr)
	log.Fatal(http.ListenAndServe(cfg.ListenAddr, handler))
}
