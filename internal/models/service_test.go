package models_test

import (
	"context"
	"errors"
	"reflect"
	"rest_waka/pkg/patch"
	"sort"
	"testing"
	"time"

	"gorm.io/datatypes"

	"rest_waka/internal/models"
	"rest_waka/pkg/modelsutil"
)

func TestServiceCreate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		req         models.CreateModelRequest
		wantErr     error
		wantName    string
		wantFlavors []string
	}{
		{
			name: "success trims name and defaults nil flavors",
			req: models.CreateModelRequest{
				Name:     "  Waka X ",
				PuffsMax: 600,
			},
			wantName:    "Waka X",
			wantFlavors: []string{},
		},
		{
			name: "rejects empty name",
			req: models.CreateModelRequest{
				Name:     "   ",
				PuffsMax: 600,
			},
			wantErr: models.ErrInvalidArgument,
		},
		{
			name: "rejects non-positive puffs",
			req: models.CreateModelRequest{
				Name:     "Waka",
				PuffsMax: 0,
			},
			wantErr: models.ErrInvalidArgument,
		},
		{
			name: "rejects negative price",
			req: models.CreateModelRequest{
				Name:       "Waka",
				PuffsMax:   600,
				PriceCents: ptr(int64(-1)),
			},
			wantErr: models.ErrInvalidArgument,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := newFakeRepo()
			svc := models.NewService(repo)

			got, err := svc.Create(context.Background(), tc.req)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Create() error = %v, want %v", err, tc.wantErr)
			}
			if tc.wantErr != nil {
				return
			}
			if got.Name != tc.wantName {
				t.Fatalf("Create() name = %q, want %q", got.Name, tc.wantName)
			}
			if !reflect.DeepEqual(got.Flavors, tc.wantFlavors) {
				t.Fatalf("Create() flavors = %#v, want %#v", got.Flavors, tc.wantFlavors)
			}
		})
	}
}

func TestServiceGet(t *testing.T) {
	t.Parallel()

	repo := newFakeRepo()
	rec := repo.seedModel("Waka", 600, []string{"mint"})
	svc := models.NewService(repo)

	tests := []struct {
		name    string
		id      uint64
		wantErr error
	}{
		{
			name:    "invalid id",
			id:      0,
			wantErr: models.ErrInvalidArgument,
		},
		{
			name:    "not found",
			id:      9999,
			wantErr: models.ErrNotFound,
		},
		{
			name: "found",
			id:   rec.ID,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := svc.Get(context.Background(), tc.id)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Get() error = %v, want %v", err, tc.wantErr)
			}
			if tc.wantErr != nil {
				return
			}
			if got.ID != rec.ID || got.Name != rec.Name {
				t.Fatalf("Get() got = %#v, want id=%d name=%q", got, rec.ID, rec.Name)
			}
		})
	}
}

func TestServiceList(t *testing.T) {
	t.Parallel()

	repo := newFakeRepo()
	first := repo.seedModel("First", 500, []string{"a"})
	second := repo.seedModel("Second", 600, []string{"b"})
	third := repo.seedModel("Third", 700, []string{"c"})

	svc := models.NewService(repo)

	tests := []struct {
		name       string
		limit      int
		offset     int
		wantLimit  int
		wantOffset int
		wantIDs    []uint64
	}{
		{
			name:       "defaults invalid limit and offset",
			limit:      -1,
			offset:     -5,
			wantLimit:  50,
			wantOffset: 0,
			wantIDs:    []uint64{third.ID, second.ID, first.ID},
		},
		{
			name:       "applies limit and offset",
			limit:      1,
			offset:     1,
			wantLimit:  1,
			wantOffset: 1,
			wantIDs:    []uint64{second.ID},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := svc.List(context.Background(), tc.limit, tc.offset)
			if err != nil {
				t.Fatalf("List() error = %v", err)
			}
			if got.Limit != tc.wantLimit || got.Offset != tc.wantOffset {
				t.Fatalf("List() paging = (%d,%d), want (%d,%d)", got.Limit, got.Offset, tc.wantLimit, tc.wantOffset)
			}

			ids := make([]uint64, 0, len(got.Items))
			for _, item := range got.Items {
				ids = append(ids, item.ID)
			}
			if !reflect.DeepEqual(ids, tc.wantIDs) {
				t.Fatalf("List() ids = %#v, want %#v", ids, tc.wantIDs)
			}
		})
	}
}

