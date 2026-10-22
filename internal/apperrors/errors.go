package apperrors

import "errors"

var (
	//user
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidCredentials=errors.New("invalid credentials")
	// Google Auth
    ErrInvalidGoogleToken   = errors.New("invalid or expired google token")
    ErrGoogleEmailNotFound  = errors.New("email not found in token claims")
    ErrFirebaseClientFailed = errors.New("failed to get firebase auth client")
	ErrEmailNotVerified=errors.New("Google email is not verified")
	ErrGoogleEmailNotVerified = errors.New("Google email is not verified")
    ErrGoogleAccountNotLinked = errors.New("This email is already registered. Please sign in with your password first.",)
    ErrGoogleAccountMismatch = errors.New("This Google account does not match the linked account.",)
	//image
	ErrInvalidSize=errors.New("file size exceeds the 5MB limit")
	ErrInvalidImageType=errors.New("invalid file format. Only JPEG, PNG, and WebP are allowed")
	ErrImageRequest=errors.New("image file is required under form key 'image'")
	ErrImageNotFound=errors.New("Image Not Found")
	ErrImageTransformFail=errors.New("failed to transform image")
	ErrInvalidTransformation=errors.New("invalid transformation request")
)
