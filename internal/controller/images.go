package controller

import (
	"image_processing_service/internal/apperrors"
	"image_processing_service/internal/middleware"
	"image_processing_service/internal/service"
	"image_processing_service/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

type imageController struct {
	ImageService service.ImageService
	tokenMaker	 token.Maker
}

type ImageControlller interface{
	RegisterRoutes(g *gin.RouterGroup)
}


func NewImageController(ImageService service.ImageService,tokenMaker token.Maker)ImageControlller{
	return &imageController{
		ImageService: ImageService,
		tokenMaker: tokenMaker,
	}
}


func(ctl *imageController)RegisterRoutes(g *gin.RouterGroup){
	g.POST("/images",middleware.AuthMiddleware(ctl.tokenMaker),ctl.UploadImage)
}


func(ctl *imageController)UploadImage(g *gin.Context){

	payload,ok:=middleware.GetTokenPayload(g)
	if!ok{
		g.JSON(http.StatusUnauthorized,gin.H{"error":"unauthorized request"})
		return
	}
	userID:=payload.UserID

	fileHeader,err:=g.FormFile("image")

	if err!=nil{
		g.JSON(http.StatusBadRequest,gin.H{"error":apperrors.ErrImageRequest.Error(),
	})
		return
	}

	image,err:=ctl.ImageService.UploadAndProcessImage(fileHeader,userID)
	if err!=nil{
		g.JSON(http.StatusBadRequest,gin.H{"error":err.Error(),
	})
	return
	}

	g.JSON(http.StatusCreated,gin.H{
		"message":"Image uploaded successfully",
		"data":image,
	})
}
