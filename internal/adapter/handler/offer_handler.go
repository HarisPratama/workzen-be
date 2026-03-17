package handler

import (
	"bwanews/internal/adapter/handler/request"
	"bwanews/internal/adapter/handler/response"
	"bwanews/internal/core/domain/entity"
	"bwanews/internal/core/service"
	"bwanews/lib/validator"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

type OfferHandler interface {
	CreateOffer(c *fiber.Ctx) error
	UpdateOffer(c *fiber.Ctx) error
	DeleteOffer(c *fiber.Ctx) error
	GetOfferByID(c *fiber.Ctx) error
	GetOffersByTenant(c *fiber.Ctx) error
	GetOffersByCandidate(c *fiber.Ctx) error
	SendOffer(c *fiber.Ctx) error
	WithdrawOffer(c *fiber.Ctx) error
	AcceptOffer(c *fiber.Ctx) error
	RejectOffer(c *fiber.Ctx) error
	NegotiateOffer(c *fiber.Ctx) error
	GetOfferMetrics(c *fiber.Ctx) error
}

type offerHandler struct {
	offerService service.OfferService
}

func NewOfferHandler(offerService service.OfferService) OfferHandler {
	return &offerHandler{
		offerService: offerService,
	}
}

func (h *offerHandler) CreateOffer(c *fiber.Ctx) error {
	var req request.OfferRequest
	claims := c.Locals("user").(*entity.JwtData)
	tenantID := claims.TenantID

	if tenantID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Unauthorized access",
			},
		})
	}

	if err := c.BodyParser(&req); err != nil {
		log.Errorw("[Handler] CreateOffer - BodyParser", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid request body",
			},
		})
	}

	if err := validator.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		})
	}

	offer := entity.Offer{
		TenantID:                 uuid.UUID{},
		CandidateApplicationID:   req.CandidateApplicationID,
		JobTitle:                 req.JobTitle,
		Department:              req.Department,
		OfferType:               entity.OfferType(req.OfferType),
		EmploymentLevel:         req.EmploymentLevel,
		BaseSalary:              req.BaseSalary,
		Currency:                req.Currency,
		SignOnBonus:             req.SignOnBonus,
		AnnualBonus:             req.AnnualBonus,
		BenefitsPackage:         req.BenefitsPackage,
		StockOptions:            req.StockOptions,
		VestingSchedule:         req.VestingSchedule,
		ProbationPeriodDays:    req.ProbationPeriodDays,
		NoticePeriodDays:       req.NoticePeriodDays,
		PaidTimeOffDays:        req.PaidTimeOffDays,
		Status:                  entity.OfferStatusDraft,
		ResponseDeadline:        req.ResponseDeadline,
		InternalNotes:          req.InternalNotes,
	}

	createdOffer, err := h.offerService.CreateOffer(c.Context(), offer)
	if err != nil {
		log.Errorw("[Handler] CreateOffer - Service", err)
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessResponse{
		Meta: response.Meta{
			Status:  true,
			Message: "Offer created successfully",
		},
		Data: createdOffer,
	})
}

func (h *offerHandler) UpdateOffer(c *fiber.Ctx) error {
	id := c.Params("id")
	offerID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid offer ID",
			},
		})
	}

	var req request.OfferUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		log.Errorw("[Handler] UpdateOffer - BodyParser", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid request body",
			},
		})
	}

	if err := validator.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		})
	}

	offer := entity.Offer{
		ID:                   offerID,
		JobTitle:             req.JobTitle,
		Department:           req.Department,
		OfferType:            entity.OfferType(req.OfferType),
		EmploymentLevel:      req.EmploymentLevel,
		BaseSalary:           req.BaseSalary,
		Currency:             req.Currency,
		SignOnBonus:          req.SignOnBonus,
		AnnualBonus:          req.AnnualBonus,
		BenefitsPackage:     req.BenefitsPackage,
		StockOptions:         req.StockOptions,
		VestingSchedule:      req.VestingSchedule,
		ProbationPeriodDays: req.ProbationPeriodDays,
		NoticePeriodDays:    req.NoticePeriodDays,
		PaidTimeOffDays:     req.PaidTimeOffDays,
		InternalNotes:       req.InternalNotes,
	}

	updatedOffer, err := h.offerService.UpdateOffer(c.Context(), offerID, offer)
	if err != nil {
		log.Errorw("[Handler] UpdateOffer - Service", err)
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
		Meta: response.Meta{
			Status:  true,
			Message: "Offer updated successfully",
		},
		Data: updatedOffer,
	})
}

