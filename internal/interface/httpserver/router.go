package httpserver

import (
	"net/http"

	"github.com/Desnn1ch/pr-reviewer-service/internal/interface/httpserver/handler"
	"github.com/go-chi/chi/v5"
)

func NewRouter(
	team *handler.TeamHandler,
	user *handler.UserHandler,
	pr *handler.PRHandler,
	st *handler.StatsHandler,
) http.Handler {
	r := chi.NewRouter()

	UseMiddlewares(r)

	registerTeamRoutes(r, team)
	registerUserRoutes(r, user)
	registerPRRoutes(r, pr)
	registerStatsRoutes(r, st)

	return r
}

func registerTeamRoutes(r chi.Router, h *handler.TeamHandler) {
	r.Route("/team", func(r chi.Router) {
		r.Post("/add", h.Add)
		r.Get("/get", h.Get)
	})
}

func registerUserRoutes(r chi.Router, h *handler.UserHandler) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/setIsActive", h.SetIsActive)
		r.Get("/getReview", h.GetReview)
	})
}

func registerPRRoutes(r chi.Router, h *handler.PRHandler) {
	r.Route("/pullRequest", func(r chi.Router) {
		r.Post("/create", h.Create)
		r.Post("/merge", h.Merge)
		r.Post("/reassign", h.Reassign)
	})
}

func registerStatsRoutes(r chi.Router, h *handler.StatsHandler) {
	r.Route("/stats", func(r chi.Router) {
		r.Get("/reviewers", h.GetReviewerStats)
	})
}
