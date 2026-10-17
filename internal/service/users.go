package service

import (
	"fmt"
	"image_processing_service/internal/apperrors"
	"image_processing_service/internal/models"
	"image_processing_service/internal/repository"
	"image_processing_service/internal/util"
	"image_processing_service/token"
	"time"
)

type userService struct {
	UserRepository repository.UserRepository
	tokenMaker	   token.Maker
	tokenDuration  time.Duration
}

type UserService interface{
	CreateUser(*models.User)error
	Login(email,password string)(string,error)
}

func NewUserService(userRepository repository.UserRepository,tokenMaker token.Maker,tokenDuration time.Duration)UserService{
	return &userService{
		UserRepository: userRepository,
		tokenMaker: tokenMaker,
		tokenDuration: tokenDuration,
	}
}

func(s *userService)CreateUser(user *models.User)error{

	existingUser,err:=s.UserRepository.GetUserByEmail(user.Email)
	if err !=nil{
		return err
	}
	if existingUser!=nil{
		return apperrors.ErrUserAlreadyExists

	}
	
	hashedpassword,err:=util.HashPassword(user.Password)
	if err!=nil{
		return err
	}
	user.Password=hashedpassword

	return s.UserRepository.CreateUser(user)
}

func(s *userService)Login(email,password string)(string,error){
	user,err:=s.UserRepository.GetUserByEmail(email)
	if err!=nil{
		return "",err
	}
	if user ==nil{
		return "",apperrors.ErrInvalidCredentials
	}
	err=util.VerifyPassword(user.Password,password)
	if err!=nil{
		return "",apperrors.ErrInvalidCredentials
	}
	accessToken,err:=s.tokenMaker.CreateToken(user.Email,user.ID,s.tokenDuration)
	if err!=nil{
		return "",fmt.Errorf("failed to generate tokens")
	}
	return accessToken,nil
}

