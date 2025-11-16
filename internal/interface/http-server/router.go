package http_server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/Desnn1ch/pr-reviewer-service/internal/interface/http-server/handler"
)

func NewRouter(
	team *handler.TeamHandler,
	user *handler.UserHandler,
	pr *handler.PRHandler,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8081"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Route("/team", func(r chi.Router) {
		r.Post("/add", team.Add)
		r.Get("/get", team.Get)
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/setIsActive", user.SetIsActive)
		r.Get("/getReview", user.GetReview)
	})

	r.Route("/pullRequest", func(r chi.Router) {
		r.Post("/create", pr.Create)
		r.Post("/merge", pr.Merge)
		r.Post("/reassign", pr.Reassign)
	})

	return r
}
