package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (hashedpassword string,err error) {
	byte, err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	if err!=nil{
		return "",fmt.Errorf("failed to hash password:%w",err)
	}
	return string(byte),nil
}

func VerifyPassword(hashedpassword,password string)error{
	return bcrypt.CompareHashAndPassword([]byte(hashedpassword),[]byte(password))
}
