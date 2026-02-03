package routes

import (
	"net/http"

	"github.com/Saikatdeb12/TodoApp2/internal/handlers"
	middlewares "github.com/Saikatdeb12/TodoApp2/internal/middleware"
	"github.com/Saikatdeb12/TodoApp2/internal/utils"
	"github.com/go-chi/chi/v5"
)

func SetupRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Route("/v1", func(v1 chi.Router) {
		v1.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			utils.RespondJSON(w, http.StatusOK, map[string]string{
				"status": "server is running",
			})
		})
		v1.Post("/auth/register", handlers.Register)
		v1.Post("/auth/login", handlers.Login)
		router.Post("/auth/logout", handlers.Logout) // to make it inside the middleware
		router.Group(func(r chi.Router) {
			r.Use(middlewares.Authenticate)
			r.Route("/todos", func(r chi.Router) {
				r.Get("/", handlers.GetTodos)
				r.Get("/upcoming", handlers.UpcomingTodosByDate)
				r.Post("/", handlers.CreateTodo)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", handlers.GetTodoByID)
					r.Patch("/", handlers.UpdateTodoByID)
					r.Delete("/", handlers.DeleteTodoByID)
				})
			})
			r.Route("/user", handlers.DeleteUser)
		})
	})

	return router
}
