package service

import (
	"context"
	"math"
	"time"

	"elearning-api/apperror"
	"elearning-api/dto"
	"elearning-api/repository"

	"github.com/google/uuid"
)

type RevenueService interface {
	GetInstructorStatistics(ctx context.Context, userID uuid.UUID) (*dto.RevenueStatisticsResponse, error)
	GetInstructorStatisticsByDay(ctx context.Context, userID uuid.UUID, date time.Time) (*dto.RevenueStatisticsResponse, error)
	GetInstructorStatisticsByMonth(ctx context.Context, userID uuid.UUID, year, month int) (*dto.RevenueStatisticsResponse, error)
	GetInstructorStatisticsByYear(ctx context.Context, userID uuid.UUID, year int) (*dto.RevenueStatisticsResponse, error)
	GetInstructorRevenueOverview(ctx context.Context, userID uuid.UUID) (*dto.RevenueOverviewResponse, error)

	GetAdminStatistics(ctx context.Context) (*dto.RevenueStatisticsResponse, error)
	GetAdminStatisticsByDay(ctx context.Context, date time.Time) (*dto.RevenueStatisticsResponse, error)
	GetAdminStatisticsByMonth(ctx context.Context, year, month int) (*dto.RevenueStatisticsResponse, error)
	GetAdminStatisticsByYear(ctx context.Context, year int) (*dto.RevenueStatisticsResponse, error)
	GetAdminRevenueOverview(ctx context.Context) (*dto.RevenueOverviewResponse, error)

	GetAdminSalesSegmentation(ctx context.Context) (*dto.SalesSegmentationResponse, error)
	GetAdminTransactions(ctx context.Context, limit, offset int, sortBy, sortOrder string) ([]*dto.RevenueTransactionItemResponse, int64, error)
}

type revenueService struct {
	revenueRepository repository.RevenueRepository
}

func NewRevenueService(revenueRepository repository.RevenueRepository) RevenueService {
	return &revenueService{revenueRepository: revenueRepository}
}

func (s *revenueService) GetInstructorStatistics(ctx context.Context, userID uuid.UUID) (*dto.RevenueStatisticsResponse, error) {
	coursePurchase, subscription, err := s.revenueRepository.GetInstructorRevenue(ctx, userID)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve instructor revenue statistics")
	}
	return buildRevenueStatistics(coursePurchase, subscription), nil
}

func (s *revenueService) GetInstructorStatisticsByDay(ctx context.Context, userID uuid.UUID, date time.Time) (*dto.RevenueStatisticsResponse, error) {
	coursePurchase, subscription, err := s.revenueRepository.GetInstructorRevenueByDay(ctx, userID, date)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve instructor daily revenue statistics")
	}
	return buildRevenueStatistics(coursePurchase, subscription), nil
}

func (s *revenueService) GetInstructorStatisticsByMonth(ctx context.Context, userID uuid.UUID, year, month int) (*dto.RevenueStatisticsResponse, error) {
	coursePurchase, subscription, err := s.revenueRepository.GetInstructorRevenueByMonth(ctx, userID, year, month)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve instructor monthly revenue statistics")
	}
	return buildRevenueStatistics(coursePurchase, subscription), nil
}

func (s *revenueService) GetInstructorStatisticsByYear(ctx context.Context, userID uuid.UUID, year int) (*dto.RevenueStatisticsResponse, error) {
	coursePurchase, subscription, err := s.revenueRepository.GetInstructorRevenueByYear(ctx, userID, year)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve instructor yearly revenue statistics")
	}
	return buildRevenueStatistics(coursePurchase, subscription), nil
}

