package service

import (
	"context"

	"elearning-api/apperror"
	"elearning-api/dto"
	"elearning-api/model"
	"elearning-api/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetDetail(ctx context.Context, params dto.GetUserDetailRequest) (*dto.UserResponse, error)
	Create(ctx context.Context, data dto.CreateUserRequest) (*dto.UserResponse, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (s *userService) GetDetail(
	ctx context.Context,
	params dto.GetUserDetailRequest,
) (*dto.UserResponse, error) {

	userID, err := uuid.Parse(params.ID)
	if err != nil {
		return nil, apperror.NewBadRequestError("Invalid user ID format")
	}

	user, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, apperror.NewNotFoundError("User not found")
	}

	response := dto.NewUserDetailResponse(*user)
	return &response, nil
}

func (s *userService) Create(
	ctx context.Context,
	data dto.CreateUserRequest,
) (*dto.UserResponse, error) {

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(data.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	newUser := &model.User{
		Email:    data.Email,
		Username: data.Username,
		Password: string(hashedPassword),
		FullName: data.Username,
	}

	// Check email
	user, err := s.userRepository.GetByEmail(ctx, data.Email)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to check existing user", err)
	}
	if user != nil {
		return nil, apperror.NewBadRequestError("Email already exists")
	}

	// Check username
	user, err = s.userRepository.GetByUsername(ctx, data.Username)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to check existing user", err)
	}
	if user != nil {
		return nil, apperror.NewBadRequestError("Username already exists")
	}

	// Create user
	user, err = s.userRepository.Create(ctx, newUser)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to create user", err)
	}

	response := dto.NewUserDetailResponse(*user)

	return &response, nil
}
