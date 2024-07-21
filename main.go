package main

import (
	"fmt"
	"net/http"
	"nrs16/cme/config"
	"nrs16/cme/server"

	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func main() {

	/// get full path to pass to load conf function
	executablePath, err := os.Executable()
	if err != nil {
		log.Fatalf("error getting file executable path: %s", err.Error())
	}
	executableDir := filepath.Dir(executablePath)

	///load configuretaion
	config, err := config.LoadConfig(executableDir + "/config.toml")
	if err != nil {
		log.Fatalf("could not load configs: %s", err.Error())
	}
	///// get information database and insert in redis:

	////initialise app with router and DB
	app, err := server.InitialiseApp(config)
	if err != nil {
		log.Fatalf("could not build app: %s", err.Error())
	}

	log.Infof("starting server on port: %d", config.App.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", config.App.Port), app.Router)
	if err != nil {
		log.Fatalf("could not start server on port: %d, %s", config.App.Port, err.Error())
	}
}
