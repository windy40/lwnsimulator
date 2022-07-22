package main

import (
	"log"

	cnt "github.com/windy40/lwnsimulator/controllers"
	"github.com/windy40/lwnsimulator/models"
	repo "github.com/windy40/lwnsimulator/repositories"
	ws "github.com/windy40/lwnsimulator/webserver"
)

func main() {

	var info *models.ServerConfig
	var err error

	simulatorRepository := repo.NewSimulatorRepository()
	simulatorController := cnt.NewSimulatorController(simulatorRepository)
	simulatorController.GetIstance()

	log.Println("LWN Simulator is online...")

	info, err = models.GetConfigFile("config.json")
	if err != nil {
		log.Fatal(err)
	}

	WebServer := ws.NewWebServer(info, simulatorController)
	WebServer.Run()

}
