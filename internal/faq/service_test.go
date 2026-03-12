package faq

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"gorm.io/datatypes"
)

type fakeRepository struct {
	listTopicsFn              func(ctx context.Context, activeOnly bool) ([]Topic, error)
	getTopicFn                func(ctx context.Context, id uint64) (Topic, error)
	createTopicFn             func(ctx context.Context, t *Topic) error
	updateTopicFn             func(ctx context.Context, id uint64, patch UpdateTopicRequest) (Topic, error)
	listArticlesByTopicFn     func(ctx context.Context, topicID uint64, channel string) ([]ArticleSummary, error)
	getArticleFn              func(ctx context.Context, id uint64) (Article, []Block, error)
	searchArticlesFn          func(ctx context.Context, q string, channel string, limit, offset int) ([]ArticleSummary, error)
	listArticlesAdminFn       func(ctx context.Context, filter AdminArticleFilter) ([]Article, int64, error)
	getArticleAnyStatusFn     func(ctx context.Context, id uint64) (Article, []Block, error)
	createArticleFn           func(ctx context.Context, a *Article) error
	updateArticleFn           func(ctx context.Context, id uint64, patch UpdateArticleRequest) (Article, error)
	replaceBlocksFn           func(ctx context.Context, articleID uint64, blocks []Block) ([]Block, error)
	createBlockFn             func(ctx context.Context, b *Block) error
	updateBlockFn             func(ctx context.Context, id uint64, patch UpdateBlockRequest) (Block, error)
	deleteBlockFn             func(ctx context.Context, id uint64) (uint64, error)
	updateArticleSearchTextFn func(ctx context.Context, articleID uint64, searchText *string) error

	listTopicsCalls              []bool
	listArticlesByTopicCalls     []listArticlesByTopicCall
	getArticleCalls              []uint64
	searchArticlesCalls          []searchArticlesCall
	listArticlesAdminCalls       []AdminArticleFilter
	getArticleAnyStatusCalls     []uint64
	createTopicCalls             []Topic
	updateTopicCalls             []updateTopicCall
	createArticleCalls           []Article
	updateArticleCalls           []updateArticleCall
	replaceBlocksCalls           []replaceBlocksCall
	createBlockCalls             []Block
	updateBlockCalls             []updateBlockCall
	deleteBlockCalls             []uint64
	updateArticleSearchTextCalls []updateArticleSearchTextCall
}

type listArticlesByTopicCall struct {
	topicID uint64
	channel string
}

type searchArticlesCall struct {
	q       string
	channel string
	limit   int
	offset  int
}

type updateTopicCall struct {
	id    uint64
	patch UpdateTopicRequest
}

type updateArticleCall struct {
	id    uint64
	patch UpdateArticleRequest
}

type replaceBlocksCall struct {
	articleID uint64
	blocks    []Block
}

type updateBlockCall struct {
	id    uint64
	patch UpdateBlockRequest
}

type updateArticleSearchTextCall struct {
	articleID  uint64
	searchText *string
}

func (f *fakeRepository) ListTopics(ctx context.Context, activeOnly bool) ([]Topic, error) {
	f.listTopicsCalls = append(f.listTopicsCalls, activeOnly)
	if f.listTopicsFn != nil {
		return f.listTopicsFn(ctx, activeOnly)
	}
	return nil, nil
}

func (f *fakeRepository) GetTopic(ctx context.Context, id uint64) (Topic, error) {
	if f.getTopicFn != nil {
		return f.getTopicFn(ctx, id)
	}
	return Topic{}, nil
}

func (f *fakeRepository) CreateTopic(ctx context.Context, t *Topic) error {
	if t != nil {
		f.createTopicCalls = append(f.createTopicCalls, *t)
	}
	if f.createTopicFn != nil {
		return f.createTopicFn(ctx, t)
	}
	return nil
}

func (f *fakeRepository) UpdateTopic(ctx context.Context, id uint64, patch UpdateTopicRequest) (Topic, error) {
	f.updateTopicCalls = append(f.updateTopicCalls, updateTopicCall{id: id, patch: patch})
	if f.updateTopicFn != nil {
		return f.updateTopicFn(ctx, id, patch)
	}
	return Topic{}, nil
}

func (f *fakeRepository) ListArticlesByTopic(ctx context.Context, topicID uint64, channel string) ([]ArticleSummary, error) {
	f.listArticlesByTopicCalls = append(f.listArticlesByTopicCalls, listArticlesByTopicCall{topicID: topicID, channel: channel})
	if f.listArticlesByTopicFn != nil {
		return f.listArticlesByTopicFn(ctx, topicID, channel)
	}
	return nil, nil
}

func (f *fakeRepository) GetArticle(ctx context.Context, id uint64) (Article, []Block, error) {
	f.getArticleCalls = append(f.getArticleCalls, id)
	if f.getArticleFn != nil {
		return f.getArticleFn(ctx, id)
	}
	return Article{}, nil, nil
}

func (f *fakeRepository) SearchArticles(ctx context.Context, q string, channel string, limit, offset int) ([]ArticleSummary, error) {
	f.searchArticlesCalls = append(f.searchArticlesCalls, searchArticlesCall{
		q:       q,
		channel: channel,
		limit:   limit,
		offset:  offset,
	})
	if f.searchArticlesFn != nil {
		return f.searchArticlesFn(ctx, q, channel, limit, offset)
	}
	return nil, nil
}

func (f *fakeRepository) ListArticlesAdmin(ctx context.Context, filter AdminArticleFilter) ([]Article, int64, error) {
	f.listArticlesAdminCalls = append(f.listArticlesAdminCalls, filter)
	if f.listArticlesAdminFn != nil {
		return f.listArticlesAdminFn(ctx, filter)
	}
	return nil, 0, nil
}

func (f *fakeRepository) GetArticleAnyStatus(ctx context.Context, id uint64) (Article, []Block, error) {
	f.getArticleAnyStatusCalls = append(f.getArticleAnyStatusCalls, id)
	if f.getArticleAnyStatusFn != nil {
		return f.getArticleAnyStatusFn(ctx, id)
	}
	return Article{}, nil, nil
}

func (f *fakeRepository) CreateArticle(ctx context.Context, a *Article) error {
	if a != nil {
		f.createArticleCalls = append(f.createArticleCalls, *a)
	}
	if f.createArticleFn != nil {
		return f.createArticleFn(ctx, a)
	}
	return nil
}

func (f *fakeRepository) UpdateArticle(ctx context.Context, id uint64, patch UpdateArticleRequest) (Article, error) {
	f.updateArticleCalls = append(f.updateArticleCalls, updateArticleCall{id: id, patch: patch})
	if f.updateArticleFn != nil {
		return f.updateArticleFn(ctx, id, patch)
	}
	return Article{}, nil
}

func (f *fakeRepository) ReplaceBlocks(ctx context.Context, articleID uint64, blocks []Block) ([]Block, error) {
	f.replaceBlocksCalls = append(f.replaceBlocksCalls, replaceBlocksCall{
		articleID: articleID,
		blocks:    cloneBlocks(blocks),
	})
	if f.replaceBlocksFn != nil {
		return f.replaceBlocksFn(ctx, articleID, blocks)
	}
	return blocks, nil
}

