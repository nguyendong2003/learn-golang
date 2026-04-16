package handler

import (
	"net/http"

	"elearning-api/apperror"
	"elearning-api/consts"
	"elearning-api/dto"
	"elearning-api/service"
	"elearning-api/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CourseHandler interface {
	Create() gin.HandlerFunc
	Update() gin.HandlerFunc
	AssignCoupons() gin.HandlerFunc
	UpdateStatus() gin.HandlerFunc
	Delete() gin.HandlerFunc

	GetByID() gin.HandlerFunc
	GetBySlug() gin.HandlerFunc
	GetList() gin.HandlerFunc
	GetMyTaughtCourses() gin.HandlerFunc

	PublishCourse() gin.HandlerFunc

	CreatePurchaseCheckoutSession() gin.HandlerFunc
	PreviewPurchaseCheckout() gin.HandlerFunc

	CreateEvent() gin.HandlerFunc
	GetEvents() gin.HandlerFunc
	DeleteEvent() gin.HandlerFunc
	UpdateEvent() gin.HandlerFunc
	SubmitForReview() gin.HandlerFunc
	ApproveCourse() gin.HandlerFunc
	RejectCourse() gin.HandlerFunc

	GetCoursePurchaseBySessionID() gin.HandlerFunc
	// Only student who is in plan HIGHER than free plan can enroll DIRECTLY in course
	// and only if they are not enrolled in the course yet
	EnrollInCourse() gin.HandlerFunc

	GetStatistics() gin.HandlerFunc
	GetNewCoursesLast30Days() gin.HandlerFunc
}

type courseHandler struct {
	courseService       service.CourseService
	subscriptionService service.SubscriptionService
	uploadService       service.UploadService
	enrollmentService   service.EnrollmentService
}

func NewCourseHandler(
	courseService service.CourseService,
	subscriptionService service.SubscriptionService,
	instructorProfileService service.InstructorProfileService,
	uploadService service.UploadService,
	enrollmentService service.EnrollmentService,
) CourseHandler {
	return &courseHandler{
		courseService:       courseService,
		subscriptionService: subscriptionService,
		uploadService:       uploadService,
		enrollmentService:   enrollmentService,
	}
}

// Create godoc
// @Summary Create a new course
// @Description Create a new course (instructor only)
// @Tags courses
// @Accept json
// @Produce json
// @Param payload body dto.CreateCourseRequest true "Course create request"
// @Success 201 {object} dto.ApiResponse{data=dto.CourseResponse}
// @Failure 400 {object} dto.ApiResponse "Invalid request"
// @Failure 401 {object} dto.ApiResponse "Unauthorized"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /api/v1/courses [post]
// @Security BearerAuth
func (h *courseHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}
		var request dto.CreateCourseRequest
		if err := util.BindAndValidateJSON(c, &request); err != nil {
			_ = c.Error(err)
			return
		}
		createdCourse, err := h.courseService.Create(c.Request.Context(), userID, request)
		if err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Status = dto.NewResponseStatus(http.StatusCreated)
		res.Request = dto.GetRequestClient(c)
		res.Data = createdCourse

		c.JSON(http.StatusCreated, res)
	}
}

// Update godoc
// @Summary Update an existing course which status is not published (instructor only)
// @Description Update a unpublished course by ID (instructor only) using JSON body.
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID (UUID format)"
// @Param payload body dto.UpdateCourseRequest true "Update payload"
// @Success 200 {object} dto.ApiResponse{data=dto.CourseResponse}
// @Failure 400 {object} dto.ApiResponse "Invalid UUID or payload"
// @Failure 401 {object} dto.ApiResponse "Unauthorized"
// @Failure 404 {object} dto.ApiResponse "Course not found"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /api/v1/courses/{id} [put]
// @Security BearerAuth
func (h *courseHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID format"))
			return
		}
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}
		var request dto.UpdateCourseRequest
		if err := util.BindAndValidateJSON(c, &request); err != nil {
			_ = c.Error(err)
			return
		}
		updatedCourse, err := h.courseService.Update(c.Request.Context(), userID, id, request)
		if err != nil {
			_ = c.Error(err)
			return
		}
		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = updatedCourse
		c.JSON(http.StatusOK, res)
	}
}

