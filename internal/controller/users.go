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
	g.POST("/login",ctl.Login)
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

func(ctl *userController)Login(g *gin.Context){
	req:=presenters.LoginReq{}

	if err:=g.ShouldBindJSON(&req);err!=nil{
		g.JSON(http.StatusBadRequest,gin.H{"error":err.Error(),
	})
	return
	}
	
	accessToken,err:=ctl.userService.Login(req.Email,req.Password)
	if err!=nil{
		if errors.Is(err,apperrors.ErrInvalidCredentials){
			g.JSON(http.StatusUnauthorized,gin.H{"error":apperrors.ErrInvalidCredentials.Error(),
		})
		return
		}
		g.JSON(http.StatusInternalServerError,gin.H{"error":err.Error(),
	})
	return
	}
	g.JSON(http.StatusOK,gin.H{
		"access_token":accessToken,
	})
}
