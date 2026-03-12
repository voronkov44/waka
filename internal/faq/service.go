package faq

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"gorm.io/datatypes"
)

type Service interface {
	// public
	ListTopics(ctx context.Context) ([]Topic, error)
	ListArticlesByTopic(ctx context.Context, topicID uint64, channel string) ([]ArticleSummary, error)
	GetArticle(ctx context.Context, id uint64) (ArticleDetail, error)
	Search(ctx context.Context, q, channel string, limit, offset int) ([]ArticleSummary, error)

	// admin read
	ListTopicsAdmin(ctx context.Context) ([]Topic, error)
	ListArticlesAdmin(ctx context.Context, filter AdminArticleFilter) (ListAdminArticlesResponse, error)
	ListArticlesAdminWithBlocks(ctx context.Context, filter AdminArticleFilter) (ListAdminArticleDetailsResponse, error)
	GetArticleAdmin(ctx context.Context, id uint64) (ArticleDetail, error)

	// admin write
	CreateTopic(ctx context.Context, req CreateTopicRequest) (Topic, error)
	UpdateTopic(ctx context.Context, id uint64, req UpdateTopicRequest) (Topic, error)

	CreateArticle(ctx context.Context, req CreateArticleRequest) (Article, error)
	UpdateArticle(ctx context.Context, id uint64, req UpdateArticleRequest) (Article, error)

	PutBlocks(ctx context.Context, articleID uint64, req PutBlocksRequest) ([]Block, error)
	CreateBlock(ctx context.Context, articleID uint64, req CreateBlockRequest) (Block, error)
	UpdateBlock(ctx context.Context, blockID uint64, req UpdateBlockRequest) (Block, error)
	DeleteBlock(ctx context.Context, blockID uint64) error
}

type service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
		now:  time.Now,
	}
}

// public

func (s *service) ListTopics(ctx context.Context) ([]Topic, error) {
	return s.repo.ListTopics(ctx, true)
}

func (s *service) ListArticlesByTopic(ctx context.Context, topicID uint64, channel string) ([]ArticleSummary, error) {
	ch, ok := normalizeChannel(channel)
	if !ok {
		return nil, ErrInvalidArgument
	}
	if topicID == 0 {
		return nil, ErrInvalidArgument
	}
	return s.repo.ListArticlesByTopic(ctx, topicID, ch)
}

func (s *service) GetArticle(ctx context.Context, id uint64) (ArticleDetail, error) {
	if id == 0 {
		return ArticleDetail{}, ErrInvalidArgument
	}
	a, blocks, err := s.repo.GetArticle(ctx, id)
	if err != nil {
		return ArticleDetail{}, err
	}
	return ArticleDetail{Article: a, Blocks: blocks}, nil
}

func (s *service) Search(ctx context.Context, q, channel string, limit, offset int) ([]ArticleSummary, error) {
	ch, ok := normalizeChannel(channel)
	if !ok {
		return nil, ErrInvalidArgument
	}
	return s.repo.SearchArticles(ctx, q, ch, limit, offset)
}

// admin read

func (s *service) ListTopicsAdmin(ctx context.Context) ([]Topic, error) {
	return s.repo.ListTopics(ctx, false)
}

func (s *service) ListArticlesAdmin(ctx context.Context, filter AdminArticleFilter) (ListAdminArticlesResponse, error) {
	norm, err := normalizeAdminFilter(filter)
	if err != nil {
		return ListAdminArticlesResponse{}, err
	}

	items, total, err := s.repo.ListArticlesAdmin(ctx, norm)
	if err != nil {
		return ListAdminArticlesResponse{}, err
	}

	out := make([]AdminArticleSummary, 0, len(items))
	for _, a := range items {
		out = append(out, AdminArticleSummary{
			ID:          a.ID,
			TopicID:     a.TopicID,
			Slug:        a.Slug,
			Title:       a.Title,
			Status:      a.Status,
			Channel:     a.Channel,
			PublishedAt: a.PublishedAt,
			UpdatedAt:   a.UpdatedAt,
		})
	}

	return ListAdminArticlesResponse{
		Items:  out,
		Limit:  norm.Limit,
		Offset: norm.Offset,
		Total:  total,
	}, nil
}

func (s *service) ListArticlesAdminWithBlocks(ctx context.Context, filter AdminArticleFilter) (ListAdminArticleDetailsResponse, error) {
	norm, err := normalizeAdminFilter(filter)
	if err != nil {
		return ListAdminArticleDetailsResponse{}, err
	}

	items, total, err := s.repo.ListArticlesAdmin(ctx, norm)
	if err != nil {
		return ListAdminArticleDetailsResponse{}, err
	}

	out := make([]ArticleDetail, 0, len(items))
	for _, a := range items {
		full, blocks, err := s.repo.GetArticleAnyStatus(ctx, a.ID)
		if err != nil {
			return ListAdminArticleDetailsResponse{}, err
		}
		out = append(out, ArticleDetail{
			Article: full,
			Blocks:  blocks,
		})
	}

	return ListAdminArticleDetailsResponse{
		Items:  out,
		Limit:  norm.Limit,
		Offset: norm.Offset,
		Total:  total,
	}, nil
}

