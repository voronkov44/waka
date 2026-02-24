package favorites

import (
	"rest_waka/internal/models"
	"time"
)

type Favorite struct {
	UserID  uint64 `gorm:"primaryKey"`
	ModelID uint64 `gorm:"primaryKey"`

	CreatedAt time.Time
}

func (Favorite) TableName() string {
	return "favorites"
}

type ListFavoritesResponse struct {
	Items  []models.Model `json:"items"`
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
}
