package models_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"rest_waka/internal/models"
)

func TestModelsHandlerCreateComputesPhotoURL(t *testing.T) {
	t.Parallel()

	svc := &fakeModelsService{
		createFn: func(_ context.Context, req models.CreateModelRequest) (models.Model, error) {
			if req.Name != "Waka" {
				t.Fatalf("Create() request name = %q, want Waka", req.Name)
			}
			return models.Model{
				ID:       1,
				Name:     "Waka",
				Status:   models.StatusHidden,
				PuffsMax: 600,
				Flavors:  []string{"mint"},
				PhotoKey: strPtr("models/1/photo.jpg"),
			}, nil
		},
	}
	s3 := &fakeS3{publicBase: "https://cdn.example/bucket"}
	mux := newModelsMux(svc, s3, false)

	rr := doJSONRequest(mux, http.MethodPost, "/api/models", `{"name":"Waka","puffs_max":600,"flavors":["mint"]}`)
	if rr.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusCreated)
	}

	var got models.Model
	decodeJSON(t, rr, &got)
	if got.PhotoURL == nil || *got.PhotoURL != "https://cdn.example/bucket/models/1/photo.jpg" {
		t.Fatalf("photo_url = %#v, want computed public url", got.PhotoURL)
	}
	if len(s3.publicCalls) != 1 || s3.publicCalls[0] != "models/1/photo.jpg" {
		t.Fatalf("public url calls = %#v, want one call with model key", s3.publicCalls)
	}
}

func TestModelsHandlerListComputesPhotoURL(t *testing.T) {
	t.Parallel()

	svc := &fakeModelsService{
		listFn: func(_ context.Context, limit, offset int) (models.ListModelsResponse, error) {
			if limit != 20 || offset != 5 {
				t.Fatalf("List() paging = (%d,%d), want (20,5)", limit, offset)
			}
			return models.ListModelsResponse{
				Items: []models.Model{
					{ID: 1, Name: "One", Status: models.StatusActive, PuffsMax: 500, Flavors: []string{"mint"}, PhotoKey: strPtr("models/1/a.jpg")},
					{ID: 2, Name: "Two", Status: models.StatusHidden, PuffsMax: 800, Flavors: []string{"cola"}, PhotoKey: nil},
					{ID: 3, Name: "Three", Status: models.StatusArchive, PuffsMax: 1200, Flavors: []string{"berry"}, PhotoKey: strPtr("models/3/c.jpg")},
				},
				Limit:  limit,
				Offset: offset,
			}, nil
		},
	}
	s3 := &fakeS3{publicBase: "https://cdn.example/media"}
	mux := newModelsMux(svc, s3, false)

	rr := doRequest(mux, http.MethodGet, "/api/models?limit=20&offset=5", nil, "")
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}

	var got models.ListModelsResponse
	decodeJSON(t, rr, &got)
	if len(got.Items) != 3 {
		t.Fatalf("items = %d, want 3", len(got.Items))
	}
	if got.Items[0].PhotoURL == nil || *got.Items[0].PhotoURL != "https://cdn.example/media/models/1/a.jpg" {
		t.Fatalf("item[0].photo_url = %#v, want computed url", got.Items[0].PhotoURL)
	}
	if got.Items[1].PhotoURL != nil {
		t.Fatalf("item[1].photo_url = %#v, want nil", got.Items[1].PhotoURL)
	}
	if got.Items[2].PhotoURL == nil || *got.Items[2].PhotoURL != "https://cdn.example/media/models/3/c.jpg" {
		t.Fatalf("item[2].photo_url = %#v, want computed url", got.Items[2].PhotoURL)
	}
	if len(s3.publicCalls) != 2 {
		t.Fatalf("PublicURL calls = %d, want 2", len(s3.publicCalls))
	}
}

func TestModelsHandlerUploadPhotoValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		requestBuilder   func() *http.Request
		wantStatus       int
		wantBodyContains string
	}{
		{
			name: "rejects non multipart",
			requestBuilder: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/api/models/7/photo", strings.NewReader(`{"x":1}`))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus:       http.StatusBadRequest,
			wantBodyContains: "invalid multipart form",
		},
		{
			name: "rejects missing file",
			requestBuilder: func() *http.Request {
				return newMultipartRequest(t, http.MethodPost, "/api/models/7/photo", false, "", nil)
			},
			wantStatus:       http.StatusBadRequest,
			wantBodyContains: "file is required",
		},
		{
			name: "rejects non image",
			requestBuilder: func() *http.Request {
				return newMultipartRequest(t, http.MethodPost, "/api/models/7/photo", true, "note.txt", []byte("plain text"))
			},
			wantStatus:       http.StatusBadRequest,
			wantBodyContains: "only images are allowed",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := &fakeModelsService{
				getFn: func(_ context.Context, id uint64) (models.Model, error) {
					if id != 7 {
						t.Fatalf("Get() id = %d, want 7", id)
					}
					return models.Model{ID: 7, Name: "W", Status: models.StatusHidden, PuffsMax: 600, Flavors: []string{}}, nil
				},
				setPhotoKeyFn: func(_ context.Context, _ uint64, _ *string) (models.Model, error) {
					t.Fatal("SetPhotoKey() should not be called")
					return models.Model{}, nil
				},
			}
			s3 := &fakeS3{publicBase: "https://cdn.example/b"}
			mux := newModelsMux(svc, s3, false)

			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, tc.requestBuilder())

			if rr.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d", rr.Code, tc.wantStatus)
			}
			if !strings.Contains(rr.Body.String(), tc.wantBodyContains) {
				t.Fatalf("body = %q, want contains %q", rr.Body.String(), tc.wantBodyContains)
			}
			if len(s3.putCalls) != 0 {
				t.Fatalf("Put() calls = %d, want 0", len(s3.putCalls))
			}
		})
	}
}

func TestModelsHandlerUploadPhotoSuccess(t *testing.T) {
	t.Parallel()

	var setPhotoKeyCalls int
	var newKey string

	svc := &fakeModelsService{
		getFn: func(_ context.Context, id uint64) (models.Model, error) {
			if id != 7 {
				t.Fatalf("Get() id = %d, want 7", id)
			}
			return models.Model{
				ID:       7,
				Name:     "W",
				Status:   models.StatusHidden,
				PuffsMax: 600,
				Flavors:  []string{"mint"},
				PhotoKey: strPtr("models/7/old.jpg"),
			}, nil
		},
		setPhotoKeyFn: func(_ context.Context, id uint64, key *string) (models.Model, error) {
			setPhotoKeyCalls++
			if id != 7 {
				t.Fatalf("SetPhotoKey() id = %d, want 7", id)
			}
			if key == nil || *key == "" {
				t.Fatalf("SetPhotoKey() key = %#v, want non-empty", key)
			}
			newKey = *key
			return models.Model{
				ID:       7,
				Name:     "W",
				Status:   models.StatusHidden,
				PuffsMax: 600,
				Flavors:  []string{"mint"},
				PhotoKey: strPtr(*key),
			}, nil
		},
	}
	s3 := &fakeS3{publicBase: "https://cdn.example/files"}
	mux := newModelsMux(svc, s3, false)

	jpg := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46, 0x00}
	req := newMultipartRequest(t, http.MethodPost, "/api/models/7/photo", true, "avatar.jpg", jpg)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
	if setPhotoKeyCalls != 1 {
		t.Fatalf("SetPhotoKey() calls = %d, want 1", setPhotoKeyCalls)
	}
	if !strings.HasPrefix(newKey, "models/7/") {
		t.Fatalf("new photo key = %q, want prefix models/7/", newKey)
	}
	if !strings.HasSuffix(newKey, ".jpg") {
		t.Fatalf("new photo key = %q, want .jpg suffix", newKey)
	}
	if len(s3.putCalls) != 1 {
		t.Fatalf("Put() calls = %d, want 1", len(s3.putCalls))
	}
	if s3.putCalls[0].key != newKey {
		t.Fatalf("Put() key = %q, want %q", s3.putCalls[0].key, newKey)
	}
	if s3.putCalls[0].contentType != "image/jpeg" {
		t.Fatalf("Put() content type = %q, want image/jpeg", s3.putCalls[0].contentType)
	}
	if len(s3.deleteCalls) != 1 || s3.deleteCalls[0] != "models/7/old.jpg" {
		t.Fatalf("Delete() calls = %#v, want delete old key", s3.deleteCalls)
	}

	var got models.Model
	decodeJSON(t, rr, &got)
	if got.PhotoKey == nil || *got.PhotoKey != newKey {
		t.Fatalf("photo_key = %#v, want %q", got.PhotoKey, newKey)
	}
	if got.PhotoURL == nil || *got.PhotoURL != "https://cdn.example/files/"+newKey {
		t.Fatalf("photo_url = %#v, want computed url", got.PhotoURL)
	}
}

