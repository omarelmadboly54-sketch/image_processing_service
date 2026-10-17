package postgres

import (
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
	result:=r.db.Create(user)
	if result.Error!=nil{
		return result.Error
	}
	return nil
}

func(r *userRepo)GetUserByEmail(email string)(user models.User,err error){
	result:=r.db.Where("email=?",email).First(&user)
	if result.Error!=nil{
		return models.User{},result.Error
	}
	return user,nil
}
