package favorites

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"rest_waka/internal/models"
	"rest_waka/pkg/middleware"
)

func TestFavoritesHandlerListRequiresUser(t *testing.T) {
	t.Parallel()

	handler := &Handler{svc: &fakeFavoritesService{}}
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/favorites", nil)

	handler.List().ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusUnauthorized)
	}
	if !strings.Contains(rr.Body.String(), "unauthorized") {
		t.Fatalf("body = %q, want unauthorized", rr.Body.String())
	}
}

func TestFavoritesHandlerListComputesPhotoURL(t *testing.T) {
	t.Parallel()

	svc := &fakeFavoritesService{
		listFn: func(_ context.Context, userID uint64, limit, offset int) (ListFavoritesResponse, error) {
			if userID != 77 {
				t.Fatalf("List() userID = %d, want 77", userID)
			}
			if limit != 20 || offset != 4 {
				t.Fatalf("List() paging = (%d,%d), want (20,4)", limit, offset)
			}
			return ListFavoritesResponse{
				Items: []models.Model{
					{ID: 1, Name: "One", Status: models.StatusActive, PuffsMax: 500, Flavors: []string{"mint"}, PhotoKey: strPtr("models/1/a.jpg")},
					{ID: 2, Name: "Two", Status: models.StatusHidden, PuffsMax: 900, Flavors: []string{"cola"}},
				},
				Limit:  limit,
				Offset: offset,
			}, nil
		},
	}
	s3 := &fakePhotoResolver{publicBase: "https://cdn.example/favs"}

	handler := &Handler{
		svc:          svc,
		s3:           s3,
		usePresigned: false,
		presignTTL:   time.Minute,
	}

	req := httptest.NewRequest(http.MethodGet, "/api/favorites?limit=20&offset=4", nil)
	req = req.WithContext(middleware.ContextWithUserID(req.Context(), 77))
	rr := httptest.NewRecorder()

	handler.List().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}

	var got ListFavoritesResponse
	decodeJSON(t, rr, &got)
	if len(got.Items) != 2 {
		t.Fatalf("items = %d, want 2", len(got.Items))
	}
	if got.Items[0].PhotoURL == nil || *got.Items[0].PhotoURL != "https://cdn.example/favs/models/1/a.jpg" {
		t.Fatalf("item[0].photo_url = %#v, want computed url", got.Items[0].PhotoURL)
	}
	if got.Items[1].PhotoURL != nil {
		t.Fatalf("item[1].photo_url = %#v, want nil", got.Items[1].PhotoURL)
	}
	if len(s3.publicCalls) != 1 || s3.publicCalls[0] != "models/1/a.jpg" {
		t.Fatalf("PublicURL calls = %#v, want one call with key", s3.publicCalls)
	}
}

type fakeFavoritesService struct {
	addFn    func(ctx context.Context, userID, modelID uint64) error
	removeFn func(ctx context.Context, userID, modelID uint64) error
	listFn   func(ctx context.Context, userID uint64, limit, offset int) (ListFavoritesResponse, error)
}

func (f *fakeFavoritesService) Add(ctx context.Context, userID, modelID uint64) error {
	if f.addFn == nil {
		return nil
	}
	return f.addFn(ctx, userID, modelID)
}

func (f *fakeFavoritesService) Remove(ctx context.Context, userID, modelID uint64) error {
	if f.removeFn == nil {
		return nil
	}
	return f.removeFn(ctx, userID, modelID)
}

func (f *fakeFavoritesService) List(ctx context.Context, userID uint64, limit, offset int) (ListFavoritesResponse, error) {
	if f.listFn == nil {
		return ListFavoritesResponse{}, nil
	}
	return f.listFn(ctx, userID, limit, offset)
}

type fakePhotoResolver struct {
	publicBase  string
	publicCalls []string
}

func (f *fakePhotoResolver) PublicURL(key string) string {
	f.publicCalls = append(f.publicCalls, key)
	return strings.TrimRight(f.publicBase, "/") + "/" + strings.TrimLeft(key, "/")
}

func (f *fakePhotoResolver) PresignedGetURL(_ context.Context, key string, _ time.Duration) (string, error) {
	return "https://signed.example/" + key, nil
}

func decodeJSON(t *testing.T, rr *httptest.ResponseRecorder, dst any) {
	t.Helper()
	if err := json.Unmarshal(rr.Body.Bytes(), dst); err != nil {
		t.Fatalf("json.Unmarshal() error = %v, body=%q", err, rr.Body.String())
	}
}

func strPtr(v string) *string {
	return &v
}
