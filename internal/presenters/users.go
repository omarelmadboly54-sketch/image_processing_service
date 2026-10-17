package presenters

import "image_processing_service/internal/models"

type CreateUserReq struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginReq struct{
	Email		string		`json:"email"`
	Password	string		`json:"password"`
}

func(c *CreateUserReq) ToUser() models.User{
	return models.User{
		Name: c.Name,
		Email: c.Email,
		Password: c.Password,
	}
}

