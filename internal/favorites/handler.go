package favorites

import (
	"context"
	"net/http"
	"rest_waka/pkg/httpx"
	"rest_waka/pkg/middleware"
	"rest_waka/pkg/photourl"
	"rest_waka/pkg/res"
	"time"
)

type favoritesService interface {
	Add(ctx context.Context, userID, modelID uint64) error
	Remove(ctx context.Context, userID, modelID uint64) error
	List(ctx context.Context, userID uint64, limit, offset int) (ListFavoritesResponse, error)
}

type HandlerDeps struct {
	Service      favoritesService
	JWTSecret    string
	S3           photourl.Resolver
	UsePresigned bool
	PresignTTL   time.Duration
}

type Handler struct {
	svc          favoritesService
	s3           photourl.Resolver
	usePresigned bool
	presignTTL   time.Duration
}

func NewFavoritesHandler(router *http.ServeMux, deps HandlerDeps) {
	handler := &Handler{
		svc:          deps.Service,
		s3:           deps.S3,
		usePresigned: deps.UsePresigned,
		presignTTL:   deps.PresignTTL,
	}

	router.Handle("GET /api/favorites", middleware.RequireUser(handler.List(), deps.JWTSecret))
	router.Handle("POST /api/favorites/{model_id}", middleware.RequireUser(handler.Add(), deps.JWTSecret))
	router.Handle("DELETE /api/favorites/{model_id}", middleware.RequireUser(handler.Remove(), deps.JWTSecret))

}

func (handler *Handler) Add() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid, ok := middleware.UserIDFromContext(r.Context())
		if !ok || uid == 0 {
			res.Json(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		modelID, err := httpx.PathUint64(r, "model_id")
		if err != nil || modelID == 0 {
			res.Json(w, "invalid id", http.StatusBadRequest)
			return
		}

		if err := handler.svc.Add(r.Context(), uid, modelID); err != nil {
			writeFavErr(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (handler *Handler) Remove() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid, ok := middleware.UserIDFromContext(r.Context())
		if !ok || uid == 0 {
			res.Json(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		modelID, err := httpx.PathUint64(r, "model_id")
		if err != nil || modelID == 0 {
			res.Json(w, "invalid id", http.StatusBadRequest)
			return
		}

		if err := handler.svc.Remove(r.Context(), uid, modelID); err != nil {
			writeFavErr(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (handler *Handler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid, ok := middleware.UserIDFromContext(r.Context())
		if !ok || uid == 0 {
			res.Json(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		limit := httpx.QueryInt(r, "limit", 50)
		offset := httpx.QueryInt(r, "offset", 0)

		data, err := handler.svc.List(r.Context(), uid, limit, offset)
		if err != nil {
			writeFavErr(w, err)
			return
		}
		for i := range data.Items {
			data.Items[i].PhotoURL = handler.resolvePhotoURL(r.Context(), data.Items[i].PhotoKey)
		}
		res.Json(w, data, http.StatusOK)
	}
}

// resolvePhotoURL - приватный метод для того, чтобы укоротить код
func (handler *Handler) resolvePhotoURL(ctx context.Context, key *string) *string {
	return photourl.Resolve(ctx, handler.s3, key, photourl.Options{
		UsePresigned: handler.usePresigned,
		PresignTTL:   handler.presignTTL,
	})
}

func writeFavErr(w http.ResponseWriter, err error) {
	httpx.WriteMappedError(w, err,
		[]httpx.ErrMap{
			{Err: ErrInvalidArgument, Status: http.StatusBadRequest, Message: "invalid argument"},
			{Err: ErrNotFound, Status: http.StatusNotFound, Message: "not found"},
			{Err: ErrAlreadyExists, Status: http.StatusConflict, Message: "already exists"},
		},
		http.StatusInternalServerError,
		"internal error",
	)
}
