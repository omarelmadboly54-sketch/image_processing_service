package postgres

import (
	"image_processing_service/internal/models"
	"image_processing_service/internal/repository"

	"gorm.io/gorm"
)

type imageRepo struct {
	db *gorm.DB
}


func NewImageRepo(db *gorm.DB)repository.ImageRepository{
	return &imageRepo{
		db: db,
	}

}

func(r *imageRepo)CreateImage(image *models.Image)error{
	return r.db.Create(image).Error
}
