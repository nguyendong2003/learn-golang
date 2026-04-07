package handler

import (
	"net/http"
	"time"

	"elearning-api/apperror"
	"elearning-api/dto"
	"elearning-api/service"
	"elearning-api/util"

	"github.com/gin-gonic/gin"
)

type RevenueHandler interface {
	GetInstructorStatistics() gin.HandlerFunc
	GetInstructorStatisticsByDay() gin.HandlerFunc
	GetInstructorStatisticsByMonth() gin.HandlerFunc
	GetInstructorStatisticsByYear() gin.HandlerFunc
	GetInstructorRevenueOverview() gin.HandlerFunc

	GetAdminStatistics() gin.HandlerFunc
	GetAdminStatisticsByDay() gin.HandlerFunc
	GetAdminStatisticsByMonth() gin.HandlerFunc
	GetAdminStatisticsByYear() gin.HandlerFunc
	GetAdminRevenueOverview() gin.HandlerFunc

	GetAdminSalesSegmentation() gin.HandlerFunc
	GetAdminTransactions() gin.HandlerFunc

	GetAllTeachersRevenue() gin.HandlerFunc
}

type revenueHandler struct {
	revenueService service.RevenueService
}

func NewRevenueHandler(revenueService service.RevenueService) RevenueHandler {
	return &revenueHandler{revenueService: revenueService}
}

// GetInstructorStatistics godoc
// @Summary Get instructor revenue statistics
// @Description Get all-time revenue statistics for the authenticated instructor, including course purchases and subscriptions.
// @Tags statistics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.ApiResponse{data=dto.RevenueStatisticsResponse}
// @Failure 401 {object} dto.ApiResponse
// @Failure 403 {object} dto.ApiResponse
// @Failure 500 {object} dto.ApiResponse
// @Router /api/v1/users/me/revenue/statistics [get]
func (h *revenueHandler) GetInstructorStatistics() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}

		stats, err := h.revenueService.GetInstructorStatistics(c.Request.Context(), userID)
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

// GetInstructorStatisticsByDay godoc
// @Summary Get instructor revenue statistics by day
// @Description Get daily revenue statistics for the authenticated instructor for a specific date, including course purchases and subscriptions.
// @Tags statistics
// @Accept json
// @Produce json
// @Param date query string true "Date in YYYY-MM-DD format"
// @Security BearerAuth
// @Success 200 {object} dto.ApiResponse{data=dto.RevenueStatisticsResponse}
// @Failure 400 {object} dto.ApiResponse "Invalid date format, expected YYYY-MM-DD"
// @Failure 401 {object} dto.ApiResponse
// @Failure 403 {object} dto.ApiResponse
// @Failure 500 {object} dto.ApiResponse
// @Router /api/v1/users/me/revenue/statistics/day [get]
func (h *revenueHandler) GetInstructorStatisticsByDay() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}

		var query dto.RevenueByDayQuery
		if err := util.BindAndValidateQuery(c, &query); err != nil {
			_ = c.Error(err)
			return
		}

		date, err := time.Parse("2006-01-02", query.Date)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid date format, expected YYYY-MM-DD"))
			return
		}

		stats, err := h.revenueService.GetInstructorStatisticsByDay(c.Request.Context(), userID, date)
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

// GetInstructorStatisticsByMonth godoc
// @Summary Get instructor revenue statistics by month
// @Description Get monthly revenue statistics for the authenticated instructor for a specific year and month, including course purchases and subscriptions.
// @Tags statistics
// @Accept json
// @Produce json
// @Param year query integer true "Year (e.g., 2026)"
// @Param month query integer true "Month (1-12)"
// @Security BearerAuth
// @Success 200 {object} dto.ApiResponse{data=dto.RevenueStatisticsResponse}
// @Failure 400 {object} dto.ApiResponse "Invalid year or month values"
// @Failure 401 {object} dto.ApiResponse
// @Failure 403 {object} dto.ApiResponse
// @Failure 500 {object} dto.ApiResponse
// @Router /api/v1/users/me/revenue/statistics/month [get]
func (h *revenueHandler) GetInstructorStatisticsByMonth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}

		var query dto.RevenueByMonthQuery
		if err := util.BindAndValidateQuery(c, &query); err != nil {
			_ = c.Error(err)
			return
		}

		stats, err := h.revenueService.GetInstructorStatisticsByMonth(c.Request.Context(), userID, query.Year, query.Month)
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

