package models

import (
	"context"
	"errors"
	"rest_waka/pkg/modelsutil"
	"strings"
)

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

	flvJson, err := modelsutil.MarshalFlavors(req.Flavors)
	if err != nil {
		return Model{}, err
	}

	rec := ModelRecord{
		Name:        name,
		Description: req.Description,
		PhotoURL:    req.PhotoURL,
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

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return Model{}, ErrInvalidArgument
		}
		rec.Name = name
	}

	// description (**string)
	if req.Description != nil {
		if *req.Description == nil {
			rec.Description = nil
		} else {
			v := **req.Description
			rec.Description = &v
		}
	}

	// photo_url (**string)
	if req.PhotoURL != nil {
		if *req.PhotoURL == nil {
			rec.PhotoURL = nil
		} else {
			v := **req.PhotoURL
			rec.PhotoURL = &v
		}
	}

	if req.PuffsMax != nil {
		if *req.PuffsMax <= 0 {
			return Model{}, ErrInvalidArgument
		}
		rec.PuffsMax = *req.PuffsMax
	}

	// flavors: пока patch заменяет список целиком
	if req.Flavors != nil {
		flvJSON, err := modelsutil.MarshalFlavors(*req.Flavors)
		if err != nil {
			return Model{}, err
		}
		rec.Flavors = flvJSON
	}

	// price (**int64)
	if req.PriceCents != nil {
		if *req.PriceCents == nil {
			rec.PriceCents = nil
		} else {
			v := **req.PriceCents
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

func isEmptyPatch(req UpdateModelRequest) bool {
	return req.Name == nil &&
		req.Description == nil &&
		req.PhotoURL == nil &&
		req.PuffsMax == nil &&
		req.Flavors == nil &&
		req.PriceCents == nil
}

func (s *Service) toAPI(rec ModelRecord) (Model, error) {
	flv, err := modelsutil.UnmarshalFlavors(rec.Flavors)
	if err != nil {
		return Model{}, err
	}

	return Model{
		ID:          rec.ID,
		Name:        rec.Name,
		Description: rec.Description,
		PhotoURL:    rec.PhotoURL,
		PuffsMax:    rec.PuffsMax,
		Flavors:     flv,
		PriceCents:  rec.PriceCents,
		CreatedAt:   rec.CreatedAt.Unix(),
		UpdatedAt:   rec.UpdatedAt.Unix(),
	}, nil
}
