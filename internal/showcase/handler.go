package showcase

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path"
	"strings"
	"time"

	"rest_waka/pkg/httpx"
	"rest_waka/pkg/middleware"
	"rest_waka/pkg/photourl"
	"rest_waka/pkg/randHex"
	"rest_waka/pkg/req"
	"rest_waka/pkg/res"
)

type showcaseService interface {
	Create(ctx context.Context, req CreateItemRequest) (Item, error)
	List(ctx context.Context, limit, offset int) (ListItemsResponse, error)
	Get(ctx context.Context, id uint64) (Item, error)
	Update(ctx context.Context, id uint64, req UpdateItemRequest) (Item, error)
	Delete(ctx context.Context, id uint64) error
	SetPhotoKey(ctx context.Context, id uint64, key *string) (Item, error)

	ListActive(ctx context.Context, limit, offset int) (ListItemsResponse, error)
	GetActive(ctx context.Context, id uint64) (Item, error)
}

type photoStore interface {
	Put(ctx context.Context, key string, body io.Reader, contentType string) error
	Delete(ctx context.Context, key string) error
	photourl.Resolver
}

type HandlerDeps struct {
	Service      showcaseService
	JWTSecret    string
	S3           photoStore
	UsePresigned bool
	PresignTTL   time.Duration
}

type Handler struct {
	svc          showcaseService
	s3           photoStore
	usePresigned bool
	presignTTL   time.Duration
}

func NewShowcaseHandler(router *http.ServeMux, deps HandlerDeps) {
	handler := &Handler{
		svc:          deps.Service,
		s3:           deps.S3,
		usePresigned: deps.UsePresigned,
		presignTTL:   deps.PresignTTL,
	}

	// public
	router.HandleFunc("GET /api/showcase", handler.ListActive())
	router.HandleFunc("GET /api/showcase/{id}", handler.GetActive())

	// admin
	router.Handle("POST /api/admin/showcase", middleware.RequireAdmin(handler.Create(), deps.JWTSecret))
	router.Handle("GET /api/admin/showcase", middleware.RequireAdmin(handler.List(), deps.JWTSecret))
	router.Handle("GET /api/admin/showcase/{id}", middleware.RequireAdmin(handler.Get(), deps.JWTSecret))
	router.Handle("PATCH /api/admin/showcase/{id}", middleware.RequireAdmin(handler.Update(), deps.JWTSecret))
	router.Handle("DELETE /api/admin/showcase/{id}", middleware.RequireAdmin(handler.Delete(), deps.JWTSecret))

	router.Handle("POST /api/admin/showcase/{id}/photo", middleware.RequireAdmin(handler.UploadPhoto(), deps.JWTSecret))
	router.Handle("DELETE /api/admin/showcase/{id}/photo", middleware.RequireAdmin(handler.DeletePhoto(), deps.JWTSecret))
}

func (handler *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.Decode[CreateItemRequest](r.Body)
		if err != nil {
			res.Json(w, "invalid json body", http.StatusBadRequest)
			return
		}

		item, err := handler.svc.Create(r.Context(), body)
		if err != nil {
			writeShowcaseErr(w, err)
			return
		}

		item.PhotoURL = handler.resolvePhotoURL(r.Context(), item.PhotoKey)
		res.Json(w, item, http.StatusCreated)
	}
}

func (handler *Handler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit := httpx.QueryInt(r, "limit", 50)
		offset := httpx.QueryInt(r, "offset", 0)

		data, err := handler.svc.List(r.Context(), limit, offset)
		if err != nil {
			res.Json(w, "internal error", http.StatusInternalServerError)
			return
		}

		for i := range data.Items {
			data.Items[i].PhotoURL = handler.resolvePhotoURL(r.Context(), data.Items[i].PhotoKey)
		}

		res.Json(w, data, http.StatusOK)
	}
}

