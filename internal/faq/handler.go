package faq

import (
	"net/http"
	"rest_waka/pkg/httpx"
	"rest_waka/pkg/req"
	"rest_waka/pkg/res"
)

type HandlerDeps struct {
	Service Service
}

type Handler struct {
	svc Service
}

func NewFaqHandler(router *http.ServeMux, deps HandlerDeps) {
	handler := &Handler{svc: deps.Service}

	// public
	router.HandleFunc("GET /api/faq/topics", handler.ListTopics())
	router.HandleFunc("GET /api/faq/topics/{topicID}/articles", handler.ListArticlesByTopic())
	router.HandleFunc("GET /api/faq/articles/{id}", handler.GetArticle())
	router.HandleFunc("GET /api/faq/search", handler.Search())

	// admin
	router.HandleFunc("POST /api/admin/faq/topics", handler.CreateTopic())
	router.HandleFunc("PATCH /api/admin/faq/topics/{id}", handler.UpdateTopic())

	router.HandleFunc("POST /api/admin/faq/articles", handler.CreateArticle())
	router.HandleFunc("PATCH /api/admin/faq/articles/{id}", handler.UpdateArticle())
	router.HandleFunc("PUT /api/admin/faq/articles/{id}/blocks", handler.PutBlocks())
}

// public

func (handler *Handler) ListTopics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := handler.svc.ListTopics(r.Context())
		if err != nil {
			writeFaqErr(w, err)
			return
		}
		res.Json(w, data, http.StatusOK)
	}
}

func (handler *Handler) ListArticlesByTopic() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		topicID, err := httpx.PathUint64(r, "topicID")
		if err != nil {
			res.Json(w, "invalid topicID", http.StatusBadRequest)
			return
		}

		channel := r.URL.Query().Get("channel") // all|tg|miniapp (пусто -> all)
		data, err := handler.svc.ListArticlesByTopic(r.Context(), topicID, channel)
		if err != nil {
			writeFaqErr(w, err)
			return
		}
		res.Json(w, data, http.StatusOK)
	}
}

func (handler *Handler) GetArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := httpx.PathUint64(r, "id")
		if err != nil {
			res.Json(w, "invalid id", http.StatusBadRequest)
			return
		}

		data, err := handler.svc.GetArticle(r.Context(), id)
		if err != nil {
			writeFaqErr(w, err)
			return
		}
		res.Json(w, data, http.StatusOK)
	}
}

func (handler *Handler) Search() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		channel := r.URL.Query().Get("channel")
		limit := httpx.QueryInt(r, "limit", 20)
		offset := httpx.QueryInt(r, "offset", 0)

		data, err := handler.svc.Search(r.Context(), q, channel, limit, offset)
		if err != nil {
			writeFaqErr(w, err)
			return
		}
		res.Json(w, data, http.StatusOK)
	}
}

// admin

func (handler *Handler) CreateTopic() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := req.Decode[CreateTopicRequest](r.Body)
		if err != nil {
			res.Json(w, "invalid body", http.StatusBadRequest)
			return
		}

		data, err := handler.svc.CreateTopic(r.Context(), payload)
		if err != nil {
			writeFaqErr(w, err)
			return
		}
		res.Json(w, data, http.StatusCreated)
	}
}

func (handler *Handler) UpdateTopic() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := httpx.PathUint64(r, "id")
		if err != nil {
			res.Json(w, "invalid id", http.StatusBadRequest)
			return
		}

		payload, err := req.Decode[UpdateTopicRequest](r.Body)
		if err != nil {
			res.Json(w, "invalid body", http.StatusBadRequest)
			return
		}

		data, err := handler.svc.UpdateTopic(r.Context(), id, payload)
		if err != nil {
			writeFaqErr(w, err)
			return
		}
		res.Json(w, data, http.StatusOK)
	}
}

func (handler *Handler) CreateArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := req.Decode[CreateArticleRequest](r.Body)
		if err != nil {
			res.Json(w, "invalid body", http.StatusBadRequest)
			return
		}

		data, err := handler.svc.CreateArticle(r.Context(), payload)
		if err != nil {
			writeFaqErr(w, err)
			return
		}
		res.Json(w, data, http.StatusCreated)
	}
}

func (handler *Handler) UpdateArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := httpx.PathUint64(r, "id")
		if err != nil {
			res.Json(w, "invalid id", http.StatusBadRequest)
			return
		}

		payload, err := req.Decode[UpdateArticleRequest](r.Body)
		if err != nil {
			res.Json(w, "invalid body", http.StatusBadRequest)
			return
		}

		data, err := handler.svc.UpdateArticle(r.Context(), id, payload)
		if err != nil {
			writeFaqErr(w, err)
			return
		}
		res.Json(w, data, http.StatusOK)
	}
}

func (handler *Handler) PutBlocks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		articleID, err := httpx.PathUint64(r, "id")
		if err != nil {
			res.Json(w, "invalid id", http.StatusBadRequest)
			return
		}

		payload, err := req.Decode[PutBlocksRequest](r.Body)
		if err != nil {
			res.Json(w, "invalid body", http.StatusBadRequest)
			return
		}

		data, err := handler.svc.PutBlocks(r.Context(), articleID, payload)
		if err != nil {
			writeFaqErr(w, err)
			return
		}
		res.Json(w, data, http.StatusOK)
	}
}

func writeFaqErr(w http.ResponseWriter, err error) {
	httpx.WriteMappedError(w, err,
		[]httpx.ErrMap{
			{Err: ErrInvalidArgument, Status: http.StatusBadRequest, Message: "invalid argument"},
			{Err: ErrNotFound, Status: http.StatusNotFound, Message: "not found"},
			{Err: ErrConflict, Status: http.StatusConflict, Message: "conflict"},
		},
		http.StatusInternalServerError,
		"internal error",
	)
}
