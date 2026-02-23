package auth

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type RepositoryGorm interface {
	UpsertTelegram(ctx context.Context, tg TelegramProfile) (User, error)
	Get(ctx context.Context, id uint64) (User, error)
}

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) UpsertTelegram(ctx context.Context, tg TelegramProfile) (User, error) {
	if tg.TgID <= 0 {
		return User{}, ErrInvalidArgument
	}

	u := User{
		TgID:      tg.TgID,
		Username:  tg.Username,
		FirstName: tg.FirstName,
		LastName:  tg.LastName,
		PhotoURL:  tg.PhotoURL,
	}

	now := time.Now()

	err := r.db.WithContext(ctx).
		Clauses(
			clause.OnConflict{
				Columns: []clause.Column{{Name: "tg_id"}},
				DoUpdates: clause.Assignments(map[string]any{
					"username":   u.Username,
					"first_name": u.FirstName,
					"last_name":  u.LastName,
					"photo_url":  u.PhotoURL,
					"updated_at": now,
				}),
			},
			clause.Returning{}, // чтобы вернулся id и при update тоже
		).
		Create(&u).Error

	if err != nil {
		return User{}, err
	}

	return u, nil
}

func (r *GormRepository) Get(ctx context.Context, id uint64) (User, error) {
	var u User
	res := r.db.WithContext(ctx).First(&u, id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return User{}, ErrNotFound
	}
	return u, nil
}