func (h *offerHandler) DeleteOffer(c *fiber.Ctx) error {
	id := c.Params("id")
	offerID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid offer ID",
			},
		})
	}

	if err := h.offerService.DeleteOffer(c.Context(), offerID); err != nil {
		log.Errorw("[Handler] DeleteOffer - Service", err)
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
		Meta: response.Meta{
			Status:  true,
			Message: "Offer deleted successfully",
		},
	})
}

func (h *offerHandler) GetOfferByID(c *fiber.Ctx) error {
	id := c.Params("id")
	offerID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid offer ID",
			},
		})
	}

	offer, err := h.offerService.GetOfferByID(c.Context(), offerID)
	if err != nil {
		log.Errorw("[Handler] GetOfferByID - Service", err)
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
		Meta: response.Meta{
			Status:  true,
			Message: "Success",
		},
		Data: offer,
	})
}

func (h *offerHandler) GetOffersByTenant(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	tenantID := claims.TenantID

	if tenantID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Unauthorized access",
			},
		})
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	offers, total, err := h.offerService.GetOffersByTenant(c.Context(), uuid.UUID{}, page, limit)
	if err != nil {
		log.Errorw("[Handler] GetOffersByTenant - Service", err)
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
		Meta: response.Meta{
			Status:  true,
			Message: "Success",
		},
		Data: response.PaginationData{
			Data:  offers,
			Total: total,
			Page:  page,
			Limit: limit,
		},
	})
}

func (h *offerHandler) GetOffersByCandidate(c *fiber.Ctx) error {
	candidateApplicationID := c.Params("candidateApplicationId")
	candidateUUID, err := uuid.Parse(candidateApplicationID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid candidate application ID",
			},
		})
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	offers, total, err := h.offerService.GetOffersByCandidate(c.Context(), candidateUUID, page, limit)
	if err != nil {
		log.Errorw("[Handler] GetOffersByCandidate - Service", err)
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
		Meta: response.Meta{
			Status:  true,
			Message: "Success",
		},
		Data: response.PaginationData{
			Data:  offers,
			Total: total,
			Page:  page,
			Limit: limit,
		},
	})
}

func (h *offerHandler) SendOffer(c *fiber.Ctx) error {
	id := c.Params("id")
	offerID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid offer ID",
			},
		})
	}

	if err := h.offerService.SendOffer(c.Context(), offerID); err != nil {
		log.Errorw("[Handler] SendOffer - Service", err)
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
		Meta: response.Meta{
			Status:  true,
			Message: "Offer sent successfully",
		},
	})
}

func (h *offerHandler) WithdrawOffer(c *fiber.Ctx) error {
	id := c.Params("id")
	offerID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid offer ID",
			},
		})
	}

	var req request.WithdrawRequest
	if err := c.BodyParser(&req); err != nil {
		log.Errorw("[Handler] WithdrawOffer - BodyParser", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid request body",
			},
		})
	}

	if err := h.offerService.WithdrawOffer(c.Context(), offerID, req.Reason); err != nil {
		log.Errorw("[Handler] WithdrawOffer - Service", err)
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
		Meta: response.Meta{
			Status:  true,
			Message: "Offer withdrawn successfully",
		},
	})
}

func (h *offerHandler) AcceptOffer(c *fiber.Ctx) error {
	id := c.Params("id")
	offerID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid offer ID",
			},
		})
	}

	if err := h.offerService.AcceptOffer(c.Context(), offerID); err != nil {
		log.Errorw("[Handler] AcceptOffer - Service", err)
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
		Meta: response.Meta{
			Status:  true,
			Message: "Offer accepted successfully",
		},
	})
}

func (h *offerHandler) RejectOffer(c *fiber.Ctx) error {
	id := c.Params("id")
	offerID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid offer ID",
			},
		})
	}

	var req request.RejectRequest
	if err := c.BodyParser(&req); err != nil {
		log.Errorw("[Handler] RejectOffer - BodyParser", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				