func (f *fakeRepository) CreateBlock(ctx context.Context, b *Block) error {
	if b != nil {
		f.createBlockCalls = append(f.createBlockCalls, cloneBlock(*b))
	}
	if f.createBlockFn != nil {
		return f.createBlockFn(ctx, b)
	}
	return nil
}

func (f *fakeRepository) UpdateBlock(ctx context.Context, id uint64, patch UpdateBlockRequest) (Block, error) {
	f.updateBlockCalls = append(f.updateBlockCalls, updateBlockCall{id: id, patch: patch})
	if f.updateBlockFn != nil {
		return f.updateBlockFn(ctx, id, patch)
	}
	return Block{}, nil
}

func (f *fakeRepository) DeleteBlock(ctx context.Context, id uint64) (uint64, error) {
	f.deleteBlockCalls = append(f.deleteBlockCalls, id)
	if f.deleteBlockFn != nil {
		return f.deleteBlockFn(ctx, id)
	}
	return 0, nil
}

func (f *fakeRepository) UpdateArticleSearchText(ctx context.Context, articleID uint64, searchText *string) error {
	var copied *string
	if searchText != nil {
		v := *searchText
		copied = &v
	}
	f.updateArticleSearchTextCalls = append(f.updateArticleSearchTextCalls, updateArticleSearchTextCall{
		articleID:  articleID,
		searchText: copied,
	})
	if f.updateArticleSearchTextFn != nil {
		return f.updateArticleSearchTextFn(ctx, articleID, searchText)
	}
	return nil
}

func cloneBlocks(in []Block) []Block {
	out := make([]Block, len(in))
	for i := range in {
		out[i] = cloneBlock(in[i])
	}
	return out
}

func cloneBlock(in Block) Block {
	out := in
	if in.Data != nil {
		out.Data = append(datatypes.JSON(nil), in.Data...)
	}
	return out
}

func stringPtr(s string) *string { return &s }
func intPtr(v int) *int          { return &v }
func boolPtr(v bool) *bool       { return &v }
func uint64Ptr(v uint64) *uint64 { return &v }

func mustJSON(s string) datatypes.JSON {
	return datatypes.JSON([]byte(s))
}

func TestService_ListTopics(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	now := time.Now

	t.Run("active only", func(t *testing.T) {
		t.Parallel()

		repo := &fakeRepository{
			listTopicsFn: func(ctx context.Context, activeOnly bool) ([]Topic, error) {
				return []Topic{{ID: 1, Title: "A", IsActive: true}}, nil
			},
		}
		svc := &service{repo: repo, now: now}

		got, err := svc.ListTopics(ctx)
		if err != nil {
			t.Fatalf("ListTopics() error = %v", err)
		}
		if len(repo.listTopicsCalls) != 1 || !repo.listTopicsCalls[0] {
			t.Fatalf("ListTopics() should call repo with activeOnly=true, calls=%v", repo.listTopicsCalls)
		}
		want := []Topic{{ID: 1, Title: "A", IsActive: true}}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("ListTopics() got=%v want=%v", got, want)
		}
	})

	t.Run("repo error", func(t *testing.T) {
		t.Parallel()

		repo := &fakeRepository{
			listTopicsFn: func(ctx context.Context, activeOnly bool) ([]Topic, error) {
				return nil, ErrNotFound
			},
		}
		svc := &service{repo: repo, now: now}

		_, err := svc.ListTopics(ctx)
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("ListTopics() error = %v, want %v", err, ErrNotFound)
		}
	})
}

func TestService_ListArticlesByTopic(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name        string
		topicID     uint64
		channel     string
		repoFn      func(ctx context.Context, topicID uint64, channel string) ([]ArticleSummary, error)
		wantErr     error
		wantCalled  bool
		wantChannel string
	}{
		{
			name:    "invalid channel",
			topicID: 10,
			channel: "email",
			wantErr: ErrInvalidArgument,
		},
		{
			name:    "zero topic id",
			topicID: 0,
			channel: ChannelTG,
			wantErr: ErrInvalidArgument,
		},
		{
			name:        "success with normalization",
			topicID:     42,
			channel:     " TG ",
			wantCalled:  true,
			wantChannel: ChannelTG,
			repoFn: func(ctx context.Context, topicID uint64, channel string) ([]ArticleSummary, error) {
				return []ArticleSummary{{ID: 1, TopicID: topicID, Title: "X"}}, nil
			},
		},
		{
			name:       "repo error",
			topicID:    42,
			channel:    ChannelAll,
			wantErr:    ErrNotFound,
			wantCalled: true,
			repoFn: func(ctx context.Context, topicID uint64, channel string) ([]ArticleSummary, error) {
				return nil, ErrNotFound
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &fakeRepository{listArticlesByTopicFn: tt.repoFn}
			svc := &service{repo: repo, now: time.Now}

			got, err := svc.ListArticlesByTopic(ctx, tt.topicID, tt.channel)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("ListArticlesByTopic() err=%v wantErr=%v", err, tt.wantErr)
			}
			if !tt.wantCalled {
				if len(repo.listArticlesByTopicCalls) != 0 {
					t.Fatalf("repo should not be called, got calls=%v", repo.listArticlesByTopicCalls)
				}
				return
			}

			if len(repo.listArticlesByTopicCalls) != 1 {
				t.Fatalf("expected one repo call, got=%d", len(repo.listArticlesByTopicCalls))
			}
			call := repo.listArticlesByTopicCalls[0]
			if call.topicID != tt.topicID {
				t.Fatalf("topicID passed=%d want=%d", call.topicID, tt.topicID)
			}
			if tt.wantChannel != "" && call.channel != tt.wantChannel {
				t.Fatalf("channel passed=%q want=%q", call.channel, tt.wantChannel)
			}
			if tt.wantErr == nil && len(got) == 0 {
				t.Fatalf("expected non-empty result on success")
			}
		})
	}
}

func TestService_GetArticle(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("invalid id", func(t *testing.T) {
		t.Parallel()

		repo := &fakeRepository{}
		svc := &service{repo: repo, now: time.Now}

		_, err := svc.GetArticle(ctx, 0)
		if !errors.Is(err, ErrInvalidArgument) {
			t.Fatalf("GetArticle() err=%v want=%v", err, ErrInvalidArgument)
		}
		if len(repo.getArticleCalls) != 0 {
			t.Fatalf("repo should not be called on invalid id")
		}
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		wantArticle := Article{ID: 7, Title: "A"}
		wantBlocks := []Block{{ID: 1, ArticleID: 7, Type: BlockText, Data: mustJSON(`{"text":"hello"}`)}}

		repo := &fakeRepository{
			getArticleFn: func(ctx context.Context, id uint64) (Article, []Block, error) {
				return wantArticle, wantBlocks, nil
			},
		}
		svc := &service{repo: repo, now: time.Now}

		got, err := svc.GetArticle(ctx, 7)
		if err != nil {
			t.Fatalf("GetArticle() err=%v", err)
		}
		if got.Article.ID != wantArticle.ID || len(got.Blocks) != len(wantBlocks) {
			t.Fatalf("GetArticle() got=%+v want article=%+v blocks=%+v", got, wantArticle, wantBlocks)
		}
	})

	t.Run("repo not found", func(t *testing.T) {
		t.Parallel()

		repo := &fakeRepository{
			getArticleFn: func(ctx context.Context, id uint64) (Article, []Block, error) {
				return Article{}, nil, ErrNotFound
			},
		}
		svc := &service{repo: repo, now: time.Now}

		_, err := svc.GetArticle(ctx, 1)
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("GetArticle() err=%v want=%v", err, ErrNotFound)
		}
	})
}

