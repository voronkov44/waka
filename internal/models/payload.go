package models

import (
	"gorm.io/datatypes"
	"rest_waka/pkg/patch"
	"time"
)

type Model struct {
	ID          uint64   `json:"id"`
	Name        string   `json:"name"`
	Description *string  `json:"description,omitempty"`
	PhotoURL    *string  `json:"photo_url,omitempty"`
	PuffsMax    int      `json:"puffs_max"`
	Flavors     []string `json:"flavors"`
	PriceCents  *int64   `json:"price_cents,omitempty"`
	CreatedAt   int64    `json:"created_at"`
	UpdatedAt   int64    `json:"updated_at"`
}

type CreateModelRequest struct {
	Name        string   `json:"name"`
	Description *string  `json:"description"`
	PhotoURL    *string  `json:"photo_url"`
	PuffsMax    int      `json:"puffs_max"`
	Flavors     []string `json:"flavors"`
	PriceCents  *int64   `json:"price_cents"`
}

type UpdateModelRequest struct {
	Name        patch.Field[string]   `json:"name"`
	Description patch.Field[string]   `json:"description"`
	PhotoURL    patch.Field[string]   `json:"photo_url"`
	PuffsMax    patch.Field[int]      `json:"puffs_max"`
	Flavors     patch.Field[[]string] `json:"flavors"`
	PriceCents  patch.Field[int64]    `json:"price_cents"`
}

type ListModelsResponse struct {
	Items  []Model `json:"items"`
	Limit  int     `json:"limit"`
	Offset int     `json:"offset"`
}

type FlavorRequest struct {
	Value string `json:"value"`
}

type ModelRecord struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement"`
	Name        string         `gorm:"not null"`
	Description *string        `gorm:""`
	PhotoURL    *string        `gorm:""`
	PuffsMax    int            `gorm:"not null"`
	Flavors     datatypes.JSON `gorm:"not null"` // JSON slice
	PriceCents  *int64         `gorm:""`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
