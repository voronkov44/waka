package auth

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

func NewAuthHandler(router *http.ServeMux, deps HandlerDeps) {
	handler := &Handler{
		svc: deps.Service,
	}

	router.HandleFunc("POST /api/auth/telegram", handler.LoginTelegram())
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
