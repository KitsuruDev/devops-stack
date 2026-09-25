// backend/handlers.go
package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Note struct {
	ID int64 `json:"id"`
	Title string `json:"title"`
	Body string `json:"body"`
}

func health(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	dbOK := pool.Ping(ctx) == nil
	status := "ok"
	if !dbOK {
		status = "degraded"
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{"status": status, "db": dbOK})
}

func listNotes(w http.ResponseWriter, r *http.Request) {
	rows, err := pool.Query(r.Context(), "SELECT id, title, body FROM notes ORDER BY id")
	if err != nil {
		http.Error(w, `{"error":"db error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	notes := []Note{}
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.Title, &n.Body); err != nil {
			http.Error(w, `{"error":"scan error"}`, http.StatusInternalServerError)
			return
		}
		notes = append(notes, n)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

func addNote(w http.ResponseWriter, r *http.Request) {
	var n Note
	if err := json.NewDecoder(io.LimitReader(r.Body, 1<<20)).Decode(&n); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}
	id, err := pool.QueryRow(r.Context(), "INSERT INTO notes (title, body) VALUES ($1, $2) RETURNING id", n.Title, n.Body).Scan()
	if err != nil {
		http.Error(w, `{"error":"db error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"id": id})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

func getNote(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	var n Note
	if err := pool.QueryRow(r.Context(), "SELECT id, title, body FROM notes WHERE id = $1",
	id).Scan(&n.ID, &n.Title, &n.Body); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	writeJSON(w, http.StatusOK, n)
}
