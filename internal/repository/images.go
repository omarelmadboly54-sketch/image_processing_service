package repository

import "image_processing_service/internal/models"

type ImageRepository interface {
	CreateImage(image *models.Image)error
	GetImageById(id int64,user_id int64)(*models.Image,error)
	GetAll(user_id int64,page int,limit int)([]models.Image,int64,error)
	DeleteImage(id int64,user_id int64)error
}
