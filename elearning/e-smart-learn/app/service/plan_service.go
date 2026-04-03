package service

import (
	"context"
	"elearning-api/apperror"
	"elearning-api/config"
	"elearning-api/consts"
	"elearning-api/dto"
	"elearning-api/model"
	"elearning-api/repository"
	"elearning-api/util"
	"encoding/json"
	"math"
	"strings"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v83"
	stripeprice "github.com/stripe/stripe-go/v83/price"
	stripeproduct "github.com/stripe/stripe-go/v83/product"
	stripesubscription "github.com/stripe/stripe-go/v83/subscription"
	"gorm.io/gorm"
)

type PlanService interface {
	Create(ctx context.Context, req dto.CreatePlanRequest) (*dto.PlanResponse, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdatePlanRequest) (*dto.PlanResponse, error)
	Activate(ctx context.Context, id uuid.UUID) (*dto.PlanResponse, error)
	Deactivate(ctx context.Context, id uuid.UUID) (*dto.PlanResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error

	GetByID(ctx context.Context, id uuid.UUID) (*dto.PlanResponse, error)
	GetActivePlans(ctx context.Context) ([]*dto.PlanResponse, error)
	GetList(ctx context.Context, request dto.ListPlanRequest) ([]*dto.PlanResponse, int64, error)
}

type planService struct {
	planRepository         repository.PlanRepository
	subscriptionRepository repository.SubscriptionRepository
	stripeConfig           *config.StripeConfig
}

func NewPlanService(
	planRepository repository.PlanRepository,
	subscriptionRepository repository.SubscriptionRepository,
	stripeConfig *config.StripeConfig,
) PlanService {
	stripe.Key = stripeConfig.SecretKey
	return &planService{
		planRepository:         planRepository,
		subscriptionRepository: subscriptionRepository,
		stripeConfig:           stripeConfig,
	}
}

func (s *planService) Create(ctx context.Context, req dto.CreatePlanRequest) (*dto.PlanResponse, error) {
	if req.Price < 0 {
		return nil, apperror.NewBadRequestError("price must be greater than or equal to 0")
	}
	exists, err := s.planRepository.CheckExists(ctx, "name = ?", req.Name)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to validate plan name")
	}
	if exists {
		return nil, apperror.NewDuplicateEntryError("plan with this name already exists")
	}

	currency := strings.ToLower(strings.TrimSpace(s.stripeConfig.DefaultCurrency))
	if req.Currency != nil && strings.TrimSpace(*req.Currency) != "" {
		currency = strings.ToLower(strings.TrimSpace(*req.Currency))
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	productParams := &stripe.ProductParams{
		Name:        stripe.String(req.Name),
		Description: stripe.String(""),
		Active:      stripe.Bool(isActive),
	}
	if req.Description != nil {
		productParams.Description = stripe.String(strings.TrimSpace(*req.Description))
	}

	stripeProduct, err := stripeproduct.New(productParams)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to create Stripe product")
	}

	unitAmount := int64(math.Round(req.Price * 100))
	if unitAmount <= 0 {
		return nil, apperror.NewBadRequestError("price must be greater than 0")
	}

	interval := "month"
	if req.BillingCycle == consts.BillingCycleYearly {
		interval = "year"
	}

	priceParams := &stripe.PriceParams{
		Currency:   stripe.String(currency),
		UnitAmount: stripe.Int64(unitAmount),
		Product:    stripe.String(stripeProduct.ID),
		Recurring: &stripe.PriceRecurringParams{
			Interval: stripe.String(interval),
		},
		Active: stripe.Bool(isActive),
	}

	stripePrice, err := stripeprice.New(priceParams)
	if err != nil {
		_, _ = stripeproduct.Update(stripeProduct.ID, &stripe.ProductParams{Active: stripe.Bool(false)})
		return nil, apperror.NewInternalServerError("Failed to create Stripe price")
	}

	plan := &model.Plan{
		Name:            strings.TrimSpace(req.Name),
		Description:     "",
		BillingCycle:    string(req.BillingCycle),
		Price:           req.Price,
		StripePriceID:   stripePrice.ID,
		StripeProductID: stripeProduct.ID,
		Currency:        currency,
		IsActive:        isActive,
	}
	if req.Description != nil {
		plan.Description = strings.TrimSpace(*req.Description)
	}

	// Handle AccessFeatures
	if len(req.AccessFeatures) > 0 {
		featuresJSON, _ := json.Marshal(req.AccessFeatures)
		plan.AccessFeatures = string(featuresJSON)
	} else {
		plan.AccessFeatures = "[]"
	}

	created, err := s.planRepository.Create(ctx, plan)
	if err != nil {
		_, _ = stripeprice.Update(stripePrice.ID, &stripe.PriceParams{Active: stripe.Bool(false)})
		_, _ = stripeproduct.Update(stripeProduct.ID, &stripe.ProductParams{Active: stripe.Bool(false)})
		return nil, apperror.NewInternalServerError("Failed to save plan")
	}

	return dto.NewPlanDetailResponse(created), nil
}

