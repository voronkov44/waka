package auth

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"rest_waka/internal/users"
	"time"
)

type RepositoryGorm interface {
	UpsertTelegram(ctx context.Context, tg TelegramProfile) (users.User, error)
	Get(ctx context.Context, id uint64) (users.User, error)
}

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) UpsertTelegram(ctx context.Context, tg TelegramProfile) (users.User, error) {
	if tg.TgID <= 0 {
		return users.User{}, ErrInvalidArgument
	}

	u := users.User{
		TgID:      tg.TgID,
		Username:  tg.Username,
		FirstName: tg.FirstName,
		LastName:  tg.LastName,
		PhotoURL:  tg.PhotoURL,
	}

	now := time.Now()

	updates := map[string]any{
		"username":   tg.Username,
		"first_name": tg.FirstName,
		"last_name":  tg.LastName,
		"photo_url":  tg.PhotoURL,
		"updated_at": now,
	}

	err := r.db.WithContext(ctx).
		Clauses(
			clause.OnConflict{
				Columns:   []clause.Column{{Name: "tg_id"}},
				DoUpdates: clause.Assignments(updates),
			},
			clause.Returning{},
		).
		Create(&u).Error
	if err != nil {
		return users.User{}, err
	}

	return u, nil
}

func (r *GormRepository) Get(ctx context.Context, id uint64) (users.User, error) {
	var u users.User
	res := r.db.WithContext(ctx).First(&u, id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return users.User{}, ErrNotFound
	}
	return u, res.Error
}
