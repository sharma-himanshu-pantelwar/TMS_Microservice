package routes

import (
	"net/http"
	sessionclient "task_service/src/internal/adaptors/grpcclient"
	"task_service/src/internal/interfaces/input/api/rest/handler"
	samw "task_service/src/internal/interfaces/input/api/rest/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func InitRoutes(taskHandler *handler.TaskHandler, sessionClient *sessionclient.Client) http.Handler {

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http:*", "https:*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))
	router.Route("/v1/tasks", func(r chi.Router) {
		r.Use(samw.SessionAuthMiddleware(sessionClient))

		r.Post("/", taskHandler.RegisterTaskHandler)
		r.Post("/usrstatus", taskHandler.CheckAssignedUserStatus)
		r.Put("/{id}", taskHandler.UpdateTaskHandler)
		r.Get("/all", taskHandler.GetAllTasksHandler)
		r.Get("/", taskHandler.GetMyTasksHandler)
		r.Delete("/{id}", taskHandler.DeleteTaskHandler)
		r.Get("/bin/", taskHandler.GetTaskBinHandler)
		r.Patch("/bin/restore/{id}", taskHandler.RestoreTaskFromBinHandler)
		r.Delete("/bin/delete/{id}", taskHandler.DeleteTaskFromBinHandler)
		r.Delete("/permanent/{id}", taskHandler.DeleteTaskPermanentlyHandler)

	})

	return router
}
