package service

import (
	"context"
	"strings"

	"elearning-api/apperror"
	"elearning-api/dto"
	"elearning-api/model"
	"elearning-api/repository"
	"elearning-api/util"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Create(ctx context.Context, data dto.CreateUserRequest) (*dto.UserResponse, error)
	Update(ctx context.Context, id uuid.UUID, data dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error

	GetByID(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error)
	GetList(ctx context.Context, limit, offset int) ([]*dto.UserResponse, int64, error)
	FilterAndPaginateAndSort(ctx context.Context, filter dto.FilterUserRequest) ([]*dto.UserResponse, int64, error)
}

type userService struct {
	userRepository repository.UserRepository
	roleRepository repository.RoleRepository
}

func NewUserService(
	userRepository repository.UserRepository,
	roleRepository repository.RoleRepository) UserService {
	return &userService{
		userRepository: userRepository,
		roleRepository: roleRepository,
	}
}

func (s *userService) Create(ctx context.Context, request dto.CreateUserRequest) (*dto.UserResponse, error) {
	existingUser, err := s.userRepository.GetByEmailOrUsername(ctx, request.Email, request.Username)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to check existing user")
	}

	if existingUser != nil {
		validationErrors := map[string]string{}

		if existingUser.Email == request.Email {
			validationErrors["email"] = "Email already exists"
		}

		if existingUser.Username == request.Username {
			validationErrors["username"] = "Username already exists"
		}

		return nil, apperror.NewValidationError(validationErrors)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to hash password")
	}

	roleID, err := uuid.Parse(request.RoleID)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to parse uuid")
	}

	role, err := s.roleRepository.FindByID(ctx, roleID)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get role")
	}

	if role == nil {
		return nil, apperror.NewNotFoundError("Role not found")
	}

	newUser := &model.User{
		Email:    request.Email,
		Username: request.Username,
		Password: string(hashedPassword),
		Name:     request.Username,
		Role:     role,
	}

	createdUser, err := s.userRepository.Create(ctx, newUser)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to create user")
	}

	return dto.NewUserDetailResponse(createdUser), nil
}

func (s *userService) Update(ctx context.Context, id uuid.UUID, request dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.userRepository.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get user")
	}

	if user == nil {
		return nil, apperror.NewNotFoundError("User not found")
	}

	updates := map[string]any{}

	if request.Email != nil {
		updates["email"] = *request.Email
	}

	if request.Username != nil {
		updates["username"] = *request.Username
	}

	if request.Name != nil {
		updates["name"] = *request.Name
	}

	if request.Password != nil {

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*request.Password), bcrypt.DefaultCost)

		if err != nil {
			return nil, apperror.NewInternalServerError("Failed to hash password")
		}

		updates["password"] = string(hashedPassword)
	}

	updatedUser, err := s.userRepository.Update(ctx, id, updates)

	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to update user")
	}

	return dto.NewUserDetailResponse(updatedUser), nil
}

func (s *userService) DeleteByID(ctx context.Context, id uuid.UUID) error {
	user, err := s.userRepository.FindByID(ctx, id)
	if err != nil {
		return apperror.NewInternalServerError("Failed to get user")
	}

	if user == nil {
		return apperror.NewNotFoundError("User not found")
	}

	if err := s.userRepository.Delete(ctx, id); err != nil {
		return apperror.NewInternalServerError("Failed to delete user")
	}

	return nil
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepository.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get user detail")
	}

	if user == nil {
		return nil, apperror.NewNotFoundError("User not found")
	}

	return dto.NewUserDetailResponse(user), nil
}

func (s *userService) GetList(ctx context.Context, limit int, offset int) ([]*dto.UserResponse, int64, error) {
	// default pagination
	if limit <= 0 {
		limit = 10
	}

	if offset < 0 {
		offset = 0
	}

	users, total, err := s.userRepository.List(
		ctx,
		limit,
		offset,
		"created_at DESC",
		"",
	)

	if err != nil {
		return nil, 0, apperror.NewInternalServerError("Failed to get users")
	}

	return dto.NewListUserResponse(users), total, nil
}

func (s *userService) FilterAndPaginateAndSort(ctx context.Context, filter dto.FilterUserRequest) ([]*dto.UserResponse, int64, error) {

	limit := filter.Limit
	offset := filter.Offset

	// Build sort query
	orderQuery := buildSortQuery(filter.SortBy, filter.SortOrder)

	// Build filter query
	var conditions []string
	var args []any

	filters := map[string]*string{
		"username": filter.Username,
		"name":     filter.Name,
	}

	for column, value := range filters {
		util.AddILIKECondition(&conditions, &args, column, value)
	}

	query := strings.Join(conditions, " AND ")

	// call repository
	users, total, err := s.userRepository.List(
		ctx,
		limit,
		offset,
		orderQuery,
		query,
		args...,
	)

	if err != nil {
		return nil, 0, apperror.NewInternalServerError("Failed to get users")
	}

	return dto.NewListUserResponse(users), total, nil
}

func buildSortQuery(sortBy string, sortOrder string) string {
	defaultSort := "created_at DESC"

	if sortBy == "" {
		return defaultSort
	}

	if sortOrder == "" {
		sortOrder = "DESC"
	}

	allowedSort := map[string]bool{
		"created_at": true,
		"updated_at": true,
		"username":   true,
		"email":      true,
		"name":       true,
	}

	if !allowedSort[sortBy] {
		return defaultSort
	}

	return sortBy + " " + strings.ToUpper(sortOrder)
}
