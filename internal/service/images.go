package service

import (
	"context"
	"errors"
	"image_processing_service/internal/apperrors"
	"image_processing_service/internal/models"
	"image_processing_service/internal/presenters"
	"image_processing_service/internal/repository"
	"image_processing_service/internal/util"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"gorm.io/gorm"
)

type imageService struct {
	ImageRepository repository.ImageRepository
	cld     		*cloudinary.Cloudinary
}

type ImageService interface{
	UploadAndProcessImage(fileHeader *multipart.FileHeader,userID int64)(*models.Image,error)
	GetImageById(id int64,user_id int64)(image *models.Image,err error)
	GetAll(user_id int64,page int,limit int)([]models.Image,int64,error)
	DeleteImage(id int64,user_id int64)error
	TransformAndSaveImage(id int64,user_id int64,req *presenters.TransformImageRequest)(*models.Image,error)
}

func NewImageService(ImageRepository repository.ImageRepository,cld *cloudinary.Cloudinary)ImageService{
	return &imageService{
		ImageRepository: ImageRepository,
		cld: cld,
	}
}


func(s *imageService)UploadAndProcessImage(fileHeader *multipart.FileHeader,userID int64)(*models.Image,error){

	//check file size
	var maxfilesize int64=5*1024*1024

	if fileHeader.Size>maxfilesize{
		return nil,apperrors.ErrInvalidSize
	}

	//check file format
	ext:=strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext!=".jpg" &&ext !=".jpeg" &&ext!=".png" &&ext!=".webq"{
		return nil,apperrors.ErrInvalidImageType
	}

	file,err:=fileHeader.Open()
	if err!=nil{
		return nil,err
	}

	defer file.Close()

	ctx:=context.Background()
	uploadResult,err:=s.cld.Upload.Upload(ctx,file,uploader.UploadParams{
		Folder:"omarix",
	})
	if err!=nil{
		return nil,err
	}

	image:=&models.Image{
		FileUrl: uploadResult.SecureURL,
		CloudinaryID:uploadResult.PublicID,
		OriginalName: fileHeader.Filename,
		FileSize: fileHeader.Size,
		Format: uploadResult.Format,
		Width: uploadResult.Width,
		Height: uploadResult.Height,
		UserID: userID,
	}

	err=s.ImageRepository.CreateImage(image)

	if err!=nil{
		return nil,err
	}
	return image,nil
}


func(s *imageService)GetImageById(id int64,user_id int64)(image *models.Image,err error){
	
	image,err=s.ImageRepository.GetImageById(id,user_id)

	if err!=nil{
		if errors.Is(err,gorm.ErrRecordNotFound){
			return nil,apperrors.ErrImageNotFound
		}
		return nil,err
	}
	return image,nil
}



func(s *imageService)GetAll(user_id int64,page int,limit int)([]models.Image,int64,error){
	return s.ImageRepository.GetAll(user_id,page,limit)
}

func(s *imageService)DeleteImage(id int64,user_id int64)error{
	image,err:=s.ImageRepository.GetImageById(id,user_id)

	if err!=nil{
		if errors.Is(err,gorm.ErrRecordNotFound){
			return apperrors.ErrImageNotFound
		}
		return err
	}

	ctx:=context.Background()
	_,err=s.cld.Upload.Destroy(ctx,uploader.DestroyParams{
		PublicID:image.CloudinaryID ,
	})

	if err!=nil{
		return err
	}

	err=s.ImageRepository.DeleteImage(image.ID,image.UserID)
	if err!=nil{
		return err
	}
	return nil
}


func(s *imageService)TransformAndSaveImage(id int64,user_id int64,req *presenters.TransformImageRequest)(*models.Image,error){
	image,err:=s.ImageRepository.GetImageById(id,user_id)
	if err!=nil{
		if errors.Is(err,gorm.ErrRecordNotFound){
			return nil,apperrors.ErrImageNotFound
		}
		return nil,err
	}

	ctx:=context.Background()

	transformString:=util.BuildCloudinaryTransformationString(req)

	uploadParams:=uploader.UploadParams{
		Transformation: transformString,
	}

	uploadResult,err:=s.cld.Upload.Upload(ctx,image.FileUrl,uploadParams)
	if err!=nil{
		return nil,apperrors.ErrImageTransformFail
	}
	newImage:=&models.Image{
		FileUrl: uploadResult.SecureURL,
		CloudinaryID: uploadResult.PublicID,
		OriginalName: "transformed_" + image.OriginalName,
		FileSize: int64(uploadResult.Bytes),
		Format: uploadResult.Format,
		Width: uploadResult.Width,
		Height: uploadResult.Height,
		UserID: user_id,
	}

	err=s.ImageRepository.CreateImage(newImage)
	if err!=nil{
		return nil,err
	}

	return newImage,nil

}
