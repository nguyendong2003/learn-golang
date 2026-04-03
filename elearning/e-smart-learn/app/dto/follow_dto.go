package dto

import (
	"time"

	"elearning-api/model"
)

type FollowActionResponse struct {
	FollowerID  string `json:"follower_id"`
	FolloweeID  string `json:"followee_id"`
	IsFollowing bool   `json:"is_following"`
}

type FollowUserResponse struct {
	ID         string    `json:"id"`
	Username   string    `json:"username"`
	Name       string    `json:"name"`
	Avatar     string    `json:"avatar"`
	FollowedAt time.Time `json:"followed_at"`
}

func NewFollowUserResponse(user *model.User, followedAt time.Time) *FollowUserResponse {
	if user == nil {
		return nil
	}

	return &FollowUserResponse{
		ID:         user.ID.String(),
		Username:   user.Username,
		Name:       user.Name,
		Avatar:     user.Avatar,
		FollowedAt: followedAt,
	}
}