// AssignCoupons godoc
// @Summary Assign coupons to a course
// @Description Manage coupons for a course (add, update, or remove).
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID (UUID format)"
// @Param payload body dto.AssignCourseCouponsRequest true "Coupon assignment payload"
// @Success 200 {object} dto.ApiResponse{data=dto.CourseResponse}
// @Failure 400 {object} dto.ApiResponse "Invalid UUID or payload"
// @Failure 401 {object} dto.ApiResponse "Unauthorized"
// @Failure 403 {object} dto.ApiResponse "Forbidden"
// @Failure 404 {object} dto.ApiResponse "Course not found"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /api/v1/courses/{id}/coupons [put]
// @Security BearerAuth
func (h *courseHandler) AssignCoupons() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID format"))
			return
		}
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}
		var request dto.AssignCourseCouponsRequest
		if err := util.BindAndValidateJSON(c, &request); err != nil {
			_ = c.Error(err)
			return
		}
		updatedCourse, err := h.courseService.AssignCoupons(c.Request.Context(), userID, id, request)
		if err != nil {
			_ = c.Error(err)
			return
		}
		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = updatedCourse
		c.JSON(http.StatusOK, res)
	}
}

// UpdateStatus godoc
// @Summary Update course status
// @Description Update course status. If status is `published`, system will publish and sync product catalog on Stripe.
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID (UUID format)"
// @Param payload body dto.UpdateCourseStatusRequest true "Update status payload"
// @Success 200 {object} dto.ApiResponse{data=dto.CourseResponse}
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Router /api/v1/courses/{id}/update-status [put]
// @Security BearerAuth
func (h *courseHandler) UpdateStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		var courseIDRequest dto.UUIDRequest

		if err := c.ShouldBindUri(&courseIDRequest); err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID in URI"))
			return
		}

		id, err := uuid.Parse(courseIDRequest.ID)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID format"))
			return
		}

		var request dto.UpdateCourseStatusRequest
		if err := util.BindAndValidateJSON(c, &request); err != nil {
			_ = c.Error(err)
			return
		}

		var updatedCourse *dto.CourseResponse
		if request.Status == consts.CoursePublished {
			isAlreadyPublished := false
			if v, ok := c.Get(consts.ContextCourse); ok {
				if course, ok := v.(*dto.CourseResponse); ok && course != nil && course.Status == consts.CoursePublished {
					isAlreadyPublished = true
				}
			}

			if isAlreadyPublished {
				updatedCourse, err = h.subscriptionService.SyncPublishedCourseCatalog(c.Request.Context(), id)
			} else {
				updatedCourse, err = h.subscriptionService.PublishCourseAndCreateProductCatalog(c.Request.Context(), id)
			}
			if err != nil {
				_ = c.Error(err)
				return
			}
		} else {
			updatedCourse, err = h.courseService.UpdateStatus(c.Request.Context(), id, request.Status)
			if err != nil {
				_ = c.Error(err)
				return
			}
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = updatedCourse

		c.JSON(http.StatusOK, res)
	}
}

// Delete godoc
// @Summary Delete a course
// @Description Delete a course by ID
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID (UUID format)"
// @Success 200 {object} dto.ApiResponse
// @Failure 400 {object} dto.ApiResponse "Invalid UUID format"
// @Failure 401 {object} dto.ApiResponse "Unauthorized"
// @Failure 404 {object} dto.ApiResponse "Course not found"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /api/v1/courses/{id} [delete]
// @Security BearerAuth
func (h *courseHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		var courseIDRequest dto.UUIDRequest

		if err := c.ShouldBindUri(&courseIDRequest); err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID in URI"))
			return
		}

		id, err := uuid.Parse(courseIDRequest.ID)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID format"))
			return
		}

		// Call service
		if err := h.courseService.Delete(c.Request.Context(), id); err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)

		c.JSON(http.StatusOK, res)
	}
}

