package controller

import (
	"errors"
	"image_processing_service/internal/apperrors"
	"image_processing_service/internal/middleware"
	"image_processing_service/internal/presenters"
	"image_processing_service/internal/service"
	"image_processing_service/internal/util"
	"image_processing_service/token"
	"net/http"
	"strconv"

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
	g.GET("images/:id",middleware.AuthMiddleware(ctl.tokenMaker),ctl.GetImage)
	g.GET("/images",middleware.AuthMiddleware(ctl.tokenMaker),ctl.GetAll)
	g.DELETE("/images/:id",middleware.AuthMiddleware(ctl.tokenMaker),ctl.DeleteImage)
	g.POST("/images/:id/transform",middleware.AuthMiddleware(ctl.tokenMaker),ctl.TransformAndSaveImage)
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

func(ctl *imageController)GetImage(g *gin.Context){
	id,err:=strconv.ParseInt(g.Param("id"),10,64)
	if err!=nil{
		g.JSON(http.StatusBadRequest,gin.H{"error":err.Error(),
	})
	return
	}

	payload,ok:=middleware.GetTokenPayload(g)
	if !ok{
		g.JSON(http.StatusUnauthorized,gin.H{"error":"unauthorized request"})
		return
	}

	userID:=payload.UserID

	image,err:=ctl.ImageService.GetImageById(id,userID)
	if err!=nil{
		if errors.Is(err,apperrors.ErrImageNotFound){
			g.JSON(http.StatusNotFound,gin.H{"error":apperrors.ErrImageNotFound.Error(),
		})
		return
		}
		g.JSON(http.StatusInternalServerError,gin.H{"error":err.Error(),
	})
	return
	}

	g.JSON(http.StatusOK,gin.H{"image":image})
}


func(ctl *imageController)GetAll(g *gin.Context){

	page,_:=strconv.Atoi(g.DefaultQuery("page","1"))
	if page<1{
		page =1
	}
	limit,_:=strconv.Atoi(g.DefaultQuery("limit","5"))
	if limit<1{
		limit =1
	}

	payload,ok:=middleware.GetTokenPayload(g)
	if !ok{
		g.JSON(http.StatusUnauthorized,gin.H{"error":"unauthorized request"})
		return
	}
	userID:=payload.UserID

	images,total,err:=ctl.ImageService.GetAll(userID,page,limit)
	if err!=nil{
		g.JSON(http.StatusInternalServerError,gin.H{"error":err.Error(),
	})
	return
	}
	g.JSON(http.StatusOK,gin.H{
		"images":images,
		"total":total,
	})
}


func(ctl *imageController)DeleteImage(g *gin.Context){
	id,err:=strconv.ParseInt(g.Param("id"),10,64)
	if err!=nil{
		g.JSON(http.StatusBadRequest,gin.H{"error":err.Error(),
	})
	return
	}
	payload,ok:=middleware.GetTokenPayload(g)
	if !ok{
		g.JSON(http.StatusUnauthorized,gin.H{"error":"unauthorized request"})
		return
	}
	userID:=payload.UserID

	err=ctl.ImageService.DeleteImage(id,userID)
	if err!=nil{
		if errors.Is(err,apperrors.ErrImageNotFound){
			g.JSON(http.StatusNotFound,gin.H{"error":apperrors.ErrImageNotFound.Error()})
			return
		}
		g.JSON(http.StatusInternalServerError,gin.H{"error":err.Error(),
	})
	return
	}

	g.JSON(http.StatusOK,gin.H{"message":"Image is deleted successfully"})
}


func(ctl *imageController)TransformAndSaveImage(g *gin.Context){
	id,err:=strconv.ParseInt(g.Param("id"),10,64)
	if err!=nil{
		g.JSON(http.StatusBadRequest,gin.H{"error":err.Error(),
	})
	return
	}
	payload,ok:=middleware.GetTokenPayload(g)
	if !ok{
		g.JSON(http.StatusUnauthorized,gin.H{"error":"unauthorized request"})
		return
	}
	userID:=payload.UserID

	req:=presenters.TransformImageRequest{}
	if err=g.ShouldBindJSON(&req);err!=nil{
		g.JSON(http.StatusBadRequest,gin.H{"error":err.Error(),
	})
	return
	}

	if err:=util.ValidateTransformRequest(&req);err!=nil{
		g.JSON(http.StatusBadRequest,gin.H{
			"error":err.Error(),
		})
		return
	}
	
	transformedImage,err:=ctl.ImageService.TransformAndSaveImage(id,userID,&req)

	if err!=nil{
	  switch {
    case errors.Is(err, apperrors.ErrInvalidTransformation):
        g.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })

    case errors.Is(err, apperrors.ErrImageNotFound):
        g.JSON(http.StatusNotFound, gin.H{
            "error": "image not found",
        })

    case errors.Is(err, apperrors.ErrImageTransformFail):
        g.JSON(http.StatusBadGateway, gin.H{
            "error": "image transformation service failed",
        })

    default:
        g.JSON(http.StatusInternalServerError, gin.H{
            "error": "internal server error",
        })
    }
    return
}

	g.JSON(http.StatusCreated,gin.H{
		"message":"Image transformed successfully",
		"Image":transformedImage,
})
}