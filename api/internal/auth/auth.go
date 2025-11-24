package auth

import (
	"context"
	"log"
	"net/http"
	"strings"
)

// Auth enforces both API-key and OIDC/JWT based authentication
type Auth struct {
	AdminKey     string
	DeviceKey    string
	OIDCEnabled  bool
	OIDCVerifier *OIDCVerifier
}

func (a Auth) RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// If OIDC is enabled, try JWT first, then fall back to API key
		if a.OIDCEnabled && a.OIDCVerifier != nil {
			if a.verifyJWT(w, r, a.OIDCVerifier.adminRole) {
				next(w, r)
				return
			}
		}

		// Fall back to API key authentication
		if a.AdminKey != "" && r.Header.Get("X-Admin-Key") == a.AdminKey {
			next(w, r)
			return
		}

		http.Error(w, "unauthorized (admin)", http.StatusUnauthorized)
	}
}

func (a Auth) RequireDevice(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// If OIDC is enabled, try JWT first, then fall back to API key
		if a.OIDCEnabled && a.OIDCVerifier != nil {
			if a.verifyJWT(w, r, a.OIDCVerifier.deviceRole) {
				next(w, r)
				return
			}
		}

		// Fall back to API key authentication
		if a.DeviceKey != "" && r.Header.Get("X-Device-Key") == a.DeviceKey {
			next(w, r)
			return
		}

		http.Error(w, "unauthorized (device)", http.StatusUnauthorized)
	}
}

// verifyJWT validates the JWT token from the Authorization header and checks for the required role
func (a Auth) verifyJWT(w http.ResponseWriter, r *http.Request, requiredRole string) bool {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return false
	}

	token := ExtractBearerToken(authHeader)
	if token == "" {
		return false
	}

	ctx := context.Background()
	idToken, err := a.OIDCVerifier.VerifyToken(ctx, token)
	if err != nil {
		log.Printf("JWT verification failed: %v", err)
		return false
	}

	hasRole, err := a.OIDCVerifier.HasRole(idToken, requiredRole)
	if err != nil {
		log.Printf("Failed to check role: %v", err)
		return false
	}

	if !hasRole && requiredRole != "" {
		log.Printf("User missing required role: %s", requiredRole)
		return false
	}

	return true
}

// ExtractBearerToken is a helper to extract Bearer token from Authorization header
func ExtractBearerToken(authHeader string) string {
	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
		return parts[1]
	}
	return ""
}
