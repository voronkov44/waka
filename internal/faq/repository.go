package faq

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type Repository interface {
	// topics
	ListTopics(ctx context.Context, activeOnly bool) ([]Topic, error)
	GetTopic(ctx context.Context, id uint64) (Topic, error)
	CreateTopic(ctx context.Context, t *Topic) error
	UpdateTopic(ctx context.Context, id uint64, patch UpdateTopicRequest) (Topic, error)
	DeleteTopic(ctx context.Context, id uint64) error

	// public articles
	ListArticlesByTopic(ctx context.Context, topicID uint64, channel string) ([]ArticleSummary, error)
	GetArticle(ctx context.Context, id uint64) (Article, []Block, error)
	SearchArticles(ctx context.Context, q string, channel string, limit, offset int) ([]ArticleSummary, error)

	// admin articles
	ListArticlesAdmin(ctx context.Context, filter AdminArticleFilter) ([]Article, int64, error)
	GetArticleAnyStatus(ctx context.Context, id uint64) (Article, []Block, error)

	CreateArticle(ctx context.Context, a *Article) error
	UpdateArticle(ctx context.Context, id uint64, patch UpdateArticleRequest) (Article, error)
	DeleteArticle(ctx context.Context, id uint64) error

	ReplaceBlocks(ctx context.Context, articleID uint64, blocks []Block) ([]Block, error)
	CreateBlock(ctx context.Context, b *Block) error
	UpdateBlock(ctx context.Context, id uint64, patch UpdateBlockRequest) (Block, error)
	DeleteBlock(ctx context.Context, id uint64) (uint64, error)

	UpdateArticleSearchText(ctx context.Context, articleID uint64, searchText *string) error
}

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) ListTopics(ctx context.Context, activeOnly bool) ([]Topic, error) {
	var out []Topic
	q := r.db.WithContext(ctx).Model(&Topic{})
	if activeOnly {
		q = q.Where("is_active = TRUE")
	}
	err := q.Order("sort asc, title asc").Find(&out).Error
	return out, err
}

func (r *GormRepository) GetTopic(ctx context.Context, id uint64) (Topic, error) {
	var t Topic
	err := r.db.WithContext(ctx).First(&t, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return Topic{}, ErrNotFound
	}
	return t, err
}

func (r *GormRepository) CreateTopic(ctx context.Context, t *Topic) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *GormRepository) UpdateTopic(ctx context.Context, id uint64, patch UpdateTopicRequest) (Topic, error) {
	var t Topic
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&t, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}

		updates := map[string]any{}
		if patch.Title != nil {
			updates["title"] = strings.TrimSpace(*patch.Title)
		}
		if patch.Sort != nil {
			updates["sort"] = *patch.Sort
		}
		if patch.IsActive != nil {
			updates["is_active"] = *patch.IsActive
		}
		if len(updates) == 0 {
			return nil
		}

		if err := tx.Model(&Topic{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			return err
		}
		return tx.First(&t, "id = ?", id).Error
	})
	if err != nil {
		return Topic{}, err
	}
	return t, nil
}

func (r *GormRepository) DeleteTopic(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var t Topic
		if err := tx.First(&t, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}

		var count int64
		if err := tx.Model(&Article{}).Where("topic_id = ?", id).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return ErrConflict
		}

		if err := tx.Delete(&Topic{}, "id = ?", id).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *GormRepository) ListArticlesByTopic(ctx context.Context, topicID uint64, channel string) ([]ArticleSummary, error) {
	var out []ArticleSummary

	q := r.db.WithContext(ctx).
		Model(&Article{}).
		Select("id, topic_id, slug, title, updated_at").
		Where("topic_id = ?", topicID).
		Where("status = ?", StatusPublished)

	if channel != ChannelAll {
		q = q.Where("channel IN (?, ?)", ChannelAll, channel)
	}

	err := q.Order("updated_at desc, id desc").Find(&out).Error
	return out, err
}

func (r *GormRepository) GetArticle(ctx context.Context, id uint64) (Article, []Block, error) {
	var a Article
	err := r.db.WithContext(ctx).First(&a, "id = ? AND status = ?", id, StatusPublished).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return Article{}, nil, ErrNotFound
	}
	if err != nil {
		return Article{}, nil, err
	}

	var blocks []Block
	if err := r.db.WithContext(ctx).
		Model(&Block{}).
		Where("article_id = ?", a.ID).
		Order("sort asc, id asc").
		Find(&blocks).Error; err != nil {
		return Article{}, nil, err
	}

	return a, blocks, nil
}

func (r *GormRepository) GetArticleAnyStatus(ctx context.Context, id uint64) (Article, []Block, error) {
	var a Article
	err := r.db.WithContext(ctx).First(&a, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return Article{}, nil, ErrNotFound
	}
	if err != nil {
		return Article{}, nil, err
	}

	var blocks []Block
	if err := r.db.WithContext(ctx).
		Model(&Block{}).
		Where("article_id = ?", a.ID).
		Order("sort asc, id asc").
		Find(&blocks).Error; err != nil {
		return Article{}, nil, err
	}

	return a, blocks, nil
}

