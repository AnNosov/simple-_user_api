package http

import (
	"github.com/AnNosov/simple_user_api/internal/usecase"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func NewRouter(r *chi.Mux, t usecase.UserActionUseCase) {
	r.Use(middleware.Logger)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))
	r.Use(middleware.SetHeader("charset", "utf-8"))

	NewUserActionRoutes(r, t)

}