func TestService_Search(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("invalid channel", func(t *testing.T) {
		t.Parallel()

		repo := &fakeRepository{}
		svc := &service{repo: repo, now: time.Now}

		_, err := svc.Search(ctx, "q", "email", 10, 0)
		if !errors.Is(err, ErrInvalidArgument) {
			t.Fatalf("Search() err=%v want=%v", err, ErrInvalidArgument)
		}
		if len(repo.searchArticlesCalls) != 0 {
			t.Fatalf("repo should not be called on invalid channel")
		}
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		repo := &fakeRepository{
			searchArticlesFn: func(ctx context.Context, q string, channel string, limit, offset int) ([]ArticleSummary, error) {
				return []ArticleSummary{{ID: 10, Title: "Found"}}, nil
			},
		}
		svc := &service{repo: repo, now: time.Now}

		got, err := svc.Search(ctx, " hello ", " TG ", 25, 2)
		if err != nil {
			t.Fatalf("Search() err=%v", err)
		}
		if len(got) != 1 || got[0].ID != 10 {
			t.Fatalf("Search() unexpected result: %+v", got)
		}
		if len(repo.searchArticlesCalls) != 1 {
			t.Fatalf("expected one repo call")
		}
		call := repo.searchArticlesCalls[0]
		if call.channel != ChannelTG || call.limit != 25 || call.offset != 2 || call.q != " hello " {
			t.Fatalf("repo call mismatch: %+v", call)
		}
	})
}

func TestService_ListTopicsAdmin(t *testing.T) {
	t.Parallel()

	repo := &fakeRepository{
		listTopicsFn: func(ctx context.Context, activeOnly bool) ([]Topic, error) {
			return []Topic{{ID: 1, IsActive: false}}, nil
		},
	}
	svc := &service{repo: repo, now: time.Now}

	got, err := svc.ListTopicsAdmin(context.Background())
	if err != nil {
		t.Fatalf("ListTopicsAdmin() err=%v", err)
	}
	if len(got) != 1 || got[0].ID != 1 {
		t.Fatalf("ListTopicsAdmin() unexpected result: %+v", got)
	}
	if len(repo.listTopicsCalls) != 1 || repo.listTopicsCalls[0] {
		t.Fatalf("ListTopicsAdmin() must call repo with activeOnly=false, calls=%v", repo.listTopicsCalls)
	}
}

func TestNormalizeAdminFilter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      AdminArticleFilter
		want    AdminArticleFilter
		wantErr error
	}{
		{
			name: "default pagination",
			in: AdminArticleFilter{
				Limit:  0,
				Offset: -10,
			},
			want: AdminArticleFilter{
				Limit:  20,
				Offset: 0,
			},
		},
		{
			name: "normalize channel and status",
			in: AdminArticleFilter{
				Channel: " TG ",
				Status:  " PUBLISHED ",
				Limit:   50,
				Offset:  5,
			},
			want: AdminArticleFilter{
				Channel: ChannelTG,
				Status:  StatusPublished,
				Limit:   50,
				Offset:  5,
			},
		},
		{
			name: "invalid channel",
			in: AdminArticleFilter{
				Channel: "email",
			},
			wantErr: ErrInvalidArgument,
		},
		{
			name: "invalid status",
			in: AdminArticleFilter{
				Status: "new",
			},
			wantErr: ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := normalizeAdminFilter(tt.in)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("normalizeAdminFilter() err=%v wantErr=%v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				return
			}
			if got.Channel != tt.want.Channel || got.Status != tt.want.Status || got.Limit != tt.want.Limit || got.Offset != tt.want.Offset {
				t.Fatalf("normalizeAdminFilter() got=%+v want=%+v", got, tt.want)
			}
		})
	}
}

func TestService_ListArticlesAdmin(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("invalid filter", func(t *testing.T) {
		t.Parallel()

		repo := &fakeRepository{}
		svc := &service{repo: repo, now: time.Now}

		_, err := svc.ListArticlesAdmin(ctx, AdminArticleFilter{Channel: "smtp"})
		if !errors.Is(err, ErrInvalidArgument) {
			t.Fatalf("ListArticlesAdmin() err=%v want=%v", err, ErrInvalidArgument)
		}
		if len(repo.listArticlesAdminCalls) != 0 {
			t.Fatalf("repo should not be called on invalid filter")
		}
	})

	t.Run("normalize filter and map response", func(t *testing.T) {
		t.Parallel()

		publishedAt := time.Date(2025, 2, 1, 10, 0, 0, 0, time.UTC)
		repo := &fakeRepository{
			listArticlesAdminFn: func(ctx context.Context, filter AdminArticleFilter) ([]Article, int64, error) {
				return []Article{
					{
						ID:          11,
						TopicID:     3,
						Slug:        "faq-slug",
						Title:       "FAQ",
						Status:      StatusPublished,
						Channel:     ChannelTG,
						PublishedAt: &publishedAt,
					},
				}, 77, nil
			},
		}
		svc := &service{repo: repo, now: time.Now}

		resp, err := svc.ListArticlesAdmin(ctx, AdminArticleFilter{
			Channel: " TG ",
			Status:  " PUBLISHED ",
			Limit:   0,
			Offset:  -1,
		})
		if err != nil {
			t.Fatalf("ListArticlesAdmin() err=%v", err)
		}
		if len(repo.listArticlesAdminCalls) != 1 {
			t.Fatalf("expected one repo call")
		}
		call := repo.listArticlesAdminCalls[0]
		if call.Channel != ChannelTG || call.Status != StatusPublished || call.Limit != 20 || call.Offset != 0 {
			t.Fatalf("filter passed to repo mismatch: %+v", call)
		}
		if resp.Limit != 20 || resp.Offset != 0 || resp.Total != 77 {
			t.Fatalf("response pagination mismatch: %+v", resp)
		}
		if len(resp.Items) != 1 {
			t.Fatalf("expected one item, got=%d", len(resp.Items))
		}
		item := resp.Items[0]
		if item.ID != 11 || item.TopicID != 3 || item.Slug != "faq-slug" || item.Status != StatusPublished || item.Channel != ChannelTG {
			t.Fatalf("mapped item mismatch: %+v", item)
		}
	})
}

