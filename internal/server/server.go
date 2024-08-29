package server

import (
	"errors"
	"fmt"
	"net/http"
	"wb_l0/internal/repository"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	router *chi.Mux
	repo   repository.Repo
	cache  repository.Repo
}

func New(repo repository.Repo, cache repository.Repo) *Server {
	server := Server{
		repo:  repo,
		cache: cache,
	}
	server.initRoutes()
	return &server
}

func (s *Server) Run() error {
	return http.ListenAndServe("localhost:8081", s.router)
}

func (s *Server) initRoutes() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/{orderUID}", s.getOrder)
	s.router = router
}

func (s *Server) getOrder(w http.ResponseWriter, r *http.Request) {
	orderUID := chi.URLParam(r, "orderUID")

	order, err := s.cache.GetOrderByID(r.Context(), orderUID)

	if err != nil {
		fmt.Fprint(w, order)
		return
	}

	// internal error
	if !errors.Is(err, repository.ErrorNotFound) {
		fmt.Fprintf(w, "error getting order from cache: %v", err)
		return
	}

	//
	order, err = s.repo.GetOrderByID(r.Context(), orderUID)
	if err != nil {
		fmt.Fprintf(w, "error getting order from repo: %v", err)
		return

	}

	err = s.cache.AddOrder(r.Context(), order)
	if err != nil {
		fmt.Fprintf(w, "error adding order to cache: %v", err)
		return
	}
	fmt.Fprint(w, order)
}
