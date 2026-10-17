package service

import (
	"image_processing_service/internal/apperrors"
	"image_processing_service/internal/models"
	"image_processing_service/internal/repository"
	"image_processing_service/internal/util"
)

type userService struct {
	UserRepository repository.UserRepository
}

type UserService interface{
	CreateUser(*models.User)error
}

func NewUserService(userRepository repository.UserRepository)UserService{
	return &userService{
		UserRepository: userRepository,
	}
}

func(s *userService)CreateUser(user *models.User)error{
	
	_,err:=s.UserRepository.GetUserByEmail(user.Email)
	if err == nil{
		return apperrors.ErrUserAlreadyExists
	}
	hashedpassword,err:=util.HashPassword(user.Password)
	if err!=nil{
		return err
	}
	user.Password=hashedpassword

	return s.UserRepository.CreateUser(user)
}

