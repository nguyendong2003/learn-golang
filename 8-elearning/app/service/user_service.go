package service

import (
	"context"
	"fmt"
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
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (s *userService) Create(
	ctx context.Context,
	data dto.CreateUserRequest,
) (*dto.UserResponse, error) {

	existingUser, err := s.userRepository.GetByEmailOrUsername(
		ctx,
		data.Email,
		data.Username,
	)

	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to check existing user", err)
	}

	if existingUser != nil {

		validationErrors := map[string]string{}

		if existingUser.Email == data.Email {
			validationErrors["email"] = "Email already exists"
		}

		if existingUser.Username == data.Username {
			validationErrors["username"] = "Username already exists"
		}

		return nil, apperror.NewValidationError(validationErrors)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(data.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to hash password", err)
	}

	newUser := &model.User{
		Email:    data.Email,
		Username: data.Username,
		Password: string(hashedPassword),
		Name:     data.Username,
	}

	createdUser, err := s.userRepository.Create(ctx, newUser)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to create user", err)
	}

	return dto.NewUserDetailResponse(createdUser), nil
}

func (s *userService) Update(
	ctx context.Context,
	id uuid.UUID,
	data dto.UpdateUserRequest,
) (*dto.UserResponse, error) {

	user, err := s.userRepository.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get user", err)
	}

	if user == nil {
		return nil, apperror.NewNotFoundError("User not found")
	}

	updates := map[string]any{}

	if data.Email != nil {
		updates["email"] = *data.Email
	}

	if data.Username != nil {
		updates["username"] = *data.Username
	}

	if data.Name != nil {
		updates["name"] = *data.Name
	}

	if data.Password != nil {

		hashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(*data.Password),
			bcrypt.DefaultCost,
		)

		if err != nil {
			return nil, apperror.NewInternalServerError("Failed to hash password", err)
		}

		updates["password"] = string(hashedPassword)
	}

	updatedUser, err := s.userRepository.Update(ctx, id, updates)

	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to update user", err)
	}

	return dto.NewUserDetailResponse(updatedUser), nil
}

func (s *userService) DeleteByID(
	ctx context.Context,
	id uuid.UUID,
) error {

	user, err := s.userRepository.FindByID(ctx, id)
	if err != nil {
		return apperror.NewInternalServerError("Failed to get user", err)
	}

	if user == nil {
		return apperror.NewNotFoundError("User not found")
	}

	if err := s.userRepository.Delete(ctx, id); err != nil {
		return apperror.NewInternalServerError("Failed to delete user", err)
	}

	return nil
}

func (s *userService) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*dto.UserResponse, error) {

	user, err := s.userRepository.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get user detail", err)
	}

	if user == nil {
		return nil, apperror.NewNotFoundError("User not found")
	}

	return dto.NewUserDetailResponse(user), nil
}

func (s *userService) GetList(
	ctx context.Context,
	limit int,
	offset int,
) ([]*dto.UserResponse, int64, error) {

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
		return nil, 0, apperror.NewInternalServerError("Failed to get users", err)
	}

	return dto.NewListUserResponse(users), total, nil
}

func (s *userService) FilterAndPaginateAndSort(
	ctx context.Context,
	filter dto.FilterUserRequest,
) ([]*dto.UserResponse, int64, error) {

	limit := filter.Limit
	offset := filter.GetOffset()

	// Build sort query
	orderQuery := buildSortQuery(filter.Sort)

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
		return nil, 0, apperror.NewInternalServerError("Failed to get users", err)
	}

	return dto.NewListUserResponse(users), total, nil
}

func buildSortQuery(sort *string) string {

	defaultSort := "created_at DESC"

	if sort == nil || *sort == "" {
		return defaultSort
	}

	allowedSort := map[string]bool{
		"created_at": true,
		"updated_at": true,
		"username":   true,
		"email":      true,
		"name":       true,
	}

	var orders []string

	items := strings.Split(*sort, ",")

	for _, item := range items {

		parts := strings.Split(item, ":")

		column := parts[0]
		direction := "ASC"

		if len(parts) > 1 {
			direction = strings.ToUpper(parts[1])
		}

		if !allowedSort[column] {
			continue
		}

		if direction != "ASC" && direction != "DESC" {
			direction = "ASC"
		}

		orders = append(orders, fmt.Sprintf("%s %s", column, direction))
	}

	if len(orders) == 0 {
		return defaultSort
	}

	return strings.Join(orders, ", ")
}