func TestServiceUpdate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		id        uint64
		req       models.UpdateModelRequest
		wantErr   error
		assertion func(t *testing.T, m models.Model, repo *fakeRepo)
	}{
		{
			name:    "invalid id",
			id:      0,
			req:     models.UpdateModelRequest{Name: patchVal("x")},
			wantErr: models.ErrInvalidArgument,
		},
		{
			name:    "empty patch rejected",
			id:      1,
			req:     models.UpdateModelRequest{},
			wantErr: models.ErrInvalidArgument,
		},
		{
			name:    "not found",
			id:      9999,
			req:     models.UpdateModelRequest{Name: patchVal("x")},
			wantErr: models.ErrNotFound,
		},
		{
			name: "name gets trimmed",
			id:   1,
			req:  models.UpdateModelRequest{Name: patchVal("  New Name  ")},
			assertion: func(t *testing.T, m models.Model, _ *fakeRepo) {
				t.Helper()
				if m.Name != "New Name" {
					t.Fatalf("Update() name = %q, want %q", m.Name, "New Name")
				}
			},
		},
		{
			name:    "invalid empty name",
			id:      1,
			req:     models.UpdateModelRequest{Name: patchVal("   ")},
			wantErr: models.ErrInvalidArgument,
		},
		{
			name: "description absent keeps old value",
			id:   1,
			req:  models.UpdateModelRequest{Name: patchVal("Name2")},
			assertion: func(t *testing.T, m models.Model, _ *fakeRepo) {
				t.Helper()
				if m.Description == nil || *m.Description != "old description" {
					t.Fatalf("Update() description = %#v, want old description", m.Description)
				}
			},
		},
		{
			name: "description null clears value",
			id:   1,
			req:  models.UpdateModelRequest{Description: patchNull[string]()},
			assertion: func(t *testing.T, m models.Model, _ *fakeRepo) {
				t.Helper()
				if m.Description != nil {
					t.Fatalf("Update() description = %#v, want nil", m.Description)
				}
			},
		},
		{
			name: "description value sets value",
			id:   1,
			req:  models.UpdateModelRequest{Description: patchVal("new description")},
			assertion: func(t *testing.T, m models.Model, _ *fakeRepo) {
				t.Helper()
				if m.Description == nil || *m.Description != "new description" {
					t.Fatalf("Update() description = %#v, want new description", m.Description)
				}
			},
		},
		{
			name: "price null clears value",
			id:   1,
			req:  models.UpdateModelRequest{PriceCents: patchNull[int64]()},
			assertion: func(t *testing.T, m models.Model, _ *fakeRepo) {
				t.Helper()
				if m.PriceCents != nil {
					t.Fatalf("Update() price = %#v, want nil", m.PriceCents)
				}
			},
		},
		{
			name: "price value sets value",
			id:   1,
			req:  models.UpdateModelRequest{PriceCents: patchVal(int64(1999))},
			assertion: func(t *testing.T, m models.Model, _ *fakeRepo) {
				t.Helper()
				if m.PriceCents == nil || *m.PriceCents != 1999 {
					t.Fatalf("Update() price = %#v, want 1999", m.PriceCents)
				}
			},
		},
		{
			name:    "negative price rejected",
			id:      1,
			req:     models.UpdateModelRequest{PriceCents: patchVal(int64(-1))},
			wantErr: models.ErrInvalidArgument,
		},
		{
			name: "flavors replaces entire list",
			id:   1,
			req:  models.UpdateModelRequest{Flavors: patchVal([]string{"cola", "berry"})},
			assertion: func(t *testing.T, m models.Model, _ *fakeRepo) {
				t.Helper()
				want := []string{"cola", "berry"}
				if !reflect.DeepEqual(m.Flavors, want) {
					t.Fatalf("Update() flavors = %#v, want %#v", m.Flavors, want)
				}
			},
		},
		{
			name: "flavors null clears list",
			id:   1,
			req:  models.UpdateModelRequest{Flavors: patchNull[[]string]()},
			assertion: func(t *testing.T, m models.Model, _ *fakeRepo) {
				t.Helper()
				if len(m.Flavors) != 0 {
					t.Fatalf("Update() flavors = %#v, want empty slice", m.Flavors)
				}
			},
		},
		{
			name:    "puffs max must be >0",
			id:      1,
			req:     models.UpdateModelRequest{PuffsMax: patchVal(0)},
			wantErr: models.ErrInvalidArgument,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := newFakeRepo()
			repo.seedModelWithID(1, "Old Name", 600, []string{"mint"})
			old := repo.items[1]
			old.Description = ptr("old description")
			old.PhotoURL = ptr("https://old.local/photo.jpg")
			old.PriceCents = ptr(int64(1499))
			repo.items[1] = old

			svc := models.NewService(repo)
			got, err := svc.Update(context.Background(), tc.id, tc.req)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Update() error = %v, want %v", err, tc.wantErr)
			}
			if tc.wantErr != nil {
				return
			}
			if tc.assertion != nil {
				tc.assertion(t, got, repo)
			}
		})
	}
}

