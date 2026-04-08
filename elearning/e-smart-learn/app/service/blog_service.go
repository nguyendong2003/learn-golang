package service

import (
	"context"
	"math"
	"strings"
	"time"

	"elearning-api/apperror"
	"elearning-api/consts"
	"elearning-api/dto"
	"elearning-api/job"
	"elearning-api/model"
	"elearning-api/repository"
	"elearning-api/util"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

type BlogService interface {
	GetPublishedBlogs(ctx context.Context, filter dto.SearchBlogRequest) ([]*dto.BlogResponse, int64, error)
	GetBlogBySlug(ctx context.Context, slug string, userId uuid.UUID) (*dto.BlogResponse, error)
	GetBlogs(ctx context.Context, filter dto.SearchBlogRequest) ([]*dto.BlogResponse, int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dto.BlogResponse, error)
	Create(ctx context.Context, creatorId uuid.UUID, data dto.CreateBlogRequest) (*dto.BlogResponse, error)
	DeleteBlog(ctx context.Context, authorId, blogId uuid.UUID) error
	UpdateBlogs(ctx context.Context, authorId, blogId uuid.UUID, data dto.UpdateBlogRequest) (*dto.BlogResponse, error)
	GetBlogStatistics(ctx context.Context) (*dto.BlogStatisticsResponse, error)
}

type blogService struct {
	db                 repository.DbRepository
	blogRepository     repository.BlogRepository
	categoryRepository repository.CategoryRepository
	userRepository     repository.UserRepository
	uploadService      UploadService
	asynqClient        *asynq.Client
	followRepository   repository.FollowRepository
}

func NewBlogService(
	db repository.DbRepository,
	blogRepository repository.BlogRepository,
	categoryRepository repository.CategoryRepository,
	userRepository repository.UserRepository,
	asynqClient *asynq.Client,
	followRepository repository.FollowRepository,
	uploadService UploadService,
) BlogService {
	return &blogService{
		db:                 db,
		blogRepository:     blogRepository,
		categoryRepository: categoryRepository,
		userRepository:     userRepository,
		asynqClient:        asynqClient,
		followRepository:   followRepository,
		uploadService:      uploadService,
	}
}

func (s *blogService) Create(
	ctx context.Context,
	creatorId uuid.UUID,
	data dto.CreateBlogRequest,
) (*dto.BlogResponse, error) {
	if err := validateBlogSchedule(consts.BlogStatus(data.Status), data.ScheduledAt); err != nil {
		return nil, err
	}

	var (
		createdBlog *model.Blog
		author      *model.User
		category    *model.Category
	)

	err := s.db.Transaction(ctx, func(txDb repository.DbRepository) error {
		txBlogRepository := repository.NewBlogRepository(txDb)
		txCategoryRepository := repository.NewCategoryRepository(txDb)
		txUserRepository := repository.NewUserRepository(txDb)

		titleSlug := util.GenerateSlug(data.Title)
		var err error

		category, err = txCategoryRepository.FindByID(ctx, data.CategoryID, nil)
		if err != nil {
			return apperror.NewInternalServerError("Failed to retrieve category")
		}
		if category == nil {
			return apperror.NewNotFoundError("Category not found")
		}

		author, err = txUserRepository.FindByID(ctx, creatorId, nil)
		if err != nil {
			return apperror.NewInternalServerError("Failed to retrieve author")
		}
		if author == nil {
			return apperror.NewNotFoundError("Author not found")
		}

		if data.ImageURL != "" {
			isValid, err := s.uploadService.ValidateImageURL(ctx, data.ImageURL)
			if !isValid {
				return apperror.NewBadRequestError("Image URL is not valid")
			}
			if err != nil {
				return apperror.NewInternalServerError("Failed to validate image URL")
			}
			trackingRepo := repository.NewPresignedUploadTrackingRepository(txDb)
			presignUrl, err := trackingRepo.Find(ctx, "object_url = ?", nil, data.ImageURL)
			if err != nil {
				return apperror.NewInternalServerError("Failed to retrieve presigned URL")
			}
			if presignUrl == nil {
				return apperror.NewNotFoundError("Presigned URL not found")
			}
			if presignUrl.Status == consts.PresignedUploadStatusConfirmed {
				return apperror.NewBadRequestError("Image URL already used")
			}
			if err := trackingRepo.ConfirmByObjectURL(ctx, data.ImageURL); err != nil {
				return apperror.NewInternalServerError("Failed to confirm image URL")
			}
		}

		blog := &model.Blog{
			Slug:        titleSlug,
			Title:       data.Title,
			CategoryID:  data.CategoryID,
			Content:     data.Content,
			ViewTotal:   0,
			AuthorID:    creatorId,
			ImageURL:    data.ImageURL,
			Tags:        data.Tags,
			Status:      consts.BlogStatus(data.Status),
			ScheduledAt: data.ScheduledAt,
		}
		if blog.Status == consts.BlogStatusPublished {
			now := time.Now().UTC()
			blog.PublishedAt = &now
		}

		if createdBlog, err = txBlogRepository.Create(ctx, blog); err != nil {
			return apperror.NewInternalServerError("Failed to create blog")
		}

		if createdBlog.Status == consts.BlogStatusScheduled {
			task, opts := job.NewBlogPublishTask(createdBlog.ID, *createdBlog.ScheduledAt)
			if _, err := s.asynqClient.Enqueue(task, opts...); err != nil {
				return apperror.NewInternalServerError("Failed to schedule blog publish task")
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	createdBlog.Author = author
	createdBlog.Category = category

	return dto.NewBlogDetailResponse(createdBlog), nil
}

func (s *blogService) GetPublishedBlogs(ctx context.Context, filter dto.SearchBlogRequest) ([]*dto.BlogResponse, int64, error) {
	allowedSortColumns := map[string]bool{
		"title":      true,
		"created_at": true,
		"view_total": true,
	}
	sortBy := "created_at"
	if allowedSortColumns[filter.SortBy] {
		sortBy = filter.SortBy
	}

	sortOrder := "desc"
	if filter.SortOrder == "asc" || filter.SortOrder == "desc" {
		sortOrder = filter.SortOrder
	}
	orderClause := sortBy + " " + sortOrder

	var query string
	var args []any
	var conditions []string
	if filter.CategoryID != "" {
		conditions = append(conditions, "category_id = ?")
		args = append(args, filter.CategoryID)
	}
	conditions = append(conditions, "status = ?")
	args = append(args, consts.BlogStatusPublished)
	query = strings.Join(conditions, " AND ")
	blogs, total, err := s.blogRepository.List(
		ctx,
		filter.Limit,
		filter.Offset,
		orderClause,
		query,
		[]repository.Preload{
			repository.Author,
			repository.Category,
		},
		args...,
	)
	if err != nil {
		return nil, 0, apperror.NewInternalServerError("Failed to retrieve blogs")
	}

	return dto.NewListBlogResponse(blogs), total, nil
}

func (s *blogService) GetBlogStatistics(ctx context.Context) (*dto.BlogStatisticsResponse, error) {
	stats, err := s.blogRepository.GetStats(ctx)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve blog statistics")
	}

	avgEngagement := float64(0)
	if stats.TotalArticles > 0 {
		avgEngagement = math.Round(((float64(stats.Engaged)/float64(stats.TotalArticles))*100)*100) / 100
	}

	return &dto.BlogStatisticsResponse{
		TotalArticles: stats.TotalArticles,
		TotalViews:    stats.TotalViews,
		AvgEngagement: avgEngagement,
		Published:     stats.Published,
		Drafts:        stats.Drafts,
		Scheduled:     stats.Scheduled,
	}, nil
}

func (s *blogService) GetBlogBySlug(ctx context.Context, slug string, userId uuid.UUID) (*dto.BlogResponse, error) {
	blog, err := s.blogRepository.FindBySlug(
		ctx,
		slug,
		[]repository.Preload{
			repository.Author,
			repository.Category,
		},
	)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve blog")
	}
	if blog == nil {
		return nil, apperror.NewNotFoundError("Blog not found")
	}
	if err := s.blogRepository.IncrementView(ctx, blog.ID); err != nil {
		return nil, apperror.NewInternalServerError("Failed to increment blog view")
	}

	result := dto.NewBlogDetailResponse(blog)
	// Check if user is following the author
	isFollowing, err := s.followRepository.Exists(ctx, userId, blog.AuthorID)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to check follow status")
	}
	result.Author.IsFollowing = isFollowing

	return result, nil
}

func (s *blogService) DeleteBlog(ctx context.Context, authorId, blogId uuid.UUID) error {
	blog, err := s.blogRepository.FindByID(ctx, blogId, nil)
	if err != nil {
		return apperror.NewInternalServerError("Failed to retrieve blog")
	}
	if blog == nil {
		return apperror.NewNotFoundError("Blog not found")
	}
	if blog.AuthorID != authorId {
		return apperror.NewForbiddenError("You do not have permission to delete this blog")
	}
	if err := s.blogRepository.Delete(ctx, blogId); err != nil {
		return apperror.NewInternalServerError("Failed to delete blog")
	}

	return nil
}

func (s *blogService) UpdateBlogs(ctx context.Context, authorId, blogId uuid.UUID, data dto.UpdateBlogRequest) (*dto.BlogResponse, error) {
	if err := validateBlogSchedule(consts.BlogStatus(data.Status), data.ScheduledAt); err != nil {
		return nil, err
	}

	var blog *model.Blog
	err := s.db.Transaction(ctx, func(txDb repository.DbRepository) error {
		txBlogRepository := repository.NewBlogRepository(txDb)
		txCategoryRepository := repository.NewCategoryRepository(txDb)
		var err error

		blog, err = txBlogRepository.FindByID(ctx, blogId, nil)
		if err != nil {
			return apperror.NewInternalServerError("Failed to retrieve blog")
		}
		if blog == nil {
			return apperror.NewNotFoundError("Blog not found")
		}
		if blog.AuthorID != authorId {
			return apperror.NewForbiddenError("You do not have permission to update this blog")
		}

		category, err := txCategoryRepository.FindByID(ctx, data.CategoryID, nil)
		if err != nil {
			return apperror.NewInternalServerError("Failed to retrieve category")
		}
		if category == nil {
			return apperror.NewNotFoundError("Category not found")
		}

		if data.ImageURL != "" && blog.ImageURL != data.ImageURL {
			isValid, err := s.uploadService.ValidateImageURL(ctx, data.ImageURL)
			if !isValid {
				return apperror.NewBadRequestError("Image URL is not valid")
			}
			if err != nil {
				return apperror.NewInternalServerError("Failed to validate image URL")
			}
			trackingRepo := repository.NewPresignedUploadTrackingRepository(txDb)
			presignUrl, err := trackingRepo.Find(ctx, "object_url = ?", nil, data.ImageURL)
			if err != nil {
				return apperror.NewInternalServerError("Failed to retrieve presigned URL")
			}
			if presignUrl == nil {
				return apperror.NewNotFoundError("Presigned URL not found")
			}
			if presignUrl.Status == consts.PresignedUploadStatusConfirmed {
				return apperror.NewBadRequestError("Image URL already used")
			}
			if err := trackingRepo.ConfirmByObjectURL(ctx, data.ImageURL); err != nil {
				return apperror.NewInternalServerError("Failed to confirm image URL")
			}
		}
		

		blog.Title = data.Title
		blog.CategoryID = data.CategoryID
		blog.Content = data.Content
		blog.ImageURL = data.ImageURL
		blog.Tags = data.Tags
		if blog.Status != consts.BlogStatusPublished {
			blog.Status = consts.BlogStatus(data.Status)
			blog.ScheduledAt = data.ScheduledAt
			if blog.Status == consts.BlogStatusPublished {
				now := time.Now().UTC()
				blog.PublishedAt = &now
			}
		}
		if blog, err = txBlogRepository.Updates(ctx, blog); err != nil {
			return apperror.NewInternalServerError("Failed to update blog")
		}
		if blog.Status == consts.BlogStatusScheduled {
			task, opts := job.NewBlogPublishTask(blog.ID, *blog.ScheduledAt)
			if _, err := s.asynqClient.Enqueue(task, opts...); err != nil {
				return apperror.NewInternalServerError("Failed to schedule blog publish task")
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	updatedBlog, _ := s.blogRepository.FindByID(ctx, blog.ID, []repository.Preload{
		repository.PreloadPath(repository.Author, repository.User),
		repository.Category,
	})
	if updatedBlog == nil {
		return dto.NewBlogDetailResponse(blog), nil
	}

	return dto.NewBlogDetailResponse(updatedBlog), nil
}

func (s *blogService) GetBlogs(ctx context.Context, filter dto.SearchBlogRequest) ([]*dto.BlogResponse, int64, error) {
	allowedSortColumns := map[string]bool{
		"title":      true,
		"created_at": true,
		"view_total": true,
	}
	sortBy := "created_at"
	if allowedSortColumns[filter.SortBy] {
		sortBy = filter.SortBy
	}

	sortOrder := "desc"
	if filter.SortOrder == "asc" || filter.SortOrder == "desc" {
		sortOrder = filter.SortOrder
	}
	orderClause := sortBy + " " + sortOrder

	var query string
	var args []any
	var conditions []string
	if filter.CategoryID != "" {
		conditions = append(conditions, "category_id = ?")
		args = append(args, filter.CategoryID)
	}
	if filter.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, filter.Status)
	}
	filter.Keyword = strings.TrimSpace(filter.Keyword)
	if filter.Keyword != "" {
		conditions = append(conditions, "title ILIKE ?")
		args = append(args, "%"+filter.Keyword+"%")
	}
	query = strings.Join(conditions, " AND ")
	blogs, total, err := s.blogRepository.List(
		ctx,
		filter.Limit,
		filter.Offset,
		orderClause,
		query,
		[]repository.Preload{
			repository.Author,
			repository.Category,
		},
		args...,
	)
	if err != nil {
		return nil, 0, apperror.NewInternalServerError("Failed to retrieve blogs")
	}

	return dto.NewListBlogResponse(blogs), total, nil
}

func (s *blogService) GetByID(ctx context.Context, id uuid.UUID) (*dto.BlogResponse, error) {
	blog, err := s.blogRepository.FindByID(
		ctx,
		id,
		[]repository.Preload{
			repository.Author,
			repository.Category,
		},
	)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve blog")
	}
	if blog == nil {
		return nil, apperror.NewNotFoundError("Blog not found")
	}

	return dto.NewBlogDetailResponse(blog), nil
}

func validateBlogSchedule(status consts.BlogStatus, scheduledAt *time.Time) error {
	if status != consts.BlogStatusScheduled {
		return nil
	}

	if scheduledAt == nil {
		return apperror.NewBadRequestError("ScheduledAt is required when status is scheduled")
	}

	if !scheduledAt.UTC().After(time.Now().UTC()) {
		return apperror.NewBadRequestError("ScheduledAt must be in the future")
	}

	return nil
}
