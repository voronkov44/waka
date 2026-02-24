package users

import (
	"context"
	"errors"
	"gorm.io/gorm"
)

type RepositoryGorm interface {
	Get(ctx context.Context, id uint64) (User, error)
	List(ctx context.Context, limit, offset int) ([]User, error)
}

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Get(ctx context.Context, id uint64) (User, error) {
	var user User
	res := r.db.WithContext(ctx).First(&user, id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return User{}, ErrNotFound
	}
	return user, res.Error
}

func (r *GormRepository) List(ctx context.Context, limit, offset int) ([]User, error) {
	var users []User
	query := r.db.WithContext(ctx).Model(&User{}).Order("id DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
