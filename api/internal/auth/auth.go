package auth

import "net/http"

// Auth enforces API-key based auth now,
// and keeps an extension point for OIDC later.
type Auth struct {
	AdminKey    string
	DeviceKey   string
	OIDCEnabled bool
}

func (a Auth) RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.OIDCEnabled {
			if r.Header.Get("X-Admin-Key") != a.AdminKey {
				http.Error(w, "unauthorized (admin)", http.StatusUnauthorized)
				return
			}
			next(w, r)
			return
		}
		http.Error(w, "OIDC not enabled in this build", http.StatusNotImplemented)
	}
}

func (a Auth) RequireDevice(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.OIDCEnabled {
			if r.Header.Get("X-Device-Key") != a.DeviceKey {
				http.Error(w, "unauthorized (device)", http.StatusUnauthorized)
				return
			}
			next(w, r)
			return
		}
		http.Error(w, "OIDC not enabled in this build", http.StatusNotImplemented)
	}
}
