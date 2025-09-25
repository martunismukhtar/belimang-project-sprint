package main

import (
	"belimang/src/api/routes"
	"belimang/src/config"
	"log"
)

// @title           Fitbyte API
// @version         1.0
// @description     This is a sample server for a fitness tracking application.

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	v := config.NewViper()
	db := config.NewGorm(v)
	app := config.NewFiber(v)

	if err := config.NewSwagger(app); err != nil {
		log.Printf("Failed to initialize Swagger: %v", err)
	}

	services := config.InitServices(db)
	routes.SetupRoutes(app, v, db, services)

	// Run server
	port := v.GetString("SERVER_PORT")
	if port == "" {
		port = "7880"
	}
	log.Fatal(app.Listen(":" + port))
}