// GetInstructorStatisticsByYear godoc
// @Summary Get instructor revenue statistics by year
// @Description Get yearly revenue statistics for the authenticated instructor for a specific year, including course purchases and subscriptions.
// @Tags statistics
// @Accept json
// @Produce json
// @Param year query integer true "Year (e.g., 2026)"
// @Security BearerAuth
// @Success 200 {object} dto.ApiResponse{data=dto.RevenueStatisticsResponse}
// @Failure 400 {object} dto.ApiResponse "Invalid year value"
// @Failure 401 {object} dto.ApiResponse
// @Failure 403 {object} dto.ApiResponse
// @Failure 500 {object} dto.ApiResponse
// @Router /api/v1/users/me/revenue/statistics/year [get]
func (h *revenueHandler) GetInstructorStatisticsByYear() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}

		var query dto.RevenueByYearQuery
		if err := util.BindAndValidateQuery(c, &query); err != nil {
			_ = c.Error(err)
			return
		}

		stats, err := h.revenueService.GetInstructorStatisticsByYear(c.Request.Context(), userID, query.Year)
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

// GetInstructorRevenueOverview godoc
// @Summary Get instructor revenue overview card
// @Description Get current month instructor revenue and growth versus previous month for dashboard card.
// @Tags statistics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.ApiResponse{data=dto.RevenueOverviewResponse}
// @Failure 401 {object} dto.ApiResponse
// @Failure 403 {object} dto.ApiResponse
// @Failure 500 {object} dto.ApiResponse
// @Router /api/v1/users/me/revenue/overview [get]
func (h *revenueHandler) GetInstructorRevenueOverview() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}

		stats, err := h.revenueService.GetInstructorRevenueOverview(c.Request.Context(), userID)
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

