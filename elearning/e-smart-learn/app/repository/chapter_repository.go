package repository

import (
	"elearning-api/model"
)

type ChapterRepository interface {
	Repository[model.Chapter]
}

type chapterRepository struct {
	*repository[model.Chapter]
}

func NewChapterRepository(db DbRepository) ChapterRepository {
	return &chapterRepository{
		repository: NewBaseRepository[model.Chapter](db),
	}
}
