package photourl

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeResolver struct {
	publicOut    string
	presignedOut string
	presignedErr error

	publicCalls    int
	presignedCalls int
	lastKey        string
	lastTTL        time.Duration
}

func (f *fakeResolver) PublicURL(key string) string {
	f.publicCalls++
	f.lastKey = key
	if f.publicOut != "" {
		return f.publicOut
	}
	return "https://public.example/" + key
}

func (f *fakeResolver) PresignedGetURL(_ context.Context, key string, ttl time.Duration) (string, error) {
	f.presignedCalls++
	f.lastKey = key
	f.lastTTL = ttl
	if f.presignedErr != nil {
		return "", f.presignedErr
	}
	if f.presignedOut != "" {
		return f.presignedOut, nil
	}
	return "https://signed.example/" + key, nil
}

func TestResolve(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		resolver       *fakeResolver
		photoKey       *string
		opt            Options
		wantURL        *string
		wantPublic     int
		wantPresigned  int
		wantLastKey    string
		wantPresignTTL time.Duration
	}{
		{
			name:     "returns nil for nil key",
			resolver: &fakeResolver{},
			photoKey: nil,
			opt:      Options{UsePresigned: false},
			wantURL:  nil,
		},
		{
			name:       "returns nil for empty key",
			resolver:   &fakeResolver{},
			photoKey:   strPtr(""),
			opt:        Options{UsePresigned: true, PresignTTL: time.Minute},
			wantURL:    nil,
			wantPublic: 0,
		},
		{
			name:        "returns public url when presign disabled",
			resolver:    &fakeResolver{publicOut: "https://cdn.example/bucket/key.jpg"},
			photoKey:    strPtr("models/1/key.jpg"),
			opt:         Options{UsePresigned: false},
			wantURL:     strPtr("https://cdn.example/bucket/key.jpg"),
			wantPublic:  1,
			wantLastKey: "models/1/key.jpg",
		},
		{
			name:           "returns presigned url when enabled",
			resolver:       &fakeResolver{presignedOut: "https://signed.example/get?x=1"},
			photoKey:       strPtr("models/2/key.jpg"),
			opt:            Options{UsePresigned: true, PresignTTL: 5 * time.Minute},
			wantURL:        strPtr("https://signed.example/get?x=1"),
			wantPresigned:  1,
			wantLastKey:    "models/2/key.jpg",
			wantPresignTTL: 5 * time.Minute,
		},
		{
			name:          "returns nil when presigned resolver errors",
			resolver:      &fakeResolver{presignedErr: errors.New("boom")},
			photoKey:      strPtr("models/3/key.jpg"),
			opt:           Options{UsePresigned: true, PresignTTL: time.Minute},
			wantURL:       nil,
			wantPresigned: 1,
			wantLastKey:   "models/3/key.jpg",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := Resolve(context.Background(), tc.resolver, tc.photoKey, tc.opt)

			if !equalStringPtr(got, tc.wantURL) {
				t.Fatalf("Resolve() url = %#v, want %#v", got, tc.wantURL)
			}
			if tc.resolver.publicCalls != tc.wantPublic {
				t.Fatalf("Resolve() public calls = %d, want %d", tc.resolver.publicCalls, tc.wantPublic)
			}
			if tc.resolver.presignedCalls != tc.wantPresigned {
				t.Fatalf("Resolve() presigned calls = %d, want %d", tc.resolver.presignedCalls, tc.wantPresigned)
			}
			if tc.wantLastKey != "" && tc.resolver.lastKey != tc.wantLastKey {
				t.Fatalf("Resolve() last key = %q, want %q", tc.resolver.lastKey, tc.wantLastKey)
			}
			if tc.wantPresignTTL != 0 && tc.resolver.lastTTL != tc.wantPresignTTL {
				t.Fatalf("Resolve() ttl = %s, want %s", tc.resolver.lastTTL, tc.wantPresignTTL)
			}
		})
	}
}

func strPtr(v string) *string {
	return &v
}

func equalStringPtr(a, b *string) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}
