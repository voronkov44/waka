package showcase

import (
	"rest_waka/pkg/patch"
	"time"
)

type ItemTag struct {
	Label     string `json:"label"`
	BgColor   string `json:"bg_color"`
	TextColor string `json:"text_color"`
	Outlined  bool   `json:"outlined"`
}

type Item struct {
	ID          uint64  `json:"id"`
	Tag         ItemTag `json:"tag"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	ModelID     uint64  `json:"model_id"`

	PhotoKey *string `json:"photo_key,omitempty"`
	PhotoURL *string `json:"photo_url,omitempty"`

	Sort      int       `json:"sort"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PublicItem struct {
	ID          uint64    `json:"id"`
	Tag         ItemTag   `json:"tag"`
	Title       string    `json:"title"`
	Description *string   `json:"description,omitempty"`
	ModelID     uint64    `json:"model_id"`
	PhotoURL    *string   `json:"photo_url,omitempty"`
	Sort        int       `json:"sort"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateItemRequest struct {
	Tag         ItemTag `json:"tag"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	ModelID     uint64  `json:"model_id"`
	Sort        *int    `json:"sort,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

type UpdateItemRequest struct {
	Tag         patch.Field[ItemTag] `json:"tag"`
	Title       patch.Field[string]  `json:"title"`
	Description patch.Field[string]  `json:"description"`
	ModelID     patch.Field[uint64]  `json:"model_id"`
	Sort        patch.Field[int]     `json:"sort"`
	IsActive    patch.Field[bool]    `json:"is_active"`
}

type ListItemsResponse struct {
	Items  []Item `json:"items"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

type ListPublicItemsResponse struct {
	Items  []PublicItem `json:"items"`
	Limit  int          `json:"limit"`
	Offset int          `json:"offset"`
}

// Showcase - gorm table
type Showcase struct {
	ID           uint64 `gorm:"primaryKey;autoIncrement"`
	TagLabel     string `gorm:"not null"`
	TagBgColor   string `gorm:"not null"`
	TagTextColor string `gorm:"not null"`
	TagOutlined  bool   `gorm:"not null;default:false"`

	Title       string  `gorm:"not null"`
	Description *string `gorm:""`

	ModelID uint64 `gorm:"not null;index"`

	PhotoKey *string `gorm:""`

	Sort     int  `gorm:"not null;default:0;index"`
	IsActive bool `gorm:"not null;default:true;index"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Showcase) TableName() string {
	return "showcase_items"
}
