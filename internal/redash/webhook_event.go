package redash

import (
	"encoding/json"
	"io"
)

type WebhookEvent struct {
	URLBase string `json:"url_base"`
	Event   string `json:"event"`
	Alert   Alert  `json:"alert"`
}

type Alert struct {
	State           string  `json:"state"`
	QueryID         int64   `json:"query_id"`
	Name            string  `json:"name"`
	Rearm           string  `json:"rearm"`
	UpdatedAt       string  `json:"updated_at"`
	UserID          int64   `json:"user_id"`
	CreatedAt       string  `json:"created_at"`
	LastTriggeredAt string  `json:"last_triggered_at"`
	ID              int64   `json:"id"`
	Options         Options `json:"options"`
}

type Options map[string]interface{}

func Parse(r io.Reader) (*WebhookEvent, error) {
	event := WebhookEvent{}
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&event); err != nil {
		return nil, err
	}

	return &event, nil
}
