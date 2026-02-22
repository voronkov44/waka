package middleware

import (
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"
)

//TODO X-Forwarded-For and X-Real-IP can be faked - danger

func Logging(log *slog.Logger) Middleware {
	if log == nil {
		log = slog.Default()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrapper := &WrapperWriter{
				ResponseWriter: w,
				StatusCode:     http.StatusOK,
			}
			next.ServeHTTP(wrapper, r)

			duration := time.Since(start)

			attributes := []any{
				"status", wrapper.StatusCode,
				"method", r.Method,
				"path", r.URL.Path,
				"bytes", wrapper.Bytes,
				"duration_ms", duration.Milliseconds(),
				"ip", clientIP(r),
			}
			if query := r.URL.RawQuery; query != "" {
				attributes = append(attributes, "query", query)
			}

			switch {
			case wrapper.StatusCode >= 500:
				log.Error("http request failed", attributes...)
			case wrapper.StatusCode >= 400:
				log.Warn("http request failed", attributes...)
			default:
				log.Info("http request succeeded", attributes...)
			}
		})
	}
}

func clientIP(r *http.Request) string {
	if xff := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); xff != "" {
		for _, part := range strings.Split(xff, ",") {
			ip := strings.TrimSpace(part)
			if ip != "" {
				return ip
			}
		}
	}
	if xr := strings.TrimSpace(r.Header.Get("X-Real-IP")); xr != "" {
		return xr
	}

	ra := strings.TrimSpace(r.RemoteAddr)
	if ra == "" {
		return "unknown"
	}
	host, _, err := net.SplitHostPort(ra)
	if err == nil && host != "" {
		return host
	}
	return ra
}