func TestService_ListArticlesAdminWithBlocks(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("invalid filter", func(t *testing.T) {
		t.Parallel()

		repo := &fakeRepository{}
		svc := &service{repo: repo, now: time.Now}

		_, err := svc.ListArticlesAdminWithBlocks(ctx, AdminArticleFilter{Status: "new"})
		if !errors.Is(err, ErrInvalidArgument) {
			t.Fatalf("ListArticlesAdminWithBlocks() err=%v want=%v", err, ErrInvalidArgument)
		}
	})

	t.Run("fetch details through get article any status", func(t *testing.T) {
		t.Parallel()

		repo := &fakeRepository{
			listArticlesAdminFn: func(ctx context.Context, filter AdminArticleFilter) ([]Article, int64, error) {
				return []Article{
					{ID: 10, Title: "A"},
					{ID: 20, Title: "B"},
				}, 2, nil
			},
			getArticleAnyStatusFn: func(ctx context.Context, id uint64) (Article, []Block, error) {
				switch id {
				case 10:
					return Article{ID: 10, Title: "A-full"}, []Block{{ID: 1, ArticleID: 10, Type: BlockText, Data: mustJSON(`{"text":"a"}`)}}, nil
				case 20:
					return Article{ID: 20, Title: "B-full"}, []Block{{ID: 2, ArticleID: 20, Type: BlockText, Data: mustJSON(`{"text":"b"}`)}}, nil
				default:
					return Article{}, nil, ErrNotFound
				}
			},
		}
		svc := &service{repo: repo, now: time.Now}

		resp, err := svc.ListArticlesAdminWithBlocks(ctx, AdminArticleFilter{Limit: 10, Offset: 0})
		if err != nil {
			t.Fatalf("ListArticlesAdminWithBlocks() err=%v", err)
		}
		if len(repo.listArticlesAdminCalls) != 1 {
			t.Fatalf("expected one ListArticlesAdmin call")
		}
		wantIDs := []uint64{10, 20}
		if !reflect.DeepEqual(repo.getArticleAnyStatusCalls, wantIDs) {
			t.Fatalf("GetArticleAnyStatus call ids=%v want=%v", repo.getArticleAnyStatusCalls, wantIDs)
		}
		if resp.Total != 2 || resp.Limit != 10 || resp.Offset != 0 {
			t.Fatalf("response metadata mismatch: %+v", resp)
		}
		if len(resp.Items) != 2 {
			t.Fatalf("expected 2 items, got=%d", len(resp.Items))
		}
		if resp.Items[0].Article.Title != "A-full" || len(resp.Items[0].Blocks) != 1 {
			t.Fatalf("first item mismatch: %+v", resp.Items[0])
		}
	})

	t.Run("detail fetch error is returned", func(t *testing.T) {
		t.Parallel()

		repo := &fakeRepository{
			listArticlesAdminFn: func(ctx context.Context, filter AdminArticleFilter) ([]Article, int64, error) {
				return []Article{{ID: 1}}, 1, nil
			},
			getArticleAnyStatusFn: func(ctx context.Context, id uint64) (Article, []Block, error) {
				return Article{}, nil, ErrNotFound
			},
		}
		svc := &service{repo: repo, now: time.Now}

		_, err := svc.ListArticlesAdminWithBlocks(ctx, AdminArticleFilter{})
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("ListArticlesAdminWithBlocks() err=%v want=%v", err, ErrNotFound)
		}
	})
}

func TestService_GetArticleAdmin(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("invalid id", func(t *testing.T) {
		t.Parallel()

		svc := &service{repo: &fakeRepository{}, now: time.Now}
		_, err := svc.GetArticleAdmin(ctx, 0)
		if !errors.Is(err, ErrInvalidArgument) {
			t.Fatalf("GetArticleAdmin() err=%v want=%v", err, ErrInvalidArgument)
		}
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		repo := &fakeRepository{
			getArticleAnyStatusFn: func(ctx context.Context, id uint64) (Article, []Block, error) {
				return Article{ID: id, Title: "admin"}, []Block{{ID: 1, ArticleID: id, Type: BlockText}}, nil
			},
		}
		svc := &service{repo: repo, now: time.Now}

		got, err := svc.GetArticleAdmin(ctx, 5)
		if err != nil {
			t.Fatalf("GetArticleAdmin() err=%v", err)
		}
		if got.Article.ID != 5 || len(got.Blocks) != 1 {
			t.Fatalf("GetArticleAdmin() unexpected result: %+v", got)
		}
	})
}

func TestService_CreateTopic(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("invalid empty title", func(t *testing.T) {
		t.Parallel()

		svc := &service{repo: &fakeRepository{}, now: time.Now}
		_, err := svc.CreateTopic(ctx, CreateTopicRequest{Title: "   "})
		if !errors.Is(err, ErrInvalidArgument) {
			t.Fatalf("CreateTopic() err=%v want=%v", err, ErrInvalidArgument)
		}
	})

	t.Run("defaults and trim", func(t *testing.T) {
		t.Parallel()

		repo := &fakeRepository{
			createTopicFn: func(ctx context.Context, t *Topic) error {
				t.ID = 100
				return nil
			},
		}
		svc := &service{repo: repo, now: time.Now}

		got, err := svc.CreateTopic(ctx, CreateTopicRequest{Title: "  Basics  "})
		if err != nil {
			t.Fatalf("CreateTopic() err=%v", err)
		}
		if len(repo.createTopicCalls) != 1 {
			t.Fatalf("expected one create topic call")
		}
		sent := repo.createTopicCalls[0]
		if sent.Title != "Basics" || sent.Sort != 0 || !sent.IsActive {
			t.Fatalf("topic sent to repo mismatch: %+v", sent)
		}
		if got.ID != 100 || got.Title != "Basics" {
			t.Fatalf("CreateTopic() result mismatch: %+v", got)
		}
	})

	t.Run("custom sort and is_active", func(t *testing.T) {
		t.Parallel()

		repo := &fakeRepository{}
		svc := &service{repo: repo, now: time.Now}
		sort := 5
		active := false

		got, err := svc.CreateTopic(ctx, CreateTopicRequest{
			Title:    "Topic",
			Sort:     &sort,
			IsActive: &active,
		})
		if err != nil {
			t.Fatalf("CreateTopic() err=%v", err)
		}
		if got.Sort != 5 || got.IsActive {
			t.Fatalf("CreateTopic() result mismatch: %+v", got)
		}
	})
}

