package handler

import (
	"errors"
	"io"
	"strings"
	"unicode/utf8"
	"workzen-be/internal/adapter/handler/response"
	"workzen-be/internal/core/domain/entity"
	"workzen-be/internal/core/service"
	"workzen-be/lib/conv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"

	"github.com/lu4p/cat"
)

type JobPostingHandler struct {
	manpowerReqService service.ManpowerReqService
	jobPostingService  service.JobPostingService
}

func NewJobPostingHandler(manpowerReqService service.ManpowerReqService, jobPostingService service.JobPostingService) *JobPostingHandler {
	return &JobPostingHandler{
		manpowerReqService: manpowerReqService,
		jobPostingService:  jobPostingService,
	}
}

func (h *JobPostingHandler) GenerateLink(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	tenantID := claims.TenantID

	if tenantID == 0 {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	idParam := c.Params("manpowerRequestID")
	manpowerRequestID, err := conv.StringToInt64(idParam)
	if err != nil {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid manpower request ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	token, err := h.manpowerReqService.GeneratePublicLink(c.Context(), int64(tenantID), manpowerRequestID)
	if err != nil {
		code := "[HANDLER] GenerateLink - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"
	defaultSuccessResponse.Data = fiber.Map{
		"public_token": token,
	}

	return c.JSON(defaultSuccessResponse)
}

func (h *JobPostingHandler) GetJobPosting(c *fiber.Ctx) error {
	token := c.Params("token")
	if token == "" {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid token"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	manpowerReq, err := h.manpowerReqService.GetManpowerReqByPublicToken(c.Context(), token)
	if err != nil {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Job posting not found"
		return c.Status(fiber.StatusNotFound).JSON(errorResp)
	}

	if manpowerReq.Status != "OPEN" {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "This job posting is no longer accepting applications"
		return c.Status(fiber.StatusGone).JSON(errorResp)
	}

	jobPosting := response.JobPostingResponse{
		Position:       manpowerReq.Position,
		CompanyName:    manpowerReq.Client.CompanyName,
		WorkLocation:   manpowerReq.WorkLocation,
		SalaryMin:      manpowerReq.SalaryMin,
		SalaryMax:      manpowerReq.SalaryMax,
		JobDescription: manpowerReq.JobDescription,
		DeadlineDate:   manpowerReq.DeadlineDate.In(jakartaTZ).Format("02 January 2006"),
		Status:         manpowerReq.Status,
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"
	defaultSuccessResponse.Data = jobPosting

	return c.JSON(defaultSuccessResponse)
}

func (h *JobPostingHandler) ApplyToJob(c *fiber.Ctx) error {
	token := c.Params("token")
	if token == "" {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid token"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	manpowerReq, err := h.manpowerReqService.GetManpowerReqByPublicToken(c.Context(), token)
	if err != nil {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Job posting not found"
		return c.Status(fiber.StatusNotFound).JSON(errorResp)
	}

	if manpowerReq.Status != "OPEN" {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "This job posting is no longer accepting applications"
		return c.Status(fiber.StatusGone).JSON(errorResp)
	}

	fullName := c.FormValue("full_name")
	email := c.FormValue("email")
	phone := c.FormValue("phone")

	if fullName == "" || email == "" {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "full_name and email are required"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	file, err := c.FormFile("cv")
	if err != nil {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "CV file is required"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	ext := strings.ToLower(file.Filename[strings.LastIndex(file.Filename, "."):])
	if ext != ".pdf" && ext != ".docx" && ext != ".doc" && ext != ".txt" {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Only PDF, DOCX, DOC, or TXT files are accepted"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	f, err := file.Open()
	if err != nil {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Failed to read CV file"
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}
	defer f.Close()

	fileBytes, err := io.ReadAll(f)
	if err != nil {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Failed to read CV file"
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	var cvText string
	if ext == ".txt" {
		cvText = string(fileBytes)
	} else {
		cvText, err = cat.FromBytes(fileBytes)
		if err != nil {
			code := "[HANDLER] ApplyToJob - ExtractText"
			log.Errorw(code, err)
			errorResp.Meta.Status = false
			errorResp.Meta.Message = "Failed to extract text from CV"
			return c.Status(fiber.StatusBadRequest).JSON(errorResp)
		}
	}

	if !utf8.ValidString(cvText) {
		cvText = strings.ToValidUTF8(cvText, "")
	}
	cvText = strings.TrimSpace(cvText)

	// Limit CV text to ~8000 chars to avoid AI timeout on large documents
	if len(cvText) > 8000 {
		cvText = cvText[:8000]
	}
	if cvText == "" {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Could not extract any text from the CV file"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	result, err := h.jobPostingService.ApplyToJob(c.Context(), manpowerReq, fullName, email, phone, cvText)
	if err != nil {
		if errors.Is(err, service.ErrScoreTooLow) {
			return c.Status(fiber.StatusOK).JSON(response.DefaultSuccessResponse{
				Meta: response.Meta{
					Status:  true,
					Message: "Your CV has been analyzed but does not meet the minimum match score (70) for this position",
				},
				Data: response.JobApplyResponse{
					Message:       "CV does not meet minimum match requirements",
					Score:         result.Score,
					Verdict:       result.Verdict,
					MatchedSkills: result.MatchedSkills,
					MissingSkills: result.MissingSkills,
					Explanation:   result.Explanation,
				},
			})
		}

		if errors.Is(err, service.ErrAIServiceUnavailable) {
			errorResp.Meta.Status = false
			errorResp.Meta.Message = "AI service is currently unavailable, please try again later"
			return c.Status(fiber.StatusServiceUnavailable).JSON(errorResp)
		}

		code := "[HANDLER] ApplyToJob - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Failed to process application"
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	return c.Status(fiber.StatusCreated).JSON(response.DefaultSuccessResponse{
		Meta: response.Meta{
			Status:  true,
			Message: "Application submitted successfully",
		},
		Data: response.JobApplyResponse{
			Message:       "Your application has been submitted successfully",
			Score:         result.Score,
			Verdict:       result.Verdict,
			MatchedSkills: result.MatchedSkills,
			MissingSkills: result.MissingSkills,
			Explanation:   result.Explanation,
		},
	})
}
