package service

import (
	"context"

	"elearning-api/apperror"
	"elearning-api/consts"
	"elearning-api/dto"
	"elearning-api/model"
	"elearning-api/repository"
	"elearning-api/util"

	"github.com/google/uuid"
)

type LessonService interface {
	Create(ctx context.Context, userID, courseID uuid.UUID, chapterRequests []dto.CreateLessonRequest) ([]*dto.ChapterResponse, error)
	UpdateCourseWithChapters(ctx context.Context, userID, courseID uuid.UUID, request dto.UpdateCourseWithChaptersRequest) ([]*dto.ChapterResponse, error)
	GetByCourseID(ctx context.Context, courseID uuid.UUID) ([]*dto.ChapterResponse, error)
}

type lessonService struct {
	lessonRepository  repository.LessonRepository
	db                repository.DbRepository
	courseRepository  repository.CourseRepository
	chapterRepository repository.ChapterRepository
	uploadService     UploadService
}

func NewLessonService(
	lessonRepository repository.LessonRepository,
	db repository.DbRepository,
	courseRepository repository.CourseRepository,
	chapterRepository repository.ChapterRepository,
	uploadService UploadService,
) LessonService {
	return &lessonService{
		lessonRepository:  lessonRepository,
		db:                db,
		courseRepository:  courseRepository,
		chapterRepository: chapterRepository,
		uploadService:     uploadService,
	}
}

