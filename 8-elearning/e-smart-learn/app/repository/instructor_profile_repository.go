package repository

import (
	"elearning-api/model"
)

type InstructorProfileRepository interface {
	Repository[model.InstructorProfile]
}

type instructorProfileRepository struct {
	*repository[model.InstructorProfile]
}

func NewInstructorProfileRepository(db DbRepository) InstructorProfileRepository {
	return &instructorProfileRepository{
		repository: NewBaseRepository[model.InstructorProfile](db),
	}
}
