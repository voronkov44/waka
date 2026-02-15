package models_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"rest_waka/internal/models"
)

func TestModelsHandlerBadJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		method string
		path   string
		body   string
	}{
		{
			name:   "create bad json",
			method: http.MethodPost,
			path:   "/api/models",
			body:   "{",
		},
		{
			name:   "update bad json",
			method: http.MethodPatch,
			path:   "/api/models/1",
			body:   "{",
		},
		{
			name:   "add flavor bad json",
			method: http.MethodPost,
			path:   "/api/models/1/flavors",
			body:   "{",
		},
		{
			name:   "remove flavor bad json",
			method: http.MethodDelete,
			path:   "/api/models/1/flavors",
			body:   "{",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mux := newModelsTestMux(newFakeRepo())
			rr := doRequest(t, mux, tc.method, tc.path, tc.body)
			if rr.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
			}
			if !strings.Contains(rr.Body.String(), "invalid json body") {
				t.Fatalf("body = %q, want invalid json body", rr.Body.String())
			}
		})
	}
}

func TestModelsHandlerInvalidID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		method string
		path   string
		body   string
	}{
		{name: "get invalid id", method: http.MethodGet, path: "/api/models/not-a-number"},
		{name: "get must read id from path not query", method: http.MethodGet, path: "/api/models/not-a-number?id=1"},
		{name: "patch invalid id", method: http.MethodPatch, path: "/api/models/not-a-number", body: `{"name":"x"}`},
		{name: "delete invalid id", method: http.MethodDelete, path: "/api/models/not-a-number"},
		{name: "add flavor invalid id", method: http.MethodPost, path: "/api/models/not-a-number/flavors", body: `{"value":"mint"}`},
		{name: "remove flavor invalid id", method: http.MethodDelete, path: "/api/models/not-a-number/flavors", body: `{"value":"mint"}`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mux := newModelsTestMux(newFakeRepo())
			rr := doRequest(t, mux, tc.method, tc.path, tc.body)
			if rr.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
			}
			if !strings.Contains(rr.Body.String(), "invalid id") {
				t.Fatalf("body = %q, want invalid id", rr.Body.String())
			}
		})
	}
}

func TestModelsHandlerNotFound(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		method string
		path   string
		body   string
	}{
		{name: "get not found", method: http.MethodGet, path: "/api/models/404"},
		{name: "patch not found", method: http.MethodPatch, path: "/api/models/404", body: `{"name":"new"}`},
		{name: "delete not found", method: http.MethodDelete, path: "/api/models/404"},
		{name: "add flavor not found", method: http.MethodPost, path: "/api/models/404/flavors", body: `{"value":"mint"}`},
		{name: "remove flavor not found", method: http.MethodDelete, path: "/api/models/404/flavors", body: `{"value":"mint"}`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mux := newModelsTestMux(newFakeRepo())
			rr := doRequest(t, mux, tc.method, tc.path, tc.body)
			if rr.Code != http.StatusNotFound {
				t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
			}
			if !strings.Contains(rr.Body.String(), "not found") {
				t.Fatalf("body = %q, want not found", rr.Body.String())
			}
		})
	}
}

