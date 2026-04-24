package service

import (
	"context"
	"math"
	"slices"
	"time"

	"elearning-api/apperror"
	"elearning-api/consts"
	"elearning-api/dto"
	"elearning-api/model"
	"elearning-api/pkg"
	"elearning-api/repository"
	"elearning-api/util"

	"github.com/google/uuid"
)

type UserService interface {
	GetUserWithRole(ctx context.Context, id uuid.UUID) (*model.User, error)

	GetByID(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error)
	UpdateUser(ctx context.Context, id uuid.UUID, request dto.UpdateUserRequest) (*dto.UserResponse, error)
	GetList(ctx context.Context, limit, offset int) ([]*dto.UserListResponse, int64, error)

	// Blogs interactions
	MarkBlogAsRead(ctx context.Context, userID uuid.UUID, request dto.UpdateReadBlogHistoryRequest) (*dto.UserResponse, error)
	GetReadBlogs(ctx context.Context, userID uuid.UUID, request dto.ViewReadBlogHistoryRequest) (*dto.GetReadBlogsResponse, error)

	// Instructor application
	ApplyToInstructor(ctx context.Context, userID uuid.UUID, request dto.ApplyInstructorRequest) error
	GetPendingInstructorApplications(ctx context.Context, request dto.PagingRequest) ([]*dto.InstructorApplicationResponse, error)
	UpdateInstructorApplicationStatus(ctx context.Context, applicationID string, status consts.InstructorProfileStatus) (*dto.UserResponse, error)
	GetStatistics(ctx context.Context) (*dto.UserStatisticsResponse, error)
	GetActiveStudentStatistics(ctx context.Context) (*dto.ActiveStudentsStatisticsResponse, error)
}

type userService struct {
	db                          repository.DbRepository
	userRepository              repository.UserRepository
	roleRepository              repository.RoleRepository
	blogRepository              repository.BlogRepository
	instructorProfileRepository repository.InstructorProfileRepository
	enrollmentRepository        repository.EnrollmentRepository
	storageProvider             pkg.StorageProvider
	categoryRepository          repository.CategoryRepository
}

func NewUserService(
	db repository.DbRepository,
	userRepository repository.UserRepository,
	roleRepository repository.RoleRepository,
	blogRepository repository.BlogRepository,
	instructorProfileRepository repository.InstructorProfileRepository,
	enrollmentRepository repository.EnrollmentRepository,
	storageProvider pkg.StorageProvider,
	categoryRepository repository.CategoryRepository,
) UserService {
	return &userService{
		db:                          db,
		userRepository:              userRepository,
		roleRepository:              roleRepository,
		blogRepository:              blogRepository,
		instructorProfileRepository: instructorProfileRepository,
		enrollmentRepository:        enrollmentRepository,
		storageProvider:             storageProvider,
		categoryRepository:          categoryRepository,
	}
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepository.FindByID(ctx, id,
		[]repository.Preload{
			repository.Role,
			repository.PreloadPath(repository.InstructorProfile, repository.Category),
		})
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get user detail")
	}

	if user == nil {
		return nil, apperror.NewNotFoundError("User not found")
	}
	totalCourseEnrolled, totalCourseCompleted, err := s.enrollmentRepository.GetCourseEnrollmentCount(ctx, user.ID)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get user enrollment stats")
	}
	userResponse := dto.NewUserDetailResponse(user)
	userResponse.TotalCoursesEnrolled = totalCourseEnrolled
	userResponse.TotalCourseInProgress = totalCourseEnrolled - totalCourseCompleted

	return userResponse, nil
}

func (s *userService) UpdateUser(ctx context.Context, id uuid.UUID, request dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.userRepository.FindByID(ctx, id, nil)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get user")
	}

	if user == nil {
		return nil, apperror.NewNotFoundError("User not found")
	}
	user.Name = request.Name
	user.Email = request.Email
	user.PhoneNumber = request.PhoneNumber
	user.Address = request.Address
	if user.Avatar != request.Avatar {
		if user.Avatar != "" {
			if err := s.storageProvider.Delete(ctx, user.Avatar); err != nil {
				logger := util.LoggerFromContextWithLayer(ctx, util.LayerService)
				logger.Warn("failed to delete old avatar", "user_id", user.ID.String(), "avatar", user.Avatar, "error", err)
			}
		}
		user.Avatar = request.Avatar
	}
	if _, err := s.userRepository.Updates(ctx, user); err != nil {
		return nil, apperror.NewInternalServerError("Failed to update user")
	}
	return dto.NewUserDetailResponse(user), nil
}

