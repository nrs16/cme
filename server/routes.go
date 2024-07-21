package server

import (
	"nrs16/cme/middleware"
)

func ConfigureRoutes(app *App) {

	app.Router.Use(middleware.AuthenticateUser)
	app.Router.HandleFunc("/api/v1/register", app.Register).Methods("POST")
	app.Router.HandleFunc("/api/v1/login", app.Login).Methods("POST")
	app.Router.HandleFunc("/api/v1/message", app.SendMessage).Methods("POST")
	app.Router.HandleFunc("/api/v1/chat", app.GetChats).Methods("GET")
	app.Router.HandleFunc("/api/v1/chat/{chatId}/message", app.GetChatMessages).Methods("GET")
}
