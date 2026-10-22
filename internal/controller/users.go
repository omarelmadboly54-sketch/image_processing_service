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
	g.POST("/login/google",ctl.GoogleLogin)
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

func (ctl *userController) GoogleLogin(g *gin.Context) {
    var req presenters.GoogleLoginRequest

    if err := g.ShouldBindJSON(&req); err != nil {
        g.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid Google login request",
        })
        return
    }

    accessToken, err := ctl.userService.LoginWithGoogle(
        g.Request.Context(),
        req.IdToken,
    )

    if err != nil {
        switch {

        case errors.Is(err, apperrors.ErrInvalidGoogleToken),
            errors.Is(err, apperrors.ErrGoogleEmailNotFound),
            errors.Is(err, apperrors.ErrGoogleEmailNotVerified),
            errors.Is(err, apperrors.ErrGoogleAccountMismatch):

            g.JSON(http.StatusUnauthorized, gin.H{
                "error": err.Error(),
            })

        case errors.Is(err, apperrors.ErrGoogleAccountNotLinked):

            g.JSON(http.StatusConflict, gin.H{
                "error": err.Error(),
            })

        case errors.Is(err, apperrors.ErrFirebaseClientFailed):

            g.JSON(http.StatusInternalServerError, gin.H{
                "error": "Google authentication is temporarily unavailable",
            })

        default:

            g.JSON(http.StatusInternalServerError, gin.H{
                "error": "An unexpected error occurred",
            })
        }

        return
    }

    g.JSON(http.StatusOK, gin.H{
        "access_token": accessToken,
    })
}