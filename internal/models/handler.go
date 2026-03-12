package models

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path"
	"rest_waka/pkg/httpx"
	"rest_waka/pkg/photourl"
	"rest_waka/pkg/randHex"
	"rest_waka/pkg/req"
	"rest_waka/pkg/res"
	"strings"
	"time"
)

type modelsService interface {
	Create(ctx context.Context, req CreateModelRequest) (Model, error)
	List(ctx context.Context, limit, offset int) (ListModelsResponse, error)
	Get(ctx context.Context, id uint64) (Model, error)
	Update(ctx context.Context, id uint64, req UpdateModelRequest) (Model, error)
	Delete(ctx context.Context, id uint64) error
	AddFlavor(ctx context.Context, id uint64, value string) (Model, error)
	RemoveFlavor(ctx context.Context, id uint64, value string) (Model, error)
	SetPhotoKey(ctx context.Context, id uint64, key *string) (Model, error)

	ListActive(ctx context.Context, limit, offset int) (ListModelsResponse, error)
	GetActive(ctx context.Context, id uint64) (Model, error)
}

type photoStore interface {
	Put(ctx context.Context, key string, body io.Reader, contentType string) error
	Delete(ctx context.Context, key string) error
	photourl.Resolver
}

type HandlerDeps struct {
	Service      modelsService
	S3           photoStore
	UsePresigned bool
	PresignTTL   time.Duration
}
type Handler struct {
	svc          modelsService
	s3           photoStore
	usePresigned bool
	presignTTL   time.Duration
}

func NewModelsHandler(router *http.ServeMux, deps HandlerDeps) {
	handler := &Handler{
		svc:          deps.Service,
		s3:           deps.S3,
		usePresigned: deps.UsePresigned,
		presignTTL:   deps.PresignTTL,
	}

	router.HandleFunc("POST /api/models", handler.CreateModels())
	router.HandleFunc("GET /api/models", handler.ListModels())
	router.HandleFunc("GET /api/models/{id}", handler.GetModels())
	router.HandleFunc("PATCH /api/models/{id}", handler.UpdateModels())
	router.HandleFunc("DELETE /api/models/{id}", handler.DeleteModels())

	// flavors pen
	router.HandleFunc("POST /api/models/{id}/flavors", handler.AddFlavor())
	router.HandleFunc("DELETE /api/models/{id}/flavors", handler.RemoveFlavor())

	//	upload photo
	router.HandleFunc("POST /api/models/{id}/photo", handler.UploadPhoto())
	router.HandleFunc("DELETE /api/models/{id}/photo", handler.DeletePhoto())

	// public pen
	router.HandleFunc("GET /api/catalog/models", handler.ActiveListModels())
	router.HandleFunc("GET /api/catalog/models/{id}", handler.ActiveGetModels())
}

func (handler *Handler) CreateModels() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.Decode[CreateModelRequest](r.Body)
		if err != nil {
			res.Json(w, "invalid json body", http.StatusBadRequest)
			return
		}

		m, err := handler.svc.Create(r.Context(), body)
		if err != nil {
			writeModelErr(w, err)
			return
		}
		m.PhotoURL = handler.resolvePhotoURL(r.Context(), m.PhotoKey)
		res.Json(w, m, http.StatusCreated)
	}
}

func (handler *Handler) ListModels() http.HandlerFunc {
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

func (handler *Handler) GetModels() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := httpx.PathUint64(r, "id")
		if err != nil {
			res.Json(w, "invalid id", http.StatusBadRequest)
			return
		}

		m, err := handler.svc.Get(r.Context(), id)
		if err != nil {
			writeModelErr(w, err)
			return
		}
		m.PhotoURL = handler.resolvePhotoURL(r.Context(), m.PhotoKey)
		res.Json(w, m, http.StatusOK)
	}
}

func (handler *Handler) UpdateModels() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := httpx.PathUint64(r, "id")
		if err != nil {
			res.Json(w, "invalid id", http.StatusBadRequest)
			return
		}

		body, err := req.Decode[UpdateModelRequest](r.Body)
		if err != nil {
			res.Json(w, "invalid json body", http.StatusBadRequest)
			return
		}

		m, err := handler.svc.Update(r.Context(), id, body)
		if err != nil {
			writeModelErr(w, err)
			return
		}

		m.PhotoURL = handler.resolvePhotoURL(r.Context(), m.PhotoKey)
		res.Json(w, m, http.StatusOK)
	}
}

func (handler *Handler) DeleteModels() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := httpx.PathUint64(r, "id")
		if err != nil || id == 0 {
			res.Json(w, "invalid id", http.StatusBadRequest)
			return
		}

		// узнаем photo_key
		var key string
		if handler.s3 != nil {
			m, err := handler.svc.Get(r.Context(), id)
			if err == nil && m.PhotoKey != nil {
				key = *m.PhotoKey
			}
		}

		// удаляем запись
		if err := handler.svc.Delete(r.Context(), id); err != nil {
			writeModelErr(w, err)
			return
		}

		// удаляем файл из s3
		if handler.s3 != nil && key != "" {
			if err := handler.s3.Delete(r.Context(), key); err != nil {
				log.Printf("WARN: failed to delete model photo: model_id=%d key=%s err=%v", id, key, err)
			}
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// Flavor pens

func (handler *Handler) AddFlavor() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := httpx.PathUint64(r, "id")
		if err != nil {
			res.Json(w, "invalid id", http.StatusBadRequest)
			return
		}

		body, err := req.Decode[FlavorRequest](r.Body)
		if err != nil {
			res.Json(w, "invalid json body", http.StatusBadRequest)
			return
		}

		m, err := handler.svc.AddFlavor(r.Context(), id, body.Value)
		if err != nil {
			writeModelErr(w, err)
			return
		}

		res.Json(w, m, http.StatusOK)
	}
}

