package httpx

import (
	"net/http"
	"strconv"
)

func QueryInt(r *http.Request, key string, def int) int {
	s := r.URL.Query().Get(key)
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}

func PathUint64(r *http.Request, key string) (uint64, error) {
	return strconv.ParseUint(r.PathValue(key), 10, 64)
}
