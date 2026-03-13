package repository

import "elearning-api/model"

type CourseRepository interface {
	Repository[model.Course]
}

type courseRepository struct {
	*repository[model.Course]
}

func NewCourseRepository(db DbRepository) CourseRepository {
	return &courseRepository{
		repository: NewBaseRepository[model.Course](db),
	}
}
