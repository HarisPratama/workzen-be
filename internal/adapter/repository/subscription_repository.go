package repository

import (
	"context"
	"fmt"
	"math"
	"time"
	"workzen-be/internal/core/domain/entity"
	"workzen-be/internal/core/domain/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SubscriptionRepository interface {
	// Subscription Plans
	GetSubscriptionPlans(ctx context.Context, query entity.SubscriptionPlanQueryString) ([]entity.SubscriptionPlanEntity, int64, int64, error)
	GetSubscriptionPlanByID(ctx context.Context, id int64) (*entity.SubscriptionPlanEntity, error)
	CreateSubscriptionPlan(ctx context.Context, plan entity.SubscriptionPlanEntity) error
	UpdateSubscriptionPlan(ctx context.Context, id int64, req entity.SubscriptionPlanUpdateRequest) error
	DeleteSubscriptionPlan(ctx context.Context, id int64) error

	// Tenant Subscriptions
	GetTenantSubscriptions(ctx context.Context, tenantID int64, query entity.TenantSubscriptionQueryString) ([]entity.TenantSubscriptionEntity, int64, int64, error)
	GetActiveTenantSubscription(ctx context.Context, tenantID int64) (*entity.TenantSubscriptionEntity, error)
	SubscribeTenant(ctx context.Context, subscription entity.TenantSubscriptionEntity) error
	CancelSubscription(ctx context.Context, subscriptionID int64, reason string) error
	ChangePlan(ctx context.Context, tenantID int64, newPlanID int64, paymentMethod string) error
}

type subscriptionRepository struct {
	db *gorm.DB
}

// ======================== SUBSCRIPTION PLANS ========================

func (r *subscriptionRepository) GetSubscriptionPlans(ctx context.Context, query entity.SubscriptionPlanQueryString) ([]entity.SubscriptionPlanEntity, int64, int64, error) {
	var models []model.SubscriptionPlan
	var countData int64

	orderBy := "created_at"
	if query.OrderBy != "" {
		orderBy = query.OrderBy
	}
	orderType := "desc"
	if query.OrderType != "" {
		orderType = query.OrderType
	}
	order := fmt.Sprintf("%s %s", orderBy, orderType)
	offset := (query.Page - 1) * query.Limit

	sqlMain := r.db.WithContext(ctx).Model(&model.SubscriptionPlan{})

	if query.Search != "" {
		sqlMain = sqlMain.Where("name ILIKE ? OR description ILIKE ?", "%"+query.Search+"%", "%"+query.Search+"%")
	}

	if query.Tier != "" {
		sqlMain = sqlMain.Where("tier = ?", query.Tier)
	}

	if query.IsActive != nil {
		sqlMain = sqlMain.Where("is_active = ?", *query.IsActive)
	}

	if err := sqlMain.Count(&countData).Error; err != nil {
		code = "[REPOSITORY] GetSubscriptionPlans - 1"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}

	totalPages := int64(math.Ceil(float64(countData) / float64(query.Limit)))

	if err := sqlMain.Order(order).Limit(query.Limit).Offset(offset).Find(&models).Error; err != nil {
		code = "[REPOSITORY] GetSubscriptionPlans - 2"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}

	results := make([]entity.SubscriptionPlanEntity, 0, len(models))
	for _, item := range models {
		results = append(results, mapPlanModelToEntity(item))
	}

	return results, countData, totalPages, nil
}

func (r *subscriptionRepository) GetSubscriptionPlanByID(ctx context.Context, id int64) (*entity.SubscriptionPlanEntity, error) {
	var m model.SubscriptionPlan
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		code = "[REPOSITORY] GetSubscriptionPlanByID - 1"
		log.Errorw(code, err)
		return nil, err
	}

	result := mapPlanModelToEntity(m)
	return &result, nil
}

func (r *subscriptionRepository) CreateSubscriptionPlan(ctx context.Context, plan entity.SubscriptionPlanEntity) error {
	m := model.SubscriptionPlan{
		Name:                plan.Name,
		Tier:                plan.Tier,
		Description:         plan.Description,
		Price:               plan.Price,
		BillingCycle:        plan.BillingCycle,
		MaxEmployees:        plan.MaxEmployees,
		MaxClients:          plan.MaxClients,
		MaxManpowerRequests: plan.MaxManpowerRequests,
		Features:            plan.Features,
		IsActive:            plan.IsActive,
	}

	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		code = "[REPOSITORY] CreateSubscriptionPlan - 1"
		log.Errorw(code, err)
		return err
	}
	return nil
}