func (s *service) GetArticleAdmin(ctx context.Context, id uint64) (ArticleDetail, error) {
	if id == 0 {
		return ArticleDetail{}, ErrInvalidArgument
	}

	a, blocks, err := s.repo.GetArticleAnyStatus(ctx, id)
	if err != nil {
		return ArticleDetail{}, err
	}

	return ArticleDetail{Article: a, Blocks: blocks}, nil
}

// admin write

func (s *service) CreateTopic(ctx context.Context, req CreateTopicRequest) (Topic, error) {
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return Topic{}, ErrInvalidArgument
	}

	sort := 0
	if req.Sort != nil {
		sort = *req.Sort
	}
	active := true
	if req.IsActive != nil {
		active = *req.IsActive
	}

	t := Topic{
		Title:    title,
		Sort:     sort,
		IsActive: active,
	}
	if err := s.repo.CreateTopic(ctx, &t); err != nil {
		return Topic{}, err
	}
	return t, nil
}

func (s *service) UpdateTopic(ctx context.Context, id uint64, req UpdateTopicRequest) (Topic, error) {
	if id == 0 {
		return Topic{}, ErrInvalidArgument
	}
	if req.Title != nil && strings.TrimSpace(*req.Title) == "" {
		return Topic{}, ErrInvalidArgument
	}
	return s.repo.UpdateTopic(ctx, id, req)
}

func (s *service) CreateArticle(ctx context.Context, req CreateArticleRequest) (Article, error) {
	if req.TopicID == 0 {
		return Article{}, ErrInvalidArgument
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		return Article{}, ErrInvalidArgument
	}

	ch := ChannelAll
	if req.Channel != nil {
		v, ok := normalizeChannel(*req.Channel)
		if !ok {
			return Article{}, ErrInvalidArgument
		}
		ch = v
	}

	st := StatusDraft
	if req.Status != nil {
		v, ok := normalizeStatus(*req.Status)
		if !ok {
			return Article{}, ErrInvalidArgument
		}
		st = v
	}

	slug := ""
	if req.Slug != nil {
		slug = strings.TrimSpace(*req.Slug)
	}
	if slug == "" {
		slug = makeSlug(title)
	}

	a := Article{
		TopicID: req.TopicID,
		Title:   title,
		Slug:    slug,
		Status:  st,
		Channel: ch,
	}
	if st == StatusPublished {
		now := s.now()
		a.PublishedAt = &now
	}

	if err := s.repo.CreateArticle(ctx, &a); err != nil {
		return Article{}, err
	}
	return a, nil
}

func (s *service) UpdateArticle(ctx context.Context, id uint64, req UpdateArticleRequest) (Article, error) {
	if id == 0 {
		return Article{}, ErrInvalidArgument
	}

	patch := req

	if patch.Title != nil {
		v := strings.TrimSpace(*patch.Title)
		if v == "" {
			return Article{}, ErrInvalidArgument
		}
		patch.Title = &v
	}
	if patch.Slug != nil {
		v := strings.TrimSpace(*patch.Slug)
		if v == "" {
			return Article{}, ErrInvalidArgument
		}
		patch.Slug = &v
	}
	if patch.Channel != nil {
		v, ok := normalizeChannel(*patch.Channel)
		if !ok {
			return Article{}, ErrInvalidArgument
		}
		patch.Channel = &v
	}
	if patch.Status != nil {
		v, ok := normalizeStatus(*patch.Status)
		if !ok {
			return Article{}, ErrInvalidArgument
		}
		patch.Status = &v
	}

	return s.repo.UpdateArticle(ctx, id, patch)
}

func (s *service) PutBlocks(ctx context.Context, articleID uint64, req PutBlocksRequest) ([]Block, error) {
	if articleID == 0 {
		return nil, ErrInvalidArgument
	}

	blocks := make([]Block, 0, len(req.Blocks))

	for _, b := range req.Blocks {
		if !isValidBlockType(b.Type) {
			return nil, ErrInvalidArgument
		}
		if !json.Valid(b.Data) {
			return nil, ErrInvalidArgument
		}

		bt := strings.ToLower(strings.TrimSpace(b.Type))

		blocks = append(blocks, Block{
			Sort: b.Sort,
			Type: bt,
			Data: b.Data,
		})
	}

	created, err := s.repo.ReplaceBlocks(ctx, articleID, blocks)
	if err != nil {
		return nil, err
	}

	if err := s.refreshArticleSearchText(ctx, articleID); err != nil {
		return nil, err
	}

	return created, nil
}

