package handler

import (
	"context"
	"workzen-be/internal/adapter/handler/response"
	"workzen-be/internal/core/domain/entity"
	"workzen-be/internal/core/service"
	"workzen-be/lib/conv"
	validatorLib "workzen-be/lib/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type SubscriptionHandler interface {
	// Plans (SUPER_ADMIN)
	GetPlans(c *fiber.Ctx) error
	GetPlanByID(c *fiber.Ctx) error
	CreatePlan(c *fiber.Ctx) error
	UpdatePlan(c *fiber.Ctx) error
	DeletePlan(c *fiber.Ctx) error

	// Tenant Subscriptions
	GetMySubscription(c *fiber.Ctx) error
	GetSubscriptionHistory(c *fiber.Ctx) error
	Subscribe(c *fiber.Ctx) error
	CancelSubscription(c *fiber.Ctx) error
	ChangePlan(c *fiber.Ctx) error
}

type subscriptionHandler struct {
	subscriptionService service.SubscriptionService
}

// ======================== SUBSCRIPTION PLANS (SUPER_ADMIN) ========================

func (h *subscriptionHandler) GetPlans(c *fiber.Ctx) error {
	query := entity.SubscriptionPlanQueryString{
		Limit:  10,
		Page:   1,
		Search: c.Query("search"),
		Tier:   c.Query("tier"),
	}

	if c.Query("limit") != "" {
		if val, err := conv.StringToInt(c.Query("limit")); err == nil {
			query.Limit = val
		}
	}

	if c.Query("page") != "" {
		if val, err := conv.StringToInt(c.Query("page")); err == nil {
			query.Page = val
		}
	}

	if c.Query("order_by") != "" {
		query.OrderBy = c.Query("order_by")
	}

	if c.Query("order_type") != "" {
		query.OrderType = c.Query("order_type")
	}

	if c.Query("is_active") != "" {
		isActive := c.Query("is_active") == "true"
		query.IsActive = &isActive
	}

	results, totalData, totalPages, err := h.subscriptionService.GetSubscriptionPlans(context.Background(), query)
	if err != nil {
		code := "[HANDLER] GetPlans - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	var respData []response.SubscriptionPlanResponse
	for _, item := range results {
		respData = append(respData, response.SubscriptionPlanResponse{
			ID:                  item.ID,
			Name:                item.Name,
			Tier:                item.Tier,
			Description:         item.Description,
			Price:               item.Price,
			BillingCycle:        item.BillingCycle,
			MaxEmployees:        item.MaxEmployees,
			MaxClients:          item.MaxClients,
			MaxManpowerRequests: item.MaxManpowerRequests,
			Features:            item.Features,
			IsActive:            item.IsActive,
		})
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"
	defaultSuccessResponse.Data = respData
	defaultSuccessResponse.Pagination = &response.PaginationResponse{
		TotalRecords: int(totalData),
		Page:         query.Page,
		PerPage:      query.Limit,
		TotalPages:   int(totalPages),
	}

	return c.JSON(defaultSuccessResponse)
}

func (h *subscriptionHandler) GetPlanByID(c *fiber.Ctx) error {
	idParam := c.Params("planID")
	planID, err := conv.StringToInt64(idParam)
	if err != nil {
		code := "[HANDLER] GetPlanByID - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid plan ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	result, err := h.subscriptionService.GetSubscriptionPlanByID(context.Background(), planID)
	if err != nil {
		code := "[HANDLER] GetPlanByID - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusNotFound).JSON(errorResp)
	}

	respData := response.SubscriptionPlanResponse{
		ID:                  result.ID,
		Name:                result.Name,
		Tier:                result.Tier,
		Description:         result.Description,
		Price:               result.Price,
		BillingCycle:        result.BillingCycle,
		MaxEmployees:        result.MaxEmployees,
		MaxClients:          result.MaxClients,
		MaxManpowerRequests: result.MaxManpowerRequests,
		Features:            result.Features,
		IsActive:            result.IsActive,
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"
	defaultSuccessResponse.Data = respData

	return c.JSON(defaultSuccessResponse)
}

func (h *subscriptionHandler) CreatePlan(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.Role != "SUPER_ADMIN" {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Only SUPER_ADMIN can manage subscription plans"
		return c.Status(fiber.StatusForbidden).JSON(errorResp)
	}

	var req entity.SubscriptionPlanRequest
	if err := c.BodyParser(&req); err != nil {
		code := "[HANDLER] CreatePlan - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid request body"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err := validatorLib.ValidateStruct(req); err != nil {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err := h.subscriptionService.CreateSubscriptionPlan(context.Background(), req); err != nil {
		code := "[HANDLER] CreatePlan - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Subscription plan created successfully"
	defaultSuccessResponse.Data = nil

	return c.Status(fiber.StatusCreated).JSON(defaultSuccessResponse)
}

func (h *subscriptionHandler) UpdatePlan(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.Role != "SUPER_ADMIN" {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Only SUPER_ADMIN can manage subscription plans"
		return c.Status(fiber.StatusForbidden).JSON(errorResp)
	}

	idParam := c.Params("planID")
	planID, err := conv.StringToInt64(idParam)
	if err != nil {
		code := "[HANDLER] UpdatePlan - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid plan ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	var req entity.SubscriptionPlanUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		code := "[HANDLER] UpdatePlan - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid request body"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err := h.subscriptionService.UpdateSubscriptionPlan(context.Background(), planID, req); err != nil {
		code := "[HANDLER] UpdatePlan - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Subscription plan updated successfully"
	defaultSuccessResponse.Data = nil

	return c.JSON(defaultSuccessResponse)
}

func (h *subscriptionHandler) DeletePlan(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.Role != "SUPER_ADMIN" {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Only SUPER_ADMIN can manage subscription plans"
		return c.Status(fiber.StatusForbidden).JSON(errorResp)
	}

	idParam := c.Params("planID")
	planID, err := conv.StringToInt64(idParam)
	if err != nil {
		code := "[HANDLER] DeletePlan - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid plan ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err := h.subscriptionService.DeleteSubscriptionPlan(context.Background(), planID); err != nil {
		code := "[HANDLER] DeletePlan - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Subscription plan deleted successfully"
	defaultSuccessResponse.Data = nil

	return c.JSON(defaultSuccessResponse)
}

// ======================== TENANT SUBSCRIPTIONS ========================

func (h *subscriptionHandler) GetMySubscription(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.TenantID == 0 {
		code := "[HANDLER] GetMySubscription - 1"
		log.Errorw(code, "tenant id is empty")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	tenantID := claims.TenantID

	result, err := h.subscriptionService.GetActiveTenantSubscription(context.Background(), int64(tenantID))
	if err != nil {
		code := "[HANDLER] GetMySubscription - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "No active subscription found"
		return c.Status(fiber.StatusNotFound).JSON(errorResp)
	}

	respData := mapSubscriptionToResponse(*result)

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"
	defaultSuccessResponse.Data = respData

	return c.JSON(defaultSuccessResponse)
}

func (h *subscriptionHandler) GetSubscriptionHistory(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.TenantID == 0 {
		code := "[HANDLER] GetSubscriptionHistory - 1"
		log.Errorw(code, "tenant id is empty")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	tenantID := claims.TenantID

	query := entity.TenantSubscriptionQueryString{
		Limit:  10,
		Page:   1,
		Status: c.Query("status"),
	}

	if c.Query("limit") != "" {
		if val, err := conv.StringToInt(c.Query("limit")); err == nil {
			query.Limit = val
		}
	}

	if c.Query("page") != "" {
		if val, err := conv.StringToInt(c.Query("page")); err == nil {
			query.Page = val
		}
	}

	if c.Query("order_by") != "" {
		query.OrderBy = c.Query("order_by")
	}

	if c.Query("order_type") != "" {
		query.OrderType = c.Query("order_type")
	}

	results, totalData, totalPages, err := h.subscriptionService.GetTenantSubscriptions(context.Background(), int64(tenantID), query)
	if err != nil {
		code := "[HANDLER] GetSubscriptionHistory - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	var respData []response.TenantSubscriptionResponse
	for _, item := range results {
		respData = append(respData, mapSubscriptionToResponse(item))
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"
	defaultSuccessResponse.Data = respData
	defaultSuccessResponse.Pagination = &response.PaginationResponse{
		TotalRecords: int(totalData),
		Page:         query.Page,
		PerPage:      query.Limit,
		TotalPages:   int(totalPages),
	}

	return c.JSON(defaultSuccessResponse)
}

func (h *subscriptionHandler) Subscribe(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.TenantID == 0 {
		code := "[HANDLER] Subscribe - 1"
		log.Errorw(code, "tenant id is empty")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	tenantID := claims.TenantID

	var req entity.SubscribeTenantRequest
	if err := c.BodyParser(&req); err != nil {
		code := "[HANDLER] Subscribe - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid request body"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err := validatorLib.ValidateStruct(req); err != nil {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err := h.subscriptionService.SubscribeTenant(context.Background(), int64(tenantID), req); err != nil {
		code := "[HANDLER] Subscribe - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Subscription activated successfully"
	defaultSuccessResponse.Data = nil

	return c.Status(fiber.StatusCreated).JSON(defaultSuccessResponse)
}

func (h *subscriptionHandler) CancelSubscription(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.TenantID == 0 {
		code := "[HANDLER] CancelSubscription - 1"
		log.Errorw(code, "tenant id is empty")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	tenantID := claims.TenantID

	idParam := c.Params("subscriptionID")
	subscriptionID, err := conv.StringToInt64(idParam)
	if err != nil {
		code := "[HANDLER] CancelSubscription - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid subscription ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	var req entity.CancelSubscriptionRequest
	if err := c.BodyParser(&req); err != nil {
		code := "[HANDLER] CancelSubscription - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid request body"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err := validatorLib.ValidateStruct(req); err != nil {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err := h.subscriptionService.CancelSubscription(context.Background(), int64(tenantID), subscriptionID, req); err != nil {
		code := "[HANDLER] CancelSubscription - 5"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Subscription cancelled successfully"
	defaultSuccessResponse.Data = nil

	return c.JSON(defaultSuccessResponse)
}

func (h *subscriptionHandler) ChangePlan(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.TenantID == 0 {
		code := "[HANDLER] ChangePlan - 1"
		log.Errorw(code, "tenant id is empty")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	tenantID := claims.TenantID

	var req entity.ChangeSubscriptionPlanRequest
	if err := c.BodyParser(&req); err != nil {
		code := "[HANDLER] ChangePlan - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid request body"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err := validatorLib.ValidateStruct(req); err != nil {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err := h.subscriptionService.ChangePlan(context.Background(), int64(tenantID), req); err != nil {
		code := "[HANDLER] ChangePlan - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Plan changed successfully"
	defaultSuccessResponse.Data = nil

	return c.JSON(defaultSuccessResponse)
}

// ======================== HELPER ========================

func mapSubscriptionToResponse(item entity.TenantSubscriptionEntity) response.TenantSubscriptionResponse {
	resp := response.TenantSubscriptionResponse{
		ID:            item.ID,
		TenantID:      item.TenantID,
		Status:        item.Status,
		StartDate:     item.StartDate.In(jakartaTZ).Format("2006-01-02"),
		EndDate:       item.EndDate.In(jakartaTZ).Format("2006-01-02"),
		AutoRenew:     item.AutoRenew,
		PaymentMethod: item.PaymentMethod,
		Plan: response.SubscriptionPlanResponse{
			ID:                  item.Plan.ID,
			Name:                item.Plan.Name,
			Tier:                item.Plan.Tier,
			Description:         item.Plan.Description,
			Price:               item.Plan.Price,
			BillingCycle:        item.Plan.BillingCycle,
			MaxEmployees:        item.Plan.MaxEmployees,
			MaxClients:          item.Plan.MaxClients,
			MaxManpowerRequests: item.Plan.MaxManpowerRequests,
		},
	}

	if item.LastPaymentAt != nil {
		resp.LastPaymentAt = item.LastPaymentAt.In(jakartaTZ).Format("2006-01-02")
	}
	if item.CancelledAt != nil {
		resp.CancelledAt = item.CancelledAt.In(jakartaTZ).Format("2006-01-02")
	}
	if item.CancelReason != "" {
		resp.CancelReason = item.CancelReason
	}

	return resp
}

func NewSubscriptionHandler(subscriptionService service.SubscriptionService) SubscriptionHandler {
	return &subscriptionHandler{subscriptionService: subscriptionService}
}