// GetByID godoc
// @Summary Get course by ID
// @Description Retrieve a course by its ID
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID (UUID format)"
// @Success 200 {object} dto.ApiResponse{data=dto.CourseResponse}
// @Failure 400 {object} dto.ApiResponse "Invalid UUID format"
// @Failure 404 {object} dto.ApiResponse "Course not found"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /api/v1/courses/{id} [get]
func (h *courseHandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var idRequest dto.UUIDRequest

		if err := c.ShouldBindUri(&idRequest); err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID in URI"))
			return
		}

		id, err := uuid.Parse(idRequest.ID)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID format"))
			return
		}

		// Call service
		data, err := h.courseService.GetByID(c.Request.Context(), id)
		if err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = data

		c.JSON(http.StatusOK, res)
	}
}

// GetBySlug godoc
// @Summary Get course by slug
// @Description Retrieve a course by its slug
// @Tags courses
// @Accept json
// @Produce json
// @Param slug path string true "Course slug"
// @Success 200 {object} dto.ApiResponse{data=dto.CourseResponse}
// @Failure 400 {object} dto.ApiResponse "Invalid slug"
// @Failure 404 {object} dto.ApiResponse "Course not found"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /api/v1/courses/slug/{slug} [get]
func (h *courseHandler) GetBySlug() gin.HandlerFunc {
	return func(c *gin.Context) {
		var slugRequest dto.CourseSlugRequest

		if err := c.ShouldBindUri(&slugRequest); err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid slug in URI"))
			return
		}

		userRole := util.GetRole(c)
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}

		data, err := h.courseService.GetBySlug(c.Request.Context(), slugRequest.Slug, userID, userRole)
		if err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = data

		c.JSON(http.StatusOK, res)
	}
}

// GetList godoc
// @Summary Get list of courses
// @Description Retrieve paginated list of courses with filters
// @Tags courses
// @Accept json
// @Produce json
// @Param limit query int false "Items per page" default(10)
// @Param offset query int false "Offset" default(0)
// @Param sortBy query string false "Sort field" default(created_at)
// @Param sortOrder query string false "Sort order" Enums(asc,desc) default(desc)
// @Param title query string false "Filter by title"
// @Param category_id query string false "Filter by category ID (UUID)"
// @Param status query string false "Filter by status" Enums(draft,pending_review,published,rejected,archived)
// @Success 200 {object} dto.ApiResponse{data=[]dto.CourseResponse,metadata=dto.Pagination}
// @Failure 400 {object} dto.ApiResponse "Invalid query"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /api/v1/courses [get]
func (h *courseHandler) GetList() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dto.ListCourseRequest

		if err := c.ShouldBindQuery(&request); err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid query parameters"))
			return
		}

		userRole := util.GetRole(c)
		// Only admin and instructor can see all courses
		// Students can only see published courses
		// Instructors can see their own courses in any status + all published courses
		if userRole != string(consts.RoleAdmin) {
			publishedStatus := consts.CoursePublished
			request.Status = &publishedStatus
		} else if userRole == string(consts.RoleInstructor) {
			userID, err := util.GetRequestUserID(c)
			if err != nil {
				_ = c.Error(err)
				return
			}
			userIDStr := userID.String()
			request.UserID = &userIDStr
		}

		// Process default pagination
		request.Process()

		limit := request.Limit
		offset := request.Offset
		sortBy := request.SortBy
		sortOrder := request.SortOrder

		// Call service
		data, total, err := h.courseService.GetList(c.Request.Context(), request)
		if err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = data
		res.Metadata = dto.NewPagination(limit, offset, int(total), sortBy, sortOrder)

		c.JSON(http.StatusOK, res)
	}
}