func TestService_UpdateTopic(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		id        uint64
		req       UpdateTopicRequest
		repoFn    func(ctx context.Context, id uint64, patch UpdateTopicRequest) (Topic, error)
		wantErr   error
		wantCalls int
	}{
		{
			name:    "zero id",
			id:      0,
			req:     UpdateTopicRequest{},
			wantErr: ErrInvalidArgument,
		},
		{
			name:    "empty title",
			id:      5,
			req:     UpdateTopicRequest{Title: stringPtr("   ")},
			wantErr: ErrInvalidArgument,
		},
		{
			name:      "repo error not found",
			id:        7,
			req:       UpdateTopicRequest{Sort: intPtr(2)},
			wantErr:   ErrNotFound,
			wantCalls: 1,
			repoFn: func(ctx context.Context, id uint64, patch UpdateTopicRequest) (Topic, error) {
				return Topic{}, ErrNotFound
			},
		},
		{
			name:      "success",
			id:        8,
			req:       UpdateTopicRequest{Sort: intPtr(3)},
			wantCalls: 1,
			repoFn: func(ctx context.Context, id uint64, patch UpdateTopicRequest) (Topic, error) {
				return Topic{ID: id, Sort: *patch.Sort}, nil
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &fakeRepository{updateTopicFn: tt.repoFn}
			svc := &service{repo: repo, now: time.Now}

			got, err := svc.UpdateTopic(ctx, tt.id, tt.req)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("UpdateTopic() err=%v wantErr=%v", err, tt.wantErr)
			}
			if len(repo.updateTopicCalls) != tt.wantCalls {
				t.Fatalf("UpdateTopic() repo calls=%d want=%d", len(repo.updateTopicCalls), tt.wantCalls)
			}
			if tt.wantErr == nil && got.ID != tt.id {
				t.Fatalf("UpdateTopic() result mismatch: %+v", got)
			}
		})
	}
}

func TestService_CreateArticle(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	fixedNow := time.Date(2026, 3, 13, 10, 30, 0, 0, time.UTC)

	t.Run("validation errors", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			req  CreateArticleRequest
		}{
			{
				name: "zero topic id",
				req:  CreateArticleRequest{TopicID: 0, Title: "X"},
			},
			{
				name: "empty title",
				req:  CreateArticleRequest{TopicID: 1, Title: "   "},
			},
			{
				name: "invalid channel",
				req:  CreateArticleRequest{TopicID: 1, Title: "X", Channel: stringPtr("email")},
			},
			{
				name: "invalid status",
				req:  CreateArticleRequest{TopicID: 1, Title: "X", Status: stringPtr("new")},
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				repo := &fakeRepository{}
				svc := &service{repo: repo, now: func() time.Time { return fixedNow }}

				_, err := svc.CreateArticle(ctx, tt.req)
				if !errors.Is(err, ErrInvalidArgument) {
					t.Fatalf("CreateArticle() err=%v want=%v", err, ErrInvalidArgument)
				}
				if len(repo.createArticleCalls) != 0 {
					t.Fatalf("repo should not be called on invalid input")
				}
			})
		}
	})

	t.Run("defaults and slug generation", func(t *testing.T) {
		t.Parallel()

		repo := &fakeRepository{
			createArticleFn: func(ctx context.Context, a *Article) error {
				a.ID = 55
				return nil
			},
		}
		svc := &service{repo: repo, now: func() time.Time { return fixedNow }}

		got, err := svc.CreateArticle(ctx, CreateArticleRequest{
			TopicID: 10,
			Title:   "  Hello World  ",
		})
		if err != nil {
			t.Fatalf("CreateArticle() err=%v", err)
		}
		if len(repo.createArticleCalls) != 1 {
			t.Fatalf("expected one repo call")
		}
		sent := repo.createArticleCalls[0]
		if sent.TopicID != 10 || sent.Title != "Hello World" || sent.Channel != ChannelAll || sent.Status != StatusDraft || sent.Slug != "hello-world" || sent.PublishedAt != nil {
			t.Fatalf("article sent to repo mismatch: %+v", sent)
		}
		if got.ID != 55 {
			t.Fatalf("expected repo-mutated id in result, got=%d", got.ID)
		}
	})

	t.Run("published sets published_at and trims custom slug", func(t *testing.T) {
		t.Parallel()

		repo := &fakeRepository{}
		svc := &service{repo: repo, now: func() time.Time { return fixedNow }}

		got, err := svc.CreateArticle(ctx, CreateArticleRequest{
			TopicID:  2,
			Title:    "Title",
			Slug:     stringPtr("  Custom-Slug  "),
			Status:   stringPtr(" PUBLISHED "),
			Channel:  stringPtr(" TG "),
		})
		if err != nil {
			t.Fatalf("CreateArticle() err=%v", err)
		}
		sent := repo.createArticleCalls[0]
		if sent.Slug != "Custom-Slug" || sent.Channel != ChannelTG || sent.Status != StatusPublished {
			t.Fatalf("article normalization mismatch: %+v", sent)
		}
		if sent.PublishedAt == nil || !sent.PublishedAt.Equal(fixedNow) {
			t.Fatalf("expected published_at=%v, got=%v", fixedNow, sent.PublishedAt)
		}
		if got.PublishedAt == nil || !got.PublishedAt.Equal(fixedNow) {
			t.Fatalf("result PublishedAt mismatch: %+v", got)
		}
	})

	t.Run("empty provided slug falls back to generated", func(t *testing.T) {
		t.Parallel()

		repo := &fakeRepository{}
		svc := &service{repo: repo, now: func() time.Time { return fixedNow }}

		_, err := svc.CreateArticle(ctx, CreateArticleRequest{
			TopicID: 1,
			Title:   "  Go FAQ  ",
			Slug:    stringPtr("    "),
		})
		if err != nil {
			t.Fatalf("CreateArticle() err=%v", err)
		}
		if repo.createArticleCalls[0].Slug != "go-faq" {
			t.Fatalf("expected generated slug, got=%q", repo.createArticleCalls[0].Slug)
		}
	})

	t.Run("repo conflict", func(t *testing.T) {
		t.Parallel()

		repo := &fakeRepository{
			createArticleFn: func(ctx context.Context, a *Article) error {
				return ErrConflict
			},
		}
		svc := &service{repo: repo, now: func() time.Time { return fixedNow }}

		_, err := svc.CreateArticle(ctx, CreateArticleRequest{TopicID: 1, Title: "x"})
		if !errors.Is(err, ErrConflict) {
			t.Fatalf("CreateArticle() err=%v want=%v", err, ErrConflict)
		}
	})
}

