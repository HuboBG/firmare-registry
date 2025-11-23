package webhook

import (
	"database/sql"
	"encoding/json"
)

// Repository persists webhooks.
type Repository interface {
	List() ([]Webhook, error)
	Create(h Webhook) (int64, error)
	Update(id int64, h Webhook) error
	Delete(id int64) error
}

// SQLiteRepo implements Repository over SQLite.
type SQLiteRepo struct {
	DB *sql.DB
}

func (r *SQLiteRepo) List() ([]Webhook, error) {
	rows, err := r.DB.Query(`SELECT id, url, events, enabled FROM webhooks`)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var out []Webhook
	for rows.Next() {
		var h Webhook
		var eventsJSON string
		if err := rows.Scan(&h.ID, &h.URL, &eventsJSON, &h.Enabled); err != nil {
			continue
		}
		_ = json.Unmarshal([]byte(eventsJSON), &h.Events)
		out = append(out, h)
	}
	return out, nil
}

func (r *SQLiteRepo) Create(h Webhook) (int64, error) {
	ev, _ := json.Marshal(h.Events)
	res, err := r.DB.Exec(
		`INSERT INTO webhooks(url, events, enabled) VALUES(?,?,?)`,
		h.URL, string(ev), h.Enabled,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *SQLiteRepo) Update(id int64, h Webhook) error {
	ev, _ := json.Marshal(h.Events)
	_, err := r.DB.Exec(
		`UPDATE webhooks SET url=?, events=?, enabled=? WHERE id=?`,
		h.URL, string(ev), h.Enabled, id,
	)
	return err
}

func (r *SQLiteRepo) Delete(id int64) error {
	_, err := r.DB.Exec(`DELETE FROM webhooks WHERE id=?`, id)
	return err
}
