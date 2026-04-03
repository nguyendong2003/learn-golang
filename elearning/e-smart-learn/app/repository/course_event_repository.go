package repository

import "elearning-api/model"

type CourseEventRepository interface {
	Repository[model.CourseEvent]
}

type courseEventRepository struct {
	*repository[model.CourseEvent]
}

func NewCourseEventRepository(db DbRepository) CourseEventRepository {
	return &courseEventRepository{
		repository: NewBaseRepository[model.CourseEvent](db),
	}
}