func (r *subscriptionRepository) UpdateSubscriptionPlan(ctx context.Context, id int64, req entity.SubscriptionPlanUpdateRequest) error {
	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Price > 0 {
		updates["price"] = req.Price
	}
	if req.BillingCycle != "" {
		updates["billing_cycle"] = req.BillingCycle
	}
	if req.MaxEmployees > 0 {
		updates["max_employees"] = req.MaxEmployees
	}
	if req.MaxClients > 0 {
		updates["max_clients"] = req.MaxClients
	}
	if req.MaxManpowerRequests > 0 {
		updates["max_manpower_requests"] = req.MaxManpowerRequests
	}
	if req.Features != "" {
		updates["features"] = req.Features
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if err := r.db.WithContext(ctx).Model(&model.SubscriptionPlan{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		code = "[REPOSITORY] UpdateSubscriptionPlan - 1"
		log.Errorw(code, err)
		return err
	}
	return nil
}

func (r *subscriptionRepository) DeleteSubscriptionPlan(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.SubscriptionPlan{}).Error; err != nil {
		code = "[REPOSITORY] DeleteSubscriptionPlan - 1"
		log.Errorw(code, err)
		return err
	}
	return nil
}

// ======================== TENANT SUBSCRIPTIONS ========================

func (r *subscriptionRepository) GetTenantSubscriptions(ctx context.Context, tenantID int64, query entity.TenantSubscriptionQueryString) ([]entity.TenantSubscriptionEntity, int64, int64, error) {
	var models []model.TenantSubscription
	var countData int64

	orderBy := "created_at"
	if query.OrderBy != "" {
		orderBy = query.OrderBy
	}
	orderType := "desc"
	if query.OrderType != "" {
		orderType = query.OrderType
	}
	order := fmt.Sprintf("%s %s", orderBy, orderType)
	offset := (query.Page - 1) * query.Limit

	sqlMain := r.db.WithContext(ctx).
		Preload(clause.Associations).
		Where("tenant_id = ?", tenantID)

	if query.Status != "" {
		sqlMain = sqlMain.Where("status = ?", query.Status)
	}

	if err := sqlMain.Model(&model.TenantSubscription{}).Count(&countData).Error; err != nil {
		code = "[REPOSITORY] GetTenantSubscriptions - 1"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}

	totalPages := int64(math.Ceil(float64(countData) / float64(query.Limit)))

	if err := sqlMain.Order(order).Limit(query.Limit).Offset(offset).Find(&models).Error; err != nil {
		code = "[REPOSITORY] GetTenantSubscriptions - 2"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}

	results := make([]entity.TenantSubscriptionEntity, 0, len(models))
	for _, item := range models {
		results = append(results, mapSubscriptionModelToEntity(item))
	}

	return results, countData, totalPages, nil
}

func (r *subscriptionRepository) GetActiveTenantSubscription(ctx context.Context, tenantID int64) (*entity.TenantSubscriptionEntity, error) {
	var m model.TenantSubscription
	if err := r.db.WithContext(ctx).
		Preload(clause.Associations).
		Where("tenant_id = ? AND status = ?", tenantID, "ACTIVE").
		Order("created_at DESC").
		First(&m).Error; err != nil {
		code = "[REPOSITORY] GetActiveTenantSubscription - 1"
		log.Errorw(code, err)
		return nil, err
	}

	result := mapSubscriptionModelToEntity(m)
	return &result, nil
}

func (r *subscriptionRepository) SubscribeTenant(ctx context.Context, subscription entity.TenantSubscriptionEntity) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Deactivate any existing active subscription
		tx.Model(&model.TenantSubscription{}).
			Where("tenant_id = ? AND status = ?", subscription.TenantID, "ACTIVE").
			Updates(map[string]interface{}{
				"status":     "EXPIRED",
				"updated_at": time.Now(),
			})

		m := model.TenantSubscription{
			TenantID:      subscription.TenantID,
			PlanID:        subscription.PlanID,
			Status:        "ACTIVE",
			StartDate:     subscription.StartDate,
			EndDate:       subscription.EndDate,
			AutoRenew:     subscription.AutoRenew,
			PaymentMethod: subscription.PaymentMethod,
			LastPaymentAt: subscription.LastPaymentAt,
		}

		if err := tx.Create(&m).Error; err != nil {
			code = "[REPOSITORY] SubscribeTenant - 1"
			log.Errorw(code, err)
			return err
		}

		// Update tenant plan
		var plan model.SubscriptionPlan
		if err := tx.Where("id = ?", subscription.PlanID).First(&plan).Error; err != nil {
			code = "[REPOSITORY] SubscribeTenant - 2"
			log.Errorw(code, err)
			return err
		}

		if err := tx.Model(&model.Tenant{}).Where("id = ?", subscription.TenantID).
			Update("plan", plan.Tier).Error; err != nil {
			code = "[REPOSITORY] SubscribeTenant - 3"
			log.Errorw(code, err)
			return err
		}

		return nil
	})
}