func (s *revenueService) GetInstructorRevenueOverview(ctx context.Context, userID uuid.UUID) (*dto.RevenueOverviewResponse, error) {
	currentStart, previousStart := monthComparisonRange(time.Now().UTC())

	currentPurchase, currentSubscription, err := s.revenueRepository.GetInstructorRevenueByMonth(ctx, userID, currentStart.Year(), int(currentStart.Month()))
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve instructor revenue overview")
	}

	previousPurchase, previousSubscription, err := s.revenueRepository.GetInstructorRevenueByMonth(ctx, userID, previousStart.Year(), int(previousStart.Month()))
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve instructor revenue overview")
	}

	currentTotal := currentPurchase.TotalAmount + currentSubscription.TotalAmount
	previousTotal := previousPurchase.TotalAmount + previousSubscription.TotalAmount

	return buildRevenueOverview(currentTotal, previousTotal), nil
}

func (s *revenueService) GetAdminStatistics(ctx context.Context) (*dto.RevenueStatisticsResponse, error) {
	coursePurchase, subscription, err := s.revenueRepository.GetAdminRevenue(ctx)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve admin revenue statistics")
	}
	return buildRevenueStatistics(coursePurchase, subscription), nil
}

func (s *revenueService) GetAdminStatisticsByDay(ctx context.Context, date time.Time) (*dto.RevenueStatisticsResponse, error) {
	coursePurchase, subscription, err := s.revenueRepository.GetAdminRevenueByDay(ctx, date)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve admin daily revenue statistics")
	}
	return buildRevenueStatistics(coursePurchase, subscription), nil
}

func (s *revenueService) GetAdminStatisticsByMonth(ctx context.Context, year, month int) (*dto.RevenueStatisticsResponse, error) {
	coursePurchase, subscription, err := s.revenueRepository.GetAdminRevenueByMonth(ctx, year, month)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve admin monthly revenue statistics")
	}
	return buildRevenueStatistics(coursePurchase, subscription), nil
}

func (s *revenueService) GetAdminStatisticsByYear(ctx context.Context, year int) (*dto.RevenueStatisticsResponse, error) {
	coursePurchase, subscription, err := s.revenueRepository.GetAdminRevenueByYear(ctx, year)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve admin yearly revenue statistics")
	}
	return buildRevenueStatistics(coursePurchase, subscription), nil
}

func (s *revenueService) GetAdminRevenueOverview(ctx context.Context) (*dto.RevenueOverviewResponse, error) {
	currentStart, previousStart := monthComparisonRange(time.Now().UTC())

	currentPurchase, currentSubscription, err := s.revenueRepository.GetAdminRevenueByMonth(ctx, currentStart.Year(), int(currentStart.Month()))
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve admin revenue overview")
	}

	previousPurchase, previousSubscription, err := s.revenueRepository.GetAdminRevenueByMonth(ctx, previousStart.Year(), int(previousStart.Month()))
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve admin revenue overview")
	}

	currentTotal := currentPurchase.TotalAmount + currentSubscription.TotalAmount
	previousTotal := previousPurchase.TotalAmount + previousSubscription.TotalAmount

	return buildRevenueOverview(currentTotal, previousTotal), nil
}

func (s *revenueService) GetAdminSalesSegmentation(ctx context.Context) (*dto.SalesSegmentationResponse, error) {
	coursePurchase, subscription, err := s.revenueRepository.GetAdminRevenue(ctx)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve admin sales segmentation")
	}

	totalAmount := float64(coursePurchase.TotalAmount + subscription.TotalAmount)
	subscriptionAmount := centsToDollars(subscription.TotalAmount)
	coursePurchaseAmount := centsToDollars(coursePurchase.TotalAmount)

	subscriptionPercent := 0
	coursePurchasePercent := 0
	if totalAmount > 0 {
		subscriptionPercent = int(math.Round((float64(subscription.TotalAmount) / totalAmount) * 100))
		coursePurchasePercent = 100 - subscriptionPercent
	}

	return &dto.SalesSegmentationResponse{
		MembershipSubs: dto.SalesSegmentationItemResponse{
			Percent: subscriptionPercent,
			Amount:  subscriptionAmount,
		},
		SinglePurchases: dto.SalesSegmentationItemResponse{
			Percent: coursePurchasePercent,
			Amount:  coursePurchaseAmount,
		},
	}, nil
}