func TestService_UpdateArticle(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("validation errors", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			id   uint64
			req  UpdateArticleRequest
		}{
			{
				name: "zero id",
				id:   0,
				req:  UpdateArticleRequest{Title: stringPtr("x")},
			},
			{
				name: "empty title",
				id:   1,
				req:  UpdateArticleRequest{Title: stringPtr("   ")},
			},
			{
				name: "empty slug",
				id:   1,
				req:  UpdateArticleRequest{Slug: stringPtr("   ")},
			},
			{
				name: "invalid channel",
				id:   1,
				req:  UpdateArticleRequest{Channel: stringPtr("email")},
			},
			{
				name: "invalid status",
				id:   1,
				req:  UpdateArticleRequest{Status: stringPtr("new")},
			},
			{
				name: "zero topic id in patch",
				id:   1,
				req:  UpdateArticleRequest{TopicID: uint64Ptr(0)},
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				repo := &fakeRepository{}
				svc := &service{repo: repo, now: time.Now}

				_, err := svc.UpdateArticle(ctx, tt.id, tt.req)
				if !errors.Is(err, ErrInvalidArgument) {
					t.Fatalf("UpdateArticle() err=%v want=%v", err, ErrInvalidArgument)
				}
				if len(repo.updateArticleCalls) != 0 {
					t.Fatalf("repo should not be called on invalid input")
				}
			})
		}
	})

	t.Run("normalize patch and pass to repo", func(t *testing.T) {
		t.Parallel()

		repo := &fakeRepository{
			updateArticleFn: func(ctx context.Context, id uint64, patch UpdateArticleRequest) (Article, error) {
				return Article{
					ID:      id,
					TopicID: *patch.TopicID,
					Title:   *patch.Title,
					Slug:    *patch.Slug,
					Channel: *patch.Channel,
					Status:  *patch.Status,
				}, nil
			},
		}
		svc := &service{repo: repo, now: time.Now}

		got, err := svc.UpdateArticle(ctx, 11, UpdateArticleRequest{
			TopicID: uint64Ptr(3),
			Title:   stringPtr("  New title  "),
			Slug:    stringPtr("  New-Slug  "),
			Channel: stringPtr(" TG "),
			Status:  stringPtr(" PUBLISHED "),
		})
		if err != nil {
			t.Fatalf("UpdateArticle() err=%v", err)
		}
		if len(repo.updateArticleCalls) != 1 {
			t.Fatalf("expected one repo call")
		}
		patch := repo.updateArticleCalls[0].patch
		if patch.TopicID == nil || *patch.TopicID != 3 {
			t.Fatalf("expected topicID patch=3, got=%v", patch.TopicID)
		}
		if patch.Title == nil || *patch.Title != "New title" {
			t.Fatalf("expected trimmed title, got=%v", patch.Title)
		}
		if patch.Slug == nil || *patch.Slug != "New-Slug" {
			t.Fatalf("expected trimmed slug, got=%v", patch.Slug)
		}
		if patch.Channel == nil || *patch.Channel != ChannelTG {
			t.Fatalf("expected normalized channel, got=%v", patch.Channel)
		}
		if patch.Status == nil || *patch.Status != StatusPublished {
			t.Fatalf("expected normalized status, got=%v", patch.Status)
		}
		if got.ID != 11 || got.TopicID != 3 || got.Channel != ChannelTG || got.Status != StatusPublished {
			t.Fatalf("result mismatch: %+v", got)
		}
	})

	t.Run("repo errors", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			name    string
			repoErr error
		}{
			{name: "not found", repoErr: ErrNotFound},
			{name: "conflict", repoErr: ErrConflict},
		}

		for _, c := range cases {
			c := c
			t.Run(c.name, func(t *testing.T) {
				t.Parallel()

				repo := &fakeRepository{
					updateArticleFn: func(ctx context.Context, id uint64, patch UpdateArticleRequest) (Article, error) {
						return Article{}, c.repoErr
					},
				}
				svc := &service{repo: repo, now: time.Now}

				_, err := svc.UpdateArticle(ctx, 1, UpdateArticleRequest{Title: stringPtr("ok")})
				if !errors.Is(err, c.repoErr) {
					t.Fatalf("UpdateArticle() err=%v want=%v", err, c.repoErr)
				}
			})
		}
	})
}

func TestMakeSlug(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		title string
		want  string
	}{
		{
			name:  "basic normalization",
			title: "  Hello   world  ",
			want:  "hello-world",
		},
		{
			name:  "keeps latin numbers and dash",
			title: "Go_1 — Intro",
			want:  "go-1-intro",
		},
		{
			name:  "empty after normalization",
			title: "!!!",
			want:  "faq",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := makeSlug(tt.title)
			if got != tt.want {
				t.Fatalf("makeSlug(%q)=%q want=%q", tt.title, got, tt.want)
			}
		})
	}
}

func TestExtractTextFromJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  datatypes.JSON
		want string
	}{
		{
			name: "invalid json",
			raw:  datatypes.JSON([]byte(`{`)),
			want: "",
		},
		{
			name: "extract string",
			raw:  mustJSON(`{"text":"hello"}`),
			want: "hello",
		},
		{
			name: "extract strings slice",
			raw:  mustJSON(`{"items":["one","two"]}`),
			want: "one two",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := extractTextFromJSON(tt.raw)
			if got != tt.want {
				t.Fatalf("extractTextFromJSON(%s)=%q want=%q", string(tt.raw), got, tt.want)
			}
		})
	}
}

func TestService_PutBlocks(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("validation errors", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name      string
			articleID uint64
			req       PutBlocksRequest
		}{
			{
				name:      "zero article id",
				articleID: 0,
				req: PutBlocksRequest{
					Blocks: []PutBlock{{Type: BlockText, Data: mustJSON(`{"text":"x"}`)}},
				},
			},
			{
				name:      "invalid block type",
				articleID: 1,
				req: PutBlocksRequest{
					Blocks: []PutBlock{{Type: "video", Data: mustJSON(`{"text":"x"}`)}},
				},
			},
			{
				name:      "invalid json",
				articleID: 1,
				req: PutBlocksRequest{
					Blocks: []PutBlock{{Type: BlockText, Data: datatypes.JSON([]byte(`{`))}},
				},
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				repo := &fakeRepository{}
				svc := &service{repo: repo, now: time.Now}

				_, err := svc.PutBlocks(ctx, tt.articleID, tt.req)
				if !errors.Is(err, ErrInvalidArgument) {
					t.Fatalf("PutBlocks() err=%v want=%v", err, ErrInvalidArgument)
				}
				if len(repo.replaceBlocksCalls) != 0 || len(repo.updateArticleSearchTextCalls) != 0 {
					t.Fatalf("repo should not be called on invalid input")
				}
			})
		}
	})

	t.Run("replace blocks and set nil search text when extracted text is empty", func(t *testing.T) {
		t.Parallel()

		repo := &fakeRepository{
			replaceBlocksFn: func(ctx context.Context, articleID uint64, blocks []Block) ([]Block, error) {
				return []Block{{ID: 1, ArticleID: articleID, Sort: blocks[0].Sort, Type: blocks[0].Type, Data: blocks[0].Data}}, nil
			},
			getArticleAnyStatusFn: func(ctx context.Context, id uint64) (Article, []Block, error) {
				return Article{ID: id}, []Block{
					{ID: 99, ArticleID: id, Type: BlockImage, Data: mustJSON(`{"url":"img.png"}`)},
					{ID: 100, ArticleID: id, Type: BlockDivider, Data: mustJSON(`{}`)},
				}, nil
			},
		}
		svc := &service{repo: repo, now: time.Now}

		created, err := svc.PutBlocks(ctx, 33, PutBlocksRequest{
			Blocks: []PutBlock{
				{Sort: 1, Type: " Text ", Data: mustJSON(`{"text":"hello"}`)},
			},
		})
		if err != nil {
			t.Fatalf("PutBlocks() err=%v", err)
		}
		if len(created) != 1 || created[0].Type != BlockText {
			t.Fatalf("PutBlocks() created mismatch: %+v", created)
		}
		if len(repo.replaceBlocksCalls) != 1 {
			t.Fatalf("expected one ReplaceBlocks call")
		}
		if repo.replaceBlocksCalls[0].blocks[0].Type != BlockText {
			t.Fatalf("block type should be normalized to %q, got=%q", BlockText, repo.replaceBlocksCalls[0].blocks[0].Type)
		}
		if len(repo.getArticleAnyStatusCalls) != 1 || repo.getArticleAnyStatusCalls[0] != 33 {
			t.Fatalf("refresh should request article 33, got calls=%v", repo.getArticleAnyStatusCalls)
		}
		if len(repo.updateArticleSearchTextCalls) != 1 {
			t.Fatalf("expected one UpdateArticleSearchText call")
		}
		if repo.updateArticleSearchTextCalls[0].searchText != nil {
			t.Fatalf("expected nil search text for non-textual blocks, got=%v", repo.updateArticleSearchTextCalls[0].searchText)
		}
	})
}

