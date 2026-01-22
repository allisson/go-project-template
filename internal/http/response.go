// Package http provides HTTP server implementation and request handlers.
package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

// makeResponse writes an HTTP response with the specified body, status code, and content type.
func makeResponse(w http.ResponseWriter, body []byte, statusCode int, contentType string) {
	w.Header().Set("Content-Type", fmt.Sprintf("%s; charset=utf-8", contentType))
	w.WriteHeader(statusCode)
	if _, err := w.Write(body); err != nil {
		slog.Error("failed to write response body", slog.Any("error", err))
	}
}

// makeJSONResponse marshals the body to JSON and writes it as an HTTP response.
func makeJSONResponse(w http.ResponseWriter, statusCode int, body interface{}) {
	d, err := json.Marshal(body)
	if err != nil {
		slog.Error("failed to marshal body", slog.Any("error", err))
	}
	c := new(bytes.Buffer)
	err = json.Compact(c, d)
	if err != nil {
		slog.Error("failed to compact json", slog.Any("error", err))
	}
	makeResponse(w, c.Bytes(), statusCode, "application/json")
}