func (handler *Handler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := httpx.PathUint64(r, "id")
		if err != nil {
			res.Json(w, "invalid id", http.StatusBadRequest)
			return
		}

		item, err := handler.svc.Get(r.Context(), id)
		if err != nil {
			writeShowcaseErr(w, err)
			return
		}

		item.PhotoURL = handler.resolvePhotoURL(r.Context(), item.PhotoKey)
		res.Json(w, item, http.StatusOK)
	}
}

func (handler *Handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := httpx.PathUint64(r, "id")
		if err != nil {
			res.Json(w, "invalid id", http.StatusBadRequest)
			return
		}

		body, err := req.Decode[UpdateItemRequest](r.Body)
		if err != nil {
			res.Json(w, "invalid json body", http.StatusBadRequest)
			return
		}

		item, err := handler.svc.Update(r.Context(), id, body)
		if err != nil {
			writeShowcaseErr(w, err)
			return
		}

		item.PhotoURL = handler.resolvePhotoURL(r.Context(), item.PhotoKey)
		res.Json(w, item, http.StatusOK)
	}
}

func (handler *Handler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := httpx.PathUint64(r, "id")
		if err != nil || id == 0 {
			res.Json(w, "invalid id", http.StatusBadRequest)
			return
		}

		var key string
		if handler.s3 != nil {
			item, err := handler.svc.Get(r.Context(), id)
			if err == nil && item.PhotoKey != nil {
				key = *item.PhotoKey
			}
		}

		if err := handler.svc.Delete(r.Context(), id); err != nil {
			writeShowcaseErr(w, err)
			return
		}

		if handler.s3 != nil && key != "" {
			if err := handler.s3.Delete(r.Context(), key); err != nil {
				log.Printf("WARN: failed to delete showcase photo: showcase_id=%d key=%s err=%v", id, key, err)
			}
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (handler *Handler) UploadPhoto() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if handler.s3 == nil {
			res.Json(w, "s3 is not configuration", http.StatusInternalServerError)
			return
		}

		id, err := httpx.PathUint64(r, "id")
		if err != nil || id == 0 {
			res.Json(w, "invalid id", http.StatusBadRequest)
			return
		}

		cur, err := handler.svc.Get(r.Context(), id)
		if err != nil {
			writeShowcaseErr(w, err)
			return
		}

		oldKey := ""
		if cur.PhotoKey != nil {
			oldKey = *cur.PhotoKey
		}

		const maxSize = 10 << 20
		r.Body = http.MaxBytesReader(w, r.Body, maxSize)
		if err := r.ParseMultipartForm(maxSize); err != nil {
			res.Json(w, "invalid multipart form", http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			res.Json(w, "file is required", http.StatusBadRequest)
			return
		}
		defer func(file multipart.File) {
			_ = file.Close()
		}(file)

		data, err := io.ReadAll(file)
		if err != nil {
			res.Json(w, "cannot read file", http.StatusBadRequest)
			return
		}
		if len(data) == 0 {
			res.Json(w, "empty file", http.StatusBadRequest)
			return
		}

		contentType := http.DetectContentType(data[:min(512, len(data))])
		if !strings.HasPrefix(contentType, "image/") {
			res.Json(w, "only images are allowed", http.StatusBadRequest)
			return
		}

		body := bytes.NewReader(data)

		ext := strings.ToLower(path.Ext(header.Filename))
		if ext == "" {
			ext = ".jpg"
		}

		newKey := fmt.Sprintf("showcase/%d/%s%s", id, randHex.RandHex(16), ext)

		if err := handler.s3.Put(r.Context(), newKey, body, contentType); err != nil {
			res.Json(w, "internal error", http.StatusInternalServerError)
			return
		}

		item, err := handler.svc.SetPhotoKey(r.Context(), id, &newKey)
		if err != nil {
			_ = handler.s3.Delete(r.Context(), newKey)
			writeShowcaseErr(w, err)
			return
		}

		if oldKey != "" && oldKey != newKey {
			if err := handler.s3.Delete(r.Context(), oldKey); err != nil {
				log.Printf("WARN: failed to delete old showcase photo: showcase_id=%d old_key=%s err=%v", id, oldKey, err)
			}
		}

		item.PhotoURL = handler.resolvePhotoURL(r.Context(), item.PhotoKey)
		res.Json(w, item, http.StatusOK)
	}
}

func (handler *Handler) DeletePhoto() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if handler.s3 == nil {
			res.Json(w, "s3 is not configuration", http.StatusInternalServerError)
			return
		}

		id, err := httpx.PathUint64(r, "id")
		if err != nil || id == 0 {
			res.Json(w, "invalid id", http.StatusBadRequest)
			return
		}

		cur, err := handler.svc.Get(r.Context(), id)
		if err != nil {
			writeShowcaseErr(w, err)
			return
		}

		oldKey := ""
		if cur.PhotoKey != nil {
			oldKey = *cur.PhotoKey
		}

		item, err := handler.svc.SetPhotoKey(r.Context(), id, nil)
		if err != nil {
			writeShowcaseErr(w, err)
			return
		}

		if oldKey != "" {
			if err := handler.s3.Delete(r.Context(), oldKey); err != nil {
				log.Printf("WARN: failed to delete showcase photo: showcase_id=%d key=%s err=%v", id, oldKey, err)
			}
		}

		item.PhotoURL = handler.resolvePhotoURL(r.Context(), item.PhotoKey)
		res.Json(w, item, http.StatusOK)
	}
}

