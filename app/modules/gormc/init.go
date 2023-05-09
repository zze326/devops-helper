package gormc

import (
	"gorm.io/gorm"
	"time"
)

// Required for vendoring see golang.org/issue/13832
type Model struct {
	ID        int            `gorm:"primarykey" json:"id,omitempty"`
	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
