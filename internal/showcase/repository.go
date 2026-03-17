package showcase

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RepositoryGorm interface {
	Create(ctx context.Context, rec *Showcase) error
	Get(ctx context.Context, id uint64) (Showcase, error)
	List(ctx context.Context, limit, offset int) ([]Showcase, error)
	Save(ctx context.Context, rec *Showcase) error
	Delete(ctx context.Context, id uint64) error

	UpdatePhotoKey(ctx context.Context, id uint64, key *string) (Showcase, error)

	ListActive(ctx context.Context, limit, offset int) ([]Showcase, error)
	GetActive(ctx context.Context, id uint64) (Showcase, error)
}

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(ctx context.Context, rec *Showcase) error {
	return r.db.WithContext(ctx).Create(rec).Error
}

func (r *GormRepository) Get(ctx context.Context, id uint64) (Showcase, error) {
	var rec Showcase
	res := r.db.WithContext(ctx).First(&rec, id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return Showcase{}, ErrNotFound
	}
	return rec, res.Error
}

func (r *GormRepository) List(ctx context.Context, limit, offset int) ([]Showcase, error) {
	var list []Showcase

	q := r.db.WithContext(ctx).
		Model(&Showcase{}).
		Order("sort ASC, id DESC")

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

func (r *GormRepository) Save(ctx context.Context, rec *Showcase) error {
	return r.db.WithContext(ctx).Save(rec).Error
}

func (r *GormRepository) Delete(ctx context.Context, id uint64) error {
	res := r.db.WithContext(ctx).Delete(&Showcase{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *GormRepository) UpdatePhotoKey(ctx context.Context, id uint64, key *string) (Showcase, error) {
	var out Showcase

	tx := r.db.WithContext(ctx).
		Model(&Showcase{}).
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"photo_key": key,
		}).
		Scan(&out)

	if tx.Error != nil {
		return Showcase{}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return Showcase{}, ErrNotFound
	}
	return out, nil
}

func (r *GormRepository) ListActive(ctx context.Context, limit, offset int) ([]Showcase, error) {
	var list []Showcase

	q := r.db.WithContext(ctx).
		Model(&Showcase{}).
		Where("is_active = ?", true).
		Order("sort ASC, id DESC")

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

func (r *GormRepository) GetActive(ctx context.Context, id uint64) (Showcase, error) {
	var rec Showcase

	res := r.db.WithContext(ctx).
		Where("id = ? AND is_active = ?", id, true).
		First(&rec)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return Showcase{}, ErrNotFound
	}
	return rec, res.Error
}