func TestModelsHandlerDeletePhotoSuccessBestEffortS3Delete(t *testing.T) {
	t.Parallel()

	var setPhotoKeyArg *string
	svc := &fakeModelsService{
		getFn: func(_ context.Context, id uint64) (models.Model, error) {
			if id != 11 {
				t.Fatalf("Get() id = %d, want 11", id)
			}
			return models.Model{
				ID:       11,
				Name:     "W",
				Status:   models.StatusHidden,
				PuffsMax: 700,
				Flavors:  []string{"mint"},
				PhotoKey: strPtr("models/11/old.jpg"),
			}, nil
		},
		setPhotoKeyFn: func(_ context.Context, id uint64, key *string) (models.Model, error) {
			if id != 11 {
				t.Fatalf("SetPhotoKey() id = %d, want 11", id)
			}
			setPhotoKeyArg = key
			return models.Model{
				ID:       11,
				Name:     "W",
				Status:   models.StatusHidden,
				PuffsMax: 700,
				Flavors:  []string{"mint"},
				PhotoKey: nil,
			}, nil
		},
	}
	s3 := &fakeS3{
		publicBase: "https://cdn.example/files",
		deleteErr:  errors.New("s3 delete failed"),
	}
	mux := newModelsMux(svc, s3, false)

	rr := doRequest(mux, http.MethodDelete, "/api/models/11/photo", nil, "")
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
	if setPhotoKeyArg != nil {
		t.Fatalf("SetPhotoKey() arg = %#v, want nil", setPhotoKeyArg)
	}
	if len(s3.deleteCalls) != 1 || s3.deleteCalls[0] != "models/11/old.jpg" {
		t.Fatalf("Delete() calls = %#v, want one old key", s3.deleteCalls)
	}

	var got models.Model
	decodeJSON(t, rr, &got)
	if got.PhotoKey != nil {
		t.Fatalf("photo_key = %#v, want nil", got.PhotoKey)
	}
	if got.PhotoURL != nil {
		t.Fatalf("photo_url = %#v, want nil", got.PhotoURL)
	}
}

func TestModelsHandlerDeleteModelBestEffortS3Delete(t *testing.T) {
	t.Parallel()

	events := make([]string, 0, 3)
	svc := &fakeModelsService{
		getFn: func(_ context.Context, id uint64) (models.Model, error) {
			events = append(events, "svc.get")
			if id != 42 {
				t.Fatalf("Get() id = %d, want 42", id)
			}
			return models.Model{
				ID:       42,
				Name:     "W",
				Status:   models.StatusActive,
				PuffsMax: 1000,
				Flavors:  []string{"mint"},
				PhotoKey: strPtr("models/42/photo.jpg"),
			}, nil
		},
		deleteFn: func(_ context.Context, id uint64) error {
			events = append(events, "svc.delete")
			if id != 42 {
				t.Fatalf("Delete() id = %d, want 42", id)
			}
			return nil
		},
	}
	s3 := &fakeS3{
		deleteErr: errors.New("temporary s3 error"),
		events:    &events,
	}
	mux := newModelsMux(svc, s3, false)

	rr := doRequest(mux, http.MethodDelete, "/api/models/42", nil, "")
	if rr.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d body=%s", rr.Code, http.StatusNoContent, rr.Body.String())
	}
	if len(s3.deleteCalls) != 1 || s3.deleteCalls[0] != "models/42/photo.jpg" {
		t.Fatalf("Delete() calls = %#v, want one photo key", s3.deleteCalls)
	}

	wantOrder := []string{"svc.get", "svc.delete", "s3.delete"}
	if strings.Join(events, ",") != strings.Join(wantOrder, ",") {
		t.Fatalf("call order = %#v, want %#v", events, wantOrder)
	}
}

func newModelsMux(svc *fakeModelsService, s3 *fakeS3, usePresigned bool) *http.ServeMux {
	mux := http.NewServeMux()
	models.NewModelsHandler(mux, models.HandlerDeps{
		Service:      svc,
		S3:           s3,
		UsePresigned: usePresigned,
		PresignTTL:   3 * time.Minute,
	})
	return mux
}

func doJSONRequest(h http.Handler, method, path, body string) *httptest.ResponseRecorder {
	return doRequest(h, method, path, strings.NewReader(body), "application/json")
}

func doRequest(h http.Handler, method, path string, body io.Reader, contentType string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, body)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	return rr
}

