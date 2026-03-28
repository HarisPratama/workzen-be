package middleware

import (
	"fmt"
	"workzen-be/internal/adapter/handler/response"
	"workzen-be/internal/core/domain/entity"
	"workzen-be/internal/core/domain/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// QuotaResource represents a resource type that has quota limits.
type QuotaResource string

const (
	ResourceEmployee        QuotaResource = "employee"
	ResourceClient          QuotaResource = "client"
	ResourceManpowerRequest QuotaResource = "manpower_request"
)

// defaultFreeLimits are applied when a tenant has no active subscription.
var defaultFreeLimits = map[QuotaResource]int{
	ResourceEmployee:        5,
	ResourceClient:          3,
	ResourceManpowerRequest: 5,
}

// resourceConfig maps each QuotaResource to its DB table and display name.
type resourceConfig struct {
	tableName   string
	displayName string
	pluralName  string
}

var resourceConfigs = map[QuotaResource]resourceConfig{
	ResourceEmployee: {
		tableName:   "employees",
		displayName: "Employee",
		pluralName:  "employees",
	},
	ResourceClient: {
		tableName:   "clients",
		displayName: "Client",
		pluralName:  "clients",
	},
	ResourceManpowerRequest: {
		tableName:   "manpower_requests",
		displayName: "Manpower request",
		pluralName:  "manpower requests",
	},
}

// QuotaChecker provides middleware to enforce subscription-based resource limits.
type QuotaChecker struct {
	db *gorm.DB
}

// NewQuotaChecker creates a new QuotaChecker instance.
func NewQuotaChecker(db *gorm.DB) *QuotaChecker {
	return &QuotaChecker{db: db}
}

// CheckQuota returns a Fiber middleware that checks if the tenant has exceeded
// the quota for the given resource based on their active subscription plan.
func (q *QuotaChecker) CheckQuota(resource QuotaResource) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var errorResponse response.ErrorResponseDefault

		claims, ok := c.Locals("user").(*entity.JwtData)
		if !ok || claims.TenantID == 0 {
			errorResponse.Meta.Status = false
			errorResponse.Meta.Message = "Unauthorized access"
			return c.Status(fiber.StatusUnauthorized).JSON(errorResponse)
		}

		tenantID := int64(claims.TenantID)

		maxLimit, err := q.getMaxLimit(tenantID, resource)
		if err != nil {
			errorResponse.Meta.Status = false
			errorResponse.Meta.Message = "Failed to check subscription quota, please try again later"
			return c.Status(fiber.StatusInternalServerError).JSON(errorResponse)
		}

		// 0 means unlimited
		if maxLimit == 0 {
			return c.Next()
		}

		currentCount, err := q.getCurrentCount(tenantID, resource)
		if err != nil {
			errorResponse.Meta.Status = false
			errorResponse.Meta.Message = "Failed to check resource usage, please try again later"
			return c.Status(fiber.StatusInternalServerError).JSON(errorResponse)
		}

		if currentCount >= int64(maxLimit) {
			cfg := resourceConfigs[resource]
			errorResponse.Meta.Status = false
			errorResponse.Meta.Message = fmt.Sprintf(
				"%s limit reached (%d/%d). Please upgrade your subscription plan to add more %s.",
				cfg.displayName, currentCount, maxLimit, cfg.pluralName,
			)
			return c.Status(fiber.StatusForbidden).JSON(errorResponse)
		}

		return c.Next()
	}
}

// getMaxLimit retrieves the quota limit for a resource from the tenant's active subscription plan.
// Falls back to defaultFreeLimits if no active subscription exists.
func (q *QuotaChecker) getMaxLimit(tenantID int64, resource QuotaResource) (int, error) {
	var subscription model.TenantSubscription

	err := q.db.
		Preload("Plan").
		Where("tenant_id = ? AND status = ?", tenantID, "ACTIVE").
		Order("created_at DESC").
		First(&subscription).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return defaultFreeLimits[resource], nil
		}
		return 0, err
	}

	switch resource {
	case ResourceEmployee:
		return subscription.Plan.MaxEmployees, nil
	case ResourceClient:
		return subscription.Plan.MaxClients, nil
	case ResourceManpowerRequest:
		return subscription.Plan.MaxManpowerRequests, nil
	default:
		return 0, nil
	}
}

// getCurrentCount counts the current number of resources for the given tenant.
func (q *QuotaChecker) getCurrentCount(tenantID int64, resource QuotaResource) (int64, error) {
	cfg, ok := resourceConfigs[resource]
	if !ok {
		return 0, fmt.Errorf("unknown resource: %s", resource)
	}

	var count int64
	err := q.db.Table(cfg.tableName).Where("tenant_id = ?", tenantID).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}
