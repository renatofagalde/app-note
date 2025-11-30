package tests

import "encoding/json"

type noteResponse struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Content   json.RawMessage `json:"content"`
	CreatedAt string          `json:"created_at"`
	UpdatedAt string          `json:"updated_at"`
	DeletedAt *string         `json:"deleted_at"`
}
