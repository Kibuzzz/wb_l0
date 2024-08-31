package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"wb_l0/internal/repository"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type router struct {
	router *chi.Mux
	repo   repository.Repo
	cache  repository.Repo
}

func New(repo repository.Repo, cache repository.Repo, addr string) *http.Server {
	r := router{
		repo:  repo,
		cache: cache,
	}
	r.initRoutes()
	return &http.Server{Addr: addr, Handler: r.router}
}

func (s *router) initRoutes() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/{orderUID}", s.getOrder)
	s.router = router
}

func (s *router) getOrder(w http.ResponseWriter, r *http.Request) {
	orderUID := chi.URLParam(r, "orderUID")

	order, err := s.cache.GetOrderByID(r.Context(), orderUID)
	if err == nil {
		log.Println("Order found in cache:", order)
		writeJSONResponse(w, http.StatusOK, order)
		return
	}

	if !errors.Is(err, repository.ErrorNotFound) {
		log.Printf("Error fetching order from cache: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	order, err = s.repo.GetOrderByID(r.Context(), orderUID)

	if err != nil {
		log.Printf("Error fetching order from repository: %v", err)
		if errors.Is(err, repository.ErrorNotFound) {
			writeJSONResponse(w, http.StatusNotFound, fmt.Sprintf("Order with UID %s not exists", orderUID))
			return
		} else {
			http.Error(w, "Order not found", http.StatusInternalServerError)
			return
		}
	}

	if !errors.Is(err, repository.ErrorNotFound) {
		log.Printf("Error fetching order from cache: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := s.cache.AddOrder(r.Context(), order); err != nil {
		log.Printf("Error adding order to cache: %v", err)
	}

	writeJSONResponse(w, http.StatusOK, order)
}

func writeJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
