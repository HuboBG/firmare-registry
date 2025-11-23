package webhook

// Webhook is stored in DB.
type Webhook struct {
	ID      int64
	URL     string
	Events  []string
	Enabled bool
}

// WebhookDTO is sent/received over the API.
type WebhookDTO struct {
	ID      int64    `json:"id"`
	URL     string   `json:"url"`
	Events  []string `json:"events"`
	Enabled bool     `json:"enabled"`
}

type EventPayload struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
	Time  string `json:"time"`
}
