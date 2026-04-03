package handler

import (
	"net/http"

	"elearning-api/dto"
	"elearning-api/service"
	"elearning-api/util"

	"github.com/gin-gonic/gin"
)

type UserCourseHandler interface {
	UpdateCourseProgress() gin.HandlerFunc
	GetMyCourses() gin.HandlerFunc
}

type userCourseHandler struct {
	enrollmentService service.EnrollmentService
}

func NewUserCourseHandler(enrollmentService service.EnrollmentService) UserCourseHandler {
	return &userCourseHandler{
		enrollmentService: enrollmentService,
	}
}

// UpdateCourseProgress godoc
// @Summary Update course progress (learn course action)
// @Description Update progress for a lesson the user completed
// @Tags user-courses
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body dto.StudentCourseProgressRequest true "course progress payload"
// @Success 200 {object} dto.ApiResponse
// @Router /api/v1/users/me/courses/progress [post]
func (h *userCourseHandler) UpdateCourseProgress() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dto.StudentCourseProgressRequest
		if err := util.BindAndValidateJSON(c, &request); err != nil {
			_ = c.Error(err)
			return
		}
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}
		enrollmentResponse, err := h.enrollmentService.UpdateLearnedLessons(
			c.Request.Context(),
			userID,
			request,
		)
		if err != nil {
			_ = c.Error(err)
			return
		}

		response := dto.NewApiResponse(c)
		response.Status = dto.NewResponseStatus(http.StatusCreated)
		response.Data = enrollmentResponse
		response.Request = request

		c.JSON(http.StatusCreated, response)
	}
}

// GetMyCourses godoc
// @Summary Track user course progress
// @Description Return list of enrolled courses with progress
// @Tags user-courses
// @Accept json
// @Produce json
// @Param limit query int false "limit"
// @Param offset query int false "offset"
// @Param sortBy query string false "sortBy"
// @Param sortOrder query string false "sortOrder"
// @Success 200 {object} any
// @Router /api/v1/users/me/courses [get]
func (h *userCourseHandler) GetMyCourses() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dto.CourseProgresssRequest
		if err := util.BindAndValidateQuery(c, &request); err != nil {
			_ = c.Error(err)
			return
		}

		request.PagingRequest.Process()

		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}

		items, total, err := h.enrollmentService.GetMyCourses(c.Request.Context(), userID, request)
		if err != nil {
			_ = c.Error(err)
			return
		}

		resp := dto.NewApiResponse(c)
		resp.Request = request
		resp.Data = items
		resp.Metadata = dto.NewPagination(request.Limit, request.Offset, int(total), request.SortBy, request.SortOrder)

		c.JSON(http.StatusOK, resp)
	}
}