// CreatePurchaseCheckoutSession godoc
// @Summary Create checkout session for one-time course purchase
// @Description Create a Stripe Checkout session (payment mode) for purchasing course lifetime access
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID (UUID format)"
// @Param payload body dto.CreateCoursePurchaseCheckoutSessionRequest false "Checkout session request payload"
// @Success 200 {object} dto.ApiResponse
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Router /api/v1/courses/{id}/purchase-checkout-session [post]
// @Security BearerAuth
func (h *courseHandler) CreatePurchaseCheckoutSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		var idRequest dto.UUIDRequest
		if err := c.ShouldBindUri(&idRequest); err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID in URI"))
			return
		}

		var request dto.CreateCoursePurchaseCheckoutSessionRequest
		if err := util.BindAndValidateJSON(c, &request); err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid body"))
			return
		}

		courseID, err := uuid.Parse(idRequest.ID)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID format"))
			return
		}

		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}

		result, err := h.subscriptionService.CreateCourseCheckoutSession(c.Request.Context(), userID, courseID, request)
		if err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = result

		c.JSON(http.StatusOK, res)
	}
}

// PreviewPurchaseCheckout godoc
// @Summary Preview checkout for one-time course purchase
// @Description Preview the Stripe checkout amount for purchasing course lifetime access, including an applied coupon if available
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID (UUID format)"
// @Param payload body dto.CreateCoursePurchaseCheckoutSessionRequest false "Checkout preview request payload"
// @Success 200 {object} dto.ApiResponse{data=dto.CoursePurchaseCheckoutPreviewResponse}
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Router /api/v1/courses/{id}/purchase-checkout-preview [post]
// @Security BearerAuth
func (h *courseHandler) PreviewPurchaseCheckout() gin.HandlerFunc {
	return func(c *gin.Context) {
		var idRequest dto.UUIDRequest
		if err := c.ShouldBindUri(&idRequest); err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID in URI"))
			return
		}

		var request dto.CreateCoursePurchaseCheckoutSessionRequest
		if err := util.BindAndValidateJSON(c, &request); err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid body"))
			return
		}

		courseID, err := uuid.Parse(idRequest.ID)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID format"))
			return
		}

		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}

		result, err := h.subscriptionService.PreviewCourseCheckout(c.Request.Context(), userID, courseID, request)
		if err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = result

		c.JSON(http.StatusOK, res)
	}
}

// GetCoursePurchaseBySessionID godoc
// @Summary Get course purchase by checkout session ID
// @Description Retrieve course purchase record by Stripe checkout session ID (exclude payment intent id)
// @Tags courses
// @Accept json
// @Produce json
// @Param session_id path string true "Checkout session ID"
// @Success 200 {object} dto.ApiResponse
// @Failure 400 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Router /api/v1/course-purchase/{session_id} [get]
func (h *courseHandler) GetCoursePurchaseBySessionID() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.Param("session_id")
		if sessionID == "" {
			_ = c.Error(apperror.NewBadRequestError("Missing session id"))
			return
		}

		purchase, err := h.courseService.GetCoursePurchaseBySessionID(c.Request.Context(), sessionID)
		if err != nil {
			_ = c.Error(err)
			return
		}
		if purchase == nil {
			_ = c.Error(apperror.NewNotFoundError("Course purchase not found"))
			return
		}

		if purchase.Status != string(consts.CoursePurchaseStatusPaid) {
			_ = c.Error(apperror.NewBadRequestError("Course purchase not paid"))
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = purchase

		c.JSON(http.StatusOK, res)
	}
}

