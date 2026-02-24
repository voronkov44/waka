package users

import "time"

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

type UserView struct {
	ID        uint64    `json:"id"`
	TgID      int64     `json:"tg_id"`
	Username  string    `json:"username,omitempty"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	PhotoURL  string    `json:"photo_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListUsersResponse struct {
	Items  []UserView `json:"items"`
	Limit  int        `json:"limit"`
	Offset int        `json:"offset"`
}