func (s *userService) GetList(ctx context.Context, limit int, offset int) ([]*dto.UserListResponse, int64, error) {
	rows, total, err := s.userRepository.GetList(ctx, limit, offset)
	if err != nil {
		return nil, 0, apperror.NewInternalServerError("Failed to get users")
	}

	result := make([]*dto.UserListResponse, len(rows))
	for i, row := range rows {
		if row == nil {
			continue
		}

		userStatus := "inactive"
		if row.HasPendingInstructorApplication {
			userStatus = "pending_application"
		} else if row.IsActive {
			userStatus = "active"
		}

		overallProgressPct := 0.0
		if row.TotalLessons > 0 {
			overallProgressPct = (float64(row.CompletedLessons) / float64(row.TotalLessons)) * 100
		}

		result[i] = &dto.UserListResponse{
			ID:                              row.UserID,
			Name:                            row.Name,
			Email:                           row.Email,
			Avatar:                          row.Avatar,
			Role:                            row.RoleName,
			Status:                          userStatus,
			HasPendingInstructorApplication: row.HasPendingInstructorApplication,
			ActiveCourses:                   row.ActiveCourses,
			TotalCoursesTaught:              row.TotalCoursesTaught,
			CompletedLessons:                row.CompletedLessons,
			TotalLessons:                    row.TotalLessons,
			OverallProgressPct:              math.Round(overallProgressPct),
		}
	}

	return result, total, nil
}

func (s *userService) GetUserWithRole(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user, err := s.userRepository.FindByID(ctx, id, []repository.Preload{repository.Role, repository.PreloadPath(repository.Role, repository.Permissions)})
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get user detail")
	}

	if user == nil {
		return nil, apperror.NewNotFoundError("User not found")
	}

	return user, nil
}

func (s *userService) GetReadBlogs(ctx context.Context, userID uuid.UUID, request dto.ViewReadBlogHistoryRequest) (*dto.GetReadBlogsResponse, error) {
	user, err := s.userRepository.FindByID(ctx, userID, nil)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get user")
	}
	if user == nil {
		return nil, apperror.NewNotFoundError("User not found")
	}

	blogIDs := user.ReadBlogIDs

	start := request.Offset
	end := request.Offset + request.Limit

	if start > len(blogIDs) {
		start = len(blogIDs)
	}

	if end > len(blogIDs) {
		end = len(blogIDs)
	}

	paginatedBlogIDs := blogIDs[start:end]

	var readBlogs []dto.BlogResponse
	for _, blogID := range paginatedBlogIDs {
		blog, err := s.blogRepository.FindByID(ctx, blogID, nil)
		if err != nil {
			return nil, apperror.NewInternalServerError("Failed to get blog")
		}
		if blog != nil {
			readBlogs = append(readBlogs, *dto.NewBlogDetailResponse(blog))
		}
	}

	return &dto.GetReadBlogsResponse{
		ReadBlogs: readBlogs,
		Total:     int(len(blogIDs)),
	}, nil
}

