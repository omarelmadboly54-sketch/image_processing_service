package repository

import "image_processing_service/internal/models"

type ImageRepository interface {
	CreateImage(image *models.Image)error
}