func (s *lessonService) Create(ctx context.Context, userID, courseID uuid.UUID, chapterRequests []dto.CreateLessonRequest) ([]*dto.ChapterResponse, error) {
	course, err := s.courseRepository.FindByID(ctx, courseID, nil)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve course")
	}
	if course == nil {
		return nil, apperror.NewNotFoundError("Course not found")
	}
	if course.UserID != userID {
		return nil, apperror.NewForbiddenError("You do not have permission to modify this course")
	}
	if len(chapterRequests) == 0 {
		return nil, apperror.NewBadRequestError("Chapters is required")
	}
	var chapters []*model.Chapter
	orderChapter := 1
	err = s.db.Transaction(ctx, func(txDb repository.DbRepository) error {
		chapterRepo := repository.NewChapterRepository(txDb)
		lessonRepo := repository.NewLessonRepository(txDb)
		for _, chapterRequest := range chapterRequests {
			chapter, err := chapterRepo.Create(ctx, &model.Chapter{
				CourseID: courseID,
				Title:    chapterRequest.ChapterTitle,
				Order:    orderChapter,
			})
			if err != nil {
				return apperror.NewInternalServerError("Failed to create chapter")
			}
			lessons := make([]*model.Lesson, 0, len(chapterRequest.Lessons))
			orderLesson := 1
			for _, item := range chapterRequest.Lessons {
				hasVideo := item.VideoURL != "" && util.IsValidUrl(item.VideoURL)
				hasDocument := item.DocumentURL != "" && util.IsValidUrl(item.DocumentURL)
				if !hasVideo && !hasDocument {
					return apperror.NewBadRequestError("Lesson must contain at least video_url or document_url")
				}
				if hasVideo && hasDocument {
					return apperror.NewBadRequestError("Lesson cannot contain both video_url and document_url")
				}
				var lessonType consts.LessonType
				if hasVideo {
					lessonType = consts.LessonVideo
				} else {
					lessonType = consts.LessonDocument
				}
				lesson := &model.Lesson{
					ChapterID:       chapter.ID,
					Title:           item.Title,
					Duration:        item.Duration,
					IsAbleToPreview: item.IsAbleToPreview,
					Order:           orderLesson,
					Type:            lessonType,
				}
				if hasVideo {
					isVideoValid, err := s.uploadService.ValidateVideoURL(ctx, item.VideoURL)
					if err != nil {
						return apperror.NewBadRequestError("Failed to validate video url")
					}
					if !isVideoValid {
						return apperror.NewBadRequestError("Video does not exist")
					}
					lesson.VideoURL = item.VideoURL
				}
				if hasDocument {
					isDocumentValid, err := s.uploadService.ValidateDocumentURL(ctx, item.DocumentURL)
					if err != nil {
						return apperror.NewBadRequestError("Failed to validate document url")
					}
					if !isDocumentValid {
						return apperror.NewBadRequestError("Document does not exist")
					}
					lesson.DocumentURL = item.DocumentURL
				}
				lessons = append(lessons, lesson)
				orderLesson++
			}
			if err := lessonRepo.CreateBatch(ctx, lessons); err != nil {
				return apperror.NewInternalServerError("Failed to create lessons")
			}
			chapters = append(chapters, chapter)
			orderChapter++
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return dto.NewListChapterResponse(chapters), nil
}

func (s *lessonService) UpdateCourseWithChapters(ctx context.Context, userID, courseID uuid.UUID, request dto.UpdateCourseWithChaptersRequest) ([]*dto.ChapterResponse, error) {
	course, err := s.courseRepository.FindByID(ctx, courseID, nil)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve course")
	}
	if course == nil {
		return nil, apperror.NewNotFoundError("Course not found")
	}
	if course.UserID != userID {
		return nil, apperror.NewForbiddenError("You do not have permission to modify this course")
	}

	if len(request.Chapters) == 0 {
		return nil, apperror.NewBadRequestError("Chapters is required")
	}

	var resultChapters []*model.Chapter
	orderChapter := 1
	err = s.db.Transaction(ctx, func(txDb repository.DbRepository) error {
		chapterRepo := repository.NewChapterRepository(txDb)
		lessonRepo := repository.NewLessonRepository(txDb)

		// Map request chapters by ID for lookup
		requestChapterMap := make(map[string]*dto.UpdateChapterWithLessonsRequest)
		for i := range request.Chapters {
			if request.Chapters[i].ID != "" {
				requestChapterMap[request.Chapters[i].ID] = &request.Chapters[i]
			}
		}

		existingChapters, err := chapterRepo.FindAll(ctx, "course_id = ?", []repository.Preload{
			repository.PreloadPath(repository.Lessons),
		}, courseID)
		if err != nil {
			return apperror.NewInternalServerError("Failed to retrieve existing chapters")
		}

		// Map existing chapters by ID for lookup
		existingChapterMap := make(map[string]*model.Chapter)
		for _, ch := range existingChapters {
			existingChapterMap[ch.ID.String()] = ch
		}

		for _, chapterReq := range request.Chapters {
			var chapter *model.Chapter

			if chapterReq.ID == "" {
				newChapter := &model.Chapter{
					CourseID: courseID,
					Title:    chapterReq.Title,
					Order:    orderChapter,
				}
				createdChapter, err := chapterRepo.Create(ctx, newChapter)
				if err != nil {
					return apperror.NewInternalServerError("Failed to create chapter")
				}
				chapter = createdChapter
			} else {
				ch, exists := existingChapterMap[chapterReq.ID]
				if !exists {
					return apperror.NewBadRequestError("Chapter not found: " + chapterReq.ID)
				}
				ch.Title = chapterReq.Title
				ch.Order = orderChapter
				chapter, err = chapterRepo.Updates(ctx, ch)
				if err != nil {
					return apperror.NewInternalServerError("Failed to update chapter")
				}
				delete(existingChapterMap, chapterReq.ID) // Mark as processed
			}

			// Process lessons for this chapter
			requestLessonMap := make(map[string]*dto.UpdateLessonInChapterRequest)
			for i := range chapterReq.Lessons {
				if chapterReq.Lessons[i].ID != "" {
					requestLessonMap[chapterReq.Lessons[i].ID] = &chapterReq.Lessons[i]
				}
			}

			existingLessonMap := make(map[string]*model.Lesson)
			for _, lesson := range chapter.Lessons {
				existingLessonMap[lesson.ID.String()] = lesson
			}

			// Delete lessons not in request
			for _, lesson := range chapter.Lessons {
				if _, exists := requestLessonMap[lesson.ID.String()]; !exists {
					if err := lessonRepo.Delete(ctx, lesson.ID); err != nil {
						return apperror.NewInternalServerError("Failed to delete lesson")
					}
				}
			}

			// Create/update lessons
			lessons := make([]*model.Lesson, 0, len(chapterReq.Lessons))
			orderLesson := 1
			for _, lessonReq := range chapterReq.Lessons {
				if lessonReq.VideoURL != "" && lessonReq.DocumentURL != "" {
					return apperror.NewBadRequestError("Lesson cannot contain both video_url and document_url")
				}

				if lessonReq.VideoURL == "" && lessonReq.DocumentURL == "" {
					return apperror.NewBadRequestError("Lesson must contain at least video_url or document_url")
				}

				var lessonType consts.LessonType
				if lessonReq.VideoURL != "" {
					if !util.IsValidUrl(lessonReq.VideoURL) {
						return apperror.NewBadRequestError("Invalid video_url format")
					}
					isVideoValid, err := s.uploadService.ValidateVideoURL(ctx, lessonReq.VideoURL)
					if err != nil {
						return err
					}
					if !isVideoValid {
						return apperror.NewBadRequestError("Video does not exist")
					}
					lessonType = consts.LessonVideo
				} else {
					if !util.IsValidUrl(lessonReq.DocumentURL) {
						return apperror.NewBadRequestError("Invalid document_url format")
					}
					isDocumentValid, err := s.uploadService.ValidateDocumentURL(ctx, lessonReq.DocumentURL)
					if err != nil {
						return apperror.NewInternalServerError("Failed to check document url")
					}
					if !isDocumentValid {
						return apperror.NewBadRequestError("Document does not exist")
					}
					lessonType = consts.LessonDocument
				}
				if lessonReq.ID == "" {
					newLesson := &model.Lesson{
						ChapterID:       chapter.ID,
						Title:           lessonReq.Title,
						Duration:        lessonReq.Duration,
						IsAbleToPreview: lessonReq.IsAbleToPreview,
						Order:           orderLesson,
						Type:            lessonType,
						VideoURL:        lessonReq.VideoURL,
						DocumentURL:     lessonReq.DocumentURL,
					}
					lessons = append(lessons, newLesson)
				} else {
					lesson, exists := existingLessonMap[lessonReq.ID]
					if !exists {
						return apperror.NewBadRequestError("Lesson not found: " + lessonReq.ID)
					}
					lesson.Title = lessonReq.Title
					lesson.Duration = lessonReq.Duration
					lesson.IsAbleToPreview = lessonReq.IsAbleToPreview
					lesson.Order = orderLesson
					lesson.Type = lessonType
					lesson.VideoURL = lessonReq.VideoURL
					lesson.DocumentURL = lessonReq.DocumentURL
					lessons = append(lessons, lesson)
				}
				orderLesson++
			}
			if _, err := lessonRepo.SaveAll(ctx, lessons); err != nil {
				return apperror.NewInternalServerError("Failed to save lessons")
			}
			resultChapters = append(resultChapters, chapter)
			orderChapter++
		}

		// Delete chapters not in request
		for _, ch := range existingChapters {
			if _, exists := requestChapterMap[ch.ID.String()]; !exists {
				if err := chapterRepo.Delete(ctx, ch.ID); err != nil {
					return apperror.NewInternalServerError("Failed to delete chapter")
				}
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return dto.NewListChapterResponse(resultChapters), nil
}

func (s *lessonService) GetByCourseID(ctx context.Context, courseID uuid.UUID) ([]*dto.ChapterResponse, error) {
	chapters, err := s.chapterRepository.FindAll(ctx, "course_id = ?", []repository.Preload{
		repository.PreloadPath(repository.Lessons),
	}, courseID)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve chapters and lessons")
	}

	return dto.NewListChapterResponse(chapters), nil
}
