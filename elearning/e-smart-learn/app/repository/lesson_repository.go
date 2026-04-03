package repository

import (
	"context"

	"elearning-api/model"

	"gorm.io/gorm/clause"
)

type LessonRepository interface {
	Repository[model.Lesson]
	SaveAll(ctx context.Context, lessons []*model.Lesson) ([]*model.Lesson, error)
}

type lessonRepository struct {
	*repository[model.Lesson]
}

func NewLessonRepository(db DbRepository) LessonRepository {
	return &lessonRepository{
		repository: NewBaseRepository[model.Lesson](db),
	}
}
func (r *lessonRepository) SaveAll(ctx context.Context, lessons []*model.Lesson) ([]*model.Lesson, error) {
    if len(lessons) == 0 {
        return lessons, nil
    }
    err := r.db.GetDB().WithContext(ctx).
        Clauses(clause.OnConflict{
            Columns:   []clause.Column{{Name: "id"}},
            UpdateAll: true, 
        }).
        CreateInBatches(&lessons, 100).Error 

    if err != nil {
        return nil, err
    }

    return lessons, nil
}