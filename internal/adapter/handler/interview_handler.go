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

type InterviewHandler interface {
	CreateInterview(c *fiber.Ctx) error
	UpdateInterview(c *fiber.Ctx) error
	DeleteInterview(c *fiber.Ctx) error
	GetInterviewByID(c *fiber.Ctx) error
	GetInterviewsByTenant(c *fiber.Ctx) error
	GetInterviewsByCandidate(c *fiber.Ctx) error
	GetInterviewsByInterviewer(c *fiber.Ctx) error
	ScheduleInterview(c *fiber.Ctx) error
	RescheduleInterview(c *fiber.Ctx) error
	CancelInterview(c *fiber.Ctx) error
	CompleteInterview(c *fiber.Ctx) error
	SubmitFeedback(c *fiber.Ctx) error
	GetInterviewMetrics(c *fiber.Ctx) error
}

type interviewHandler struct {
	interviewService service.InterviewService
}

func NewInterviewHandler(interviewService service.InterviewService) InterviewHandler {
	return &interviewHandler{
		interviewService: interviewService,
	}
}

func (h *interviewHandler) CreateInterview(c *fiber.Ctx) error {
	var req request.InterviewRequest
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
		log.Errorw("[Handler] CreateInterview - BodyParser", err)
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

	scheduledAt, _ := time.Parse(time.RFC3339, req.ScheduledAt)

	interview := entity.Interview{
		TenantID:                 uuid.UUID{},
		CandidateApplicationID:   req.CandidateApplicationID,
		InterviewerID:           req.InterviewerID,
		ScheduledAt:             scheduledAt,
		DurationMinutes:         req.DurationMinutes,
		Type:                    entity.InterviewType(req.Type),
		Status:                  entity.InterviewStatusScheduled,
		Location:                req.Location,
		MeetingLink:             req.MeetingLink,
	}

	createdInterview, err := h.interviewService.CreateInterview(c.Context(), interview)
	if err != nil {
		log.Errorw("[Handler] CreateInterview - Service", err)
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
			Message: "Interview created successfully",
		},
		Data: createdInterview,
	})
}

func (h *interviewHandler) UpdateInterview(c *fiber.Ctx) error {
	id := c.Params("id")
	interviewID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid interview ID",
			},
		})
	}

	var req request.InterviewUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		log.Errorw("[Handler] UpdateInterview - BodyParser", err)
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

	scheduledAt, _ := time.Parse(time.RFC3339, req.ScheduledAt)

	interview := entity.Interview{
		ID:              interviewID,
		ScheduledAt:     scheduledAt,
		DurationMinutes: req.DurationMinutes,
		Type:            entity.InterviewType(req.Type),
		Location:        req.Location,
		MeetingLink:     req.MeetingLink,
	}

	updatedInterview, err := h.interviewService.UpdateInterview(c.Context(), interviewID, interview)
	if err != nil {
		log.Errorw("[Handler] UpdateInterview - Service", err)
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
			Message: "Interview updated successfully",
		},
		Data: updatedInterview,
	})
}

func (h *interviewHandler) DeleteInterview(c *fiber.Ctx) error {
	id := c.Params("id")
	interviewID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid interview ID",
			},
		})
	}

	if err := h.interviewService.DeleteInterview(c.Context(), interviewID); err != nil {
		log.Errorw("[Handler] DeleteInterview - Service", err)
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
			Message: "Interview deleted successfully",
		},
	})
}

func (h *interviewHandler) GetInterviewByID(c *fiber.Ctx) error {
	id := c.Params("id")
	interviewID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid interview ID",
			},
		})
	}

	interview, err := h.interviewService.GetInterviewByID(c.Context(), interviewID)
	if err != nil {
		log.Errorw("[Handler] GetInterviewByID - Service", err)
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
		Data: interview,
	})
}

func (h *interviewHandler) GetInterviewsByTenant(c *fiber.Ctx) error {
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

	interviews, total, err := h.interviewService.GetInterviewsByTenant(c.Context(), uuid.UUID{}, page, limit)
	if err != nil {
		log.Errorw("[Handler] GetInterviewsByTenant - Service", err)
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
			Data:  interviews,
			Total: total,
			Page:  page,
			Limit: limit,
		},
	})
}

func (h *interviewHandler) GetInterviewsByCandidate(c *fiber.Ctx) error {
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

	interviews, total, err := h.interviewService.GetInterviewsByCandidate(c.Context(), candidateUUID, page, limit)
	if err != nil {
		log.Errorw("[Handler] GetInterviewsByCandidate - Service", err)
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
			Data:  interviews,
			Total: total,
			Page:  page,
			Limit: limit,
		},
	})
}

