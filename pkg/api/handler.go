package api

import (
	"encoding/json"
	"net/http"
	"time"

	"octopus-cache/internal/database"

	"github.com/gorilla/mux"
)

type Handler struct {
	db *database.Database
}

func NewHandler(db *database.Database) *Handler {
	return &Handler{db: db}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/set/{key}", h.setHandler).Methods("POST")
	router.HandleFunc("/get/{key}", h.getHandler).Methods("GET")
	router.HandleFunc("/delete/{key}", h.deleteHandler).Methods("DELETE")
	router.HandleFunc("/query", h.queryHandler).Methods("GET")
	router.HandleFunc("/health", h.healthHandler).Methods("GET")

	// Serve the request
	router.ServeHTTP(w, r)
}

// setHandler handles setting a key-value pair with TTL
func (h *Handler) setHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	var request struct {
		Value   interface{}            `json:"value"`
		TTL     string                 `json:"ttl"`
		Indexes map[string]interface{} `json:"indexes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ttl, err := time.ParseDuration(request.TTL)
	if err != nil {
		http.Error(w, "invalid TTL duration", http.StatusBadRequest)
		return
	}

	h.db.Set(key, request.Value, ttl, request.Indexes)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "key": key})
}

// getHandler retrieves a value by key
func (h *Handler) getHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, exists := h.db.Get(key)
	if !exists {
		http.Error(w, "key not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"key":   key,
		"value": value,
	})
}

// deleteHandler deletes a key-value pair
func (h *Handler) deleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	h.db.Delete(key)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "key": key})
}

// queryHandler handles indexed queries
func (h *Handler) queryHandler(w http.ResponseWriter, r *http.Request) {
	index := r.URL.Query().Get("index")
	value := r.URL.Query().Get("value")

	if index == "" || value == "" {
		http.Error(w, "index and value parameters are required", http.StatusBadRequest)
		return
	}

	keys := h.db.Query(index, value)
	if len(keys) == 0 {
		http.Error(w, "no results found", http.StatusNotFound)
		return
	}

	results := make(map[string]interface{})
	for _, key := range keys {
		if value, exists := h.db.Get(key); exists {
			results[key] = value
		}
	}

	json.NewEncoder(w).Encode(results)
}

// healthHandler provides a health check endpoint
func (h *Handler) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
