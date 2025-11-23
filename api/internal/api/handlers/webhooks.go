package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"firmware-registry-api/internal/auth"
	"firmware-registry-api/internal/util"
	"firmware-registry-api/internal/webhook"
)

// WebhookHandler manages webhook CRUD.
type WebhookHandler struct {
	Auth auth.Auth
	Repo webhook.Repository
}

func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/webhooks" {
		switch r.Method {
		case http.MethodGet:
			h.Auth.RequireAdmin(func(w http.ResponseWriter, r *http.Request) {
				h.list(w)
			})(w, r)
		case http.MethodPost:
			h.Auth.RequireAdmin(func(w http.ResponseWriter, r *http.Request) {
				h.create(w, r)
			})(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// /api/webhooks/{id}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/webhooks/")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	h.Auth.RequireAdmin(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			h.update(w, r, id)
		case http.MethodDelete:
			h.delete(w, id)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})(w, r)
}

func (h *WebhookHandler) list(w http.ResponseWriter) {
	hooks, err := h.Repo.List()
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	out := make([]webhook.WebhookDTO, 0, len(hooks))
	for _, x := range hooks {
		out = append(out, webhook.WebhookDTO{
			ID: x.ID, URL: x.URL, Events: x.Events, Enabled: x.Enabled,
		})
	}
	util.WriteJSON(w, out)
}

func (h *WebhookHandler) create(w http.ResponseWriter, r *http.Request) {
	var dto webhook.WebhookDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	if dto.URL == "" || len(dto.Events) == 0 {
		http.Error(w, "url/events required", http.StatusBadRequest)
		return
	}
	if dto.Enabled == false {
		dto.Enabled = true
	}

	id, err := h.Repo.Create(webhook.Webhook{
		URL: dto.URL, Events: dto.Events, Enabled: dto.Enabled,
	})
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	util.WriteJSON(w, map[string]any{"id": id})
}

func (h *WebhookHandler) update(w http.ResponseWriter, r *http.Request, id int64) {
	var dto webhook.WebhookDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	if err := h.Repo.Update(id, webhook.Webhook{
		URL: dto.URL, Events: dto.Events, Enabled: dto.Enabled,
	}); err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	util.WriteJSON(w, map[string]any{"updated": true})
}

func (h *WebhookHandler) delete(w http.ResponseWriter, id int64) {
	if err := h.Repo.Delete(id); err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	util.WriteJSON(w, map[string]any{"deleted": true})
}
