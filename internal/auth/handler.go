package auth

import (
	"net/http"
	"rest_waka/pkg/httpx"
	"rest_waka/pkg/middleware"
	"rest_waka/pkg/req"
	"rest_waka/pkg/res"
)

type HandlerDeps struct {
	Service   *Service
	JWTSecret string
}

type Handler struct {
	svc *Service
}

func NewAuthHandler(router *http.ServeMux, deps HandlerDeps) {
	handler := &Handler{
		svc: deps.Service,
	}

	router.HandleFunc("POST /api/auth/telegram", handler.LoginTelegram())
	router.Handle(
		"GET /api/auth/me",
		middleware.RequireUser(handler.Me(), deps.JWTSecret),
	)
}

func (handler *Handler) LoginTelegram() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.Decode[TelegramProfile](r.Body)
		if err != nil {
			res.Json(w, "invalid json body", http.StatusBadRequest)
			return
		}
		if err := req.IsValid(body); err != nil {
			res.Json(w, "invalid json body", http.StatusBadRequest)
			return
		}

		token, err := handler.svc.LoginTelegram(r.Context(), body)
		if err != nil {
			writeAuthErr(w, err)
			return
		}

		res.Json(w, TokenResponse{Token: token}, http.StatusOK)
	}
}

func (handler *Handler) Me() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid, ok := middleware.UserIDFromContext(r.Context())
		if !ok || uid == 0 {
			res.Json(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		data, err := handler.svc.Me(r.Context(), uint64(uid))
		if err != nil {
			writeAuthErr(w, err)
			return
		}

		res.Json(w, data, http.StatusOK)
	}
}

func writeAuthErr(w http.ResponseWriter, err error) {
	httpx.WriteMappedError(w, err,
		[]httpx.ErrMap{
			{Err: ErrInvalidArgument, Status: http.StatusBadRequest, Message: "invalid argument"},
			{Err: ErrNotFound, Status: http.StatusNotFound, Message: "not found"},
		},
		http.StatusInternalServerError,
		"internal error",
	)
}
