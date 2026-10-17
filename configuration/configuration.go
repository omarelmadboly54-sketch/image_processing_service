package configuration

import (
	"image_processing_service/internal/models"
	"image_processing_service/internal/util"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgres()(*gorm.DB,models.Config,error) {
	config,err:=util.LoadConfig(".")
	if err!=nil{
		return nil,models.Config{},err
	}
	db, err := gorm.Open(postgres.Open(config.DBSource), &gorm.Config{})
	if err!=nil{
		return nil,models.Config{},err
	}

	return db,config,nil

}
