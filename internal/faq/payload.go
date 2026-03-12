package faq

import (
	"gorm.io/datatypes"
	"strings"
	"time"
)

const (
	ChannelAll     = "all"
	ChannelTG      = "tg"
	ChannelMiniApp = "miniapp"
)

const (
	StatusDraft     = "draft"
	StatusPublished = "published"
	StatusArchived  = "archived"
)

const (
	BlockText    = "text"
	BlockImage   = "image"
	BlockLink    = "link"
	BlockBullets = "bullets"
	BlockDivider = "divider"
	BlockCallout = "callout"
)

// Gorm db

type Topic struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title"`
	Sort      int       `json:"sort"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Topic) TableName() string {
	return "faq_topics"
}

type Article struct {
	ID          uint64     `json:"id" gorm:"primaryKey"`
	TopicID     uint64     `json:"topic_id" gorm:"index"`
	Slug        string     `json:"slug" gorm:"uniqueIndex;size:255"`
	Title       string     `json:"title"`
	Status      string     `json:"status" gorm:"size:16;index"`
	Channel     string     `json:"channel" gorm:"size:16;index"`
	SearchText  *string    `json:"search_text,omitempty" gorm:"type:text"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (Article) TableName() string {
	return "faq_articles"
}

type Block struct {
	ID        uint64         `json:"id" gorm:"primaryKey"`
	ArticleID uint64         `json:"article_id" gorm:"index"`
	Sort      int            `json:"sort"`
	Type      string         `json:"type" gorm:"size:16"`
	Data      datatypes.JSON `json:"data" gorm:"type:jsonb"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (Block) TableName() string {
	return "faq_blocks"
}

// DTO

type ArticleSummary struct {
	ID        uint64    `json:"id"`
	TopicID   uint64    `json:"topic_id"`
	Slug      string    `json:"slug"`
	Title     string    `json:"title"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ArticleDetail struct {
	Article
	Blocks []Block `json:"blocks"`
}

type AdminArticleSummary struct {
	ID          uint64     `json:"id"`
	TopicID     uint64     `json:"topic_id"`
	Slug        string     `json:"slug"`
	Title       string     `json:"title"`
	Status      string     `json:"status"`
	Channel     string     `json:"channel"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type ListAdminArticlesResponse struct {
	Items  []AdminArticleSummary `json:"items"`
	Limit  int                   `json:"limit"`
	Offset int                   `json:"offset"`
	Total  int64                 `json:"total"`
}

type ListAdminArticleDetailsResponse struct {
	Items  []ArticleDetail `json:"items"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
	Total  int64           `json:"total"`
}

type AdminArticleFilter struct {
	TopicID *uint64 `json:"topic_id,omitempty"`
	Channel string  `json:"channel,omitempty"` // "" = все, all/tg/miniapp = точный фильтр
	Status  string  `json:"status,omitempty"`  // "" = все, draft/published/archived = точный фильтр
	Limit   int     `json:"limit"`
	Offset  int     `json:"offset"`
}

// request

type CreateTopicRequest struct {
	Title    string `json:"title"`
	Sort     *int   `json:"sort,omitempty"`
	IsActive *bool  `json:"is_active,omitempty"`
}

type UpdateTopicRequest struct {
	Title    *string `json:"title,omitempty"`
	Sort     *int    `json:"sort,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}

type CreateArticleRequest struct {
	TopicID uint64  `json:"topic_id"`
	Title   string  `json:"title"`
	Slug    *string `json:"slug,omitempty"`
	Status  *string `json:"status,omitempty"`  // draft/published/archived
	Channel *string `json:"channel,omitempty"` // all/tg/miniapp
}

type UpdateArticleRequest struct {
	TopicID *uint64 `json:"topic_id,omitempty"`
	Title   *string `json:"title,omitempty"`
	Slug    *string `json:"slug,omitempty"`
	Status  *string `json:"status,omitempty"`
	Channel *string `json:"channel,omitempty"`
}

type PutBlocksRequest struct {
	Blocks []PutBlock `json:"blocks"`
}

type PutBlock struct {
	Sort int            `json:"sort"`
	Type string         `json:"type"`
	Data datatypes.JSON `json:"data"`
}

type CreateBlockRequest struct {
	Sort int            `json:"sort"`
	Type string         `json:"type"`
	Data datatypes.JSON `json:"data"`
}

type UpdateBlockRequest struct {
	Sort *int            `json:"sort,omitempty"`
	Type *string         `json:"type,omitempty"`
	Data *datatypes.JSON `json:"data,omitempty"`
}

// helpers

func normalizeChannel(s string) (string, bool) {
	v := strings.ToLower(strings.TrimSpace(s))
	switch v {
	case "", ChannelAll:
		return ChannelAll, true
	case ChannelTG:
		return ChannelTG, true
	case ChannelMiniApp:
		return ChannelMiniApp, true
	default:
		return "", false
	}
}

func normalizeChannelFilter(s string) (string, bool) {
	v := strings.ToLower(strings.TrimSpace(s))
	switch v {
	case "":
		return "", true
	case ChannelAll:
		return ChannelAll, true
	case ChannelTG:
		return ChannelTG, true
	case ChannelMiniApp:
		return ChannelMiniApp, true
	default:
		return "", false
	}
}

func normalizeStatus(s string) (string, bool) {
	v := strings.ToLower(strings.TrimSpace(s))
	switch v {
	case "", StatusDraft:
		return StatusDraft, true
	case StatusPublished:
		return StatusPublished, true
	case StatusArchived:
		return StatusArchived, true
	default:
		return "", false
	}
}

func normalizeStatusFilter(s string) (string, bool) {
	v := strings.ToLower(strings.TrimSpace(s))
	switch v {
	case "":
		return "", true
	case StatusDraft:
		return StatusDraft, true
	case StatusPublished:
		return StatusPublished, true
	case StatusArchived:
		return StatusArchived, true
	default:
		return "", false
	}
}

func isValidBlockType(t string) bool {
	switch strings.ToLower(strings.TrimSpace(t)) {
	case BlockText, BlockImage, BlockLink, BlockBullets, BlockDivider, BlockCallout:
		return true
	default:
		return false
	}
}