func (h *interviewHandler) GetInterviewsByInterviewer(c *fiber.Ctx) error {
	interviewerID := c.Params("interviewerId")
	interviewerUUID, err := uuid.Parse(interviewerID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid interviewer ID",
			},
		})
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	interviews, total, err := h.interviewService.GetInterviewsByInterviewer(c.Context(), interviewerUUID, page, limit)
	if err != nil {
		log.Errorw("[Handler] GetInterviewsByInterviewer - Service", err)
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
			Data:  interviews,
			Total: total,
			Page:  page,
			Limit: limit,
		},
	})
}

func (h *interviewHandler) ScheduleInterview(c *fiber.Ctx) error {
	return h.CreateInterview(c)
}

func (h *interviewHandler) RescheduleInterview(c *fiber.Ctx) error {
	id := c.Params("id")
	interviewID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid interview ID",
			},
		})
	}

	var req request.RescheduleRequest
	if err := c.BodyParser(&req); err != nil {
		log.Errorw("[Handler] RescheduleInterview - BodyParser", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid request body",
			},
		})
	}

	newScheduledAt, _ := time.Parse(time.RFC3339, req.NewScheduledAt)

	if err := h.interviewService.RescheduleInterview(c.Context(), interviewID, newScheduledAt, req.Reason); err != nil {
		log.Errorw("[Handler] RescheduleInterview - Service", err)
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
			Message: "Interview rescheduled successfully",
		},
	})
}

func (h *interviewHandler) CancelInterview(c *fiber.Ctx) error {
	id := c.Params("id")
	interviewID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid interview ID",
			},
		})
	}

	var req request.CancelRequest
	if err := c.BodyParser(&req); err != nil {
		log.Errorw("[Handler] CancelInterview - BodyParser", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid request body",
			},
		})
	}

	if err := h.interviewService.CancelInterview(c.Context(), interviewID, req.Reason); err != nil {
		log.Errorw("[Handler] CancelInterview - Service", err)
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
			Message: "Interview cancelled successfully",
		},
	})
}

func (h *interviewHandler) CompleteInterview(c *fiber.Ctx) error {
	id := c.Params("id")
	interviewID, err := uuid.Parse(id)
	if err != nil {
		return c.Status
func (h *interviewHandler) CompleteInterview(c *fiber.Ctx) error {
	id := c.Params("id")
	interviewID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid interview ID",
			},
		})
	}

	var req request.CompleteInterviewRequest
	if err := c.BodyParser(&req); err != nil {
		log.Errorw("[Handler] CompleteInterview - BodyParser", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid request body",
			},
		})
	}

	if err := h.interviewService.CompleteInterview(c.Context(), interviewID, req.Feedback); err != nil {
		log.Errorw("[Handler] CompleteInterview - Service", err)
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
			Message: "Interview completed successfully",
		},
	})
}

func (h *interviewHandler) SubmitFeedback(c *fiber.Ctx) error {
	id := c.Params("id")
	interviewID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid interview ID",
			},
		})
	}

	var req request.SubmitFeedbackRequest
	if err := c.BodyParser(&req); err != nil {
		log.Errorw("[Handler] SubmitFeedback - BodyParser", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid request body",
			},
		})
	}

	if err := h.interviewService.SubmitFeedback(c.Context(), interviewID, req.Feedback, req.Rating, req.Recommendation); err != nil {
		log.Errorw("[Handler] SubmitFeedback - Service", err)
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
			Message: "Feedback submitted successfully",
		},
	})
}

func (h *interviewHandler) GetInterviewMetrics(c *fiber.Ctx) error {
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

	metrics, err := h.interviewService.GetInterviewMetrics(c.Context(), uuid.UUID{})
	if err != nil {
		log.Errorw("[Handler] GetInterviewMetrics - Service", err)
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
		Data: metrics,
	})
}

		},
	})
}

func (h *interviewHandler) SubmitFeedback(c *fiber.Ctx) error {
	id := c.Params("id")
	interviewID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid interview ID",
			},
		})
	}

	var req request.SubmitFeedbackRequest
	if err := c.BodyParser(&req); err != nil {
		log.Errorw("[Handler] SubmitFeedback - BodyParser", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid request body",
			},
		})
	}

	if err := h.interviewService.SubmitFeedback(c.Context(), interviewID, req.OverallFeedback, req.Rating, req.Recommendation); err != nil {
		log.Errorw("[Handler] SubmitFeedback - Service", err)
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
			Message: "Feedback submitted successfully",
		},
	})
}

func (h *interviewHandler) GetInterviewMetrics(c *fiber.Ctx) error {
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

	metrics, err := h.interviewService.GetInterviewMetrics(c.Context(), uuid.UUID{})
	if err != nil {
		log.Errorw("[Handler] GetInterviewMetrics - Service", err)
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
		Data: metrics,
	})
}
