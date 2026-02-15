package models

import (
	"net/http"
	"rest_waka/pkg/httpx"
	"rest_waka/pkg/req"
	"rest_waka/pkg/res"
)

type HandlerDeps struct {
	Service *Service
}
type Handler struct {
	svc *Service
}

func NewModelsHandler(router *http.ServeMux, deps HandlerDeps) {
	handler := &Handler{
		svc: deps.Service,
	}

	router.HandleFunc("POST /api/models", handler.CreateModels())
	router.HandleFunc("GET /api/models", handler.ListModels())
	router.HandleFunc("GET /api/models/{id}", handler.GetModels())
	router.HandleFunc("PATCH /api/models/{id}", handler.UpdateModels())
	router.HandleFunc("DELETE /api/models/{id}", handler.DeleteModels())

	// flavors pen
	router.HandleFunc("POST /api/models/{id}/flavors", handler.AddFlavor())
	router.HandleFunc("DELETE /api/models/{id}/flavors", handler.RemoveFlavor())
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

		res.Json(w, m, http.StatusOK)
	}
}

func (handler *Handler) DeleteModels() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := httpx.PathUint64(r, "id")
		if err != nil {
			res.Json(w, "invalid id", http.StatusBadRequest)
			return
		}

		if err := handler.svc.Delete(r.Context(), id); err != nil {
			writeModelErr(w, err)
			return
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
