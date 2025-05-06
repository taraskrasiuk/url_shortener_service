package middlewares

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type LoggerMiddleware struct {
	out     io.Writer
	handler http.Handler
}

func NewLoggerMiddleware(next http.Handler, out io.Writer) *LoggerMiddleware {
	return &LoggerMiddleware{out, next}
}

func (l *LoggerMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	logTime := start.Format(time.RFC822)
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		contentType = r.Header.Get("content-type")
	}
	// server request
	l.handler.ServeHTTP(w, r)
	// write directly to writer
	fmt.Fprintf(l.out, "[LOG]: %s %s Content-Type: %s \t %s [%s]", r.Method, r.URL.Path, contentType, logTime, time.Since(start))
}
