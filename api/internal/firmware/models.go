package firmware

import "time"

// Firmware is the internal domain model.
type Firmware struct {
	Type      string
	Version   string
	Filename  string
	SizeBytes int64
	SHA256    string
	CreatedAt time.Time
}

// FirmwareDTO is what we expose over HTTP.
type FirmwareDTO struct {
	Type        string    `json:"type"`
	Version     string    `json:"version"`
	Filename    string    `json:"filename"`
	SizeBytes   int64     `json:"sizeBytes"`
	SHA256      string    `json:"sha256"`
	CreatedAt   time.Time `json:"createdAt"`
	DownloadURL string    `json:"downloadUrl,omitempty"`
}

func (f Firmware) ToDTO(downloadURL string) FirmwareDTO {
	return FirmwareDTO{
		Type:        f.Type,
		Version:     f.Version,
		Filename:    f.Filename,
		SizeBytes:   f.SizeBytes,
		SHA256:      f.SHA256,
		CreatedAt:   f.CreatedAt,
		DownloadURL: downloadURL,
	}
}
