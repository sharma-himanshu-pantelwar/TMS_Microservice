package routes

import (
	"net/http"
	"user_service/src/internal/interfaces/input/api/rest/handler"
	"user_service/src/internal/interfaces/input/api/rest/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func InitRoutes(userHandler *handler.UserHandler, jwtkey string) http.Handler {

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http:*", "https:*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))
	router.Route("/v1/users", func(r chi.Router) {
		r.Post("/register", userHandler.RegisterUserHandler)
		r.Post("/login", userHandler.LoginUserHandler)

		r.Group(func(protected chi.Router) {
			protected.Use(middleware.Authenticate(jwtkey))
			protected.Get("/", userHandler.GetProfileHandler)
			protected.Post("/logout", userHandler.LogoutHandler)
		})
	})

	return router
}
