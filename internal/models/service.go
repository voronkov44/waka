package models

import (
	"context"
	"errors"
	"rest_waka/pkg/modelsutil"
	"strings"
)

const (
	StatusActive  = "active"
	StatusHidden  = "hidden"
	StatusArchive = "archive"
)

func normalizeStatus(v string) (string, error) {
	s := strings.ToLower(strings.TrimSpace(v))
	if s == "" {
		return StatusHidden, nil
	}
	switch s {
	case StatusActive, StatusHidden, StatusArchive:
		return s, nil
	default:
		return "", ErrInvalidArgument
	}
}

type Service struct {
	repo RepositoryGorm
}

func NewService(repo RepositoryGorm) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req CreateModelRequest) (Model, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" || req.PuffsMax <= 0 {
		return Model{}, ErrInvalidArgument
	}
	if req.PriceCents != nil && *req.PriceCents < 0 {
		return Model{}, ErrInvalidArgument
	}
	if req.Flavors == nil {
		req.Flavors = []string{}
	}

	status, err := normalizeStatus(req.Status)
	if err != nil {
		return Model{}, err
	}

	flvJson, err := modelsutil.MarshalFlavors(req.Flavors)
	if err != nil {
		return Model{}, err
	}

	tag, err := normalizeTag(req.Tag)
	if err != nil {
		return Model{}, err
	}

	tagJSON, err := marshalTag(tag)
	if err != nil {
		return Model{}, err
	}

	rec := WakaModel{
		Name:        name,
		Status:      status,
		Description: req.Description,
		Tag:         tagJSON,
		PuffsMax:    req.PuffsMax,
		Flavors:     flvJson,
		PriceCents:  req.PriceCents,
	}

	if err := s.repo.Create(ctx, &rec); err != nil {
		return Model{}, err
	}
	return s.toAPI(rec)
}

