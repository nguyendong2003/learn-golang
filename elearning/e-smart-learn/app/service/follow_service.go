package service

import (
	"context"
	"strings"

	"elearning-api/apperror"
	"elearning-api/dto"
	"elearning-api/model"
	"elearning-api/repository"

	"github.com/google/uuid"
)

type FollowService interface {
	FollowUser(ctx context.Context, followerID, followeeID uuid.UUID) (*dto.FollowActionResponse, error)
	UnfollowUser(ctx context.Context, followerID, followeeID uuid.UUID) (*dto.FollowActionResponse, error)
	GetFollowers(ctx context.Context, userID uuid.UUID, request dto.PagingRequest) ([]*dto.FollowUserResponse, int64, error)
	GetFollowings(ctx context.Context, userID uuid.UUID, request dto.PagingRequest) ([]*dto.FollowUserResponse, int64, error)
}

type followService struct {
	followRepository repository.FollowRepository
	userRepository   repository.UserRepository
}

func NewFollowService(
	followRepository repository.FollowRepository,
	userRepository repository.UserRepository,
) FollowService {
	return &followService{
		followRepository: followRepository,
		userRepository:   userRepository,
	}
}

func (s *followService) FollowUser(ctx context.Context, followerID, followeeID uuid.UUID) (*dto.FollowActionResponse, error) {
	if followerID == followeeID {
		return nil, apperror.NewBadRequestError("You cannot follow yourself")
	}

	if err := s.ensureUsersExist(ctx, followerID, followeeID); err != nil {
		return nil, err
	}

	isFollowing, err := s.followRepository.Exists(ctx, followerID, followeeID)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to check follow relationship")
	}

	if isFollowing {
		return nil, apperror.NewBadRequestError("Already followed this user")
	}

	_, err = s.followRepository.Create(ctx, &model.Follow{
		FollowerID: followerID,
		FolloweeID: followeeID,
	})
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to follow user")
	}

	return &dto.FollowActionResponse{
		FollowerID:  followerID.String(),
		FolloweeID:  followeeID.String(),
		IsFollowing: true,
	}, nil
}

func (s *followService) UnfollowUser(ctx context.Context, followerID, followeeID uuid.UUID) (*dto.FollowActionResponse, error) {
	if followerID == followeeID {
		return nil, apperror.NewBadRequestError("You cannot unfollow yourself")
	}

	if err := s.ensureUsersExist(ctx, followerID, followeeID); err != nil {
		return nil, err
	}

	isFollowing, err := s.followRepository.Exists(ctx, followerID, followeeID)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to check follow relationship")
	}

	if !isFollowing {
		return nil, apperror.NewBadRequestError("You are not following this user")
	}

	if err := s.followRepository.DeleteByPair(ctx, followerID, followeeID); err != nil {
		return nil, apperror.NewInternalServerError("Failed to unfollow user")
	}

	return &dto.FollowActionResponse{
		FollowerID:  followerID.String(),
		FolloweeID:  followeeID.String(),
		IsFollowing: false,
	}, nil
}

func (s *followService) GetFollowers(ctx context.Context, userID uuid.UUID, request dto.PagingRequest) ([]*dto.FollowUserResponse, int64, error) {
	user, err := s.userRepository.FindByID(ctx, userID, nil)
	if err != nil {
		return nil, 0, apperror.NewInternalServerError("Failed to get user")
	}
	if user == nil {
		return nil, 0, apperror.NewNotFoundError("User not found")
	}

	follows, total, err := s.followRepository.ListFollowers(ctx, userID, request.Limit, request.Offset, buildFollowSortQuery(request.SortBy, request.SortOrder))
	if err != nil {
		return nil, 0, apperror.NewInternalServerError("Failed to get followers")
	}

	data := make([]*dto.FollowUserResponse, len(follows))
	for i, follow := range follows {
		data[i] = dto.NewFollowUserResponse(follow.Follower, follow.CreatedAt)
	}

	return data, total, nil
}

func (s *followService) GetFollowings(ctx context.Context, userID uuid.UUID, request dto.PagingRequest) ([]*dto.FollowUserResponse, int64, error) {
	user, err := s.userRepository.FindByID(ctx, userID, nil)
	if err != nil {
		return nil, 0, apperror.NewInternalServerError("Failed to get user")
	}
	if user == nil {
		return nil, 0, apperror.NewNotFoundError("User not found")
	}

	follows, total, err := s.followRepository.ListFollowings(ctx, userID, request.Limit, request.Offset, buildFollowSortQuery(request.SortBy, request.SortOrder))
	if err != nil {
		return nil, 0, apperror.NewInternalServerError("Failed to get followings")
	}

	data := make([]*dto.FollowUserResponse, len(follows))
	for i, follow := range follows {
		data[i] = dto.NewFollowUserResponse(follow.Followee, follow.CreatedAt)
	}

	return data, total, nil
}

func buildFollowSortQuery(sortBy string, sortOrder string) string {
	defaultSort := "created_at DESC"

	if sortBy == "" {
		return defaultSort
	}

	if sortOrder == "" {
		sortOrder = "DESC"
	}

	allowedSort := map[string]bool{
		"created_at": true,
	}

	if !allowedSort[sortBy] {
		return defaultSort
	}

	return sortBy + " " + strings.ToUpper(sortOrder)
}

func (s *followService) ensureUsersExist(ctx context.Context, followerID, followeeID uuid.UUID) error {
	follower, err := s.userRepository.FindByID(ctx, followerID, nil)
	if err != nil {
		return apperror.NewInternalServerError("Failed to get follower")
	}
	if follower == nil {
		return apperror.NewNotFoundError("Follower not found")
	}

	followee, err := s.userRepository.FindByID(ctx, followeeID, nil)
	if err != nil {
		return apperror.NewInternalServerError("Failed to get user to follow")
	}
	if followee == nil {
		return apperror.NewNotFoundError("User to follow not found")
	}

	return nil
}
