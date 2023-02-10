package controllers

import "golang-docker-todo/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	//Users routes
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	//Tasks routes
	s.Router.HandleFunc("/tasks", middlewares.SetMiddlewareJSON(s.CreateTask)).Methods("POST")
	s.Router.HandleFunc("/tasks", middlewares.SetMiddlewareJSON(s.GetTasks)).Methods("GET")
	s.Router.HandleFunc("/tasks/{id}", middlewares.SetMiddlewareJSON(s.GetTask)).Methods("GET")
	s.Router.HandleFunc("/tasks/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateTask))).Methods("PUT")
	s.Router.HandleFunc("/tasks/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteTask)).Methods("DELETE")
}