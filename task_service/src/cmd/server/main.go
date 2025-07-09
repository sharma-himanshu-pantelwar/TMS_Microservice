package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	sessionclient "task_service/src/internal/adaptors/grpcclient"
	"task_service/src/internal/adaptors/persistance"
	"task_service/src/internal/config"
	"task_service/src/internal/interfaces/input/api/rest/handler"
	"task_service/src/internal/interfaces/input/api/rest/routes"
	"task_service/src/internal/usecase"
	"task_service/src/pkg/migrate"
)

func main() {
	database, err := persistance.NewDatabase()
	if err != nil {
		log.Fatalf("Error connecting to db %v", err)
	}
	fmt.Println("Connected to database")
	// fetch current cwd
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error fetching cwd %v", err)
	}

	//run migrations
	migrate := migrate.NewMigrate(
		database.GetDB(),
		cwd+"/src/migrations",
	)
	err = migrate.RunMigrations()
	if err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// loadconfig
	configurations, err := config.LoadConfig()
	if err != nil {
		fmt.Println("failed to load config")
	}

	//repos
	taskRepo := persistance.NewTaskRepo(database)

	//services
	taskService := usecase.NewTaskService(taskRepo)

	// handler
	taskHandler := handler.NewTaskHandler(taskService)

	sessionClient, err := sessionclient.NewClient("localhost:50051")
	if err != nil {
		fmt.Printf("unable to connect to session client due to : %v", err)
		// os.Exit(1)
	}
	fmt.Println("?connected to session client")
	defer sessionClient.Close()
	router := routes.InitRoutes(&taskHandler, sessionClient)
	// router := routes.InitRoutes(&taskHandler)

	err = http.ListenAndServe(fmt.Sprintf(":%s", configurations.APP_PORT), router)
	if err != nil {
		fmt.Printf("failed to start server %v", err)
		os.Exit(1)
	}

}

// TODO
// create more tasks api if needed
// notification http part if any
// grpc implement in all and connect
// redis for notification