func TestServiceDelete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		id      uint64
		wantErr error
	}{
		{
			name:    "invalid id",
			id:      0,
			wantErr: models.ErrInvalidArgument,
		},
		{
			name:    "not found",
			id:      999,
			wantErr: models.ErrNotFound,
		},
		{
			name: "success",
			id:   1,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := newFakeRepo()
			repo.seedModelWithID(1, "Waka", 600, []string{"mint"})
			svc := models.NewService(repo)

			err := svc.Delete(context.Background(), tc.id)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Delete() error = %v, want %v", err, tc.wantErr)
			}
			if tc.wantErr == nil {
				if _, ok := repo.items[1]; ok {
					t.Fatal("Delete() did not remove record")
				}
			}
		})
	}
}

func TestServiceAddFlavor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		id            uint64
		value         string
		seed          bool
		seedFlavors   []string
		wantErr       error
		wantFlavors   []string
		wantSaveCalls int
	}{
		{
			name:    "invalid id",
			id:      0,
			value:   "mint",
			wantErr: models.ErrInvalidArgument,
		},
		{
			name:    "not found",
			id:      2,
			value:   "mint",
			wantErr: models.ErrNotFound,
		},
		{
			name:          "adds new flavor",
			id:            1,
			value:         "  cola ",
			seed:          true,
			seedFlavors:   []string{"mint"},
			wantFlavors:   []string{"mint", "cola"},
			wantSaveCalls: 1,
		},
		{
			name:          "duplicate is idempotent case-insensitive",
			id:            1,
			value:         "MINT",
			seed:          true,
			seedFlavors:   []string{"mint"},
			wantFlavors:   []string{"mint"},
			wantSaveCalls: 0,
		},
		{
			name:          "duplicate with spaces in stored flavor is idempotent",
			id:            1,
			value:         "mint",
			seed:          true,
			seedFlavors:   []string{"  Mint  "},
			wantFlavors:   []string{"  Mint  "},
			wantSaveCalls: 0,
		},
		{
			name:    "empty flavor rejected",
			id:      1,
			value:   " ",
			seed:    true,
			wantErr: models.ErrInvalidArgument,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := newFakeRepo()
			if tc.seed {
				repo.seedModelWithID(1, "Waka", 600, tc.seedFlavors)
			}

			svc := models.NewService(repo)
			got, err := svc.AddFlavor(context.Background(), tc.id, tc.value)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("AddFlavor() error = %v, want %v", err, tc.wantErr)
			}
			if tc.wantErr != nil {
				return
			}
			if !reflect.DeepEqual(got.Flavors, tc.wantFlavors) {
				t.Fatalf("AddFlavor() flavors = %#v, want %#v", got.Flavors, tc.wantFlavors)
			}
			if repo.saveCalls != tc.wantSaveCalls {
				t.Fatalf("AddFlavor() saveCalls = %d, want %d", repo.saveCalls, tc.wantSaveCalls)
			}
		})
	}
}

func TestServiceRemoveFlavor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		id            uint64
		value         string
		seed          bool
		seedFlavors   []string
		wantErr       error
		wantFlavors   []string
		wantSaveCalls int
	}{
		{
			name:    "invalid id",
			id:      0,
			value:   "mint",
			wantErr: models.ErrInvalidArgument,
		},
		{
			name:    "not found",
			id:      2,
			value:   "mint",
			wantErr: models.ErrNotFound,
		},
		{
			name:          "removes existing flavor",
			id:            1,
			value:         "mint",
			seed:          true,
			seedFlavors:   []string{"Mint", "cola"},
			wantFlavors:   []string{"cola"},
			wantSaveCalls: 1,
		},
		{
			name:          "remove missing flavor is idempotent",
			id:            1,
			value:         "berry",
			seed:          true,
			seedFlavors:   []string{"mint"},
			wantFlavors:   []string{"mint"},
			wantSaveCalls: 0,
		},
		{
			name:          "removes flavor even if stored with spaces",
			id:            1,
			value:         "mint",
			seed:          true,
			seedFlavors:   []string{"  mint  ", "cola"},
			wantFlavors:   []string{"cola"},
			wantSaveCalls: 1,
		},
		{
			name:    "empty flavor rejected",
			id:      1,
			value:   " ",
			seed:    true,
			wantErr: models.ErrInvalidArgument,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := newFakeRepo()
			if tc.seed {
				repo.seedModelWithID(1, "Waka", 600, tc.seedFlavors)
			}

			svc := models.NewService(repo)
			got, err := svc.RemoveFlavor(context.Background(), tc.id, tc.value)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("RemoveFlavor() error = %v, want %v", err, tc.wantErr)
			}
			if tc.wantErr != nil {
				return
			}
			if !reflect.DeepEqual(got.Flavors, tc.wantFlavors) {
				t.Fatalf("RemoveFlavor() flavors = %#v, want %#v", got.Flavors, tc.wantFlavors)
			}
			if repo.saveCalls != tc.wantSaveCalls {
				t.Fatalf("RemoveFlavor() saveCalls = %d, want %d", repo.saveCalls, tc.wantSaveCalls)
			}
		})
	}
}

