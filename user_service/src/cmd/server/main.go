package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"user_service/src/internal/adaptors/persistance"
	"user_service/src/internal/config"

	"user_service/src/internal/interfaces/input/api/rest/handler"
	"user_service/src/internal/interfaces/input/api/rest/routes"
	"user_service/src/internal/usecase"
	"user_service/src/pkg/migrate"

	pb "user_service/src/internal/interfaces/output/grpc"
	grpcserver "user_service/src/internal/interfaces/output/grpc/server"

	grpcgoogle "google.golang.org/grpc"
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
	userRepo := persistance.NewUserRepo(database)
	sessionRepo := persistance.NewSessionRepo(database)

	//services
	userService := usecase.NewUserService(userRepo, sessionRepo, configurations.JWT_SECRET)

	// handler
	userHandler := handler.NewUserHandler(configurations, userService)

	router := routes.InitRoutes(&userHandler, configurations.JWT_SECRET)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Printf("Failed to listen grpc server in user service %v", err)
	}

	grpcServer := grpcgoogle.NewServer()
	// pb.RegisterSessionValidateServer(grpcServer)
	// pb.RegisterSessionValidateServer(grpcServer, &pb.SessionValidatorServer{UserService: userService})
	pb.RegisterSessionValidatorServer(grpcServer, &grpcserver.SessionServer{DB: database.GetDB()})
	go func() {

		fmt.Println("gRPC server listening on 50051 port")
		if err := grpcServer.Serve(lis); err != nil {
			fmt.Printf("failed to serve gRPC server: %v", err)
		}
	}()

	err = http.ListenAndServe(fmt.Sprintf(":%s", configurations.APP_PORT), router)
	if err != nil {
		fmt.Printf("failed to start server %v", err)
		os.Exit(1)
	}

}
