package users

import (
	"net/http"
	"rest_waka/pkg/httpx"
	"rest_waka/pkg/res"
)

type HandlerDeps struct {
	Service *Service
}

type Handler struct {
	svc *Service
}

func NewUsersHandler(router *http.ServeMux, deps HandlerDeps) {
	handler := &Handler{svc: deps.Service}

	router.HandleFunc("GET /api/users", handler.List())
	router.HandleFunc("GET /api/users/{id}", handler.Get())
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

		user, err := handler.svc.Get(r.Context(), id)
		if err != nil {
			writeUserErr(w, err)
			return
		}
		res.Json(w, user, http.StatusOK)
	}
}

func writeUserErr(w http.ResponseWriter, err error) {
	httpx.WriteMappedError(w, err,
		[]httpx.ErrMap{
			{Err: ErrInvalidArgument, Status: http.StatusBadRequest, Message: "invalid argument"},
			{Err: ErrNotFound, Status: http.StatusNotFound, Message: "not found"},
		},
		http.StatusInternalServerError,
		"internal error",
	)
}
