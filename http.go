package main

import (
	"encoding/csv"
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"strconv"
)

type HTTPServer struct {
	port string
	db   *DB
}

func NewHTTPServer(port string, db *DB) *HTTPServer {
	return &HTTPServer{
		port: port,
		db:   db,
	}
}

func (h *HTTPServer) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /__outload", h.outloadHandler)
	mux.HandleFunc("GET /health", h.healthHandler)

	http.ListenAndServe(
		net.JoinHostPort("", h.port),
		mux,
	)
}

func (h *HTTPServer) outloadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv")
	w.WriteHeader(http.StatusOK)

	slog.Debug("exporting data")

	data, err := h.db.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cw := csv.NewWriter(w)
	cw.Write([]string{"ID", "Value", "CreatedAt"})
	for _, row := range data {
		id := strconv.Itoa(row.ID)
		cw.Write([]string{id, row.Value, row.CreatedAt.String()})
		slog.Debug("writting row", "row", row)
	}

	cw.Flush()
}

func (h *HTTPServer) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	slog.Debug("health check")
	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": "ok",
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
