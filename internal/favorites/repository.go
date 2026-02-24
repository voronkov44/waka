package favorites

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"rest_waka/internal/models"
)

type RepositoryGorm interface {
	Add(ctx context.Context, userID, modelID uint64) error
	Remove(ctx context.Context, userID, modelID uint64) error
	ListModelsFavorites(ctx context.Context, userID uint64, limit, offset int) ([]models.WakaModel, error)
}

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Add(ctx context.Context, userID, modelID uint64) error {
	if userID == 0 || modelID == 0 {
		return ErrInvalidArgument
	}

	rec := Favorite{
		UserID:  userID,
		ModelID: modelID,
	}

	err := r.db.WithContext(ctx).Create(&rec).Error
	if err == nil {
		return err
	}

	// pg errors: unique
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return ErrAlreadyExists
		case "23503":
			// fk model_id not found
			return ErrNotFound
		}
	}
	return err
}

func (r *GormRepository) Remove(ctx context.Context, userID, modelID uint64) error {
	if userID == 0 || modelID == 0 {
		return ErrInvalidArgument
	}
	res := r.db.WithContext(ctx).
		Where("user_id = ? AND model_id = ?", userID, modelID).
		Delete(&Favorite{})

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *GormRepository) ListModelsFavorites(ctx context.Context, userID uint64, limit, offset int) ([]models.WakaModel, error) {
	if userID == 0 {
		return nil, ErrInvalidArgument
	}

	var list []models.WakaModel
	q := r.db.WithContext(ctx).
		Table("waka_models").
		Joins("JOIN favorites f ON f.model_id = waka_models.id").
		Where("f.user_id = ?", userID).
		Order("f.created_at DESC")

	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}

	if err := q.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
