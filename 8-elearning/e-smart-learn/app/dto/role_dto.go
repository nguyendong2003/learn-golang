package dto

import (
	"elearning-api/model"
	"time"
)

type RoleResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewRoleResponse(role *model.Role) *RoleResponse {
	if role == nil {
		return nil
	}
	return &RoleResponse{
		ID:          role.ID.String(),
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}
