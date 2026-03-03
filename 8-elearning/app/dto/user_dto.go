package dto

import (
	"elearning-api/model"
	"time"
)

type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	FullName  string    `json:"fullName"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewUserDetailResponse(data model.User) UserResponse {
	return UserResponse{
		ID:        data.ID.String(),
		Email:     data.Email,
		Username:  data.Username,
		FullName:  data.FullName,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
}

func NewListUserResponse(users []model.User) []UserResponse {
	res := make([]UserResponse, len(users))
	for i, u := range users {
		res[i] = NewUserDetailResponse(u)
	}
	return res
}

type GetListUserRequest struct {
	PagingRequest
}

type GetUserDetailRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
}
