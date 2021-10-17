package apiserver

import "net/http"

// Тип responseWriter будет удовлетворять интерфейсу ResponseWriter
type responseWriter struct {
	http.ResponseWriter
	code int
}

// WriteHeader ...
func (w *responseWriter) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
