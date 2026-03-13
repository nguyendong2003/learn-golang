package service

import (
	"context"
	"elearning-api/apperror"
	"elearning-api/dto"
	"elearning-api/model"
	"elearning-api/repository"
	"elearning-api/util"
	"strings"

	"github.com/google/uuid"
)

type InstructorProfileService interface {
	Create(ctx context.Context, userID uuid.UUID, request dto.CreateInstructorProfileRequest) (*dto.InstructorProfileResponse, error)
	Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, request dto.UpdateInstructorProfileRequest) (*dto.InstructorProfileResponse, error)
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error

	GetByID(ctx context.Context, id uuid.UUID) (*dto.InstructorProfileResponse, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*dto.InstructorProfileResponse, error)
	GetList(ctx context.Context, request dto.ListInstructorProfileRequest) ([]*dto.InstructorProfileResponse, int64, error)
}

type instructorProfileService struct {
	instructorProfileRepository repository.InstructorProfileRepository
}

func NewInstructorProfileService(instructorProfileRepository repository.InstructorProfileRepository) InstructorProfileService {
	return &instructorProfileService{
		instructorProfileRepository: instructorProfileRepository,
	}
}

func (s *instructorProfileService) Create(ctx context.Context, userID uuid.UUID, request dto.CreateInstructorProfileRequest) (*dto.InstructorProfileResponse, error) {
	newInstructorProfile := &model.InstructorProfile{
		UserID:    userID,
		Bio:       request.Bio,
		Education: request.Education,
	}

	if request.LinkedinURL != nil {
		newInstructorProfile.LinkedinURL = *request.LinkedinURL
	}

	if request.YoutubeURL != nil {
		newInstructorProfile.YoutubeURL = *request.YoutubeURL
	}

	if request.InstagramURL != nil {
		newInstructorProfile.InstagramURL = *request.InstagramURL
	}

	createdInstructorProfile, err := s.instructorProfileRepository.Create(ctx, newInstructorProfile)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to create instructorProfile")
	}

	instructorProfile, err := s.instructorProfileRepository.FindByID(ctx, createdInstructorProfile.ID, []repository.Preload{
		repository.PreloadPath(repository.User),
	})

	return dto.NewInstructorProfileDetailResponse(instructorProfile), nil
}

func (s *instructorProfileService) Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, request dto.UpdateInstructorProfileRequest) (*dto.InstructorProfileResponse, error) {
	instructorProfile, err := s.instructorProfileRepository.FindByID(ctx, id, []repository.Preload{
		repository.PreloadPath(repository.User),
	})

	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve instructor profile")
	}

	if instructorProfile == nil {
		return nil, apperror.NewNotFoundError("Instructor profile not found")
	}

	if request.Bio != nil {
		instructorProfile.Bio = *request.Bio
	}

	if request.Education != nil {
		instructorProfile.Education = *request.Education
	}

	if request.LinkedinURL != nil {
		instructorProfile.LinkedinURL = *request.LinkedinURL
	}

	if request.YoutubeURL != nil {
		instructorProfile.YoutubeURL = *request.YoutubeURL
	}

	if request.InstagramURL != nil {
		instructorProfile.InstagramURL = *request.InstagramURL
	}

	updatedInstructorProfile, err := s.instructorProfileRepository.Updates(ctx, instructorProfile)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to update instructor profile")
	}

	return dto.NewInstructorProfileDetailResponse(updatedInstructorProfile), nil
}

func (s *instructorProfileService) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	instructorProfile, err := s.instructorProfileRepository.FindByID(ctx, id, nil)
	if err != nil {
		return apperror.NewInternalServerError("Failed to retrieve instructorProfile")
	}

	if instructorProfile == nil {
		return apperror.NewNotFoundError("InstructorProfile not found")
	}

	if err := s.instructorProfileRepository.Delete(ctx, id); err != nil {
		return apperror.NewInternalServerError("Failed to delete instructorProfile")
	}

	return nil
}

func (s *instructorProfileService) GetByID(ctx context.Context, id uuid.UUID) (*dto.InstructorProfileResponse, error) {
	instructorProfile, err := s.instructorProfileRepository.FindByID(ctx, id, []repository.Preload{
		repository.PreloadPath(repository.User),
	})

	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve instructorProfile")
	}

	if instructorProfile == nil {
		return nil, apperror.NewNotFoundError("Instructor profile not found")
	}

	return dto.NewInstructorProfileDetailResponse(instructorProfile), nil
}

func (s *instructorProfileService) GetByUserID(ctx context.Context, userID uuid.UUID) (*dto.InstructorProfileResponse, error) {
	instructorProfile, err := s.instructorProfileRepository.Find(ctx, "user_id = ?", []repository.Preload{
		repository.PreloadPath(repository.User),
	}, userID)

	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve instructorProfile")
	}

	if instructorProfile == nil {
		return nil, apperror.NewNotFoundError("Instructor profile not found")
	}

	return dto.NewInstructorProfileDetailResponse(instructorProfile), nil
}

func (s *instructorProfileService) GetList(ctx context.Context, request dto.ListInstructorProfileRequest) ([]*dto.InstructorProfileResponse, int64, error) {
	// Build sort query
	orderQuery := buildInstructorProfileSortQuery(request.SortBy, request.SortOrder)

	// Build filter query
	query, args := buildInstructorProfileQuery(request)

	categories, total, err := s.instructorProfileRepository.List(
		ctx,
		request.Limit,
		request.Offset,
		orderQuery,
		query,
		nil,
		args...,
	)

	if err != nil {
		return nil, 0, apperror.NewInternalServerError("Failed to retrieve categories")
	}

	return dto.NewInstructorProfileListResponse(categories), total, nil
}

func buildInstructorProfileSortQuery(sortBy string, sortOrder string) string {
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
	}

	if !allowedSort[sortBy] {
		return defaultSort
	}

	return sortBy + " " + strings.ToUpper(sortOrder)
}

func buildInstructorProfileQuery(request dto.ListInstructorProfileRequest) (string, []any) {
	var conditions []string
	var args []any

	filters := map[string]*string{
		"User.name": request.Name,
	}

	for column, value := range filters {
		util.AddILIKECondition(&conditions, &args, column, value)
	}

	query := strings.Join(conditions, " AND ")

	return query, args
}