func (s *userService) MarkBlogAsRead(ctx context.Context, userID uuid.UUID, request dto.UpdateReadBlogHistoryRequest) (*dto.UserResponse, error) {
	user, err := s.userRepository.FindByID(ctx, userID, nil)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get user")
	}

	if user == nil {
		return nil, apperror.NewNotFoundError("User not found")
	}
	parsedBlogID, err := uuid.Parse(request.BlogID)
	if err != nil {
		return nil, apperror.NewBadRequestError("Invalid blog ID")
	}

	var existingBLog *model.Blog
	if existingBLog, err = s.blogRepository.FindByID(ctx, parsedBlogID, nil); err != nil {
		return nil, apperror.NewInternalServerError("Failed to check blog existence")
	} else if existingBLog == nil {
		return nil, apperror.NewNotFoundError("Blog not found")
	}

	// Increment blog view
	if err := s.blogRepository.IncrementView(ctx, parsedBlogID); err != nil {
		return nil, apperror.NewInternalServerError("Failed to increment blog view")
	}

	// Check if the blog ID is already in the user's read list -> do nothing
	if slices.Contains(user.ReadBlogIDs, parsedBlogID) {
		return dto.NewUserDetailResponse(user), nil
	}

	user.ReadBlogIDs = append(user.ReadBlogIDs, parsedBlogID)
	if _, err := s.userRepository.Updates(ctx, user); err != nil {
		return nil, apperror.NewInternalServerError("Failed to update user")
	}

	return dto.NewUserDetailResponse(user), nil
}

func (s *userService) ApplyToInstructor(ctx context.Context, userID uuid.UUID, request dto.ApplyInstructorRequest) error {
	user, err := s.userRepository.FindByID(ctx, userID, []repository.Preload{repository.Role})
	if err != nil {
		return apperror.NewInternalServerError("Failed to get user")
	}
	if user == nil {
		return apperror.NewNotFoundError("User not found")
	}
	if user.Role.Name == string(consts.RoleInstructor) {
		return apperror.NewBadRequestError("User is already an instructor")
	}

	instructorRole, err := s.roleRepository.GetByName(ctx, string(consts.RoleInstructor))
	if err != nil {
		return apperror.NewInternalServerError("Failed to get instructor role")
	}
	if instructorRole == nil {
		return apperror.NewNotFoundError("Instructor role not found")
	}
	categoryIDParse := uuid.MustParse(*request.CategoryID)
	category, err := s.categoryRepository.FindByID(ctx, categoryIDParse, nil)
	if err != nil {
		return apperror.NewInternalServerError("Failed to get category")
	}
	if category == nil {
		return apperror.NewNotFoundError("Category not found")
	}

	err = s.db.Transaction(ctx, func(txDb repository.DbRepository) error {
		// Create instructor profile record with pending status
		if _, err := s.instructorProfileRepository.Create(
			ctx,
			&model.InstructorProfile{
				UserID:            user.ID,
				Bio:               request.Bio,
				Education:         request.Education,
				YearsOfExperience: request.YearsOfExperience,
				LinkedinURL:       request.LinkedinURL,
				YoutubeURL:        request.YoutubeURL,
				InstagramURL:      request.InstagramURL,
				CVURL:             request.CVURL,
				PortfolioURL:      request.PortfolioURL,
				Certifications:    request.Certifications,
				Status:            consts.InstructorProfilePending,
				CategoryID:        categoryIDParse,
			}); err != nil {
			return apperror.NewInternalServerError("Failed to create instructor profile")
		}
		return nil
	})

	return err
}

func (s *userService) GetPendingInstructorApplications(ctx context.Context, request dto.PagingRequest) ([]*dto.InstructorApplicationResponse, error) {
	applications, _, err := s.instructorProfileRepository.List(
		ctx,
		request.Limit,
		request.Offset,
		"created_at DESC",
		"status = ?",
		[]repository.Preload{repository.User},
		consts.InstructorProfilePending,
	)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get instructor applications")
	}

	return dto.NewInstructorApplicationListResponse(applications), nil
}