// SubmitForReview godoc
// @Summary Submit course for review (instructor)
// @Description Instructor submits their course for admin review
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID (UUID format)"
// @Success 200 {object} dto.ApiResponse{data=dto.CourseResponse}
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 403 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Router /api/v1/courses/{id}/submit [post]
// @Security BearerAuth
func (h *courseHandler) SubmitForReview() gin.HandlerFunc {
	return func(c *gin.Context) {
		var idRequest dto.UUIDRequest
		if err := c.ShouldBindUri(&idRequest); err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID in URI"))
			return
		}

		courseID, err := uuid.Parse(idRequest.ID)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID format"))
			return
		}

		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}

		updated, err := h.courseService.SubmitForReview(c.Request.Context(), userID, courseID)
		if err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = updated

		c.JSON(http.StatusOK, res)
	}
}

// ApproveCourse godoc
// @Summary Approve course (admin)
// @Description Admin approves a submitted course and publishes it
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID (UUID format)"
// @Success 200 {object} dto.ApiResponse{data=dto.CourseResponse}
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 403 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Router /api/v1/admin/courses/{id}/approve [post]
// @Security BearerAuth
func (h *courseHandler) ApproveCourse() gin.HandlerFunc {
	return func(c *gin.Context) {
		var idRequest dto.UUIDRequest
		if err := c.ShouldBindUri(&idRequest); err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID in URI"))
			return
		}

		courseID, err := uuid.Parse(idRequest.ID)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID format"))
			return
		}

		course, err := h.subscriptionService.PublishCourseAndCreateProductCatalog(c.Request.Context(), courseID)
		if err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = course

		c.JSON(http.StatusOK, res)
	}
}

// RejectCourse godoc
// @Summary Reject course (admin)
// @Description Admin rejects a submitted course
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID (UUID format)"
// @Success 200 {object} dto.ApiResponse
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 403 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Router /api/v1/admin/courses/{id}/reject [post]
// @Security BearerAuth
func (h *courseHandler) RejectCourse() gin.HandlerFunc {
	return func(c *gin.Context) {
		var idRequest dto.UUIDRequest
		if err := c.ShouldBindUri(&idRequest); err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID in URI"))
			return
		}

		courseID, err := uuid.Parse(idRequest.ID)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID format"))
			return
		}

		_, err = h.courseService.UpdateStatus(c.Request.Context(), courseID, consts.CourseRejected)
		if err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)

		c.JSON(http.StatusOK, res)
	}
}

// PublishCourse godoc
// @Summary Publish course
// @Description Publish course and create/update Stripe product-price catalog for one-time purchase
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID (UUID format)"
// @Success 200 {object} dto.ApiResponse{data=dto.CourseResponse}
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Router /api/v1/courses/{id}/publish [post]
// @Security BearerAuth
func (h *courseHandler) PublishCourse() gin.HandlerFunc {
	return func(c *gin.Context) {
		var idRequest dto.UUIDRequest

		if err := c.ShouldBindUri(&idRequest); err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID in URI"))
			return
		}

		courseID, err := uuid.Parse(idRequest.ID)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID format"))
			return
		}

		course, err := h.subscriptionService.PublishCourseAndCreateProductCatalog(c.Request.Context(), courseID)
		if err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = course

		c.JSON(http.StatusOK, res)
	}
}

// CreateEvent godoc
// @Summary Create a new event for a course
// @Description Create a new event (e.g. live session) for a specific course. Validates that the event time does not overlap with existing events.
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID (UUID format)"
// @Param payload body dto.CourseEventRequest true "Course event create request"
// @Success 200 {object} dto.ApiResponse{data=dto.CourseEventResponse}
// @Failure 400 {object} dto.ApiResponse "Invalid UUID format, validation failed, or event time overlaps with existing event"
// @Failure 401 {object} dto.ApiResponse "Unauthorized"
// @Failure 404 {object} dto.ApiResponse "Course not found"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /api/v1/courses/{id}/events [post]
func (h *courseHandler) CreateEvent() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		courseID, err := uuid.Parse(id)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid course ID format"))
			return
		}
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(apperror.NewUnauthorizedError("Unauthorized"))
			return
		}
		var request dto.CourseEventRequest
		if err := util.BindAndValidateJSON(c, &request); err != nil {
			_ = c.Error(err)
			return
		}

		event, err := h.courseService.CreateEvent(c.Request.Context(), userID, courseID, request)
		if err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = event

		c.JSON(http.StatusOK, res)
	}
}

