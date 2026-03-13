package faq

import (
	"net/http"
	"rest_waka/pkg/httpx"
	"rest_waka/pkg/middleware"
	"rest_waka/pkg/req"
	"rest_waka/pkg/res"
	"strconv"
	"strings"
)

type HandlerDeps struct {
	Service   Service
	JWTSecret string
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

	// admin read
	router.Handle("GET /api/admin/faq/topics", middleware.RequireAdmin(handler.AdminListTopics(), deps.JWTSecret))
	router.Handle("GET /api/admin/faq/topics/{topicID}/articles", middleware.RequireAdmin(handler.AdminListArticlesByTopic(), deps.JWTSecret))
	router.Handle("GET /api/admin/faq/articles", middleware.RequireAdmin(handler.AdminListArticles(), deps.JWTSecret))
	router.Handle("GET /api/admin/faq/articles/{id}", middleware.RequireAdmin(handler.AdminGetArticle(), deps.JWTSecret))

	// admin write
	router.Handle("POST /api/admin/faq/topics", middleware.RequireAdmin(handler.CreateTopic(), deps.JWTSecret))
	router.Handle("PATCH /api/admin/faq/topics/{id}", middleware.RequireAdmin(handler.UpdateTopic(), deps.JWTSecret))
	router.Handle("DELETE /api/admin/faq/topics/{id}", middleware.RequireAdmin(handler.DeleteTopic(), deps.JWTSecret))

	router.Handle("POST /api/admin/faq/articles", middleware.RequireAdmin(handler.CreateArticle(), deps.JWTSecret))
	router.Handle("PATCH /api/admin/faq/articles/{id}", middleware.RequireAdmin(handler.UpdateArticle(), deps.JWTSecret))
	router.Handle("DELETE /api/admin/faq/articles/{id}", middleware.RequireAdmin(handler.DeleteArticle(), deps.JWTSecret))
	router.Handle("PUT /api/admin/faq/articles/{id}/blocks", middleware.RequireAdmin(handler.PutBlocks(), deps.JWTSecret))

	router.Handle("POST /api/admin/faq/articles/{id}/blocks", middleware.RequireAdmin(handler.CreateBlock(), deps.JWTSecret))
	router.Handle("PATCH /api/admin/faq/blocks/{id}", middleware.RequireAdmin(handler.UpdateBlock(), deps.JWTSecret))
	router.Handle("DELETE /api/admin/faq/blocks/{id}", middleware.RequireAdmin(handler.DeleteBlock(), deps.JWTSecret))
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

// admin read

func (handler *Handler) AdminListTopics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := handler.svc.ListTopicsAdmin(r.Context())
		if err != nil {
			writeFaqErr(w, err)
			return
		}
		res.Json(w, data, http.StatusOK)
	}
}

func (handler *Handler) AdminListArticlesByTopic() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		topicID, err := httpx.PathUint64(r, "topicID")
		if err != nil {
			res.Json(w, "invalid topicID", http.StatusBadRequest)
			return
		}

		filter := AdminArticleFilter{
			TopicID: &topicID,
			Channel: strings.TrimSpace(r.URL.Query().Get("channel")),
			Status:  strings.TrimSpace(r.URL.Query().Get("status")),
			Limit:   httpx.QueryInt(r, "limit", 20),
			Offset:  httpx.QueryInt(r, "offset", 0),
		}

		withBlocks := queryBool(r, "with_blocks")

		if withBlocks {
			data, err := handler.svc.ListArticlesAdminWithBlocks(r.Context(), filter)
			if err != nil {
				writeFaqErr(w, err)
				return
			}
			res.Json(w, data, http.StatusOK)
			return
		}

		data, err := handler.svc.ListArticlesAdmin(r.Context(), filter)
		if err != nil {
			writeFaqErr(w, err)
			return
		}
		res.Json(w, data, http.StatusOK)
	}
}