func buildRevenueOverview(currentTotal, previousTotal int64) *dto.RevenueOverviewResponse {
	currentAmount := centsToDollars(currentTotal)
	previousAmount := centsToDollars(previousTotal)
	growthPct := 0.0

	if previousTotal > 0 {
		growthPct = ((currentAmount - previousAmount) / previousAmount) * 100
	}

	return &dto.RevenueOverviewResponse{
		TotalAmount:    currentAmount,
		PreviousAmount: previousAmount,
		GrowthPct:      math.Round(growthPct*10) / 10,
	}
}

func monthComparisonRange(now time.Time) (time.Time, time.Time) {
	currentStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	previousStart := currentStart.AddDate(0, -1, 0)
	return currentStart, previousStart
}

func (s *revenueService) GetAdminTransactions(ctx context.Context, limit, offset int, sortBy, sortOrder string) ([]*dto.RevenueTransactionItemResponse, int64, error) {
	rows, total, err := s.revenueRepository.GetAdminTransactions(ctx, limit, offset, sortBy, sortOrder)
	if err != nil {
		return nil, 0, apperror.NewInternalServerError("Failed to retrieve admin transactions")
	}

	items := make([]*dto.RevenueTransactionItemResponse, len(rows))
	for i, row := range rows {
		items[i] = &dto.RevenueTransactionItemResponse{
			TransactionID: row.TransactionID,
			User: dto.RevenueUserResponse{
				ID:       row.UserID,
				Email:    row.UserEmail,
				Username: row.UserUsername,
				Name:     row.UserName,
				Avatar:   row.UserAvatar,
			},
			Method:    row.Method,
			Amount:    row.Amount,
			Type:      row.Type,
			CreatedAt: row.CreatedAt,
		}
	}

	return items, total, nil
}

func buildRevenueStatistics(coursePurchase, subscription *repository.RevenueBreakdownCents) *dto.RevenueStatisticsResponse {
	return &dto.RevenueStatisticsResponse{
		CoursePurchase: dto.RevenueBreakdownResponse{
			TotalAmount:     centsToDollars(coursePurchase.TotalAmount),
			InstructorGross: centsToDollars(coursePurchase.InstructorGross),
			PlatformGross:   centsToDollars(coursePurchase.PlatformGross),
			StripeFee:       centsToDollars(coursePurchase.StripeFee),
			InstructorNet:   centsToDollars(coursePurchase.InstructorNet),
			PlatformNet:     centsToDollars(coursePurchase.PlatformNet),
		},
		Subscription: dto.RevenueBreakdownResponse{
			TotalAmount:     centsToDollars(subscription.TotalAmount),
			InstructorGross: centsToDollars(subscription.InstructorGross),
			PlatformGross:   centsToDollars(subscription.PlatformGross),
			StripeFee:       centsToDollars(subscription.StripeFee),
			InstructorNet:   centsToDollars(subscription.InstructorNet),
			PlatformNet:     centsToDollars(subscription.PlatformNet),
		},
		Total: dto.RevenueBreakdownResponse{
			TotalAmount:     centsToDollars(coursePurchase.TotalAmount + subscription.TotalAmount),
			InstructorGross: centsToDollars(coursePurchase.InstructorGross + subscription.InstructorGross),
			PlatformGross:   centsToDollars(coursePurchase.PlatformGross + subscription.PlatformGross),
			StripeFee:       centsToDollars(coursePurchase.StripeFee + subscription.StripeFee),
			InstructorNet:   centsToDollars(coursePurchase.InstructorNet + subscription.InstructorNet),
			PlatformNet:     centsToDollars(coursePurchase.PlatformNet + subscription.PlatformNet),
		},
	}
}

func centsToDollars(value int64) float64 {
	return float64(value) / 100.0
}
