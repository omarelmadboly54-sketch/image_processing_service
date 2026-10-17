package postgres

import (
	"errors"
	"image_processing_service/internal/models"
	"image_processing_service/internal/repository"

	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
	
}

func NewUserRepo(db *gorm.DB)repository.UserRepository{
	return &userRepo{
		db: db,
	}
}

func(r *userRepo)CreateUser(user *models.User)error{
	return r.db.Create(user).Error
}

func(r *userRepo)GetUserByEmail(email string)(*models.User,error){
	var user models.User
	err:=r.db.Where("email=?",email).First(&user).Error
	if err!=nil{
		if errors.Is(err,gorm.ErrRecordNotFound){
			return nil,nil
		}
		return nil,err
	}
	return &user,nil
}

