package main

import (
	"encoding/csv"
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

	http.ListenAndServe(":"+h.port, mux)
}

func (h *HTTPServer) outloadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv")
	w.WriteHeader(http.StatusOK)

	data, err := h.db.GetAll()
	if err != nil {
		httpError(w, err)
		return
	}

	cw := csv.NewWriter(w)
	cw.Write([]string{"ID", "Value", "CreatedAt"})
	for _, row := range data {
		id := strconv.Itoa(row.ID)
		cw.Write([]string{id, row.Value, row.CreatedAt.String()})
	}

	cw.Flush()
}

func httpError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}
