package main

import (
	"image_processing_service/configuration"
	"image_processing_service/internal/controller"
	"image_processing_service/internal/providers/postgres"
	"image_processing_service/internal/service"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	db, config, err := configuration.InitPostgres()
	if err!=nil{
		log.Fatal("Failed to connect to DataBase",err)
	}
	log.Println("DataBase Connection Successfully")

	r:=gin.Default()

	userRepo:=postgres.NewUserRepo(db)
	userService:=service.NewUserService(userRepo)
	userController:=controller.NewUserController(userService)
	userController.RegisterRoutes(r.Group("/api"))

	if err:=r.Run(config.ServerAddress);err!=nil{
		log.Fatalf("failed to start server:%v",err)
	}
}
