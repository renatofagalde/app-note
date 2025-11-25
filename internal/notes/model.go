package notes

import (
	"time"

	"gorm.io/datatypes"
)

type Note struct {
	ID        string         `gorm:"type:text;primaryKey"       json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Content   datatypes.JSON `gorm:"type:jsonb;not null"        json:"content"`
	CreatedAt time.Time      `gorm:"autoCreateTime"             json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"             json:"updated_at"`
	DeletedAt *time.Time     `gorm:"index"                      json:"deleted_at,omitempty"`
}

type CreateNoteRequest struct {
	Name    string         `json:"name"`
	Content datatypes.JSON `json:"content"`
}

type NoteResponse struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Content   datatypes.JSON `json:"content"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt *time.Time     `json:"deleted_at,omitempty"`
}

func (n *Note) ToResponse() *NoteResponse {
	return &NoteResponse{
		ID:        n.ID,
		Name:      n.Name,
		Content:   n.Content,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
		DeletedAt: n.DeletedAt,
	}
}
