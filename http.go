package main

import (
	"embed"
	"encoding/csv"
	"encoding/json"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strconv"
)

//go:embed www
var www embed.FS

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
	mux.HandleFunc("GET /", h.indexHandler("index.html", "text/html"))
	mux.HandleFunc("GET /index.js", h.indexHandler("index.js", "application/javascript"))
	mux.HandleFunc("GET /data", h.dataHandler)
	mux.HandleFunc("GET /__outload", h.outloadHandler)
	mux.HandleFunc("GET /health", h.healthHandler)

	http.ListenAndServe(
		net.JoinHostPort("", h.port),
		mux,
	)
}

func (h *HTTPServer) indexHandler(fname string, ctype string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		f, err := http.FS(www).Open("www/" + fname)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", ctype)
		w.WriteHeader(http.StatusOK)

		io.Copy(w, f)
	}
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
	cw.Write([]string{"id", "value", "created_at"})
	for _, row := range data {
		id := strconv.Itoa(row.ID)
		cw.Write([]string{id, row.Value, row.CreatedAt.String()})
		slog.Debug("writing row", "row", row)
	}

	cw.Flush()
}

func (h *HTTPServer) dataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	slog.Debug("exporting data in json")

	data, err := h.db.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
