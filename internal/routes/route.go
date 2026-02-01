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
		r.Get("/todos", handlers.GetTodos)
		r.Get("/todos/{id}", handlers.GetTodoByID)
		r.Get("/todos/complete", handlers.CompletedTodos)
		r.Get("/todos/incomplete", handlers.InCompleteTodos)
		r.Get("/todos/upcoming-todos", handlers.UpcomingTodosByDate)
		r.Post("/todo", handlers.CreateTodo)
		r.Put("/todo/{id}", handlers.UpdateTodoByID)
		r.Delete("/todo/{id}", handlers.DeleteTodoByID)
	})

	return router
}
