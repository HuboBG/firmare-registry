package config

import (
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config contains all runtime settings.
// Load order: defaults -> YAML (optional) -> env overrides.
type Config struct {
	ListenAddr    string `yaml:"listen_addr"`
	PublicBaseURL string `yaml:"public_base_url"`

	StorageDir string `yaml:"storage_dir"`
	DBPath     string `yaml:"db_path"`

	AdminKey  string `yaml:"admin_key"`
	DeviceKey string `yaml:"device_key"`

	MaxUploadMB int64 `yaml:"max_upload_mb"`

	// OIDC/Keycloak extension point. Off by default.
	OIDC struct {
		Enabled      bool   `yaml:"enabled"`
		IssuerURL    string `yaml:"issuer_url"`
		ClientID     string `yaml:"client_id"`
		Audience     string `yaml:"audience"`
		AdminRole    string `yaml:"admin_role"`
		DeviceRole   string `yaml:"device_role"`
		JWKSCacheSec int    `yaml:"jwks_cache_sec"`
	} `yaml:"oidc"`

	Webhooks struct {
		Secret     string `yaml:"secret"`
		TimeoutSec int    `yaml:"timeout_sec"`
		Retries    int    `yaml:"retries"`
	} `yaml:"webhooks"`
}

// Load reads YAML if path is non-empty, then applies env overrides.
func Load(path string) (Config, error) {
	cfg := defaults()

	if strings.TrimSpace(path) != "" {
		b, err := os.ReadFile(path)
		if err != nil {
			return cfg, err
		}
		if err := yaml.Unmarshal(b, &cfg); err != nil {
			return cfg, err
		}
	}

	applyEnv(&cfg)
	return cfg, nil
}

func defaults() Config {
	var c Config
	c.ListenAddr = ":8080"
	c.PublicBaseURL = ""
	c.StorageDir = "/data/firmware"
	c.DBPath = "/data/db/firmware-registry.db"
	c.MaxUploadMB = 50

	c.Webhooks.TimeoutSec = 5
	c.Webhooks.Retries = 3

	c.OIDC.Enabled = false
	c.OIDC.JWKSCacheSec = 300
	return c
}

func applyEnv(cfg *Config) {
	setStr(&cfg.ListenAddr, "FW_LISTEN_ADDR")
	setStr(&cfg.PublicBaseURL, "FW_PUBLIC_BASE_URL")
	setStr(&cfg.StorageDir, "FW_STORAGE_DIR")
	setStr(&cfg.DBPath, "FW_DB_PATH")
	setStr(&cfg.AdminKey, "FW_ADMIN_KEY")
	setStr(&cfg.DeviceKey, "FW_DEVICE_KEY")

	if v := os.Getenv("FW_MAX_UPLOAD_MB"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil && n > 0 {
			cfg.MaxUploadMB = n
		}
	}

	setStr(&cfg.Webhooks.Secret, "FW_WEBHOOK_SECRET")
	if v := os.Getenv("FW_WEBHOOK_TIMEOUT_SEC"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.Webhooks.TimeoutSec = n
		}
	}
	if v := os.Getenv("FW_WEBHOOK_RETRIES"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			cfg.Webhooks.Retries = n
		}
	}

	if v := os.Getenv("FW_OIDC_ENABLED"); v != "" {
		cfg.OIDC.Enabled = v == "1" || strings.ToLower(v) == "true"
	}
	setStr(&cfg.OIDC.IssuerURL, "FW_OIDC_ISSUER_URL")
	setStr(&cfg.OIDC.ClientID, "FW_OIDC_CLIENT_ID")
	setStr(&cfg.OIDC.Audience, "FW_OIDC_AUDIENCE")
	setStr(&cfg.OIDC.AdminRole, "FW_OIDC_ADMIN_ROLE")
	setStr(&cfg.OIDC.DeviceRole, "FW_OIDC_DEVICE_ROLE")
}

func setStr(dst *string, key string) {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		*dst = v
	}
}
