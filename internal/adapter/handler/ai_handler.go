package handler

import (
	"errors"
	"workzen-be/internal/adapter/handler/request"
	"workzen-be/internal/core/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AIHandler struct {
	aiService service.AIService
}

func NewAIHandler(aiService service.AIService) *AIHandler {
	return &AIHandler{
		aiService: aiService,
	}
}

func (h *AIHandler) AnalyzeCV(c *fiber.Ctx) error {
	var req request.AnalyzeCVRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	resp, err := h.aiService.AnalyzeCV(c.Context(), req.CvText)
	if err != nil {
		return c.Status(grpcToHTTPStatus(err)).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *AIHandler) MatchJob(c *fiber.Ctx) error {
	var req request.MatchJobRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	resp, err := h.aiService.MatchJob(c.Context(), req.CvText, req.JdText)
	if err != nil {
		return c.Status(grpcToHTTPStatus(err)).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func grpcToHTTPStatus(err error) int {
	if errors.Is(err, service.ErrAIServiceUnavailable) {
		return fiber.StatusServiceUnavailable
	}

	st, ok := status.FromError(err)
	if !ok {
		return fiber.StatusInternalServerError
	}

	switch st.Code() {
	case codes.InvalidArgument:
		return fiber.StatusBadRequest
	case codes.NotFound:
		return fiber.StatusNotFound
	case codes.DeadlineExceeded:
		return fiber.StatusGatewayTimeout
	case codes.Unavailable:
		return fiber.StatusServiceUnavailable
	case codes.ResourceExhausted:
		return fiber.StatusTooManyRequests
	case codes.Unauthenticated:
		return fiber.StatusUnauthorized
	case codes.PermissionDenied:
		return fiber.StatusForbidden
	default:
		return fiber.StatusInternalServerError
	}
}