func (s *service) CreateBlock(ctx context.Context, articleID uint64, req CreateBlockRequest) (Block, error) {
	if articleID == 0 {
		return Block{}, ErrInvalidArgument
	}
	if !isValidBlockType(req.Type) {
		return Block{}, ErrInvalidArgument
	}
	if !json.Valid(req.Data) {
		return Block{}, ErrInvalidArgument
	}

	b := Block{
		ArticleID: articleID,
		Sort:      req.Sort,
		Type:      strings.ToLower(strings.TrimSpace(req.Type)),
		Data:      req.Data,
	}

	if err := s.repo.CreateBlock(ctx, &b); err != nil {
		return Block{}, err
	}

	if err := s.refreshArticleSearchText(ctx, articleID); err != nil {
		return Block{}, err
	}

	return b, nil
}

func (s *service) UpdateBlock(ctx context.Context, blockID uint64, req UpdateBlockRequest) (Block, error) {
	if blockID == 0 {
		return Block{}, ErrInvalidArgument
	}

	patch := req

	if patch.Type != nil {
		v := strings.ToLower(strings.TrimSpace(*patch.Type))
		if !isValidBlockType(v) {
			return Block{}, ErrInvalidArgument
		}
		patch.Type = &v
	}
	if patch.Data != nil && !json.Valid(*patch.Data) {
		return Block{}, ErrInvalidArgument
	}

	updated, err := s.repo.UpdateBlock(ctx, blockID, patch)
	if err != nil {
		return Block{}, err
	}

	if err := s.refreshArticleSearchText(ctx, updated.ArticleID); err != nil {
		return Block{}, err
	}

	return updated, nil
}

func (s *service) DeleteBlock(ctx context.Context, blockID uint64) error {
	if blockID == 0 {
		return ErrInvalidArgument
	}

	articleID, err := s.repo.DeleteBlock(ctx, blockID)
	if err != nil {
		return err
	}

	return s.refreshArticleSearchText(ctx, articleID)
}

// helpers

func normalizeAdminFilter(filter AdminArticleFilter) (AdminArticleFilter, error) {
	out := filter

	ch, ok := normalizeChannelFilter(filter.Channel)
	if !ok {
		return AdminArticleFilter{}, ErrInvalidArgument
	}
	st, ok := normalizeStatusFilter(filter.Status)
	if !ok {
		return AdminArticleFilter{}, ErrInvalidArgument
	}

	out.Channel = ch
	out.Status = st

	if out.Limit <= 0 || out.Limit > 100 {
		out.Limit = 20
	}
	if out.Offset < 0 {
		out.Offset = 0
	}

	return out, nil
}

func (s *service) refreshArticleSearchText(ctx context.Context, articleID uint64) error {
	_, blocks, err := s.repo.GetArticleAnyStatus(ctx, articleID)
	if err != nil {
		return err
	}

	var parts []string
	for _, b := range blocks {
		bt := strings.ToLower(strings.TrimSpace(b.Type))
		switch bt {
		case BlockText, BlockCallout, BlockBullets:
			txt := strings.TrimSpace(extractTextFromJSON(b.Data))
			if txt != "" {
				parts = append(parts, txt)
			}
		}
	}

	searchText := strings.TrimSpace(strings.Join(parts, " "))
	if searchText == "" {
		return s.repo.UpdateArticleSearchText(ctx, articleID, nil)
	}

	return s.repo.UpdateArticleSearchText(ctx, articleID, &searchText)
}

func makeSlug(title string) string {
	s := strings.ToLower(strings.TrimSpace(title))
	s = strings.ReplaceAll(s, "—", "-")
	s = strings.ReplaceAll(s, "–", "-")
	s = strings.ReplaceAll(s, "_", "-")
	s = strings.ReplaceAll(s, " ", "-")

	var out []rune
	prevDash := false
	for _, r := range s {
		isAZ := r >= 'a' && r <= 'z'
		is09 := r >= '0' && r <= '9'
		if isAZ || is09 {
			out = append(out, r)
			prevDash = false
			continue
		}
		if r == '-' && !prevDash {
			out = append(out, '-')
			prevDash = true
		}
	}
	res := strings.Trim(string(out), "-")
	if res == "" {
		return "faq"
	}
	return res
}

func extractTextFromJSON(raw datatypes.JSON) string {
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return ""
	}

	var parts []string
	for _, v := range m {
		switch x := v.(type) {
		case string:
			parts = append(parts, x)
		case []any:
			for _, it := range x {
				if s, ok := it.(string); ok {
					parts = append(parts, s)
				}
			}
		}
	}

	return strings.Join(parts, " ")
}
