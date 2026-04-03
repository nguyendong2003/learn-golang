package service

import (
	"context"
	"math"
	"strings"
	"time"

	"elearning-api/apperror"
	"elearning-api/dto"
	"elearning-api/model"
	"elearning-api/repository"
	"elearning-api/util"

	"github.com/google/uuid"
)

type InstructorProfileService interface {
	Create(ctx context.Context, userID uuid.UUID, request dto.CreateInstructorProfileRequest) (*dto.InstructorProfileResponse, error)
	Update(ctx context.Context, id uuid.UUID, request dto.UpdateInstructorProfileRequest) (*dto.InstructorProfileResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error

	GetByID(ctx context.Context, id uuid.UUID) (*dto.InstructorProfileResponse, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*dto.InstructorProfileResponse, error)
	GetList(ctx context.Context, request dto.ListInstructorProfileRequest) ([]*dto.InstructorProfileResponse, int64, error)
	GetGrowthStatistics(ctx context.Context) (*dto.TeacherGrowthStatisticsResponse, error)
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
	profile, err := s.instructorProfileRepository.Find(ctx, "user_id = ?", nil, userID)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to check existing instructor profile")
	}

	if profile != nil {
		return nil, apperror.NewBadRequestError("Instructor profile already exists for this user")
	}

	newInstructorProfile := &model.InstructorProfile{
		UserID:    userID,
		Bio:       request.Bio,
		Education: request.Education,
	}

	if request.CategoryID != nil {
		if cid, err := uuid.Parse(*request.CategoryID); err == nil {
			newInstructorProfile.CategoryID = cid
		}
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
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve instructor profile")
	}

	return dto.NewInstructorProfileDetailResponse(instructorProfile), nil
}

func (s *instructorProfileService) Update(ctx context.Context, id uuid.UUID, request dto.UpdateInstructorProfileRequest) (*dto.InstructorProfileResponse, error) {
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

	if request.CategoryID != nil {
		if cid, err := uuid.Parse(*request.CategoryID); err == nil {
			instructorProfile.CategoryID = cid
		}
	}

	updatedInstructorProfile, err := s.instructorProfileRepository.Updates(ctx, instructorProfile)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to update instructor profile")
	}

	return dto.NewInstructorProfileDetailResponse(updatedInstructorProfile), nil
}

func (s *instructorProfileService) Delete(ctx context.Context, id uuid.UUID) error {
	instructorProfile, err := s.instructorProfileRepository.FindByID(ctx, id, nil)
	if err != nil {
		return apperror.NewInternalServerError("Failed to retrieve instructor profile")
	}

	if instructorProfile == nil {
		return apperror.NewNotFoundError("Instructor Profile not found")
	}

	if err := s.instructorProfileRepository.Delete(ctx, id); err != nil {
		return apperror.NewInternalServerError("Failed to delete instructor profile")
	}

	return nil
}

func (s *instructorProfileService) GetByID(ctx context.Context, id uuid.UUID) (*dto.InstructorProfileResponse, error) {
	instructorProfile, err := s.instructorProfileRepository.FindByID(ctx, id, []repository.Preload{
		repository.PreloadPath(repository.User),
		repository.PreloadPath(repository.Category),
	})
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve instructor profile")
	}

	if instructorProfile == nil {
		return nil, apperror.NewNotFoundError("Instructor profile not found")
	}

	return dto.NewInstructorProfileDetailResponse(instructorProfile), nil
}

func (s *instructorProfileService) GetByUserID(ctx context.Context, userID uuid.UUID) (*dto.InstructorProfileResponse, error) {
	instructorProfile, err := s.instructorProfileRepository.Find(ctx, "user_id = ?", []repository.Preload{
		repository.PreloadPath(repository.User),
		repository.PreloadPath(repository.Category),
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

	instructorProfiles, total, err := s.instructorProfileRepository.ListWithJoins(
		ctx,
		request.Limit,
		request.Offset,
		orderQuery,
		query,
		[]repository.Join{repository.UserJoin},
		[]repository.Preload{
			repository.PreloadPath(repository.User),
			repository.PreloadPath(repository.Category),
		},
		args...,
	)
	if err != nil {
		return nil, 0, apperror.NewInternalServerError("Failed to retrieve instructor profiles")
	}

	return dto.NewInstructorProfileListResponse(instructorProfiles), total, nil
}

func (s *instructorProfileService) GetGrowthStatistics(ctx context.Context) (*dto.TeacherGrowthStatisticsResponse, error) {
	now := time.Now()
	quarterStart := startOfQuarter(now)
	quarterEnd := quarterStart.AddDate(0, 3, 0)

	stats, err := s.instructorProfileRepository.GetGrowthStatistics(ctx, quarterStart, quarterEnd)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve teacher growth statistics")
	}

	var topCategory *dto.TeacherGrowthTopCategoryResponse
	if stats.TopCategoryName != "" && stats.TopCategoryCount > 0 {
		sharePct := 0.0
		if stats.TotalVerifiedTeachers > 0 {
			sharePct = math.Round((float64(stats.TopCategoryCount) / float64(stats.TotalVerifiedTeachers)) * 100)
		}

		topCategory = &dto.TeacherGrowthTopCategoryResponse{
			Name:     stats.TopCategoryName,
			Count:    stats.TopCategoryCount,
			SharePct: sharePct,
		}
	}

	return &dto.TeacherGrowthStatisticsResponse{
		TotalVerifiedTeachers: stats.TotalVerifiedTeachers,
		NewThisQuarter:        stats.NewThisQuarter,
		TopCategory:           topCategory,
	}, nil
}

func startOfQuarter(now time.Time) time.Time {
	quarter := (int(now.Month())-1)/3 + 1
	startMonth := time.Month((quarter-1)*3 + 1)
	return time.Date(now.Year(), startMonth, 1, 0, 0, 0, 0, now.Location())
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
		"name": request.Name,
	}

	for column, value := range filters {
		util.AddILIKECondition(&conditions, &args, column, value)
	}

	query := strings.Join(conditions, " AND ")

	return query, args
}
