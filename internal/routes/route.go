package routes

import (
	"github.com/Saikatdeb12/TodoApp2/internal/handlers"
	middlewares "github.com/Saikatdeb12/TodoApp2/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func SetupRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/auth/register", handlers.Register)
	router.Post("/auth/login", handlers.Login)
	router.Post("/auth/logout", handlers.Logout)
	router.Group(func(r chi.Router) {
		r.Use(middlewares.Auth)
		r.Route("/todos", func(r chi.Router) {
			r.Get("/", handlers.GetTodos)
			r.Post("/", handlers.CreateTodo)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", handlers.GetTodoByID)
				r.Patch("/", handlers.UpdateTodoByID)
				r.Delete("/", handlers.DeleteTodoByID)
			})
		})
	})

	return router
}
