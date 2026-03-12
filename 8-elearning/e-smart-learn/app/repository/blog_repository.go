package repository

import "elearning-api/model"

type BlogRepository interface {
	Repository[model.Blog]
}

type blogRepository struct {
	*repository[model.Blog]
}

func NewBlogRepository(db DbRepository) BlogRepository {
	return &blogRepository{
		repository: NewBaseRepository[model.Blog](db),
	}
}