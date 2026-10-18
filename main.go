package main

import (
	"image_processing_service/configuration"
	"image_processing_service/internal/controller"
	"image_processing_service/internal/providers/postgres"
	"image_processing_service/internal/service"
	"image_processing_service/token"
	"log"
	"time"
    "github.com/cloudinary/cloudinary-go/v2"
	"github.com/gin-contrib/cors"
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

	cld,err:=cloudinary.NewFromURL(config.CloudinaryUrl)
	if err!=nil{
		log.Fatalf("failed to initialize Cloudinary:%v",err)
	}
	log.Println("Cloudinary Initialized Successfully")

	r:=gin.Default()

	r.Use(cors.New(cors.Config{
    AllowOrigins: []string{"http://localhost:5173"},
    AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
    AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
}))

	//users
	userRepo:=postgres.NewUserRepo(db)
	tokenDuration:=time.Duration(config.AccessTokenMinutes)*time.Minute
	userService:=service.NewUserService(userRepo,tokenMaker,tokenDuration)
	userController:=controller.NewUserController(userService)
	userController.RegisterRoutes(r.Group("/api"))
	//images
	imageRepo:=postgres.NewImageRepo(db)
	imageService:=service.NewImageService(imageRepo,cld)
	imageController:=controller.NewImageController(imageService,tokenMaker)
	imageController.RegisterRoutes(r.Group("/api"))

	if err:=r.Run(config.ServerAddress);err!=nil{
		log.Fatalf("failed to start server:%v",err)
	}
}

