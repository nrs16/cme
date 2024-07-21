package server

import (
	"nrs16/cme/middleware"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func ConfigureRoutes(app *App) {

	////middleware
	app.Router.Use(middleware.AuthenticateUser)

	//// appliction API
	app.Router.HandleFunc("/api/v1/register", app.Register).Methods("POST")
	app.Router.HandleFunc("/api/v1/login", app.Login).Methods("POST")
	app.Router.HandleFunc("/api/v1/message", app.SendMessage).Methods("POST")
	app.Router.HandleFunc("/api/v1/chat", app.GetChats).Methods("GET")
	app.Router.HandleFunc("/api/v1/chat/{chatId}/message", app.GetChatMessages).Methods("GET")

	////prometheus
	app.Router.Handle("/metrics", promhttp.Handler()).Methods("GET")
}