type fakeRepo struct {
	nextID    uint64
	now       time.Time
	items     map[uint64]models.WakaModel
	saveCalls int
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{
		nextID: 1,
		now:    time.Unix(1_700_000_000, 0).UTC(),
		items:  map[uint64]models.WakaModel{},
	}
}

func (f *fakeRepo) Create(_ context.Context, rec *models.WakaModel) error {
	if rec.ID == 0 {
		rec.ID = f.nextID
		f.nextID++
	}
	if rec.CreatedAt.IsZero() {
		rec.CreatedAt = f.now
	}
	if rec.UpdatedAt.IsZero() {
		rec.UpdatedAt = f.now
	}
	f.items[rec.ID] = cloneModelRecord(*rec)
	return nil
}

func (f *fakeRepo) Get(_ context.Context, id uint64) (models.WakaModel, error) {
	rec, ok := f.items[id]
	if !ok {
		return models.WakaModel{}, models.ErrNotFound
	}
	return cloneModelRecord(rec), nil
}

func (f *fakeRepo) List(_ context.Context, limit, offset int) ([]models.WakaModel, error) {
	ids := make([]uint64, 0, len(f.items))
	for id := range f.items {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] > ids[j] })

	start := offset
	if start > len(ids) {
		start = len(ids)
	}
	end := start + limit
	if end > len(ids) {
		end = len(ids)
	}

	out := make([]models.WakaModel, 0, end-start)
	for _, id := range ids[start:end] {
		out = append(out, cloneModelRecord(f.items[id]))
	}
	return out, nil
}

func (f *fakeRepo) Save(_ context.Context, rec *models.WakaModel) error {
	if _, ok := f.items[rec.ID]; !ok {
		return models.ErrNotFound
	}
	f.saveCalls++
	rec.UpdatedAt = f.now.Add(time.Second * time.Duration(f.saveCalls))
	f.items[rec.ID] = cloneModelRecord(*rec)
	return nil
}

func (f *fakeRepo) Delete(_ context.Context, id uint64) error {
	if _, ok := f.items[id]; !ok {
		return models.ErrNotFound
	}
	delete(f.items, id)
	return nil
}

func (f *fakeRepo) seedModel(name string, puffs int, flavors []string) models.WakaModel {
	id := f.nextID
	f.nextID++
	return f.seedModelWithID(id, name, puffs, flavors)
}

func (f *fakeRepo) seedModelWithID(id uint64, name string, puffs int, flavors []string) models.WakaModel {
	raw, err := modelsutil.MarshalFlavors(flavors)
	if err != nil {
		panic(err)
	}
	rec := models.WakaModel{
		ID:        id,
		Name:      name,
		PuffsMax:  puffs,
		Flavors:   raw,
		CreatedAt: f.now,
		UpdatedAt: f.now,
	}
	f.items[id] = cloneModelRecord(rec)
	if id >= f.nextID {
		f.nextID = id + 1
	}
	return cloneModelRecord(rec)
}

func cloneModelRecord(rec models.WakaModel) models.WakaModel {
	out := rec
	out.Flavors = append(datatypes.JSON(nil), rec.Flavors...)
	return out
}

func ptr[T any](v T) *T {
	return &v
}

func patchVal[T any](v T) patch.Field[T] {
	return patch.Field[T]{
		Set:   true,
		Value: v,
	}
}

func patchNull[T any]() patch.Field[T] {
	return patch.Field[T]{
		Set:  true,
		Null: true,
	}
}