// GetEvents godoc
// @Summary Get events for a course
// @Description Retrieve a list of events (e.g. live sessions) for a specific course.
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID (UUID format)"
// @Success 200 {object} dto.ApiResponse{data=[]dto.CourseEventResponse}
// @Failure 400 {object} dto.ApiResponse "Invalid UUID format"
// @Failure 404 {object} dto.ApiResponse "Course not found"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /api/v1/courses/{id}/events [get]
func (h *courseHandler) GetEvents() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		courseID, err := uuid.Parse(id)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid course ID format"))
			return
		}
		events, err := h.courseService.GetEvents(c.Request.Context(), courseID)
		if err != nil {
			_ = c.Error(err)
			return
		}
		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = events

		c.JSON(http.StatusOK, res)
	}
}

// DeleteEvent godoc
// @Summary Delete an event from a course
// @Description Delete a specific event (e.g. live session) from a course. Requires instructor permissions.
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID (UUID format)"
// @Param event_id path string true "Event ID (UUID format)"
// @Success 200 {object} dto.ApiResponse
// @Failure 400 {object} dto.ApiResponse "Invalid UUID format"
// @Failure 401 {object} dto.ApiResponse "Unauthorized"
// @Failure 404 {object} dto.ApiResponse "Course or event not found"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /api/v1/courses/{id}/events/{event_id} [delete]
// @Security BearerAuth
func (h *courseHandler) DeleteEvent() gin.HandlerFunc {
	return func(c *gin.Context) {
		idCourse := c.Param("id")
		courseID, err := uuid.Parse(idCourse)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid course ID format"))
			return
		}
		idEvent := c.Param("event_id")
		eventID, err := uuid.Parse(idEvent)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid event ID format"))
			return
		}
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(apperror.NewUnauthorizedError("Unauthorized"))
			return
		}
		if err := h.courseService.DeleteEvent(c.Request.Context(), userID, courseID, eventID); err != nil {
			_ = c.Error(err)
			return
		}
		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = gin.H{"message": "Event deleted successfully"}

		c.JSON(http.StatusOK, res)
	}
}

// UpdateEvent godoc
// @Summary Update an event for a course
// @Description Update a specific event (e.g. live session) for a course. Validates that the updated event time does not overlap with existing events. Requires instructor permissions.
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID (UUID format)"
// @Param event_id path string true "Event ID (UUID format)"
// @Param payload body dto.CourseEventRequest true "Course event update request"
// @Success 200 {object} dto.ApiResponse{data=dto.CourseEventResponse}
// @Failure 400 {object} dto.ApiResponse "Invalid UUID format, validation failed, or event time overlaps with existing event"
// @Failure 401 {object} dto.ApiResponse "Unauthorized"
// @Failure 404 {object} dto.ApiResponse "Course or event not found"
// @Failure 500 {object} dto.ApiResponse
// @Router /api/v1/courses/{id}/events/{event_id} [put]
// @Security BearerAuth
func (h *courseHandler) UpdateEvent() gin.HandlerFunc {
	return func(c *gin.Context) {
		idCourse := c.Param("id")
		courseID, err := uuid.Parse(idCourse)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid course ID format"))
			return
		}
		idEvent := c.Param("event_id")
		eventID, err := uuid.Parse(idEvent)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid event ID format"))
			return
		}
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(apperror.NewUnauthorizedError("Unauthorized"))
			return
		}
		var request dto.CourseEventRequest
		if err := util.BindAndValidateJSON(c, &request); err != nil {
			_ = c.Error(err)
			return
		}
		updatedEvent, err := h.courseService.UpdateEvent(c.Request.Context(), userID, courseID, eventID, request)
		if err != nil {
			_ = c.Error(err)
			return
		}
		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = updatedEvent

		c.JSON(http.StatusOK, res)
	}
}

