package favorites

import (
	"net/http"
	"rest_waka/pkg/httpx"
	"rest_waka/pkg/middleware"
	"rest_waka/pkg/res"
)

type HandlerDeps struct {
	Service   *Service
	JWTSecret string
}

type Handler struct {
	svc *Service
}

func NewFavoritesHandler(router *http.ServeMux, deps HandlerDeps) {
	handler := &Handler{
		svc: deps.Service,
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

		if err := handler.svc.Add(r.Context(), uint64(uid), modelID); err != nil {
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

		if err := handler.svc.Remove(r.Context(), uint64(uid), modelID); err != nil {
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

		data, err := handler.svc.List(r.Context(), uint64(uid), limit, offset)
		if err != nil {
			writeFavErr(w, err)
			return
		}
		res.Json(w, data, http.StatusOK)
	}
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
