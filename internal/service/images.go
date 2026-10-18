package service

import (
	"context"
	"image_processing_service/internal/apperrors"
	"image_processing_service/internal/models"
	"image_processing_service/internal/repository"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type imageService struct {
	ImageRepository repository.ImageRepository
	cld     		*cloudinary.Cloudinary
}

type ImageService interface{
	UploadAndProcessImage(fileHeader *multipart.FileHeader,userID int64)(*models.Image,error)


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