// GetStatistics godoc
// @Summary Get course statistics
// @Description Retrieve aggregated statistics for courses, such as total courses, published courses, and total revenue.
// @Tags courses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.ApiResponse{data=dto.CourseStatisticsResponse}
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /api/v1/admin/courses/statistics [get]
func (h *courseHandler) GetStatistics() gin.HandlerFunc {
	return func(c *gin.Context) {
		stats, err := h.courseService.GetStatistics(c.Request.Context())
		if err != nil {
			_ = c.Error(err)
			return
		}
		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = stats

		c.JSON(http.StatusOK, res)
	}
}

// GetNewCoursesLast30Days godoc
// @Summary Get new courses statistics for the last 30 days
// @Description Get the number of courses created in the last 30 days for admin dashboard.
// @Tags courses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.ApiResponse{data=dto.NewCoursesLast30DaysResponse}
// @Failure 401 {object} dto.ApiResponse
// @Failure 403 {object} dto.ApiResponse
// @Failure 500 {object} dto.ApiResponse
// @Router /api/v1/admin/courses/statistics/new [get]
func (h *courseHandler) GetNewCoursesLast30Days() gin.HandlerFunc {
	return func(c *gin.Context) {
		stats, err := h.courseService.GetNewCoursesLast30Days(c.Request.Context())
		if err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = stats

		c.JSON(http.StatusOK, res)
	}
}

// EnrollInCourse godoc
// @Summary Enroll in a course
// @Description Enroll the authenticated students who is in higher subscription than free
// @Tags courses
// @Param id path string true "Course ID (UUID format)"
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.ApiResponse{data=dto.EnrollmentResponse}
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /api/v1/courses/{id}/enroll [post]
func (h *courseHandler) EnrollInCourse() gin.HandlerFunc {
	return func(c *gin.Context) {
		var idRequest dto.UUIDRequest
		if err := c.ShouldBindUri(&idRequest); err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID in URI"))
			return
		}

		courseID, err := uuid.Parse(idRequest.ID)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID format"))
			return
		}

		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}

		enrollment, err := h.enrollmentService.EnrollInCourse(c.Request.Context(), userID, courseID)
		if err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = enrollment
		c.JSON(http.StatusCreated, res)
	}
}

// GetMyTaughtCourses godoc
// @Summary Get my taught courses with revenue
// @Description Retrieve paginated courses where current user is the instructor, including revenue and student count for each course.
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Param sortBy query string false "Sort field" Enums(created_at,title,status,total_student,revenue) default(created_at)
// @Param sortOrder query string false "Sort order" Enums(asc,desc) default(desc)
// @Success 200 {object} dto.ApiResponse{data=[]dto.InstructorTaughtCourseRevenueResponse,metadata=dto.Pagination}
// @Failure 400 {object} dto.ApiResponse "Invalid pagination parameters"
// @Failure 401 {object} dto.ApiResponse "Unauthorized"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /api/v1/users/me/courses/taught [get]
func (h *courseHandler) GetMyTaughtCourses() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}
		var request dto.PagingRequest
		if err := c.ShouldBindQuery(&request); err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid pagination parameters"))
			return
		}
		request.Process()
		data, total, err := h.courseService.GetInstructorTaughtCourseRevenue(c.Request.Context(), userID, request)
		if err != nil {
			_ = c.Error(err)
			return
		}
		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = data
		res.Metadata = dto.NewPagination(request.Limit, request.Offset, int(total), request.SortBy, request.SortOrder)

		c.JSON(http.StatusOK, res)
	}
}
