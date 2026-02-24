package auth

import "time"

type TelegramProfile struct {
	TgID      int64  `json:"tg_id" validate:"required,gt=0"`
	Username  string `json:"username,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	PhotoURL  string `json:"photo_url,omitempty"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type MeResponse struct {
	ID        uint64    `json:"id"`
	TgID      int64     `json:"tg_id"`
	Username  string    `json:"username,omitempty"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	PhotoURL  string    `json:"photo_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// User - Gorm table
type User struct {
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	TgID      int64  `gorm:"not null;uniqueIndex;index"`
	Username  string `gorm:""`
	FirstName string `gorm:""`
	LastName  string `gorm:""`
	PhotoURL  string `gorm:""`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (User) TableName() string {
	return "telegram_users"
}