func newMultipartRequest(t *testing.T, method, path string, withFile bool, filename string, file []byte) *http.Request {
	t.Helper()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	if withFile {
		part, err := writer.CreateFormFile("file", filename)
		if err != nil {
			t.Fatalf("CreateFormFile() error = %v", err)
		}
		if _, err := part.Write(file); err != nil {
			t.Fatalf("Write() error = %v", err)
		}
	} else {
		if err := writer.WriteField("x", "1"); err != nil {
			t.Fatalf("WriteField() error = %v", err)
		}
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func decodeJSON(t *testing.T, rr *httptest.ResponseRecorder, dst any) {
	t.Helper()
	if err := json.Unmarshal(rr.Body.Bytes(), dst); err != nil {
		t.Fatalf("json.Unmarshal() error = %v, body=%q", err, rr.Body.String())
	}
}

type fakeModelsService struct {
	createFn       func(ctx context.Context, req models.CreateModelRequest) (models.Model, error)
	listFn         func(ctx context.Context, limit, offset int) (models.ListModelsResponse, error)
	getFn          func(ctx context.Context, id uint64) (models.Model, error)
	updateFn       func(ctx context.Context, id uint64, req models.UpdateModelRequest) (models.Model, error)
	deleteFn       func(ctx context.Context, id uint64) error
	addFlavorFn    func(ctx context.Context, id uint64, value string) (models.Model, error)
	removeFlavorFn func(ctx context.Context, id uint64, value string) (models.Model, error)
	setPhotoKeyFn  func(ctx context.Context, id uint64, key *string) (models.Model, error)
}

func (f *fakeModelsService) Create(ctx context.Context, req models.CreateModelRequest) (models.Model, error) {
	if f.createFn == nil {
		panic("unexpected Create() call")
	}
	return f.createFn(ctx, req)
}

func (f *fakeModelsService) List(ctx context.Context, limit, offset int) (models.ListModelsResponse, error) {
	if f.listFn == nil {
		panic("unexpected List() call")
	}
	return f.listFn(ctx, limit, offset)
}

func (f *fakeModelsService) Get(ctx context.Context, id uint64) (models.Model, error) {
	if f.getFn == nil {
		panic("unexpected Get() call")
	}
	return f.getFn(ctx, id)
}

func (f *fakeModelsService) Update(ctx context.Context, id uint64, req models.UpdateModelRequest) (models.Model, error) {
	if f.updateFn == nil {
		panic("unexpected Update() call")
	}
	return f.updateFn(ctx, id, req)
}

func (f *fakeModelsService) Delete(ctx context.Context, id uint64) error {
	if f.deleteFn == nil {
		panic("unexpected Delete() call")
	}
	return f.deleteFn(ctx, id)
}

func (f *fakeModelsService) AddFlavor(ctx context.Context, id uint64, value string) (models.Model, error) {
	if f.addFlavorFn == nil {
		panic("unexpected AddFlavor() call")
	}
	return f.addFlavorFn(ctx, id, value)
}

func (f *fakeModelsService) RemoveFlavor(ctx context.Context, id uint64, value string) (models.Model, error) {
	if f.removeFlavorFn == nil {
		panic("unexpected RemoveFlavor() call")
	}
	return f.removeFlavorFn(ctx, id, value)
}

func (f *fakeModelsService) SetPhotoKey(ctx context.Context, id uint64, key *string) (models.Model, error) {
	if f.setPhotoKeyFn == nil {
		panic("unexpected SetPhotoKey() call")
	}
	return f.setPhotoKeyFn(ctx, id, key)
}

type putCall struct {
	key         string
	contentType string
	body        []byte
}

type fakeS3 struct {
	publicBase string

	putCalls    []putCall
	deleteCalls []string
	publicCalls []string

	presignedURL string
	presignedErr error
	deleteErr    error

	events *[]string
}

func (f *fakeS3) Put(_ context.Context, key string, body io.Reader, contentType string) error {
	payload, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	f.putCalls = append(f.putCalls, putCall{key: key, contentType: contentType, body: payload})
	return nil
}

func (f *fakeS3) Delete(_ context.Context, key string) error {
	if f.events != nil {
		*f.events = append(*f.events, "s3.delete")
	}
	f.deleteCalls = append(f.deleteCalls, key)
	return f.deleteErr
}

func (f *fakeS3) PublicURL(key string) string {
	f.publicCalls = append(f.publicCalls, key)
	base := strings.TrimRight(f.publicBase, "/")
	if base == "" {
		base = "https://example.invalid"
	}
	return base + "/" + strings.TrimLeft(key, "/")
}

func (f *fakeS3) PresignedGetURL(_ context.Context, key string, _ time.Duration) (string, error) {
	if f.presignedErr != nil {
		return "", f.presignedErr
	}
	if f.presignedURL != "" {
		return f.presignedURL, nil
	}
	return "https://signed.example/" + strings.TrimLeft(key, "/"), nil
}

func strPtr(v string) *string {
	return &v
}