func (r *GormRepository) SearchArticles(ctx context.Context, q string, channel string, limit, offset int) ([]ArticleSummary, error) {
	var out []ArticleSummary

	dbq := r.db.WithContext(ctx).
		Model(&Article{}).
		Select("id, topic_id, slug, title, updated_at").
		Where("status = ?", StatusPublished)

	if channel != ChannelAll {
		dbq = dbq.Where("channel IN (?, ?)", ChannelAll, channel)
	}

	if strings.TrimSpace(q) != "" {
		like := "%" + strings.TrimSpace(q) + "%"
		dbq = dbq.Where("(title ILIKE ? OR search_text ILIKE ?)", like, like)
	}

	if limit <= 0 || limit > 50 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	err := dbq.Order("updated_at desc, id desc").Limit(limit).Offset(offset).Find(&out).Error
	return out, err
}

func (r *GormRepository) ListArticlesAdmin(ctx context.Context, filter AdminArticleFilter) ([]Article, int64, error) {
	var items []Article
	var total int64

	q := r.db.WithContext(ctx).Model(&Article{})

	if filter.TopicID != nil {
		q = q.Where("topic_id = ?", *filter.TopicID)
	}
	if filter.Status != "" {
		q = q.Where("status = ?", filter.Status)
	}
	if filter.Channel != "" {
		q = q.Where("channel = ?", filter.Channel)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	limit := filter.Limit
	offset := filter.Offset
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	err := q.Order("updated_at desc, id desc").
		Limit(limit).
		Offset(offset).
		Find(&items).Error
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *GormRepository) CreateArticle(ctx context.Context, a *Article) error {
	err := r.db.WithContext(ctx).Create(a).Error
	if isUniqueViolation(err) {
		return ErrConflict
	}
	return err
}

func (r *GormRepository) UpdateArticle(ctx context.Context, id uint64, patch UpdateArticleRequest) (Article, error) {
	var a Article

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&a, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}

		updates := map[string]any{}
		if patch.TopicID != nil {
			updates["topic_id"] = *patch.TopicID
		}
		if patch.Title != nil {
			updates["title"] = strings.TrimSpace(*patch.Title)
		}
		if patch.Slug != nil {
			updates["slug"] = strings.TrimSpace(*patch.Slug)
		}
		if patch.Status != nil {
			updates["status"] = *patch.Status

			if *patch.Status == StatusPublished && a.Status != StatusPublished {
				now := time.Now()
				updates["published_at"] = &now
			}
			if *patch.Status != StatusPublished {
				updates["published_at"] = nil
			}
		}
		if patch.Channel != nil {
			updates["channel"] = *patch.Channel
		}
		if len(updates) == 0 {
			return nil
		}

		if err := tx.Model(&Article{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			if isUniqueViolation(err) {
				return ErrConflict
			}
			return err
		}

		return tx.First(&a, "id = ?", id).Error
	})
	if err != nil {
		return Article{}, err
	}
	return a, nil
}

func (r *GormRepository) DeleteArticle(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var a Article
		if err := tx.First(&a, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}

		if err := tx.Where("article_id = ?", id).Delete(&Block{}).Error; err != nil {
			return err
		}

		if err := tx.Delete(&Article{}, "id = ?", id).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *GormRepository) ReplaceBlocks(ctx context.Context, articleID uint64, blocks []Block) ([]Block, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var a Article
		if err := tx.First(&a, "id = ?", articleID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}

		if err := tx.Where("article_id = ?", articleID).Delete(&Block{}).Error; err != nil {
			return err
		}

		if len(blocks) == 0 {
			return nil
		}

		for i := range blocks {
			blocks[i].ArticleID = articleID
		}

		if err := tx.Create(&blocks).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return blocks, nil
}

func (r *GormRepository) CreateBlock(ctx context.Context, b *Block) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var a Article
		if err := tx.First(&a, "id = ?", b.ArticleID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}
		return tx.Create(b).Error
	})
}

func (r *GormRepository) UpdateBlock(ctx context.Context, id uint64, patch UpdateBlockRequest) (Block, error) {
	var out Block

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&out, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}

		updates := map[string]any{}
		if patch.Sort != nil {
			updates["sort"] = *patch.Sort
		}
		if patch.Type != nil {
			updates["type"] = *patch.Type
		}
		if patch.Data != nil {
			updates["data"] = *patch.Data
		}
		if len(updates) == 0 {
			return nil
		}

		if err := tx.Model(&Block{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			return err
		}

		return tx.First(&out, "id = ?", id).Error
	})
	if err != nil {
		return Block{}, err
	}

	return out, nil
}

func (r *GormRepository) DeleteBlock(ctx context.Context, id uint64) (uint64, error) {
	var blk Block

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&blk, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}

		if err := tx.Delete(&Block{}, "id = ?", id).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	return blk.ArticleID, nil
}

func (r *GormRepository) UpdateArticleSearchText(ctx context.Context, articleID uint64, searchText *string) error {
	return r.db.WithContext(ctx).
		Model(&Article{}).
		Where("id = ?", articleID).
		Update("search_text", searchText).Error
}

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