func (s *planService) Update(ctx context.Context, id uuid.UUID, req dto.UpdatePlanRequest) (*dto.PlanResponse, error) {
	existing, err := s.planRepository.FindByID(ctx, id, nil)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get plan")
	}
	if existing == nil {
		return nil, apperror.NewNotFoundError("plan not found")
	}

	updates := map[string]any{}

	newName := existing.Name
	if req.Name != nil {
		newName = strings.TrimSpace(*req.Name)
	}
	newDescription := existing.Description
	if req.Description != nil {
		newDescription = strings.TrimSpace(*req.Description)
	}
	newBillingCycle := existing.BillingCycle
	if req.BillingCycle != nil {
		newBillingCycle = string(*req.BillingCycle)
	}
	newPrice := existing.Price
	if req.Price != nil {
		newPrice = *req.Price
	}
	newCurrency := strings.ToLower(existing.Currency)
	if req.Currency != nil && strings.TrimSpace(*req.Currency) != "" {
		newCurrency = strings.ToLower(strings.TrimSpace(*req.Currency))
	}
	billingCycleChanged := req.BillingCycle != nil && newBillingCycle != existing.BillingCycle
	priceValueChanged := req.Price != nil && math.Abs(newPrice-existing.Price) > 1e-9
	currencyChanged := req.Currency != nil && newCurrency != strings.ToLower(existing.Currency)
	if req.Name != nil {
		exists, err := s.planRepository.CheckExists(ctx, "name = ? AND id <> ?", newName, id)
		if err != nil {
			return nil, apperror.NewInternalServerError("Failed to validate plan name")
		}
		if exists {
			return nil, apperror.NewDuplicateEntryError("Plan with this name already exists")
		}
		updates["name"] = newName
	}
	if req.Description != nil {
		updates["description"] = newDescription
	}
	if billingCycleChanged {
		updates["billing_cycle"] = newBillingCycle
	}
	if priceValueChanged {
		updates["price"] = newPrice
	}
	if currencyChanged {
		updates["currency"] = newCurrency
	}

	// Handle AccessFeatures
	if len(req.AccessFeatures) > 0 {
		featuresJSON, _ := json.Marshal(req.AccessFeatures)
		updates["access_features"] = string(featuresJSON)
	}

	stripePriceObj, err := stripeprice.Get(existing.StripePriceID, nil)
	if err != nil || stripePriceObj == nil {
		return nil, apperror.NewInternalServerError("Failed to load Stripe price")
	}

	productID := existing.StripeProductID
	if productID == "" {
		return nil, apperror.NewInternalServerError("Stripe product ID not found for plan")
	}

	if req.Name != nil || req.Description != nil {
		productParams := &stripe.ProductParams{}
		if req.Name != nil {
			productParams.Name = stripe.String(newName)
		}
		if req.Description != nil {
			productParams.Description = stripe.String(newDescription)
		}
		if _, err := stripeproduct.Update(productID, productParams); err != nil {
			return nil, apperror.NewInternalServerError("Failed to update Stripe product")
		}
	}

	priceChanged := billingCycleChanged || priceValueChanged || currencyChanged
	if priceChanged {
		unitAmount := int64(math.Round(newPrice * 100))
		if unitAmount <= 0 {
			return nil, apperror.NewBadRequestError("price must be greater than 0")
		}

		interval := "month"
		if newBillingCycle == string(consts.BillingCycleYearly) {
			interval = "year"
		}

		newStripePrice, err := stripeprice.New(&stripe.PriceParams{
			Currency:   stripe.String(newCurrency),
			UnitAmount: stripe.Int64(unitAmount),
			Product:    stripe.String(productID),
			Recurring: &stripe.PriceRecurringParams{
				Interval: stripe.String(interval),
			},
			Active: stripe.Bool(existing.IsActive),
		})
		if err != nil {
			return nil, apperror.NewInternalServerError("Failed to create new Stripe price")
		}

		_, _ = stripeprice.Update(existing.StripePriceID, &stripe.PriceParams{Active: stripe.Bool(false)})
		updates["stripe_price_id"] = newStripePrice.ID
	}

	updated, err := s.planRepository.Update(ctx, id, updates)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to update plan")
	}

	// Only cancel auto-renew subscriptions if billing-related fields changed
	if priceChanged {
		if err := s.cancelPlanAutoRenewals(ctx, id); err != nil {
			return nil, err
		}
	}

	return dto.NewPlanDetailResponse(updated), nil
}

