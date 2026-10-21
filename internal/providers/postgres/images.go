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


func(r *imageRepo)GetImageById(id int64,user_id int64)(*models.Image,error){

	var image models.Image

	result:=r.db.Where("id=? AND user_id=?",id,user_id).First(&image)

	if result.Error!=nil{
		return nil,result.Error
	}
	return &image,nil
}

func(r *imageRepo)GetAll(user_id int64,page int,limit int)([]models.Image,int64,error){

	var images []models.Image
	var total int64

	if err:=r.db.Model(&images).Where("user_id=?",user_id).Count(&total).Error;err!=nil{
		return nil,0,err
	}

	offset:=(page-1)*limit
	if err:=r.db.Where("user_id=?",user_id).Limit(limit).Offset(offset).Find(&images).Error;err!=nil{
		return nil,0,err
	}

	return images,total,nil
}

func(r *imageRepo)DeleteImage(id int64,user_id int64)error{
	return r.db.Where("id=? AND user_id=?",id,user_id).Delete(&models.Image{}).Error
}
