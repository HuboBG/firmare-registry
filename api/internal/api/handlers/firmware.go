package handlers

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"firmware-registry-api/internal/auth"
	"firmware-registry-api/internal/firmware"
	"firmware-registry-api/internal/util"
	"firmware-registry-api/internal/webhook"
)

// FirmwareHandler translates HTTP to firmware service calls.
type FirmwareHandler struct {
	Auth     auth.Auth
	Service  *firmware.Service
	Webhooks *webhook.Service
	MaxBytes int64
}

func (h *FirmwareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/firmware/")
	parts := filterEmpty(strings.Split(path, "/"))
	if len(parts) == 0 {
		http.Error(w, "missing firmware type", http.StatusBadRequest)
		return
	}
	t := parts[0]

	// GET /api/firmware/{type}
	if len(parts) == 1 && r.Method == http.MethodGet {
		h.Auth.RequireDevice(func(w http.ResponseWriter, r *http.Request) {
			h.list(w, t)
		})(w, r)
		return
	}

	// GET /api/firmware/{type}/latest
	if len(parts) == 2 && parts[1] == "latest" && r.Method == http.MethodGet {
		h.Auth.RequireDevice(func(w http.ResponseWriter, r *http.Request) {
			h.latest(w, t)
		})(w, r)
		return
	}

	// /api/firmware/{type}/{version}
	if len(parts) == 2 {
		v := parts[1]
		switch r.Method {
		case http.MethodPost:
			h.Auth.RequireAdmin(func(w http.ResponseWriter, r *http.Request) {
				h.upload(w, r, t, v)
			})(w, r)
		case http.MethodGet:
			h.Auth.RequireDevice(func(w http.ResponseWriter, r *http.Request) {
				h.download(w, t, v)
			})(w, r)
		case http.MethodDelete:
			h.Auth.RequireAdmin(func(w http.ResponseWriter, r *http.Request) {
				h.delete(w, t, v)
			})(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	http.Error(w, "invalid firmware route", http.StatusNotFound)
}

func (h *FirmwareHandler) upload(w http.ResponseWriter, r *http.Request, t, v string) {
	maxN := h.MaxBytes
	r.Body = http.MaxBytesReader(w, r.Body, maxN)

	if err := r.ParseMultipartForm(maxN); err != nil {
		http.Error(w, "invalid multipart", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "missing file field", http.StatusBadRequest)
		return
	}
	defer func(file multipart.File) {
		_ = file.Close()
	}(file)

	rec, err := h.Service.SaveFirmware(t, v, header.Filename, file)
	if err != nil {
		http.Error(w, "save failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	dto := rec.ToDTO(h.Service.DownloadURL(t, v))

	if h.Webhooks != nil {
		h.Webhooks.Dispatch("firmware.uploaded", dto)
	}

	util.WriteJSON(w, dto)
}

func (h *FirmwareHandler) download(w http.ResponseWriter, t, v string) {
	rec, err := h.Service.Repo.Get(t, v)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	path := h.Service.DownloadPath(t, v)
	f, err := os.Open(path)
	if err != nil {
		http.Error(w, "missing binary", http.StatusNotFound)
		return
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.FormatInt(rec.SizeBytes, 10))
	w.Header().Set("X-Firmware-Sha256", rec.SHA256)
	w.Header().Set("X-Firmware-Version", rec.Version)

	_, _ = io.Copy(w, f)
}

func (h *FirmwareHandler) delete(w http.ResponseWriter, t, v string) {
	rec, err := h.Service.Repo.Get(t, v)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	_ = os.RemoveAll(h.Service.Storage.Dir(t, v))
	_ = h.Service.Repo.Delete(t, v)

	dto := rec.ToDTO(h.Service.DownloadURL(t, v))

	if h.Webhooks != nil {
		h.Webhooks.Dispatch("firmware.deleted", dto)
	}

	util.WriteJSON(w, map[string]any{"deleted": true})
}

func (h *FirmwareHandler) list(w http.ResponseWriter, t string) {
	list, err := h.Service.Repo.List(t)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	sort.Slice(list, func(i, j int) bool {
		return util.CompareSemver(list[i].Version, list[j].Version) > 0
	})

	out := make([]firmware.FirmwareDTO, 0, len(list))
	for _, f := range list {
		out = append(out, f.ToDTO(h.Service.DownloadURL(f.Type, f.Version)))
	}
	util.WriteJSON(w, out)
}

func (h *FirmwareHandler) latest(w http.ResponseWriter, t string) {
	list, err := h.Service.Repo.List(t)
	if err != nil || len(list) == 0 {
		http.Error(w, "no firmware", http.StatusNotFound)
		return
	}

	sort.Slice(list, func(i, j int) bool {
		return util.CompareSemver(list[i].Version, list[j].Version) > 0
	})

	f := list[0]
	util.WriteJSON(w, f.ToDTO(h.Service.DownloadURL(f.Type, f.Version)))
}

func filterEmpty(in []string) []string {
	out := make([]string, 0, len(in))
	for _, p := range in {
		if strings.TrimSpace(p) != "" {
			out = append(out, p)
		}
	}
	return out
}