func (s *planService) Activate(ctx context.Context, id uuid.UUID) (*dto.PlanResponse, error) {
	existing, err := s.planRepository.FindByID(ctx, id, nil)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get plan")
	}
	if existing == nil {
		return nil, apperror.NewNotFoundError("Plan not found")
	}
	if existing.IsActive {
		return nil, apperror.NewBadRequestError("Plan is already activated")
	}

	if existing.StripeProductID != "" {
		_, _ = stripeproduct.Update(existing.StripeProductID, &stripe.ProductParams{Active: stripe.Bool(true)})
	}
	if existing.StripePriceID != "" {
		_, _ = stripeprice.Update(existing.StripePriceID, &stripe.PriceParams{Active: stripe.Bool(true)})
	}

	updated, err := s.planRepository.Update(ctx, id, map[string]any{"is_active": true})
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to activate plan")
	}

	return dto.NewPlanDetailResponse(updated), nil
}

func (s *planService) Deactivate(ctx context.Context, id uuid.UUID) (*dto.PlanResponse, error) {
	existing, err := s.planRepository.FindByID(ctx, id, nil)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get plan")
	}
	if existing == nil {
		return nil, apperror.NewNotFoundError("Plan not found")
	}
	if !existing.IsActive {
		return nil, apperror.NewBadRequestError("Plan is already deactivated")
	}

	if existing.StripeProductID != "" {
		_, _ = stripeproduct.Update(existing.StripeProductID, &stripe.ProductParams{Active: stripe.Bool(false)})
	}
	if existing.StripePriceID != "" {
		_, _ = stripeprice.Update(existing.StripePriceID, &stripe.PriceParams{Active: stripe.Bool(false)})
	}

	updated, err := s.planRepository.Update(ctx, id, map[string]any{"is_active": false})
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to deactivate plan")
	}

	if err := s.cancelPlanAutoRenewals(ctx, id); err != nil {
		return nil, err
	}

	return dto.NewPlanDetailResponse(updated), nil
}

func (s *planService) Delete(ctx context.Context, id uuid.UUID) error {
	existing, err := s.planRepository.FindByID(ctx, id, nil)
	if err != nil {
		return apperror.NewInternalServerError("Failed to get plan")
	}
	if existing == nil {
		return apperror.NewNotFoundError("Plan not found")
	}

	subCount, err := s.subscriptionRepository.CountByPlanID(ctx, id)
	if err != nil {
		return apperror.NewInternalServerError("Failed to verify plan subscriptions")
	}
	if subCount > 0 {
		return apperror.NewBadRequestError("Plan already has subscriptions and cannot be deleted")
	}

	if existing.StripePriceID != "" {
		if _, err := stripeprice.Update(existing.StripePriceID, &stripe.PriceParams{Active: stripe.Bool(false)}); err != nil {
			return apperror.NewInternalServerError("Failed to deactivate Stripe price")
		}
	}

	if existing.StripeProductID != "" {
		if _, err := stripeproduct.Del(existing.StripeProductID, nil); err != nil {
			if _, updateErr := stripeproduct.Update(existing.StripeProductID, &stripe.ProductParams{Active: stripe.Bool(false)}); updateErr != nil {
				return apperror.NewInternalServerError("Failed to archive Stripe product")
			}
		}
	}

	err = s.planRepository.Delete(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return apperror.NewNotFoundError("Plan not found")
		}
		return apperror.NewInternalServerError("Failed to delete plan")
	}
	return nil
}