func TestService_CreateBlock(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("validation errors", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name      string
			articleID uint64
			req       CreateBlockRequest
		}{
			{
				name:      "zero article id",
				articleID: 0,
				req:       CreateBlockRequest{Type: BlockText, Data: mustJSON(`{"text":"x"}`)},
			},
			{
				name:      "invalid block type",
				articleID: 1,
				req:       CreateBlockRequest{Type: "video", Data: mustJSON(`{"text":"x"}`)},
			},
			{
				name:      "invalid json",
				articleID: 1,
				req:       CreateBlockRequest{Type: BlockText, Data: datatypes.JSON([]byte(`{`))},
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				repo := &fakeRepository{}
				svc := &service{repo: repo, now: time.Now}

				_, err := svc.CreateBlock(ctx, tt.articleID, tt.req)
				if !errors.Is(err, ErrInvalidArgument) {
					t.Fatalf("CreateBlock() err=%v want=%v", err, ErrInvalidArgument)
				}
				if len(repo.createBlockCalls) != 0 || len(repo.updateArticleSearchTextCalls) != 0 {
					t.Fatalf("repo should not be called on invalid input")
				}
			})
		}
	})

	t.Run("create block and refresh search text", func(t *testing.T) {
		t.Parallel()

		repo := &fakeRepository{
			createBlockFn: func(ctx context.Context, b *Block) error {
				b.ID = 77
				return nil
			},
			getArticleAnyStatusFn: func(ctx context.Context, id uint64) (Article, []Block, error) {
				return Article{ID: id}, []Block{
					{ID: 1, ArticleID: id, Type: BlockText, Data: mustJSON(`{"text":"hello"}`)},
					{ID: 2, ArticleID: id, Type: BlockCallout, Data: mustJSON(`{"text":"world"}`)},
				}, nil
			},
		}
		svc := &service{repo: repo, now: time.Now}

		got, err := svc.CreateBlock(ctx, 9, CreateBlockRequest{
			Sort: 3,
			Type: " CallOut ",
			Data: mustJSON(`{"text":"hello"}`),
		})
		if err != nil {
			t.Fatalf("CreateBlock() err=%v", err)
		}
		if got.ID != 77 || got.ArticleID != 9 || got.Type != BlockCallout {
			t.Fatalf("CreateBlock() result mismatch: %+v", got)
		}
		if len(repo.createBlockCalls) != 1 || repo.createBlockCalls[0].Type != BlockCallout {
			t.Fatalf("create block call mismatch: %+v", repo.createBlockCalls)
		}
		if len(repo.getArticleAnyStatusCalls) != 1 || repo.getArticleAnyStatusCalls[0] != 9 {
			t.Fatalf("refresh should request article 9, got=%v", repo.getArticleAnyStatusCalls)
		}
		if len(repo.updateArticleSearchTextCalls) != 1 || repo.updateArticleSearchTextCalls[0].searchText == nil {
			t.Fatalf("expected non-nil search text update, got=%+v", repo.updateArticleSearchTextCalls)
		}
		if !strings.Contains(*repo.updateArticleSearchTextCalls[0].searchText, "hello") {
			t.Fatalf("search text should include extracted text, got=%q", *repo.updateArticleSearchTextCalls[0].searchText)
		}
	})

	t.Run("repo and refresh errors", func(t *testing.T) {
		t.Parallel()

		t.Run("create error", func(t *testing.T) {
			t.Parallel()

			repo := &fakeRepository{
				createBlockFn: func(ctx context.Context, b *Block) error { return ErrNotFound },
			}
			svc := &service{repo: repo, now: time.Now}

			_, err := svc.CreateBlock(ctx, 1, CreateBlockRequest{Type: BlockText, Data: mustJSON(`{"text":"x"}`)})
			if !errors.Is(err, ErrNotFound) {
				t.Fatalf("CreateBlock() err=%v want=%v", err, ErrNotFound)
			}
		})

		t.Run("refresh error", func(t *testing.T) {
			t.Parallel()

			repo := &fakeRepository{
				createBlockFn:         func(ctx context.Context, b *Block) error { return nil },
				getArticleAnyStatusFn: func(ctx context.Context, id uint64) (Article, []Block, error) { return Article{}, nil, ErrNotFound },
			}
			svc := &service{repo: repo, now: time.Now}

			_, err := svc.CreateBlock(ctx, 1, CreateBlockRequest{Type: BlockText, Data: mustJSON(`{"text":"x"}`)})
			if !errors.Is(err, ErrNotFound) {
				t.Fatalf("CreateBlock() err=%v want=%v", err, ErrNotFound)
			}
		})
	})
}

