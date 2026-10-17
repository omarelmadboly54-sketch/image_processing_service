package token

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const minsecretkeysize = 32

type JWTMaker struct {
	secretKet string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minsecretkeysize {
		return nil, fmt.Errorf("invalid secret key size must be at least:%d",minsecretkeysize)
	}
	return &JWTMaker{secretKet: secretKey},nil
}

func(maker *JWTMaker)CreateToken(email string, user_id int64,duration time.Duration) (string, error){
	Payload,err:=NewPayloadToken(email,user_id,duration)
	if err!=nil{
		return "",err
	}
	jwtToken:=jwt.NewWithClaims(jwt.SigningMethodHS256,
	jwt.MapClaims{
		"email":Payload.Email,
		"user_id":Payload.UserID,
		"iat":Payload.IssuedAt.Unix(),
		"exp":Payload.ExpiredAt.Unix(),
	})
	return jwtToken.SignedString([]byte(maker.secretKet))
}
