package service

import (
	"context"
	"errors"
	"time"
	"workzen-be/internal/adapter/repository"
	"workzen-be/internal/core/domain/entity"

	"github.com/gofiber/fiber/v2/log"
)

type SubscriptionService interface {
	// Subscription Plans
	GetSubscriptionPlans(ctx context.Context, query entity.SubscriptionPlanQueryString) ([]entity.SubscriptionPlanEntity, int64, int64, error)
	GetSubscriptionPlanByID(ctx context.Context, id int64) (*entity.SubscriptionPlanEntity, error)
	CreateSubscriptionPlan(ctx context.Context, req entity.SubscriptionPlanRequest) error
	UpdateSubscriptionPlan(ctx context.Context, id int64, req entity.SubscriptionPlanUpdateRequest) error
	DeleteSubscriptionPlan(ctx context.Context, id int64) error

	// Tenant Subscriptions
	GetTenantSubscriptions(ctx context.Context, tenantID int64, query entity.TenantSubscriptionQueryString) ([]entity.TenantSubscriptionEntity, int64, int64, error)
	GetActiveTenantSubscription(ctx context.Context, tenantID int64) (*entity.TenantSubscriptionEntity, error)
	SubscribeTenant(ctx context.Context, tenantID int64, req entity.SubscribeTenantRequest) error
	CancelSubscription(ctx context.Context, tenantID int64, subscriptionID int64, req entity.CancelSubscriptionRequest) error
	ChangePlan(ctx context.Context, tenantID int64, req entity.ChangeSubscriptionPlanRequest) error
}

type subscriptionService struct {
	subscriptionRepo repository.SubscriptionRepository
}

// ======================== SUBSCRIPTION PLANS ========================

func (s *subscriptionService) GetSubscriptionPlans(ctx context.Context, query entity.SubscriptionPlanQueryString) ([]entity.SubscriptionPlanEntity, int64, int64, error) {
	results, totalData, totalPages, err := s.subscriptionRepo.GetSubscriptionPlans(ctx, query)
	if err != nil {
		code = "[SERVICE] GetSubscriptionPlans - 1"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}
	return results, totalData, totalPages, nil
}

func (s *subscriptionService) GetSubscriptionPlanByID(ctx context.Context, id int64) (*entity.SubscriptionPlanEntity, error) {
	result, err := s.subscriptionRepo.GetSubscriptionPlanByID(ctx, id)
	if err != nil {
		code = "[SERVICE] GetSubscriptionPlanByID - 1"
		log.Errorw(code, err)
		return nil, err
	}
	return result, nil
}

func (s *subscriptionService) CreateSubscriptionPlan(ctx context.Context, req entity.SubscriptionPlanRequest) error {
	planEntity := entity.SubscriptionPlanEntity{
		Name:                req.Name,
		Tier:                req.Tier,
		Description:         req.Description,
		Price:               req.Price,
		BillingCycle:        req.BillingCycle,
		MaxEmployees:        req.MaxEmployees,
		MaxClients:          req.MaxClients,
		MaxManpowerRequests: req.MaxManpowerRequests,
		Features:            req.Features,
		IsActive:            req.IsActive,
	}

	if err := s.subscriptionRepo.CreateSubscriptionPlan(ctx, planEntity); err != nil {
		code = "[SERVICE] CreateSubscriptionPlan - 1"
		log.Errorw(code, err)
		return err
	}
	return nil
}

func (s *subscriptionService) UpdateSubscriptionPlan(ctx context.Context, id int64, req entity.SubscriptionPlanUpdateRequest) error {
	if err := s.subscriptionRepo.UpdateSubscriptionPlan(ctx, id, req); err != nil {
		code = "[SERVICE] UpdateSubscriptionPlan - 1"
		log.Errorw(code, err)
		return err
	}
	return nil
}

func (s *subscriptionService) DeleteSubscriptionPlan(ctx context.Context, id int64) error {
	if err := s.subscriptionRepo.DeleteSubscriptionPlan(ctx, id); err != nil {
		code = "[SERVICE] DeleteSubscriptionPlan - 1"
		log.Errorw(code, err)
		return err
	}
	return nil
}

