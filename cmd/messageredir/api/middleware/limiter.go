package middleware

import (
	"io"
	"net/http"
)

type LimitReadCloser struct {
	io.Reader
	underlying io.ReadCloser
}

func (lrc *LimitReadCloser) Close() error {
	return lrc.underlying.Close()
}

func BodyLimit(maxBodySize int64, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = &LimitReadCloser{io.LimitReader(r.Body, maxBodySize), r.Body}
		next.ServeHTTP(w, r)
	})
}
