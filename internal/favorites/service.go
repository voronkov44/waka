package favorites

import (
	"context"
	"rest_waka/internal/models"
	"rest_waka/pkg/modelsutil"
)

type Service struct {
	repo RepositoryGorm
}

func NewService(repo RepositoryGorm) *Service {
	return &Service{repo: repo}
}

func (s *Service) Add(ctx context.Context, userID, modelID uint64) error {
	return s.repo.Add(ctx, userID, modelID)
}

func (s *Service) Remove(ctx context.Context, userID, modelID uint64) error {
	return s.repo.Remove(ctx, userID, modelID)
}

func (s *Service) List(ctx context.Context, userID uint64, limit, offset int) (ListFavoritesResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	recs, err := s.repo.ListModelsFavorites(ctx, userID, limit, offset)
	if err != nil {
		return ListFavoritesResponse{}, err
	}

	items := make([]models.Model, 0, len(recs))
	for _, rec := range recs {
		m, err := toAPIModel(rec)
		if err != nil {
			return ListFavoritesResponse{}, err
		}
		items = append(items, m)
	}

	return ListFavoritesResponse{
		Items:  items,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func toAPIModel(rec models.WakaModel) (models.Model, error) {
	flv, err := modelsutil.UnmarshalFlavors(rec.Flavors)
	if err != nil {
		return models.Model{}, err
	}

	return models.Model{
		ID:          rec.ID,
		Name:        rec.Name,
		Status:      rec.Status,
		Description: rec.Description,
		PhotoURL:    rec.PhotoURL,
		PuffsMax:    rec.PuffsMax,
		Flavors:     flv,
		PriceCents:  rec.PriceCents,
		CreatedAt:   rec.CreatedAt,
		UpdatedAt:   rec.UpdatedAt,
	}, nil
}
