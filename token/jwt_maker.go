package token

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minsecretkeysize = 32

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minsecretkeysize {
		return nil, fmt.Errorf("invalid secret key size must be at least:%d",minsecretkeysize)
	}
	return &JWTMaker{secretKey: secretKey},nil
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
	return jwtToken.SignedString([]byte(maker.secretKey))
}



func (maker *JWTMaker) VerifyToken(tokenString string) (*Payload, error) {
    claims := jwt.MapClaims{}

    parsedToken, err := jwt.ParseWithClaims(
        tokenString,
        claims,
        func(token *jwt.Token) (interface{}, error) {
            return []byte(maker.secretKey), nil
        },
        jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
        jwt.WithExpirationRequired(),
        jwt.WithIssuedAt(),
        jwt.WithJSONNumber(),
    )

    if err != nil {
        if errors.Is(err, jwt.ErrTokenExpired) {
            return nil, ErrExpiredToken
        }
        return nil, ErrInvalidToken
    }

    if !parsedToken.Valid {
        return nil, ErrInvalidToken
    }

    email, ok := claims["email"].(string)
    if !ok || email == "" {
        return nil, ErrInvalidToken
    }

    userIDNumber, ok := claims["user_id"].(json.Number)
    if !ok {
        return nil, ErrInvalidToken
    }

    userID, err := userIDNumber.Int64()
    if err != nil || userID <= 0 {
        return nil, ErrInvalidToken
    }

    exp, err := claims.GetExpirationTime()
    if err != nil || exp == nil {
        return nil, ErrInvalidToken
    }

    return &Payload{
        UserID:    userID,
        Email:     email,
        ExpiredAt: exp.Time,
    }, nil
}
