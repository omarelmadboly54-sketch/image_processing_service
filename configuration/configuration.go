package configuration

import (
	"context"
	"image_processing_service/internal/models"
	"image_processing_service/internal/util"
	"os"

	firebase "firebase.google.com/go/v4"
	"github.com/spf13/viper"
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

func InitFireBase()(*firebase.App,error){
	credPath:=viper.GetString("FIREBASE_CREDENTIALS_FILE")
	if credPath == ""{
		credPath = "serviceAccountKey.json"
	}

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS",credPath)
	app,err:=firebase.NewApp(context.Background(),nil)
	if err!=nil{
		return nil,err
	}
	return app,nil
}


