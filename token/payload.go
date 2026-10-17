package token

import (
	"errors"
	"time"
)
var(
	ErrExpiredToken=errors.New("token is expired")
	ErrInvalidToken=errors.New("token is invalid")

)

type Payload struct {
	Email    	string 		`json:"email"`
	UserID   	int64  		`json:"user_id"`
	IssuedAt  	time.Time	`json:"issued_at"`
	ExpiredAt 	time.Time	`json:"expired_at"`	
}


func NewPayloadToken(email string,user_id int64,duration time.Duration)(*Payload,error){
	payload:=&Payload{
		Email: email,
		UserID: user_id,
		IssuedAt: time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload,nil
}

func(p *Payload)Valid()error{
	if time.Now().After(p.ExpiredAt){
		return ErrExpiredToken
	}
	return nil
}
