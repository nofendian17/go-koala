package controllers

import "github.com/nofendian17/rest/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Auth Route
	s.Router.HandleFunc("/Login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")
	s.Router.HandleFunc("/Refresh", middlewares.SetMiddlewareJSON(s.Refresh)).Methods("POST")
	s.Router.HandleFunc("/LogOut", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.LogOut))).Methods("GET")

	// Customers Route
	s.Router.HandleFunc("/Customer", middlewares.SetMiddlewareJSON(s.CreateCustomer)).Methods("POST")
	s.Router.HandleFunc("/Customers", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetCustomers))).Methods("GET")

	// Orders Route
	s.Router.HandleFunc("/Order", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.CreateOrder))).Methods("POST")
}