func (s *Service) List(ctx context.Context, limit, offset int) (ListModelsResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	recs, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return ListModelsResponse{}, err
	}

	items := make([]Model, 0, len(recs))
	for _, r := range recs {
		m, err := s.toAPI(r)
		if err != nil {
			return ListModelsResponse{}, err
		}
		items = append(items, m)
	}

	return ListModelsResponse{
		Items:  items,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (s *Service) Get(ctx context.Context, id uint64) (Model, error) {
	if id == 0 {
		return Model{}, ErrInvalidArgument
	}
	rec, err := s.repo.Get(ctx, id)
	if err != nil {
		return Model{}, err
	}
	return s.toAPI(rec)
}

func (s *Service) Update(ctx context.Context, id uint64, req UpdateModelRequest) (Model, error) {
	if id == 0 {
		return Model{}, ErrInvalidArgument
	}
	if isEmptyPatch(req) {
		return Model{}, ErrInvalidArgument
	}

	rec, err := s.repo.Get(ctx, id)
	if err != nil {
		return Model{}, err
	}

	if req.Name.Set {
		if req.Name.Null {
			return Model{}, ErrInvalidArgument
		}
		name := strings.TrimSpace(req.Name.Value)
		if name == "" {
			return Model{}, ErrInvalidArgument
		}
		rec.Name = name
	}

	if req.Status.Set {
		if req.Status.Null {
			return Model{}, ErrInvalidArgument
		}
		st, err := normalizeStatus(req.Status.Value)
		if err != nil {
			return Model{}, err
		}
		rec.Status = st
	}

	if req.Description.Set {
		if req.Description.Null {
			rec.Description = nil
		} else {
			v := req.Description.Value
			rec.Description = &v
		}
	}

	if req.Tag.Set {
		if req.Tag.Null {
			rec.Tag = nil
		} else {
			tag, err := normalizeTag(&req.Tag.Value)
			if err != nil {
				return Model{}, err
			}
			tagJSON, err := marshalTag(tag)
			if err != nil {
				return Model{}, err
			}
			rec.Tag = tagJSON
		}
	}

	if req.PuffsMax.Set {
		if req.PuffsMax.Null || req.PuffsMax.Value <= 0 {
			return Model{}, ErrInvalidArgument
		}
		rec.PuffsMax = req.PuffsMax.Value
	}

	if req.Flavors.Set {
		nextFlavors := req.Flavors.Value
		if req.Flavors.Null {
			nextFlavors = []string{}
		}
		flvJSON, err := modelsutil.MarshalFlavors(nextFlavors)
		if err != nil {
			return Model{}, err
		}
		rec.Flavors = flvJSON
	}

	if req.PriceCents.Set {
		if req.PriceCents.Null {
			rec.PriceCents = nil
		} else {
			v := req.PriceCents.Value
			if v < 0 {
				return Model{}, ErrInvalidArgument
			}
			rec.PriceCents = &v
		}
	}

	if err := s.repo.Save(ctx, &rec); err != nil {
		return Model{}, err
	}
	return s.toAPI(rec)
}

func (s *Service) Delete(ctx context.Context, id uint64) error {
	if id == 0 {
		return ErrInvalidArgument
	}
	return s.repo.Delete(ctx, id)
}

// Ручки для работы со слайсом вкуса

func (s *Service) AddFlavor(ctx context.Context, id uint64, value string) (Model, error) {
	if id == 0 {
		return Model{}, ErrInvalidArgument
	}

	rec, err := s.repo.Get(ctx, id)
	if err != nil {
		return Model{}, err
	}

	flv, err := modelsutil.UnmarshalFlavors(rec.Flavors)
	if err != nil {
		return Model{}, err
	}

	updated, changed, err := modelsutil.AddFlavorUnique(flv, value)
	if err != nil {
		if errors.Is(err, modelsutil.ErrEmptyFlavor) {
			return Model{}, ErrInvalidArgument
		}
		return Model{}, err
	}
	if changed {
		raw, err := modelsutil.MarshalFlavors(updated)
		if err != nil {
			return Model{}, err
		}
		rec.Flavors = raw
		if err := s.repo.Save(ctx, &rec); err != nil {
			return Model{}, err
		}
	}
	return s.toAPI(rec)
}

func (s *Service) RemoveFlavor(ctx context.Context, id uint64, value string) (Model, error) {
	if id == 0 {
		return Model{}, ErrInvalidArgument
	}

	rec, err := s.repo.Get(ctx, id)
	if err != nil {
		return Model{}, err
	}

	flv, err := modelsutil.UnmarshalFlavors(rec.Flavors)
	if err != nil {
		return Model{}, err
	}

	updated, changed, err := modelsutil.RemoveFlavor(flv, value)
	if err != nil {
		if errors.Is(err, modelsutil.ErrEmptyFlavor) {
			return Model{}, ErrInvalidArgument
		}
		return Model{}, err
	}
	if changed {
		raw, err := modelsutil.MarshalFlavors(updated)
		if err != nil {
			return Model{}, err
		}
		rec.Flavors = raw
		if err := s.repo.Save(ctx, &rec); err != nil {
			return Model{}, err
		}
	}
	return s.toAPI(rec)
}

func (s *Service) ListActive(ctx context.Context, limit, offset int) (ListModelsResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	recs, err := s.repo.ListByStatus(ctx, StatusActive, limit, offset)
	if err != nil {
		return ListModelsResponse{}, err
	}

	items := make([]Model, 0, len(recs))
	for _, r := range recs {
		m, err := s.toAPI(r)
		if err != nil {
			return ListModelsResponse{}, err
		}
		items = append(items, m)
	}

	return ListModelsResponse{
		Items:  items,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (s *Service) GetActive(ctx context.Context, id uint64) (Model, error) {
	if id == 0 {
		return Model{}, ErrInvalidArgument
	}

	rec, err := s.repo.GetByIDAndStatus(ctx, id, StatusActive)
	if err != nil {
		return Model{}, err
	}

	return s.toAPI(rec)
}

func isEmptyPatch(req UpdateModelRequest) bool {
	return !req.Name.Set &&
		!req.Status.Set &&
		!req.Description.Set &&
		!req.Tag.Set &&
		!req.PuffsMax.Set &&
		!req.Flavors.Set &&
		!req.PriceCents.Set
}

func (s *Service) toAPI(rec WakaModel) (Model, error) {
	flv, err := modelsutil.UnmarshalFlavors(rec.Flavors)
	if err != nil {
		return Model{}, err
	}

	tag, err := unmarshalTag(rec.Tag)
	if err != nil {
		return Model{}, err
	}

	return Model{
		ID:          rec.ID,
		Name:        rec.Name,
		Status:      rec.Status,
		Description: rec.Description,
		PhotoKey:    rec.PhotoKey,
		PhotoURL:    nil, // computed handler
		Tag:         tag,
		PuffsMax:    rec.PuffsMax,
		Flavors:     flv,
		PriceCents:  rec.PriceCents,
		CreatedAt:   rec.CreatedAt,
		UpdatedAt:   rec.UpdatedAt,
	}, nil
}

func (s *Service) SetPhotoKey(ctx context.Context, id uint64, key *string) (Model, error) {
	if id == 0 {
		return Model{}, ErrInvalidArgument
	}
	rec, err := s.repo.UpdatePhotoKey(ctx, id, key)
	if err != nil {
		return Model{}, err
	}
	return s.toAPI(rec)
}