// ======================== TENANT SUBSCRIPTIONS ========================

func (s *subscriptionService) GetTenantSubscriptions(ctx context.Context, tenantID int64, query entity.TenantSubscriptionQueryString) ([]entity.TenantSubscriptionEntity, int64, int64, error) {
	results, totalData, totalPages, err := s.subscriptionRepo.GetTenantSubscriptions(ctx, tenantID, query)
	if err != nil {
		code = "[SERVICE] GetTenantSubscriptions - 1"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}
	return results, totalData, totalPages, nil
}

func (s *subscriptionService) GetActiveTenantSubscription(ctx context.Context, tenantID int64) (*entity.TenantSubscriptionEntity, error) {
	result, err := s.subscriptionRepo.GetActiveTenantSubscription(ctx, tenantID)
	if err != nil {
		code = "[SERVICE] GetActiveTenantSubscription - 1"
		log.Errorw(code, err)
		return nil, err
	}
	return result, nil
}

func (s *subscriptionService) SubscribeTenant(ctx context.Context, tenantID int64, req entity.SubscribeTenantRequest) error {
	// Get the plan to determine billing cycle
	plan, err := s.subscriptionRepo.GetSubscriptionPlanByID(ctx, req.PlanID)
	if err != nil {
		code = "[SERVICE] SubscribeTenant - 1"
		log.Errorw(code, err)
		return errors.New("subscription plan not found")
	}

	if !plan.IsActive {
		return errors.New("subscription plan is not active")
	}

	now := time.Now()
	var endDate time.Time
	switch plan.BillingCycle {
	case "yearly":
		endDate = now.AddDate(1, 0, 0)
	case "custom":
		endDate = now.AddDate(1, 0, 0)
	default:
		endDate = now.AddDate(0, 1, 0)
	}

	subscription := entity.TenantSubscriptionEntity{
		TenantID:      tenantID,
		PlanID:        req.PlanID,
		StartDate:     now,
		EndDate:       endDate,
		AutoRenew:     req.AutoRenew,
		PaymentMethod: req.PaymentMethod,
		LastPaymentAt: &now,
	}

	if err := s.subscriptionRepo.SubscribeTenant(ctx, subscription); err != nil {
		code = "[SERVICE] SubscribeTenant - 2"
		log.Errorw(code, err)
		return err
	}
	return nil
}

func (s *subscriptionService) CancelSubscription(ctx context.Context, tenantID int64, subscriptionID int64, req entity.CancelSubscriptionRequest) error {
	// Verify the subscription belongs to the tenant
	active, err := s.subscriptionRepo.GetActiveTenantSubscription(ctx, tenantID)
	if err != nil {
		code = "[SERVICE] CancelSubscription - 1"
		log.Errorw(code, err)
		return errors.New("no active subscription found")
	}

	if active.ID != subscriptionID {
		return errors.New("subscription does not belong to this tenant")
	}

	if err := s.subscriptionRepo.CancelSubscription(ctx, subscriptionID, req.CancelReason); err != nil {
		code = "[SERVICE] CancelSubscription - 2"
		log.Errorw(code, err)
		return err
	}
	return nil
}

func (s *subscriptionService) ChangePlan(ctx context.Context, tenantID int64, req entity.ChangeSubscriptionPlanRequest) error {
	// Verify the new plan exists and is active
	plan, err := s.subscriptionRepo.GetSubscriptionPlanByID(ctx, req.NewPlanID)
	if err != nil {
		code = "[SERVICE] ChangePlan - 1"
		log.Errorw(code, err)
		return errors.New("subscription plan not found")
	}

	if !plan.IsActive {
		return errors.New("subscription plan is not active")
	}

	if err := s.subscriptionRepo.ChangePlan(ctx, tenantID, req.NewPlanID, req.PaymentMethod); err != nil {
		code = "[SERVICE] ChangePlan - 2"
		log.Errorw(code, err)
		return err
	}
	return nil
}

func NewSubscriptionService(subscriptionRepo repository.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{subscriptionRepo: subscriptionRepo}
}
