package main

import (
	"image_processing_service/configuration"
	"image_processing_service/internal/controller"
	"image_processing_service/internal/providers/postgres"
	"image_processing_service/internal/service"
	"image_processing_service/token"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	db, config, err := configuration.InitPostgres()
	if err!=nil{
		log.Fatal("Failed to connect to DataBase",err)
	}
	log.Println("DataBase Connection Successfully")
	tokenMaker,err:=token.NewJWTMaker(config.SecretKey)
	if err!=nil{
		log.Fatal("failed to generate token maker:",err)
	}

	r:=gin.Default()

	userRepo:=postgres.NewUserRepo(db)
	tokenDuration:=time.Duration(config.AccessTokenMinutes)*time.Minute
	userService:=service.NewUserService(userRepo,tokenMaker,tokenDuration)
	userController:=controller.NewUserController(userService)
	userController.RegisterRoutes(r.Group("/api"))

	if err:=r.Run(config.ServerAddress);err!=nil{
		log.Fatalf("failed to start server:%v",err)
	}
}

