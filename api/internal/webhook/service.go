package webhook

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"
)

// Service dispatches webhook events to subscribed URLs.
type Service struct {
	Repo       Repository
	Secret     string
	TimeoutSec int
	Retries    int
}

func (s *Service) Dispatch(event string, data any) {
	hooks, err := s.Repo.List()
	if err != nil {
		return
	}

	payload := EventPayload{
		Event: event,
		Data:  data,
		Time:  time.Now().UTC().Format(time.RFC3339),
	}

	body, _ := json.Marshal(payload)

	for _, h := range hooks {
		if !h.Enabled || !contains(h.Events, event) {
			continue
		}
		go s.deliver(h.URL, body)
	}
}

func (s *Service) deliver(url string, body []byte) {
	timeout := time.Duration(s.TimeoutSec) * time.Second
	retries := s.Retries
	if retries < 0 {
		retries = 0
	}

	for attempt := 0; attempt <= retries; attempt++ {
		req, _ := http.NewRequest("POST", url, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		if s.Secret != "" {
			req.Header.Set("X-Webhook-Signature", hmacHex([]byte(s.Secret), body))
		}

		client := &http.Client{Timeout: timeout}
		resp, err := client.Do(req)
		if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			_ = resp.Body.Close()
			return
		}
		if resp != nil {
			_ = resp.Body.Close()
		}

		time.Sleep(time.Duration(attempt+1) * 500 * time.Millisecond)
	}
}

func hmacHex(secret, data []byte) string {
	m := hmac.New(sha256.New, secret)
	m.Write(data)
	return hex.EncodeToString(m.Sum(nil))
}

func contains(list []string, v string) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}