func (r *subscriptionRepository) CancelSubscription(ctx context.Context, subscriptionID int64, reason string) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status":        "CANCELLED",
		"cancelled_at":  &now,
		"cancel_reason": reason,
		"auto_renew":    false,
		"updated_at":    now,
	}

	if err := r.db.WithContext(ctx).Model(&model.TenantSubscription{}).
		Where("id = ?", subscriptionID).Updates(updates).Error; err != nil {
		code = "[REPOSITORY] CancelSubscription - 1"
		log.Errorw(code, err)
		return err
	}
	return nil
}

func (r *subscriptionRepository) ChangePlan(ctx context.Context, tenantID int64, newPlanID int64, paymentMethod string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Cancel current active subscription
		now := time.Now()
		tx.Model(&model.TenantSubscription{}).
			Where("tenant_id = ? AND status = ?", tenantID, "ACTIVE").
			Updates(map[string]interface{}{
				"status":     "EXPIRED",
				"updated_at": now,
			})

		// Get new plan
		var plan model.SubscriptionPlan
		if err := tx.Where("id = ?", newPlanID).First(&plan).Error; err != nil {
			code = "[REPOSITORY] ChangePlan - 1"
			log.Errorw(code, err)
			return err
		}

		// Determine end date based on billing cycle
		var endDate time.Time
		switch plan.BillingCycle {
		case "yearly":
			endDate = now.AddDate(1, 0, 0)
		case "custom":
			endDate = now.AddDate(1, 0, 0)
		default:
			endDate = now.AddDate(0, 1, 0)
		}

		// Create new subscription
		m := model.TenantSubscription{
			TenantID:      tenantID,
			PlanID:        newPlanID,
			Status:        "ACTIVE",
			StartDate:     now,
			EndDate:       endDate,
			AutoRenew:     true,
			PaymentMethod: paymentMethod,
			LastPaymentAt: &now,
		}

		if err := tx.Create(&m).Error; err != nil {
			code = "[REPOSITORY] ChangePlan - 2"
			log.Errorw(code, err)
			return err
		}

		// Update tenant plan
		if err := tx.Model(&model.Tenant{}).Where("id = ?", tenantID).
			Update("plan", plan.Tier).Error; err != nil {
			code = "[REPOSITORY] ChangePlan - 3"
			log.Errorw(code, err)
			return err
		}

		return nil
	})
}

// ======================== HELPERS ========================

func mapPlanModelToEntity(m model.SubscriptionPlan) entity.SubscriptionPlanEntity {
	return entity.SubscriptionPlanEntity{
		ID:                  m.ID,
		Name:                m.Name,
		Tier:                m.Tier,
		Description:         m.Description,
		Price:               m.Price,
		BillingCycle:        m.BillingCycle,
		MaxEmployees:        m.MaxEmployees,
		MaxClients:          m.MaxClients,
		MaxManpowerRequests: m.MaxManpowerRequests,
		Features:            m.Features,
		IsActive:            m.IsActive,
		CreatedAt:           m.CreatedAt,
		UpdatedAt:           m.UpdatedAt,
	}
}

func mapSubscriptionModelToEntity(m model.TenantSubscription) entity.TenantSubscriptionEntity {
	return entity.TenantSubscriptionEntity{
		ID:            m.ID,
		TenantID:      m.TenantID,
		PlanID:        m.PlanID,
		Status:        m.Status,
		StartDate:     m.StartDate,
		EndDate:       m.EndDate,
		AutoRenew:     m.AutoRenew,
		PaymentMethod: m.PaymentMethod,
		LastPaymentAt: m.LastPaymentAt,
		CancelledAt:   m.CancelledAt,
		CancelReason:  m.CancelReason,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
		Plan: entity.SubscriptionPlanEntity{
			ID:           m.Plan.ID,
			Name:         m.Plan.Name,
			Tier:         m.Plan.Tier,
			Description:  m.Plan.Description,
			Price:        m.Plan.Price,
			BillingCycle: m.Plan.BillingCycle,
		},
		Tenant: entity.TenantEntity{
			ID:          m.Tenant.ID,
			CompanyName: m.Tenant.CompanyName,
			Plan:        m.Tenant.Plan,
			Status:      m.Tenant.Status,
		},
	}
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &subscriptionRepository{db: db}
}