// GetAdminStatistics godoc
// @Summary Get admin revenue statistics
// @Description Get all-time revenue statistics for the admin dashboard, including course purchases and subscriptions.
// @Tags statistics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.ApiResponse{data=dto.RevenueStatisticsResponse}
// @Failure 401 {object} dto.ApiResponse
// @Failure 403 {object} dto.ApiResponse
// @Failure 500 {object} dto.ApiResponse
// @Router /api/v1/admin/revenue/statistics [get]
func (h *revenueHandler) GetAdminStatistics() gin.HandlerFunc {
	return func(c *gin.Context) {
		stats, err := h.revenueService.GetAdminStatistics(c.Request.Context())
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

// GetAdminStatisticsByDay godoc
// @Summary Get admin revenue statistics by day
// @Description Get daily revenue statistics for the admin dashboard for a specific date, including course purchases and subscriptions across all instructors.
// @Tags statistics
// @Accept json
// @Produce json
// @Param date query string true "Date in YYYY-MM-DD format"
// @Security BearerAuth
// @Success 200 {object} dto.ApiResponse{data=dto.RevenueStatisticsResponse}
// @Failure 400 {object} dto.ApiResponse "Invalid date format, expected YYYY-MM-DD"
// @Failure 401 {object} dto.ApiResponse
// @Failure 403 {object} dto.ApiResponse
// @Failure 500 {object} dto.ApiResponse
// @Router /api/v1/admin/revenue/statistics/day [get]
func (h *revenueHandler) GetAdminStatisticsByDay() gin.HandlerFunc {
	return func(c *gin.Context) {
		var query dto.RevenueByDayQuery
		if err := util.BindAndValidateQuery(c, &query); err != nil {
			_ = c.Error(err)
			return
		}

		date, err := time.Parse("2006-01-02", query.Date)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid date format, expected YYYY-MM-DD"))
			return
		}

		stats, err := h.revenueService.GetAdminStatisticsByDay(c.Request.Context(), date)
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

// GetAdminStatisticsByMonth godoc
// @Summary Get admin revenue statistics by month
// @Description Get monthly revenue statistics for the admin dashboard for a specific year and month, including course purchases and subscriptions across all instructors.
// @Tags statistics
// @Accept json
// @Produce json
// @Param year query integer true "Year (e.g., 2026)"
// @Param month query integer true "Month (1-12)"
// @Security BearerAuth
// @Success 200 {object} dto.ApiResponse{data=dto.RevenueStatisticsResponse}
// @Failure 400 {object} dto.ApiResponse "Invalid year or month values"
// @Failure 401 {object} dto.ApiResponse
// @Failure 403 {object} dto.ApiResponse
// @Failure 500 {object} dto.ApiResponse
// @Router /api/v1/admin/revenue/statistics/month [get]
func (h *revenueHandler) GetAdminStatisticsByMonth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var query dto.RevenueByMonthQuery
		if err := util.BindAndValidateQuery(c, &query); err != nil {
			_ = c.Error(err)
			return
		}

		stats, err := h.revenueService.GetAdminStatisticsByMonth(c.Request.Context(), query.Year, query.Month)
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

// GetAdminStatisticsByYear godoc
// @Summary Get admin revenue statistics by year
// @Description Get yearly revenue statistics for the admin dashboard for a specific year, including course purchases and subscriptions across all instructors.
// @Tags statistics
// @Accept json
// @Produce json
// @Param year query integer true "Year (e.g., 2026)"
// @Security BearerAuth
// @Success 200 {object} dto.ApiResponse{data=dto.RevenueStatisticsResponse}
// @Failure 400 {object} dto.ApiResponse "Invalid year value"
// @Failure 401 {object} dto.ApiResponse
// @Failure 403 {object} dto.ApiResponse
// @Failure 500 {object} dto.ApiResponse
// @Router /api/v1/admin/revenue/statistics/year [get]
func (h *revenueHandler) GetAdminStatisticsByYear() gin.HandlerFunc {
	return func(c *gin.Context) {
		var query dto.RevenueByYearQuery
		if err := util.BindAndValidateQuery(c, &query); err != nil {
			_ = c.Error(err)
			return
		}

		stats, err := h.revenueService.GetAdminStatisticsByYear(c.Request.Context(), query.Year)
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

// GetAdminRevenueOverview godoc
// @Summary Get admin revenue overview card
// @Description Get current month admin revenue and growth versus previous month for dashboard card.
// @Tags statistics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.ApiResponse{data=dto.RevenueOverviewResponse}
// @Failure 401 {object} dto.ApiResponse
// @Failure 403 {object} dto.ApiResponse
// @Failure 500 {object} dto.ApiResponse
// @Router /api/v1/admin/revenue/overview [get]
func (h *revenueHandler) GetAdminRevenueOverview() gin.HandlerFunc {
	return func(c *gin.Context) {
		stats, err := h.revenueService.GetAdminRevenueOverview(c.Request.Context())
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

// GetAdminSalesSegmentation godoc
// @Summary Get admin sales segmentation
// @Description Get sales segmentation for dashboard card between membership subscriptions and single course purchases.
// @Tags statistics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.ApiResponse{data=dto.SalesSegmentationResponse}
// @Failure 401 {object} dto.ApiResponse
// @Failure 403 {object} dto.ApiResponse
// @Failure 500 {object} dto.ApiResponse
// @Router /api/v1/admin/revenue/statistics/sales-segmentation [get]
func (h *revenueHandler) GetAdminSalesSegmentation() gin.HandlerFunc {
	return func(c *gin.Context) {
		stats, err := h.revenueService.GetAdminSalesSegmentation(c.Request.Context())
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

// GetAdminTransactions godoc
// @Summary Get admin transaction table
// @Description Get paginated transactions for admin dashboard table from course purchases and subscriptions, sorted by time descending.
// @Tags statistics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Param sortOrder query string false "Sort direction" Enums(asc,desc) default(desc)
// @Success 200 {object} dto.ApiResponse{data=[]dto.RevenueTransactionItemResponse,metadata=dto.Pagination}
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 403 {object} dto.ApiResponse
// @Failure 500 {object} dto.ApiResponse
// @Router /api/v1/admin/revenue/transactions [get]
func (h *revenueHandler) GetAdminTransactions() gin.HandlerFunc {
	return func(c *gin.Context) {
		var pagingRequest dto.PagingRequest
		if err := c.ShouldBindQuery(&pagingRequest); err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid pagination parameters"))
			return
		}

		pagingRequest.Process()

		data, total, err := h.revenueService.GetAdminTransactions(
			c.Request.Context(),
			pagingRequest.Limit,
			pagingRequest.Offset,
			"created_at",
			pagingRequest.SortOrder,
		)
		if err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = data
		res.Metadata = dto.NewPagination(pagingRequest.Limit, pagingRequest.Offset, int(total), "created_at", pagingRequest.SortOrder)

		c.JSON(http.StatusOK, res)
	}
}

// GetAllTeachersRevenue godoc
// @Summary Get revenue statistics for all teachers
// @Description Get paginated revenue breakdown (gross, stripe fee, net) per instructor. Admin only. Supports optional date range filtering.
// @Tags statistics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start_date query string  false "Start date (YYYY-MM-DD)"
// @Param end_date   query string  false "End date   (YYYY-MM-DD)"
// @Param limit      query int     false "Limit"      default(10)
// @Param offset     query int     false "Offset"     default(0)
// @Param sort_order query string  false "Sort order" Enums(asc,desc) default(desc)
// @Success 200 {object} dto.ApiResponse{data=[]dto.TeacherRevenueItemResponse,metadata=dto.Pagination}
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 403 {object} dto.ApiResponse
// @Failure 500 {object} dto.ApiResponse
// @Router /api/v1/admin/revenue/statistics/teachers/revenue [get]
func (h *revenueHandler) GetAllTeachersRevenue() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.TeacherRevenueFilterRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid query parameters"))
			return
		}
		req.Process()

		data, total, err := h.revenueService.GetAllTeachersRevenue(c.Request.Context(), req)
		if err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = data
		res.Metadata = dto.NewPagination(req.Limit, req.Offset, int(total), "total_amount", req.SortOrder)

		c.JSON(http.StatusOK, res)
	}
}