func TestModelsHandlerHappyPath(t *testing.T) {
	t.Parallel()

	repo := newFakeRepo()
	mux := newModelsTestMux(repo)

	createBody := `{"name":"Waka 7000","puffs_max":7000,"flavors":["mint"]}`
	createResp := doRequest(t, mux, http.MethodPost, "/api/models", createBody)
	if createResp.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d", createResp.Code, http.StatusCreated)
	}
	created := decodeModelResponse(t, createResp)
	if created.ID == 0 {
		t.Fatalf("created model id = %d, want > 0", created.ID)
	}

	getResp := doRequest(t, mux, http.MethodGet, "/api/models/1", "")
	if getResp.Code != http.StatusOK {
		t.Fatalf("get status = %d, want %d", getResp.Code, http.StatusOK)
	}
	got := decodeModelResponse(t, getResp)
	if got.Name != "Waka 7000" {
		t.Fatalf("get name = %q, want %q", got.Name, "Waka 7000")
	}

	listResp := doRequest(t, mux, http.MethodGet, "/api/models?limit=10&offset=0", "")
	if listResp.Code != http.StatusOK {
		t.Fatalf("list status = %d, want %d", listResp.Code, http.StatusOK)
	}
	var list models.ListModelsResponse
	decodeJSON(t, listResp, &list)
	if len(list.Items) != 1 {
		t.Fatalf("list items = %d, want 1", len(list.Items))
	}

	updateResp := doRequest(t, mux, http.MethodPatch, "/api/models/1", `{"name":"Waka Updated"}`)
	if updateResp.Code != http.StatusOK {
		t.Fatalf("update status = %d, want %d", updateResp.Code, http.StatusOK)
	}
	updated := decodeModelResponse(t, updateResp)
	if updated.Name != "Waka Updated" {
		t.Fatalf("update name = %q, want %q", updated.Name, "Waka Updated")
	}

	clearFlavorsResp := doRequest(t, mux, http.MethodPatch, "/api/models/1", `{"flavors":null}`)
	if clearFlavorsResp.Code != http.StatusOK {
		t.Fatalf("clear flavors status = %d, want %d", clearFlavorsResp.Code, http.StatusOK)
	}
	cleared := decodeModelResponse(t, clearFlavorsResp)
	if len(cleared.Flavors) != 0 {
		t.Fatalf("clear flavors result = %#v, want empty", cleared.Flavors)
	}

	addFlavorResp := doRequest(t, mux, http.MethodPost, "/api/models/1/flavors", `{"value":"cola"}`)
	if addFlavorResp.Code != http.StatusOK {
		t.Fatalf("add flavor status = %d, want %d", addFlavorResp.Code, http.StatusOK)
	}
	withFlavor := decodeModelResponse(t, addFlavorResp)
	if !contains(withFlavor.Flavors, "cola") {
		t.Fatalf("add flavor result = %#v, want flavor cola", withFlavor.Flavors)
	}

	removeFlavorResp := doRequest(t, mux, http.MethodDelete, "/api/models/1/flavors", `{"value":"cola"}`)
	if removeFlavorResp.Code != http.StatusOK {
		t.Fatalf("remove flavor status = %d, want %d", removeFlavorResp.Code, http.StatusOK)
	}
	withoutFlavor := decodeModelResponse(t, removeFlavorResp)
	if contains(withoutFlavor.Flavors, "cola") {
		t.Fatalf("remove flavor result = %#v, flavor cola should be removed", withoutFlavor.Flavors)
	}

	deleteResp := doRequest(t, mux, http.MethodDelete, "/api/models/1", "")
	if deleteResp.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, want %d", deleteResp.Code, http.StatusNoContent)
	}
}

func newModelsTestMux(repo *fakeRepo) *http.ServeMux {
	svc := models.NewService(repo)
	mux := http.NewServeMux()
	models.NewModelsHandler(mux, models.HandlerDeps{Service: svc})
	return mux
}

func doRequest(t *testing.T, h http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	return rr
}

func decodeModelResponse(t *testing.T, rr *httptest.ResponseRecorder) models.Model {
	t.Helper()

	var out models.Model
	decodeJSON(t, rr, &out)
	return out
}

func decodeJSON(t *testing.T, rr *httptest.ResponseRecorder, dst any) {
	t.Helper()
	if err := json.Unmarshal(rr.Body.Bytes(), dst); err != nil {
		t.Fatalf("json.Unmarshal() error = %v, body=%q", err, rr.Body.String())
	}
}

func contains(values []string, needle string) bool {
	for _, v := range values {
		if v == needle {
			return true
		}
	}
	return false
}
