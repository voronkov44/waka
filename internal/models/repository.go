package models

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RepositoryGorm interface {
	Create(ctx context.Context, rec *WakaModel) error
	Get(ctx context.Context, id uint64) (WakaModel, error)
	List(ctx context.Context, limit, offset int) ([]WakaModel, error)
	Save(ctx context.Context, rec *WakaModel) error
	Delete(ctx context.Context, id uint64) error

	UpdatePhotoKey(ctx context.Context, id uint64, key *string) (WakaModel, error)

	ListByStatus(ctx context.Context, status string, limit, offset int) ([]WakaModel, error)
	GetByIDAndStatus(ctx context.Context, id uint64, status string) (WakaModel, error)
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
	q := r.db.WithContext(ctx).Model(&WakaModel{}).Order("puffs_max DESC").Order("id DESC")

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

func (r *GormRepository) UpdatePhotoKey(ctx context.Context, id uint64, key *string) (WakaModel, error) {
	var out WakaModel

	tx := r.db.WithContext(ctx).
		Model(&WakaModel{}).
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"photo_key": key,
		}).
		Scan(&out)

	if tx.Error != nil {
		return WakaModel{}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return WakaModel{}, ErrNotFound
	}
	return out, nil
}

func (r *GormRepository) ListByStatus(ctx context.Context, status string, limit, offset int) ([]WakaModel, error) {
	var list []WakaModel

	q := r.db.WithContext(ctx).
		Model(&WakaModel{}).
		Where("status = ?", status).
		Order("puffs_max DESC").
		Order("id DESC")

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

func (r *GormRepository) GetByIDAndStatus(ctx context.Context, id uint64, status string) (WakaModel, error) {
	var rec WakaModel

	res := r.db.WithContext(ctx).
		Where("id = ? AND status = ?", id, status).
		First(&rec)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return WakaModel{}, ErrNotFound
	}
	return rec, res.Error
}