func (s *userService) UpdateInstructorApplicationStatus(ctx context.Context, applicationID string, status consts.InstructorProfileStatus) (*dto.UserResponse, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to start transaction")
	}

	txInstructorProfileRepo := repository.NewInstructorProfileRepository(tx)
	txRoleRepo := repository.NewRoleRepository(tx)
	txUserRepo := repository.NewUserRepository(tx)

	application, err := txInstructorProfileRepo.FindByID(ctx, uuid.MustParse(applicationID), []repository.Preload{repository.User})
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get instructor application")
	}
	if application == nil {
		return nil, apperror.NewNotFoundError("Instructor application not found")
	}

	if application.Status != consts.InstructorProfilePending {
		return nil, apperror.NewBadRequestError("Only pending applications can be updated")
	}
	requestedStatusParse := consts.InstructorProfileStatus(status)
	application.Status = requestedStatusParse

	if _, err := txInstructorProfileRepo.Update(ctx, application.ID, map[string]any{"status": status}); err != nil {
		_ = tx.Rollback()
		return nil, apperror.NewInternalServerError("Failed to update instructor application status")
	}

	// If approved, also update user role to instructor
	if requestedStatusParse == consts.InstructorProfileApproved && application.User != nil {
		instructorRole, err := txRoleRepo.GetByName(ctx, string(consts.RoleInstructor))
		if err != nil {
			return nil, apperror.NewInternalServerError("Failed to get instructor role")
		}
		if instructorRole == nil {
			_ = tx.Rollback()
			return nil, apperror.NewNotFoundError("Instructor role not found")
		}

		application.User.Role = instructorRole
		if _, err := txUserRepo.Updates(ctx, application.User); err != nil {
			_ = tx.Rollback()
			return nil, apperror.NewInternalServerError("Failed to update user role")
		}
	}
	user, err := txUserRepo.FindByID(ctx, application.UserID, []repository.Preload{repository.Role, repository.InstructorProfile})
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get user")
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return dto.NewUserDetailResponse(user), nil
}

func (s *userService) GetStatistics(ctx context.Context) (*dto.UserStatisticsResponse, error) {
	totalActiveUsers, err := s.userRepository.Count(ctx, "is_active = ?", true)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to count active users")
	}

	now := time.Now()
	startCurrentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	totalActiveUsersLastMonth, err := s.userRepository.Count(ctx, "is_active = ? AND created_at < ?", true, startCurrentMonth)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to count last month active users")
	}

	activeUsersGrowthPct := 0.0
	if totalActiveUsersLastMonth > 0 {
		activeUsersGrowthPct = (float64(totalActiveUsers-totalActiveUsersLastMonth) / float64(totalActiveUsersLastMonth)) * 100
	} else if totalActiveUsers > 0 {
		activeUsersGrowthPct = 100
	}

	pendingInstructors, err := s.instructorProfileRepository.Count(ctx, "status = ?", consts.InstructorProfilePending)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to count pending instructors")
	}

	totalCompletedLessons, totalLessons, err := s.enrollmentRepository.GetCourseCompletionTotals(ctx)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to calculate course completion rate")
	}

	courseCompletionRate := 0.0
	if totalLessons > 0 {
		courseCompletionRate = (float64(totalCompletedLessons) / float64(totalLessons)) * 100
	}

	return &dto.UserStatisticsResponse{
		TotalActiveUsers:        totalActiveUsers,
		ActiveUsersGrowthPct:    math.Round(activeUsersGrowthPct),
		PendingInstructors:      pendingInstructors,
		CourseCompletionRatePct: math.Round(courseCompletionRate),
	}, nil
}

func (s *userService) GetActiveStudentStatistics(ctx context.Context) (*dto.ActiveStudentsStatisticsResponse, error) {
	totalActiveStudents, err := s.userRepository.CountActiveByRoleName(ctx, string(consts.RoleStudent))
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to count active students")
	}

	now := time.Now()
	startCurrentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	totalActiveStudentsLastMonth, err := s.userRepository.CountActiveByRoleNameBefore(ctx, string(consts.RoleStudent), startCurrentMonth)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to count last month active students")
	}

	growthPct := 0.0
	if totalActiveStudentsLastMonth > 0 {
		growthPct = (float64(totalActiveStudents-totalActiveStudentsLastMonth) / float64(totalActiveStudentsLastMonth)) * 100
	} else if totalActiveStudents > 0 {
		growthPct = 100
	}

	roundedGrowthPct := math.Round(growthPct*100) / 100

	return &dto.ActiveStudentsStatisticsResponse{
		TotalActiveStudents:     totalActiveStudents,
		ActiveStudentsGrowthPct: roundedGrowthPct,
	}, nil
}
