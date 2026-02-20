package models

import (
	"context"
	"errors"
	"gorm.io/gorm"
)

type RepositoryGorm interface {
	Create(ctx context.Context, rec *WakaModel) error
	Get(ctx context.Context, id uint64) (WakaModel, error)
	List(ctx context.Context, limit, offset int) ([]WakaModel, error)
	Save(ctx context.Context, rec *WakaModel) error
	Delete(ctx context.Context, id uint64) error
}

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{
		db: db,
	}
}

func (r *GormRepository) Create(ctx context.Context, rec *WakaModel) error {
	return r.db.WithContext(ctx).Create(rec).Error
}

func (r *GormRepository) Get(ctx context.Context, id uint64) (WakaModel, error) {
	var rec WakaModel
	res := r.db.WithContext(ctx).First(&rec, id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return WakaModel{}, ErrNotFound
	}
	return rec, res.Error
}

func (r *GormRepository) List(ctx context.Context, limit, offset int) ([]WakaModel, error) {
	var list []WakaModel
	q := r.db.WithContext(ctx).Model(&WakaModel{}).Order("id DESC")

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

func (r *GormRepository) Save(ctx context.Context, rec *WakaModel) error {
	return r.db.WithContext(ctx).Save(rec).Error
}

func (r *GormRepository) Delete(ctx context.Context, id uint64) error {
	res := r.db.WithContext(ctx).Delete(&WakaModel{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
