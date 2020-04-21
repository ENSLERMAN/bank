package apiserver

import "net/http"

type responseWriter struct {
	http.ResponseWriter
	code int
}

// WriteHeader - здесь мы просто перезаписываем  responseWriter на отдачу статус кода
func (w *responseWriter) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