func (handler *Handler) RemoveFlavor() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := httpx.PathUint64(r, "id")
		if err != nil {
			res.Json(w, "invalid id", http.StatusBadRequest)
			return
		}

		body, err := req.Decode[FlavorRequest](r.Body)
		if err != nil {
			res.Json(w, "invalid json body", http.StatusBadRequest)
			return
		}

		m, err := handler.svc.RemoveFlavor(r.Context(), id, body.Value)
		if err != nil {
			writeModelErr(w, err)
			return
		}
		res.Json(w, m, http.StatusOK)
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

		// запоминаем старый ключ, если он есть
		cur, err := handler.svc.Get(r.Context(), id)
		if err != nil {
			writeModelErr(w, err)
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
		// defer file.Close()
		defer func(file multipart.File) {
			err := file.Close()
			if err != nil {

			}
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

		newKey := fmt.Sprintf("models/%d/%s%s", id, randHex.RandHex(16), ext)

		// грузим новое
		if err := handler.s3.Put(r.Context(), newKey, body, contentType); err != nil {
			res.Json(w, "internal error", http.StatusInternalServerError)
			return
		}

		// обновляем photo_key в бд
		m, err := handler.svc.SetPhotoKey(r.Context(), id, &newKey)
		if err != nil {
			// если произойдет ошибка - удаляем новое, чтобы не оставить сироту
			_ = handler.s3.Delete(r.Context(), newKey)
			writeModelErr(w, err)
			return
		}

		// удаляем старый объект(если был и отличается)
		if oldKey != "" && oldKey != newKey {
			if err := handler.s3.Delete(r.Context(), oldKey); err != nil {
				log.Printf("WARN: failed to delete old photo: model_id=%d old_key=%s err=%v", id, oldKey, err)
			}
		}

		m.PhotoURL = handler.resolvePhotoURL(r.Context(), m.PhotoKey)
		res.Json(w, m, http.StatusOK)
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

		// берем текущий ключ (чтобы потом удалить объект)
		cur, err := handler.svc.Get(r.Context(), id)
		if err != nil {
			writeModelErr(w, err)
			return
		}

		oldKey := ""
		if cur.PhotoKey != nil {
			oldKey = *cur.PhotoKey
		}

		// очищаем сначала бд
		m, err := handler.svc.SetPhotoKey(r.Context(), id, nil) // null в бд
		if err != nil {
			writeModelErr(w, err)
			return
		}

		// удаляем объект из s3
		if oldKey != "" {
			if err := handler.s3.Delete(r.Context(), oldKey); err != nil {
				log.Printf("WARN: failed to delete model photo: model_id=%d key=%s err=%v", id, oldKey, err)
			}
		}

		m.PhotoURL = handler.resolvePhotoURL(r.Context(), m.PhotoKey)
		res.Json(w, m, http.StatusOK)
	}
}

func (handler *Handler) ActiveListModels() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit := httpx.QueryInt(r, "limit", 50)
		offset := httpx.QueryInt(r, "offset", 0)

		data, err := handler.svc.ListActive(r.Context(), limit, offset)
		if err != nil {
			res.Json(w, "internal error", http.StatusInternalServerError)
			return
		}

		out := ListPublicModelsResponse{
			Items:  make([]PublicModel, 0, len(data.Items)),
			Limit:  data.Limit,
			Offset: data.Offset,
		}

		for i := range data.Items {
			data.Items[i].PhotoURL = handler.resolvePhotoURL(r.Context(), data.Items[i].PhotoKey)
			out.Items = append(out.Items, toPublicModel(data.Items[i]))
		}

		res.Json(w, out, http.StatusOK)
	}
}

func (handler *Handler) ActiveGetModels() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := httpx.PathUint64(r, "id")
		if err != nil {
			res.Json(w, "invalid id", http.StatusBadRequest)
			return
		}

		m, err := handler.svc.GetActive(r.Context(), id)
		if err != nil {
			writeModelErr(w, err)
			return
		}

		m.PhotoURL = handler.resolvePhotoURL(r.Context(), m.PhotoKey)

		out := toPublicModel(m)
		res.Json(w, out, http.StatusOK)
	}
}

// toPublicModel - mapper для публичной ручки
func toPublicModel(m Model) PublicModel {
	return PublicModel{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		PhotoURL:    m.PhotoURL,
		PuffsMax:    m.PuffsMax,
		Flavors:     m.Flavors,
		PriceCents:  m.PriceCents,
	}
}

// resolvePhotoURL - приватный метод для того, чтобы укоротить код
func (handler *Handler) resolvePhotoURL(ctx context.Context, key *string) *string {
	return photourl.Resolve(ctx, handler.s3, key, photourl.Options{
		UsePresigned: handler.usePresigned,
		PresignTTL:   handler.presignTTL,
	})
}

func writeModelErr(w http.ResponseWriter, err error) {
	httpx.WriteMappedError(w, err,
		[]httpx.ErrMap{
			{Err: ErrInvalidArgument, Status: http.StatusBadRequest, Message: "invalid argument"},
			{Err: ErrNotFound, Status: http.StatusNotFound, Message: "not found"},
		},
		http.StatusInternalServerError,
		"internal error",
	)
}
