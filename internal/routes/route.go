package routes

import "github.com/go-chi/chi/v5"

func SetupRouter() *chi.Mux{
	r := chi.NewRouter()
	r.Post("/auth/register", handlers.Register)
	r.Post("/auth/login", handlers.Login)
	r.Post("/auth/logout", handlers.Logout)
	r.Group(func (r chi.Router){
		r.Use(middlewares.Auth)
		r.Get("/todos", handlers.GetTodos)
		r.Get("/todos/{id}", handlers.GetTodoByID)
		r.Get("/todos/complete", handlers.CompletedTodos)
		r.Get("/todos/incomplete", handlers.InCompleteTodos)
		r.Get("/todos/upcoming-todos", handlers.UpcomingTodosByDate)
		r.Post("/todos", handlers.CreateTodo)
		r.Put("/todos/{id}", handlers.UpdateTodoByID)
		r.Delete("/todos/{id}", handlers.DeleteTodoByID)
	})

	return r
}
