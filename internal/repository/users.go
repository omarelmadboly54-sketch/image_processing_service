package repository

import "image_processing_service/internal/models"

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string)(user models.User,err error)
}
