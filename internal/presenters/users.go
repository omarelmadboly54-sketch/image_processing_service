package presenters

import "image_processing_service/internal/models"

type CreateUserReq struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginReq struct{
	Email		string		`json:"email"`
	Password	string		`json:"password"`
}

type GoogleLoginRequest struct{
	IdToken  	string		`json:"id_token" binding:"required"`
}

func(c *CreateUserReq) ToUser() models.User{
	return models.User{
		Name: c.Name,
		Email: c.Email,
		Password: c.Password,
	}
}

