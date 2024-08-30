package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"wb_l0/internal/repository"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	srv *http.Server
}

type router struct {
	router *chi.Mux
	repo   repository.Repo
	cache  repository.Repo
}

func New(repo repository.Repo, cache repository.Repo) Server {
	r := router{
		repo:  repo,
		cache: cache,
	}
	r.initRoutes()
	return Server{&http.Server{Addr: "0.0.0.0:8081", Handler: r.router}}
}

func (s *Server) Run() error {
	return http.ListenAndServe("localhost:8081", s.srv.Handler)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
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
		log.Default().Println("order from cache", order)
		fmt.Fprint(w, order)
		return
	}

	// internal error
	if !errors.Is(err, repository.ErrorNotFound) {
		fmt.Fprintf(w, "error getting order from cache: %v", err)
		return
	}

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