func (s *planService) cancelPlanAutoRenewals(ctx context.Context, planID uuid.UUID) error {
	subscriptions, err := s.subscriptionRepository.ListAutoRenewByPlanID(ctx, planID, nil)
	if err != nil {
		return apperror.NewInternalServerError("Failed to load subscriptions for plan")
	}

	for _, sub := range subscriptions {
		if sub == nil {
			continue
		}

		if sub.StripeSubscriptionID != "" {
			stripeSub, err := stripesubscription.Update(sub.StripeSubscriptionID, &stripe.SubscriptionParams{
				CancelAtPeriodEnd: stripe.Bool(true),
			})
			if err != nil {
				return apperror.NewInternalServerError("Failed to disable auto-renew on Stripe subscriptions")
			}

			sub.CancelAtPeriodEnd = true
			if stripeSub != nil {
				sub.CurrentPeriodStart = nil
				sub.CurrentPeriodEnd = nil
				if stripeSub.Items != nil && len(stripeSub.Items.Data) > 0 {
					sub.CurrentPeriodStart = stripeTimestamp(stripeSub.Items.Data[0].CurrentPeriodStart)
					sub.CurrentPeriodEnd = stripeTimestamp(stripeSub.Items.Data[0].CurrentPeriodEnd)
				}
			}
		} else {
			sub.CancelAtPeriodEnd = true
		}

		if _, err := s.subscriptionRepository.Updates(ctx, sub); err != nil {
			return apperror.NewInternalServerError("Failed to update subscription auto-renew flag")
		}
	}

	return nil
}

func (s *planService) GetByID(ctx context.Context, id uuid.UUID) (*dto.PlanResponse, error) {
	plan, err := s.planRepository.FindByID(ctx, id, nil)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get plan")
	}
	if plan == nil {
		return nil, apperror.NewNotFoundError("Plan not found")
	}
	return dto.NewPlanDetailResponse(plan), nil
}

func (s *planService) GetActivePlans(ctx context.Context) ([]*dto.PlanResponse, error) {
	plans, err := s.planRepository.ListActivePlans(ctx)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get plans")
	}

	return dto.NewListPlanResponse(plans), nil
}

func (s *planService) GetList(ctx context.Context, request dto.ListPlanRequest) ([]*dto.PlanResponse, int64, error) {
	// Build sort query
	orderQuery := buildPlanSortQuery(request.SortBy, request.SortOrder)

	// Build filter query
	query, args := buildPlanQuery(request)

	plans, total, err := s.planRepository.List(
		ctx,
		request.Limit,
		request.Offset,
		orderQuery,
		query,
		nil,
		args...,
	)

	if err != nil {
		return nil, 0, apperror.NewInternalServerError("Failed to retrieve plans")
	}

	return dto.NewListPlanResponse(plans), total, nil
}

func buildPlanSortQuery(sortBy string, sortOrder string) string {
	defaultSort := "created_at DESC"

	if sortBy == "" {
		return defaultSort
	}

	if sortOrder == "" {
		sortOrder = "DESC"
	}

	allowedSort := map[string]bool{
		"created_at": true,
		"updated_at": true,
	}

	if !allowedSort[sortBy] {
		return defaultSort
	}

	return sortBy + " " + strings.ToUpper(sortOrder)
}

func buildPlanQuery(request dto.ListPlanRequest) (string, []any) {
	var conditions []string
	var args []any

	util.AddEqualBoolCondition(&conditions, &args, "is_active", request.IsActive)
	query := strings.Join(conditions, " AND ")
	return query, args
}
