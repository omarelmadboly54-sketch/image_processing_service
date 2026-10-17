package controller

import (
	"errors"
	"image_processing_service/internal/apperrors"
	"image_processing_service/internal/presenters"
	"image_processing_service/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userController struct {
	userService service.UserService
}

type UserController interface{
	RegisterRoutes(g *gin.RouterGroup)
}

func NewUserController(userService service.UserService)UserController{
	return &userController{
		userService: userService,
	}
}

func(ctl *userController)RegisterRoutes(g *gin.RouterGroup){
	g.POST("/users",ctl.CreateUser)
}


func(ctl *userController)CreateUser(g *gin.Context){
	req:=presenters.CreateUserReq{}

	if err:=g.ShouldBindJSON(&req);err!=nil{
		g.JSON(http.StatusBadRequest,gin.H{"error":err.Error(),
	})
	return
	}

	toUser:=req.ToUser()
	
	err:=ctl.userService.CreateUser(&toUser)
	if err!=nil{
		if errors.Is(err,apperrors.ErrUserAlreadyExists){
			g.JSON(http.StatusBadRequest,gin.H{"error":apperrors.ErrUserAlreadyExists.Error(),
		})
			return
		}
		g.JSON(http.StatusInternalServerError,gin.H{"error":err.Error(),
	})
	return
	}

	g.JSON(http.StatusCreated,gin.H{"ok":"ok"})
}