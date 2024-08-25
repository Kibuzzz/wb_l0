package server

import (
	"fmt"
	"net/http"
	"wb_l0/internal/repository"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	router *chi.Mux
	repo   repository.Repo
}

func New(repo repository.Repo) *Server {
	server := Server{repo: repo}
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
	fmt.Println(orderUID)
	order, err := s.repo.GetOrderByID(r.Context(), orderUID)
	if err != nil {
		fmt.Fprintf(w, "error: %v", err)
		return
	}
	fmt.Fprint(w, order)
}
