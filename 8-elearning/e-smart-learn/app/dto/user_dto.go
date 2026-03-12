package dto

import (
	"elearning-api/model"
	"time"
)

type UserResponse struct {
	ID        string        `json:"id"`
	Email     string        `json:"email"`
	Username  string        `json:"username"`
	Name      string        `json:"name"`
	Avatar    string        `json:"avatar"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	Role      *RoleResponse `json:"role,omitempty"`
}

func NewUserDetailResponse(data *model.User) *UserResponse {
	if data == nil {
		return nil
	}
	return &UserResponse{
		ID:        data.ID.String(),
		Email:     data.Email,
		Username:  data.Username,
		Name:      data.Name,
		Avatar:    data.Avatar,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		Role:      NewRoleResponse(data.Role),
	}
}

func NewListUserResponse(users []*model.User) []*UserResponse {
	res := make([]*UserResponse, len(users))
	for i, u := range users {
		res[i] = NewUserDetailResponse(u)
	}
	return res
}

type GetUserDetailRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	RoleID   string `json:"roleID" binding:"required,uuid"`
}

type UpdateUserRequest struct {
	Email    *string `json:"email,omitempty" binding:"omitempty,email"`
	Username *string `json:"username,omitempty" binding:"omitempty,min=3,max=50"`
	Name     *string `json:"name,omitempty"`
	Password *string `json:"password,omitempty" binding:"omitempty,min=6"`
}

type FilterUserRequest struct {
	PagingRequest

	Username *string `form:"username,omitempty"`
	Name     *string `form:"name,omitempty"`
}