func TestService_UpdateBlock(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("validation errors", func(t *testing.T) {
		t.Parallel()

		validData := mustJSON(`{"text":"x"}`)
		tests := []struct {
			name    string
			blockID uint64
			req     UpdateBlockRequest
		}{
			{
				name:    "zero block id",
				blockID: 0,
				req:     UpdateBlockRequest{Type: stringPtr(BlockText)},
			},
			{
				name:    "invalid block type",
				blockID: 1,
				req:     UpdateBlockRequest{Type: stringPtr("video")},
			},
			{
				name:    "invalid json",
				blockID: 1,
				req:     UpdateBlockRequest{Data: func() *datatypes.JSON { j := datatypes.JSON([]byte(`{`)); return &j }()},
			},
			{
				name:    "valid request control",
				blockID: 1,
				req:     UpdateBlockRequest{Type: stringPtr(" text "), Data: &validData},
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				repo := &fakeRepository{
					updateBlockFn: func(ctx context.Context, id uint64, patch UpdateBlockRequest) (Block, error) {
						return Block{ID: id, ArticleID: 2, Type: BlockText, Data: mustJSON(`{"text":"x"}`)}, nil
					},
					getArticleAnyStatusFn: func(ctx context.Context, id uint64) (Article, []Block, error) {
						return Article{ID: id}, []Block{{ID: 1, ArticleID: id, Type: BlockText, Data: mustJSON(`{"text":"x"}`)}}, nil
					},
				}
				svc := &service{repo: repo, now: time.Now}

				_, err := svc.UpdateBlock(ctx, tt.blockID, tt.req)

				if tt.name == "valid request control" {
					if err != nil {
						t.Fatalf("UpdateBlock() err=%v", err)
					}
					return
				}
				if !errors.Is(err, ErrInvalidArgument) {
					t.Fatalf("UpdateBlock() err=%v want=%v", err, ErrInvalidArgument)
				}
				if len(repo.updateBlockCalls) != 0 {
					t.Fatalf("repo should not be called on invalid input")
				}
			})
		}
	})

	t.Run("normalizes type and refreshes by updated article id", func(t *testing.T) {
		t.Parallel()

		data := mustJSON(`{"items":["one","two"]}`)
		repo := &fakeRepository{
			updateBlockFn: func(ctx context.Context, id uint64, patch UpdateBlockRequest) (Block, error) {
				return Block{
					ID:        id,
					ArticleID: 99,
					Type:      BlockBullets,
					Data:      *patch.Data,
				}, nil
			},
			getArticleAnyStatusFn: func(ctx context.Context, id uint64) (Article, []Block, error) {
				return Article{ID: id}, []Block{{ID: 1, ArticleID: id, Type: BlockBullets, Data: mustJSON(`{"items":["one","two"]}`)}}, nil
			},
		}
		svc := &service{repo: repo, now: time.Now}

		got, err := svc.UpdateBlock(ctx, 15, UpdateBlockRequest{
			Sort: intPtr(5),
			Type: stringPtr(" BULLETS "),
			Data: &data,
		})
		if err != nil {
			t.Fatalf("UpdateBlock() err=%v", err)
		}
		if got.ArticleID != 99 {
			t.Fatalf("expected updated block articleID=99, got=%d", got.ArticleID)
		}
		if len(repo.updateBlockCalls) != 1 {
			t.Fatalf("expected one update block call")
		}
		patch := repo.updateBlockCalls[0].patch
		if patch.Type == nil || *patch.Type != BlockBullets {
			t.Fatalf("expected normalized type %q, got=%v", BlockBullets, patch.Type)
		}
		if len(repo.getArticleAnyStatusCalls) != 1 || repo.getArticleAnyStatusCalls[0] != 99 {
			t.Fatalf("refresh should use articleID=99, got calls=%v", repo.getArticleAnyStatusCalls)
		}
		if len(repo.updateArticleSearchTextCalls) != 1 {
			t.Fatalf("expected one search_text update")
		}
	})

	t.Run("repo and refresh errors", func(t *testing.T) {
		t.Parallel()

		t.Run("update error", func(t *testing.T) {
			t.Parallel()

			repo := &fakeRepository{
				updateBlockFn: func(ctx context.Context, id uint64, patch UpdateBlockRequest) (Block, error) {
					return Block{}, ErrNotFound
				},
			}
			svc := &service{repo: repo, now: time.Now}

			_, err := svc.UpdateBlock(ctx, 1, UpdateBlockRequest{Type: stringPtr(BlockText)})
			if !errors.Is(err, ErrNotFound) {
				t.Fatalf("UpdateBlock() err=%v want=%v", err, ErrNotFound)
			}
		})

		t.Run("refresh error", func(t *testing.T) {
			t.Parallel()

			repo := &fakeRepository{
				updateBlockFn: func(ctx context.Context, id uint64, patch UpdateBlockRequest) (Block, error) {
					return Block{ID: id, ArticleID: 5}, nil
				},
				getArticleAnyStatusFn: func(ctx context.Context, id uint64) (Article, []Block, error) {
					return Article{}, nil, ErrNotFound
				},
			}
			svc := &service{repo: repo, now: time.Now}

			_, err := svc.UpdateBlock(ctx, 1, UpdateBlockRequest{Type: stringPtr(BlockText)})
			if !errors.Is(err, ErrNotFound) {
				t.Fatalf("UpdateBlock() err=%v want=%v", err, ErrNotFound)
			}
		})
	})
}

func TestService_DeleteBlock(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("invalid id", func(t *testing.T) {
		t.Parallel()

		svc := &service{repo: &fakeRepository{}, now: time.Now}
		err := svc.DeleteBlock(ctx, 0)
		if !errors.Is(err, ErrInvalidArgument) {
			t.Fatalf("DeleteBlock() err=%v want=%v", err, ErrInvalidArgument)
		}
	})

	t.Run("success refreshes search text", func(t *testing.T) {
		t.Parallel()

		repo := &fakeRepository{
			deleteBlockFn: func(ctx context.Context, id uint64) (uint64, error) {
				return 88, nil
			},
			getArticleAnyStatusFn: func(ctx context.Context, id uint64) (Article, []Block, error) {
				return Article{ID: id}, []Block{{ID: 1, ArticleID: id, Type: BlockText, Data: mustJSON(`{"text":"x"}`)}}, nil
			},
		}
		svc := &service{repo: repo, now: time.Now}

		err := svc.DeleteBlock(ctx, 15)
		if err != nil {
			t.Fatalf("DeleteBlock() err=%v", err)
		}
		if len(repo.deleteBlockCalls) != 1 || repo.deleteBlockCalls[0] != 15 {
			t.Fatalf("delete block call mismatch: %v", repo.deleteBlockCalls)
		}
		if len(repo.getArticleAnyStatusCalls) != 1 || repo.getArticleAnyStatusCalls[0] != 88 {
			t.Fatalf("refresh call mismatch: %v", repo.getArticleAnyStatusCalls)
		}
		if len(repo.updateArticleSearchTextCalls) != 1 || repo.updateArticleSearchTextCalls[0].searchText == nil {
			t.Fatalf("expected search_text update on delete")
		}
	})

	t.Run("repo and refresh errors", func(t *testing.T) {
		t.Parallel()

		t.Run("delete error", func(t *testing.T) {
			t.Parallel()

			repo := &fakeRepository{
				deleteBlockFn: func(ctx context.Context, id uint64) (uint64, error) { return 0, ErrNotFound },
			}
			svc := &service{repo: repo, now: time.Now}

			err := svc.DeleteBlock(ctx, 1)
			if !errors.Is(err, ErrNotFound) {
				t.Fatalf("DeleteBlock() err=%v want=%v", err, ErrNotFound)
			}
		})

		t.Run("refresh error", func(t *testing.T) {
			t.Parallel()

			repo := &fakeRepository{
				deleteBlockFn: func(ctx context.Context, id uint64) (uint64, error) { return 4, nil },
				getArticleAnyStatusFn: func(ctx context.Context, id uint64) (Article, []Block, error) {
					return Article{}, nil, ErrNotFound
				},
			}
			svc := &service{repo: repo, now: time.Now}

			err := svc.DeleteBlock(ctx, 1)
			if !errors.Is(err, ErrNotFound) {
				t.Fatalf("DeleteBlock() err=%v want=%v", err, ErrNotFound)
			}
		})
	})
}
