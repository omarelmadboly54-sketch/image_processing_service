package token

import "time"

type Maker interface {
	CreateToken(email string, user_id int64,duration time.Duration) (string, error)
		VerifyToken(token string)(*Payload,error)
}
