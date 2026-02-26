package favorites_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"gorm.io/datatypes"

	"rest_waka/internal/favorites"
	"rest_waka/internal/models"
	"rest_waka/pkg/modelsutil"
)

func TestServiceList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		limit         int
		offset        int
		wantRepoLimit int
		wantRepoOff   int
		wantOutLimit  int
		wantOutOff    int
	}{
		{
			name:          "normalizes invalid limit and offset",
			limit:         -1,
			offset:        -10,
			wantRepoLimit: 50,
			wantRepoOff:   0,
			wantOutLimit:  50,
			wantOutOff:    0,
		},
		{
			name:          "uses valid limit and offset",
			limit:         20,
			offset:        3,
			wantRepoLimit: 20,
			wantRepoOff:   3,
			wantOutLimit:  20,
			wantOutOff:    3,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := &fakeFavoritesRepo{
				listFn: func(_ context.Context, userID uint64, limit, offset int) ([]models.WakaModel, error) {
					if userID != 77 {
						t.Fatalf("ListModelsFavorites() userID = %d, want 77", userID)
					}
					if limit != tc.wantRepoLimit || offset != tc.wantRepoOff {
						t.Fatalf("ListModelsFavorites() paging = (%d,%d), want (%d,%d)", limit, offset, tc.wantRepoLimit, tc.wantRepoOff)
					}
					return []models.WakaModel{
						newWakaRecord(t, 1, []string{"mint", "cola"}),
						newWakaRecord(t, 2, []string{"berry"}),
					}, nil
				},
			}
			svc := favorites.NewService(repo)

			got, err := svc.List(context.Background(), 77, tc.limit, tc.offset)
			if err != nil {
				t.Fatalf("List() error = %v", err)
			}
			if got.Limit != tc.wantOutLimit || got.Offset != tc.wantOutOff {
				t.Fatalf("List() paging = (%d,%d), want (%d,%d)", got.Limit, got.Offset, tc.wantOutLimit, tc.wantOutOff)
			}
			if len(got.Items) != 2 {
				t.Fatalf("List() items = %d, want 2", len(got.Items))
			}
			if !reflect.DeepEqual(got.Items[0].Flavors, []string{"mint", "cola"}) {
				t.Fatalf("item[0].flavors = %#v, want %#v", got.Items[0].Flavors, []string{"mint", "cola"})
			}
			if !reflect.DeepEqual(got.Items[1].Flavors, []string{"berry"}) {
				t.Fatalf("item[1].flavors = %#v, want %#v", got.Items[1].Flavors, []string{"berry"})
			}
		})
	}
}

func TestServiceAddAndRemoveRepositoryErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		call    func(svc *favorites.Service) error
		wantErr error
	}{
		{
			name: "add propagates already exists",
			call: func(svc *favorites.Service) error {
				return svc.Add(context.Background(), 1, 2)
			},
			wantErr: favorites.ErrAlreadyExists,
		},
		{
			name: "remove propagates not found",
			call: func(svc *favorites.Service) error {
				return svc.Remove(context.Background(), 1, 2)
			},
			wantErr: favorites.ErrNotFound,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := &fakeFavoritesRepo{
				addErr:    favorites.ErrAlreadyExists,
				removeErr: favorites.ErrNotFound,
			}
			svc := favorites.NewService(repo)

			err := tc.call(svc)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("service call error = %v, want %v", err, tc.wantErr)
			}
		})
	}
}

func TestServiceListRepositoryError(t *testing.T) {
	t.Parallel()

	repoErr := errors.New("db failed")
	repo := &fakeFavoritesRepo{
		listFn: func(_ context.Context, _ uint64, _ int, _ int) ([]models.WakaModel, error) {
			return nil, repoErr
		},
	}
	svc := favorites.NewService(repo)

	_, err := svc.List(context.Background(), 1, 10, 0)
	if !errors.Is(err, repoErr) {
		t.Fatalf("List() error = %v, want %v", err, repoErr)
	}
}

type fakeFavoritesRepo struct {
	addErr    error
	removeErr error
	listFn    func(ctx context.Context, userID uint64, limit, offset int) ([]models.WakaModel, error)
}

func (f *fakeFavoritesRepo) Add(_ context.Context, _ uint64, _ uint64) error {
	return f.addErr
}

func (f *fakeFavoritesRepo) Remove(_ context.Context, _ uint64, _ uint64) error {
	return f.removeErr
}

func (f *fakeFavoritesRepo) ListModelsFavorites(ctx context.Context, userID uint64, limit, offset int) ([]models.WakaModel, error) {
	if f.listFn == nil {
		return nil, nil
	}
	return f.listFn(ctx, userID, limit, offset)
}

func newWakaRecord(t *testing.T, id uint64, flavors []string) models.WakaModel {
	t.Helper()

	raw, err := modelsutil.MarshalFlavors(flavors)
	if err != nil {
		t.Fatalf("MarshalFlavors() error = %v", err)
	}

	now := time.Unix(1_700_000_000, 0).UTC()
	return models.WakaModel{
		ID:        id,
		Name:      "Model",
		Status:    models.StatusActive,
		PuffsMax:  600,
		Flavors:   datatypes.JSON(raw),
		CreatedAt: now,
		UpdatedAt: now,
	}
}
