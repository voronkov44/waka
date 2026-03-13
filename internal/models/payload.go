package models

import (
	"gorm.io/datatypes"
	"rest_waka/pkg/patch"
	"time"
)

type ModelTag struct {
	Key       string `json:"key,omitempty"`
	Label     string `json:"label"`
	BgColor   string `json:"bg_color"`
	TextColor string `json:"text_color"`
}

type Model struct {
	ID          uint64  `json:"id"`
	Name        string  `json:"name"`
	Status      string  `json:"status"`
	Description *string `json:"description,omitempty"`

	PhotoKey *string `json:"photo_key,omitempty"` // field db
	PhotoURL *string `json:"photo_url,omitempty"` // computed, dont db

	Tag        *ModelTag `json:"tag,omitempty"`
	PuffsMax   int       `json:"puffs_max"`
	Flavors    []string  `json:"flavors"`
	PriceCents *int64    `json:"price_cents,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type PublicModel struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	PhotoURL    *string   `json:"photo_url,omitempty"`
	Tag         *ModelTag `json:"tag,omitempty"`
	PuffsMax    int       `json:"puffs_max"`
	Flavors     []string  `json:"flavors"`
	PriceCents  *int64    `json:"price_cents,omitempty"`
}

type CreateModelRequest struct {
	Name        string    `json:"name"`
	Status      string    `json:"status,omitempty"`
	Description *string   `json:"description"`
	Tag         *ModelTag `json:"tag,omitempty"`
	PuffsMax    int       `json:"puffs_max"`
	Flavors     []string  `json:"flavors"`
	PriceCents  *int64    `json:"price_cents"`
}

type UpdateModelRequest struct {
	Name        patch.Field[string]   `json:"name"`
	Status      patch.Field[string]   `json:"status"`
	Description patch.Field[string]   `json:"description"`
	Tag         patch.Field[ModelTag] `json:"tag"`
	PuffsMax    patch.Field[int]      `json:"puffs_max"`
	Flavors     patch.Field[[]string] `json:"flavors"`
	PriceCents  patch.Field[int64]    `json:"price_cents"`
}

type ListModelsResponse struct {
	Items  []Model `json:"items"`
	Limit  int     `json:"limit"`
	Offset int     `json:"offset"`
}

type ListPublicModelsResponse struct {
	Items  []PublicModel `json:"items"`
	Limit  int           `json:"limit"`
	Offset int           `json:"offset"`
}

type FlavorRequest struct {
	Value string `json:"value"`
}

// WakaModel - Gorm table
type WakaModel struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement"`
	Name        string         `gorm:"not null"`
	Status      string         `gorm:"not null;default:hidden;index"`
	Description *string        `gorm:""`
	PhotoKey    *string        `gorm:""`
	Tag         datatypes.JSON `gorm:"type:jsonb"` // JSON slice
	PuffsMax    int            `gorm:"not null"`
	Flavors     datatypes.JSON `gorm:"not null"` // JSON slice
	PriceCents  *int64         `gorm:""`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
