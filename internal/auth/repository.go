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

	// update: только непустые поля, иначе не обновляем
	updates := map[string]any{
		"updated_at": now,
	}
	if tg.Username != "" {
		updates["username"] = tg.Username
	}
	if tg.FirstName != "" {
		updates["first_name"] = tg.FirstName
	}
	if tg.LastName != "" {
		updates["last_name"] = tg.LastName
	}
	if tg.PhotoURL != "" {
		updates["photo_url"] = tg.PhotoURL
	}

	err := r.db.WithContext(ctx).
		Clauses(
			clause.OnConflict{
				Columns:   []clause.Column{{Name: "tg_id"}},
				DoUpdates: clause.Assignments(updates),
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
	return u, res.Error
}
