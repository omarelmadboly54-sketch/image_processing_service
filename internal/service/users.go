package service

import (
	"context"
	"errors"
	"fmt"
	"image_processing_service/internal/apperrors"
	"image_processing_service/internal/models"
	"image_processing_service/internal/repository"
	"image_processing_service/internal/util"
	"image_processing_service/token"
	"time"

	firebase "firebase.google.com/go/v4"
)

type userService struct {
	UserRepository repository.UserRepository
	tokenMaker	   token.Maker
	tokenDuration  time.Duration
	firebaseAPP	   *firebase.App
}

type UserService interface{
	CreateUser(*models.User)error
	Login(email,password string)(string,error)
	LoginWithGoogle(ctx context.Context,idToken string)(string,error)
}

func NewUserService(userRepository repository.UserRepository,tokenMaker token.Maker,tokenDuration time.Duration,firebaseAPP *firebase.App)UserService{
	return &userService{
		UserRepository: userRepository,
		tokenMaker: tokenMaker,
		tokenDuration: tokenDuration,
		firebaseAPP: firebaseAPP,
	}
}

func(s *userService)CreateUser(user *models.User)error{

	existingUser,err:=s.UserRepository.GetUserByEmail(user.Email)
	if err !=nil{
		return err
	}
	if existingUser!=nil{
		return apperrors.ErrUserAlreadyExists

	}
	
	hashedpassword,err:=util.HashPassword(user.Password)
	if err!=nil{
		return err
	}
	user.Password=hashedpassword

	return s.UserRepository.CreateUser(user)
}

func(s *userService)Login(email,password string)(string,error){
	user,err:=s.UserRepository.GetUserByEmail(email)
	if err!=nil{
		return "",err
	}
	if user ==nil{
		return "",apperrors.ErrInvalidCredentials
	}

	if user.GoogleID!=nil{
		return "",errors.New("please sign in using your google account")
	}

	err=util.VerifyPassword(user.Password,password)
	if err!=nil{
		return "",apperrors.ErrInvalidCredentials
	}
	accessToken,err:=s.tokenMaker.CreateToken(user.Email,user.ID,s.tokenDuration)
	if err!=nil{
		return "",fmt.Errorf("failed to generate tokens")
	}
	return accessToken,nil
}

func(s *userService)LoginWithGoogle(ctx context.Context,idToken string)(string,error){

	client,err:=s.firebaseAPP.Auth(ctx)
	if err!=nil{
		return "",apperrors.ErrFirebaseClientFailed
	}

	tokenClaims,err:=client.VerifyIDToken(ctx,idToken)
	if err!=nil{
		return "",apperrors.ErrInvalidGoogleToken
	}

	email,ok:=tokenClaims.Claims["email"].(string)
	if !ok|| email == ""{
		return "",apperrors.ErrGoogleEmailNotFound
	}

	emailVerified,ok:=tokenClaims.Claims["email_verified"].(bool)
	if !ok || !emailVerified {
		return "",apperrors.ErrEmailNotVerified
	}

	uid:=tokenClaims.UID

	name,_:=tokenClaims.Claims["name"].(string)
	if name ==""{
		name="Google User"
	}

	user,err:=s.UserRepository.GetUserByEmail(email)
	if err!=nil{
		return "",err
		}
	
	if user == nil{
		newUser:=&models.User{
			Name: name,
			Email: email,
			Password: "",
			GoogleID: &uid,
		}

		if err:=s.UserRepository.CreateUser(newUser);err!=nil{
			return "",err
		}
		user=newUser
	}else{
		if user.GoogleID==nil{
			return "",apperrors.ErrGoogleAccountNotLinked
		}

		if *user.GoogleID!=uid{
			return "",apperrors.ErrGoogleAccountMismatch
		}
	}
	accessToken,err:=s.tokenMaker.CreateToken(user.Email,user.ID,s.tokenDuration)
	if err!=nil{
		return "",err
	}
	return accessToken,nil
}