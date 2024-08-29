package routes

import (
	"Blogsite/handlers"
	middleware "Blogsite/middlewares"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/register", handlers.Register).Methods("POST")
	r.HandleFunc("/api/login", handlers.Login).Methods("POST")

	s := r.PathPrefix("/api").Subrouter()
	s.Use(middleware.AuthMiddleware)

	s.HandleFunc("/user/blog", handlers.CreateBlog).Methods("POST")
	s.HandleFunc("/feed", handlers.GetAllBlogs).Methods("GET")
	s.HandleFunc("/user/blogs", handlers.GetUserBlogs).Methods("GET")
	s.HandleFunc("/blog/{id}", handlers.GetBlogById).Methods("GET")
	s.HandleFunc("/user/{id}", handlers.GetUser).Methods("GET")
	s.HandleFunc("/user/{id}", handlers.UpdateUser).Methods("PUT")
	s.HandleFunc("/user/blog/{id}", handlers.UpdateBlog).Methods("PUT")
	s.HandleFunc("/user/blog/{id}", handlers.DeleteBlog).Methods("DELETE")

	return r
}
