package models

import (
	"context"
	"errors"
	"gorm.io/gorm"
)

type RepositoryGorm interface {
	Create(ctx context.Context, rec *ModelRecord) error
	Get(ctx context.Context, id uint64) (ModelRecord, error)
	List(ctx context.Context, limit, offset int) ([]ModelRecord, error)
	Save(ctx context.Context, rec *ModelRecord) error
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

func (r *GormRepository) Create(ctx context.Context, rec *ModelRecord) error {
	return r.db.WithContext(ctx).Create(rec).Error
}

func (r *GormRepository) Get(ctx context.Context, id uint64) (ModelRecord, error) {
	var rec ModelRecord
	res := r.db.WithContext(ctx).First(&rec, id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return ModelRecord{}, ErrNotFound
	}
	return rec, res.Error
}

func (r *GormRepository) List(ctx context.Context, limit, offset int) ([]ModelRecord, error) {
	var list []ModelRecord
	q := r.db.WithContext(ctx).Model(&ModelRecord{}).Order("id DESC")

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

func (r *GormRepository) Save(ctx context.Context, rec *ModelRecord) error {
	return r.db.WithContext(ctx).Save(rec).Error
}

func (r *GormRepository) Delete(ctx context.Context, id uint64) error {
	res := r.db.WithContext(ctx).Delete(&ModelRecord{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
