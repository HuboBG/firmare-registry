package firmware

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"strings"
	"time"
)

// Repository persists firmware metadata.
type Repository interface {
	Upsert(Firmware) error
	Get(typeName, version string) (Firmware, error)
	List(typeName string) ([]Firmware, error)
	Delete(typeName, version string) error
}

// Service holds business logic only.
type Service struct {
	Repo       Repository
	Storage    Storage
	PublicBase string
}

// SaveFirmware reads the uploaded binary, computes SHA256,
// writes to disk atomically, and upserts metadata.
func (s *Service) SaveFirmware(typeName, version, filename string, r io.Reader) (Firmware, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return Firmware{}, err
	}

	sum := sha256.Sum256(data)
	shaHex := hex.EncodeToString(sum[:])

	dir := s.Storage.Dir(typeName, version)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return Firmware{}, err
	}

	dest := s.Storage.FilePath(typeName, version)
	tmp := dest + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return Firmware{}, err
	}
	if err := os.Rename(tmp, dest); err != nil {
		return Firmware{}, err
	}

	rec := Firmware{
		Type:      typeName,
		Version:   version,
		Filename:  filename,
		SizeBytes: int64(len(data)),
		SHA256:    shaHex,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.Repo.Upsert(rec); err != nil {
		return Firmware{}, err
	}
	return rec, nil
}

func (s *Service) DownloadPath(typeName, version string) string {
	return s.Storage.FilePath(typeName, version)
}

func (s *Service) DownloadURL(typeName, version string) string {
	if s.PublicBase == "" {
		return ""
	}
	base := strings.TrimRight(s.PublicBase, "/")
	return base + "/api/firmware/" + typeName + "/" + version
}
