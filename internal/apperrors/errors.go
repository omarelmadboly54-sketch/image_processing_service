package apperrors

import "errors"

var (
	//user
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidCredentials=errors.New("invalid credentials")
	//image
	ErrInvalidSize=errors.New("file size exceeds the 5MB limit")
	ErrInvalidImageType=errors.New("invalid file format. Only JPEG, PNG, and WebP are allowed")
	ErrImageRequest=errors.New("image file is required under form key 'image'")
)
