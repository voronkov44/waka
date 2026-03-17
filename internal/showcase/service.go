package showcase

import (
	"context"
	"errors"
	"strings"

	modelspkg "rest_waka/internal/models"
)

type modelLookup interface {
	Get(ctx context.Context, id uint64) (modelspkg.WakaModel, error)
}

type Service struct {
	repo   RepositoryGorm
	models modelLookup
}

func NewService(repo RepositoryGorm, modelsRepo modelLookup) *Service {
	return &Service{
		repo:   repo,
		models: modelsRepo,
	}
}

func (s *Service) Create(ctx context.Context, req CreateItemRequest) (Item, error) {
	title := strings.TrimSpace(req.Title)
	if title == "" || req.ModelID == 0 {
		return Item{}, ErrInvalidArgument
	}

	if err := s.ensureModelExists(ctx, req.ModelID); err != nil {
		return Item{}, err
	}

	tag, err := normalizeTag(req.Tag)
	if err != nil {
		return Item{}, err
	}

	sort := 0
	if req.Sort != nil {
		sort = *req.Sort
	}

	active := true
	if req.IsActive != nil {
		active = *req.IsActive
	}

	rec := Showcase{
		TagLabel:     tag.Label,
		TagBgColor:   tag.BgColor,
		TagTextColor: tag.TextColor,
		TagOutlined:  tag.Outlined,
		Title:        title,
		Description:  normalizeOptionalString(req.Description),
		ModelID:      req.ModelID,
		Sort:         sort,
		IsActive:     active,
	}

	if err := s.repo.Create(ctx, &rec); err != nil {
		return Item{}, err
	}

	return s.toAPI(rec), nil
}

func (s *Service) List(ctx context.Context, limit, offset int) (ListItemsResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	recs, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return ListItemsResponse{}, err
	}

	items := make([]Item, 0, len(recs))
	for _, rec := range recs {
		items = append(items, s.toAPI(rec))
	}

	return ListItemsResponse{
		Items:  items,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (s *Service) Get(ctx context.Context, id uint64) (Item, error) {
	if id == 0 {
		return Item{}, ErrInvalidArgument
	}

	rec, err := s.repo.Get(ctx, id)
	if err != nil {
		return Item{}, err
	}

	return s.toAPI(rec), nil
}

func (s *Service) Update(ctx context.Context, id uint64, req UpdateItemRequest) (Item, error) {
	if id == 0 {
		return Item{}, ErrInvalidArgument
	}
	if isEmptyPatch(req) {
		return Item{}, ErrInvalidArgument
	}

	rec, err := s.repo.Get(ctx, id)
	if err != nil {
		return Item{}, err
	}

	if req.Tag.Set {
		if req.Tag.Null {
			return Item{}, ErrInvalidArgument
		}
		tag, err := normalizeTag(req.Tag.Value)
		if err != nil {
			return Item{}, err
		}
		rec.TagLabel = tag.Label
		rec.TagBgColor = tag.BgColor
		rec.TagTextColor = tag.TextColor
		rec.TagOutlined = tag.Outlined
	}

	if req.Title.Set {
		if req.Title.Null {
			return Item{}, ErrInvalidArgument
		}
		title := strings.TrimSpace(req.Title.Value)
		if title == "" {
			return Item{}, ErrInvalidArgument
		}
		rec.Title = title
	}

	if req.Description.Set {
		if req.Description.Null {
			rec.Description = nil
		} else {
			rec.Description = normalizeOptionalString(&req.Description.Value)
		}
	}

	if req.ModelID.Set {
		if req.ModelID.Null || req.ModelID.Value == 0 {
			return Item{}, ErrInvalidArgument
		}
		if err := s.ensureModelExists(ctx, req.ModelID.Value); err != nil {
			return Item{}, err
		}
		rec.ModelID = req.ModelID.Value
	}

	if req.Sort.Set {
		if req.Sort.Null {
			return Item{}, ErrInvalidArgument
		}
		rec.Sort = req.Sort.Value
	}

	if req.IsActive.Set {
		if req.IsActive.Null {
			return Item{}, ErrInvalidArgument
		}
		rec.IsActive = req.IsActive.Value
	}

	if err := s.repo.Save(ctx, &rec); err != nil {
		return Item{}, err
	}

	return s.toAPI(rec), nil
}

func (s *Service) Delete(ctx context.Context, id uint64) error {
	if id == 0 {
		return ErrInvalidArgument
	}
	return s.repo.Delete(ctx, id)
}

func (s *Service) SetPhotoKey(ctx context.Context, id uint64, key *string) (Item, error) {
	if id == 0 {
		return Item{}, ErrInvalidArgument
	}

	rec, err := s.repo.UpdatePhotoKey(ctx, id, key)
	if err != nil {
		return Item{}, err
	}

	return s.toAPI(rec), nil
}

func (s *Service) ListActive(ctx context.Context, limit, offset int) (ListItemsResponse, error) {
	if limit <= 0 || limit > 20 {
		limit = 5
	}
	if offset < 0 {
		offset = 0
	}

	recs, err := s.repo.ListActive(ctx, limit, offset)
	if err != nil {
		return ListItemsResponse{}, err
	}

	items := make([]Item, 0, len(recs))
	for _, rec := range recs {
		items = append(items, s.toAPI(rec))
	}

	return ListItemsResponse{
		Items:  items,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (s *Service) GetActive(ctx context.Context, id uint64) (Item, error) {
	if id == 0 {
		return Item{}, ErrInvalidArgument
	}

	rec, err := s.repo.GetActive(ctx, id)
	if err != nil {
		return Item{}, err
	}

	return s.toAPI(rec), nil
}

func (s *Service) ensureModelExists(ctx context.Context, modelID uint64) error {
	if modelID == 0 {
		return ErrInvalidArgument
	}
	_, err := s.models.Get(ctx, modelID)
	if err != nil {
		if errors.Is(err, modelspkg.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func isEmptyPatch(req UpdateItemRequest) bool {
	return !req.Tag.Set &&
		!req.Title.Set &&
		!req.Description.Set &&
		!req.ModelID.Set &&
		!req.Sort.Set &&
		!req.IsActive.Set
}

func normalizeOptionalString(in *string) *string {
	if in == nil {
		return nil
	}
	v := strings.TrimSpace(*in)
	if v == "" {
		return nil
	}
	return &v
}

func (s *Service) toAPI(rec Showcase) Item {
	return Item{
		ID: rec.ID,
		Tag: ItemTag{
			Label:     rec.TagLabel,
			BgColor:   rec.TagBgColor,
			TextColor: rec.TagTextColor,
			Outlined:  rec.TagOutlined,
		},
		Title:       rec.Title,
		Description: rec.Description,
		ModelID:     rec.ModelID,
		PhotoKey:    rec.PhotoKey,
		PhotoURL:    nil,
		Sort:        rec.Sort,
		IsActive:    rec.IsActive,
		CreatedAt:   rec.CreatedAt,
		UpdatedAt:   rec.UpdatedAt,
	}
}