func (handler *Handler) ListActive() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit := httpx.QueryInt(r, "limit", 5)
		offset := httpx.QueryInt(r, "offset", 0)

		data, err := handler.svc.ListActive(r.Context(), limit, offset)
		if err != nil {
			res.Json(w, "internal error", http.StatusInternalServerError)
			return
		}

		out := ListPublicItemsResponse{
			Items:  make([]PublicItem, 0, len(data.Items)),
			Limit:  data.Limit,
			Offset: data.Offset,
		}

		for i := range data.Items {
			data.Items[i].PhotoURL = handler.resolvePhotoURL(r.Context(), data.Items[i].PhotoKey)
			out.Items = append(out.Items, toPublicItem(data.Items[i]))
		}

		res.Json(w, out, http.StatusOK)
	}
}

func (handler *Handler) GetActive() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := httpx.PathUint64(r, "id")
		if err != nil {
			res.Json(w, "invalid id", http.StatusBadRequest)
			return
		}

		item, err := handler.svc.GetActive(r.Context(), id)
		if err != nil {
			writeShowcaseErr(w, err)
			return
		}

		item.PhotoURL = handler.resolvePhotoURL(r.Context(), item.PhotoKey)
		res.Json(w, toPublicItem(item), http.StatusOK)
	}
}

func toPublicItem(item Item) PublicItem {
	return PublicItem{
		ID:          item.ID,
		Tag:         item.Tag,
		Title:       item.Title,
		Description: item.Description,
		ModelID:     item.ModelID,
		PhotoURL:    item.PhotoURL,
		Sort:        item.Sort,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}

func (handler *Handler) resolvePhotoURL(ctx context.Context, key *string) *string {
	return photourl.Resolve(ctx, handler.s3, key, photourl.Options{
		UsePresigned: handler.usePresigned,
		PresignTTL:   handler.presignTTL,
	})
}

func writeShowcaseErr(w http.ResponseWriter, err error) {
	httpx.WriteMappedError(w, err,
		[]httpx.ErrMap{
			{Err: ErrInvalidArgument, Status: http.StatusBadRequest, Message: "invalid argument"},
			{Err: ErrNotFound, Status: http.StatusNotFound, Message: "not found"},
		},
		http.StatusInternalServerError,
		"internal error",
	)
}