func (handler *Handler) AdminListArticles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var topicID *uint64
		if raw := strings.TrimSpace(r.URL.Query().Get("topic_id")); raw != "" {
			v, err := strconv.ParseUint(raw, 10, 64)
			if err != nil || v == 0 {
				res.Json(w, "invalid topic_id", http.StatusBadRequest)
				return
			}
			tv := v
			topicID = &tv
		}

		filter := AdminArticleFilter{
			TopicID: topicID,
			Channel: strings.TrimSpace(r.URL.Query().Get("channel")),
			Status:  strings.TrimSpace(r.URL.Query().Get("status")),
			Limit:   httpx.QueryInt(r, "limit", 20),
			Offset:  httpx.QueryInt(r, "offset", 0),
		}

		withBlocks := queryBool(r, "with_blocks")

		if withBlocks {
			data, err := handler.svc.ListArticlesAdminWithBlocks(r.Context(), filter)
			if err != nil {
				writeFaqErr(w, err)
				return
			}
			res.Json(w, data, http.StatusOK)
			return
		}

		data, err := handler.svc.ListArticlesAdmin(r.Context(), filter)
		if err != nil {
			writeFaqErr(w, err)
			return
		}
		res.Json(w, data, http.StatusOK)
	}
}

func (handler *Handler) AdminGetArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := httpx.PathUint64(r, "id")
		if err != nil {
			res.Json(w, "invalid id", http.StatusBadRequest)
			return
		}

		data, err := handler.svc.GetArticleAdmin(r.Context(), id)
		if err != nil {
			writeFaqErr(w, err)
			return
		}
		res.Json(w, data, http.StatusOK)
	}
}

// admin write

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

func (handler *Handler) DeleteTopic() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := httpx.PathUint64(r, "id")
		if err != nil {
			res.Json(w, "invalid id", http.StatusBadRequest)
			return
		}

		if err := handler.svc.DeleteTopic(r.Context(), id); err != nil {
			writeFaqErr(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
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

func (handler *Handler) DeleteArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := httpx.PathUint64(r, "id")
		if err != nil {
			res.Json(w, "invalid id", http.StatusBadRequest)
			return
		}

		if err := handler.svc.DeleteArticle(r.Context(), id); err != nil {
			writeFaqErr(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
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

func (handler *Handler) CreateBlock() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		articleID, err := httpx.PathUint64(r, "id")
		if err != nil {
			res.Json(w, "invalid id", http.StatusBadRequest)
			return
		}

		payload, err := req.Decode[CreateBlockRequest](r.Body)
		if err != nil {
			res.Json(w, "invalid body", http.StatusBadRequest)
			return
		}

		data, err := handler.svc.CreateBlock(r.Context(), articleID, payload)
		if err != nil {
			writeFaqErr(w, err)
			return
		}
		res.Json(w, data, http.StatusCreated)
	}
}

func (handler *Handler) UpdateBlock() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		blockID, err := httpx.PathUint64(r, "id")
		if err != nil {
			res.Json(w, "invalid id", http.StatusBadRequest)
			return
		}

		payload, err := req.Decode[UpdateBlockRequest](r.Body)
		if err != nil {
			res.Json(w, "invalid body", http.StatusBadRequest)
			return
		}

		data, err := handler.svc.UpdateBlock(r.Context(), blockID, payload)
		if err != nil {
			writeFaqErr(w, err)
			return
		}
		res.Json(w, data, http.StatusOK)
	}
}

func (handler *Handler) DeleteBlock() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		blockID, err := httpx.PathUint64(r, "id")
		if err != nil {
			res.Json(w, "invalid id", http.StatusBadRequest)
			return
		}

		if err := handler.svc.DeleteBlock(r.Context(), blockID); err != nil {
			writeFaqErr(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func queryBool(r *http.Request, key string) bool {
	v := strings.ToLower(strings.TrimSpace(r.URL.Query().Get(key)))
	switch v {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
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
