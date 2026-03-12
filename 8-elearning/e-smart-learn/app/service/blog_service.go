package service

import (
	"context"
	"elearning-api/apperror"
	"elearning-api/dto"
	"elearning-api/model"
	"elearning-api/repository"
	"elearning-api/util"

	"github.com/google/uuid"
)

type BlogService interface {
	Create(ctx context.Context, creatorId uuid.UUID, data dto.CreateBlogRequest) (*dto.BlogResponse, error)
}

type blogService struct {
	blogRepository repository.BlogRepository
}

func NewBlogService(blogRepository repository.BlogRepository) BlogService {
	return &blogService{
		blogRepository: blogRepository,
	}
}

func (s *blogService) Create(
	ctx context.Context,
	creatorId uuid.UUID,
	data dto.CreateBlogRequest,
) (*dto.BlogResponse, error) {
	var err error
	titleSlug := util.GenerateSlug(data.Title)

	blog := &model.Blog{
		Slug:      titleSlug,
		Title:     data.Title,
		Content:   data.Content,
		ViewTotal: 0,
		AuthorID:  creatorId,
		ImageURL:  data.ImageURL,
	}

	if blog, err = s.blogRepository.Create(ctx, blog); err != nil {
		return nil, apperror.NewInternalServerError("Failed to create blog")
	}

	insertedBlog, err := s.blogRepository.FindByID(ctx, blog.ID, []repository.Preload{
		repository.PreloadPath(repository.Author, repository.User),
	})

	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve created blog")
	}

	return dto.NewBlogDetailResponse(insertedBlog), nil
}
