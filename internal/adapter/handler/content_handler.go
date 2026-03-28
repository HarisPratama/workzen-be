package handler

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
	"workzen-be/internal/adapter/handler/request"
	"workzen-be/internal/adapter/handler/response"
	"workzen-be/internal/core/domain/entity"
	"workzen-be/internal/core/service"
	"workzen-be/lib/conv"
	validatorLib "workzen-be/lib/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type ContentHandler interface {
	GetContents(c *fiber.Ctx) error
	GetContentByID(c *fiber.Ctx) error
	CreateContent(c *fiber.Ctx) error
	UpdateContent(c *fiber.Ctx) error
	DeleteContent(c *fiber.Ctx) error
	UploadImageR2(c *fiber.Ctx) error

	//	FE
	GetContentWithQuery(c *fiber.Ctx) error
	GetContentDetail(c *fiber.Ctx) error
}

type contentHandler struct {
	contentService service.ContentService
}

func (c2 *contentHandler) GetContentByID(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code := "[HANDLER] GetContentByID - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	idParam := c.Params("contentID")
	contentID, err := conv.StringToInt64(idParam)
	if err != nil {
		code := "[HANDLER] GetContentByID - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	result, err := c2.contentService.GetContentByID(context.Background(), contentID)
	if err != nil {
		code := "[HANDLER] GetContentByID - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"

	respContent := response.ContentResponse{
		ID:           result.ID,
		Title:        result.Title,
		Excerpt:      result.Excerpt,
		Description:  result.Description,
		Image:        result.Image,
		Tags:         result.Tags,
		Status:       result.Status,
		CategoryID:   result.CategoryID,
		CreatedByID:  result.CreatedByID,
		CreatedAt:    result.CreatedAt.In(jakartaTZ).Format("02 January 2006"),
		CategoryName: result.Category.Title,
		Author:       result.User.Name,
	}

	defaultSuccessResponse.Data = respContent
	return c.JSON(defaultSuccessResponse)
}

func (c2 *contentHandler) CreateContent(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code := "[HANDLER] DeleteContent - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	userID := claims.UserID
	var req request.ContentRequest
	if err := c.BodyParser(&req); err != nil {
		code := "[HANDLER] CreateContent - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "invalid request body"

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err = validatorLib.ValidateStruct(&req); err != nil {
		code := "[HANDLER] CreateContent - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	tags := strings.Split(req.Tags, ",")
	reqEntity := entity.ContentEntity{
		Title:       req.Title,
		Excerpt:     req.Excerpt,
		Description: req.Description,
		Image:       req.Image,
		Tags:        tags,
		Status:      req.Status,
		CategoryID:  req.CategoryID,
		CreatedByID: int64(userID),
		CreatedAt:   time.Time{},
	}

	err = c2.contentService.CreateContent(c.Context(), reqEntity)
	if err != nil {
		code := "[HANDLER] CreateContent - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"
	defaultSuccessResponse.Data = nil

	return c.Status(fiber.StatusCreated).JSON(defaultSuccessResponse)
}

func (c2 *contentHandler) UpdateContent(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code := "[HANDLER] DeleteContent - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	userID := claims.UserID
	var req request.ContentRequest
	if err := c.BodyParser(&req); err != nil {
		code := "[HANDLER] UpdateContent - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "invalid request body"

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err = validatorLib.ValidateStruct(&req); err != nil {
		code := "[HANDLER] UpdateContent - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	idParam := c.Params("contentID")
	contentID, err := conv.StringToInt64(idParam)
	if err != nil {
		code := "[HANDLER] UpdateContent - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	tags := strings.Split(req.Tags, ",")
	reqEntity := entity.ContentEntity{
		ID:          contentID,
		Title:       req.Title,
		Excerpt:     req.Excerpt,
		Description: req.Description,
		Image:       req.Image,
		Tags:        tags,
		Status:      req.Status,
		CategoryID:  req.CategoryID,
		CreatedByID: int64(userID),
	}

	err = c2.contentService.UpdateContent(c.Context(), reqEntity)
	if err != nil {
		code := "[HANDLER] UpdateContent - 5"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"
	defaultSuccessResponse.Data = nil

	return c.JSON(defaultSuccessResponse)
}

func (c2 *contentHandler) DeleteContent(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code := "[HANDLER] DeleteContent - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	idParam := c.Params("contentID")
	contentID, err := conv.StringToInt64(idParam)
	if err != nil {
		code := "[HANDLER] DeleteContent - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	err = c2.contentService.DeleteContent(c.Context(), contentID)
	if err != nil {
		code := "[HANDLER] DeleteContent - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"
	defaultSuccessResponse.Data = nil

	return c.JSON(defaultSuccessResponse)
}

func (c2 *contentHandler) UploadImageR2(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code := "[HANDLER] UploadImage - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	var req request.FileUploadRequest
	file, err := c.FormFile("image")
	if err != nil {
		code := "[HANDLER] UploadImage - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "invalid request body"

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err = c.SaveFile(file, fmt.Sprintf("./temp/content/%s", file.Filename)); err != nil {
		code := "[HANDLER] UploadImage - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	req.Image = fmt.Sprintf("./temp/content/%s", file.Filename)
	reqEntity := entity.FileUploadEntity{
		Name: fmt.Sprintf("%d-%d", int64(claims.UserID), time.Now().UnixNano()),
		Path: req.Image,
	}

	imageUrl, err := c2.contentService.UploadImageR2(c.Context(), reqEntity)
	if err != nil {
		code := "[HANDLER] UploadImage - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	if req.Image != "" {
		err = os.Remove(req.Image)
		if err != nil {
			code := "[HANDLER] UploadImage - 5"
			log.Errorw(code, err)
			errorResp.Meta.Status = false
			errorResp.Meta.Message = err.Error()

			return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
		}
	}

	urlImageResp := map[string]interface{}{
		"urlImage": imageUrl,
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"
	defaultSuccessResponse.Data = urlImageResp

	return c.Status(fiber.StatusCreated).JSON(defaultSuccessResponse)
}

func (c2 *contentHandler) GetContentWithQuery(c *fiber.Ctx) error {
	page := 1
	if c.Query("page") != "" {
		page, err = conv.StringToInt(c.Query("page"))
		if err != nil {
			code := "[HANDLER] GetContentWithQuery - 1"
			log.Errorw(code, err)
			errorResp.Meta.Status = false
			errorResp.Meta.Message = "invalid page number"

			return c.Status(fiber.StatusBadRequest).JSON(errorResp)
		}
	}

	limit := 6
	if c.Query("limit") != "" {
		limit, err = conv.StringToInt(c.Query("limit"))
		if err != nil {
			code := "[HANDLER] GetContentWithQuery - 2"
			log.Errorw(code, err)
			errorResp.Meta.Status = false
			errorResp.Meta.Message = "invalid limit number"

			return c.Status(fiber.StatusBadRequest).JSON(errorResp)
		}
	}

	orderBy := "created_at"
	if c.Query("orderBy") != "" {
		orderBy = c.Query("orderBy")
	}

	orderType := "desc"
	if c.Query("orderType") != "" {
		orderType = c.Query("orderType")
	}

	search := ""
	if c.Query("search") != "" {
		search = c.Query("search")
	}

	categoryID := 0
	if c.Query("categoryID") != "" {
		categoryID, err = conv.StringToInt(c.Query("categoryID"))
		if err != nil {
			code := "[HANDLER] GetContentWithQuery - 3"
			log.Errorw(code, err)
			errorResp.Meta.Status = false
			errorResp.Meta.Message = "Invalid category ID"

			return c.Status(fiber.StatusBadRequest).JSON(errorResp)
		}
	}

	reqEntity := entity.QueryString{
		Limit:      limit,
		Page:       page,
		OrderBy:    orderBy,
		Search:     search,
		OrderType:  orderType,
		Status:     "PUBLISH",
		CategoryID: int64(categoryID),
	}

	results, totalData, totalPages, err := c2.contentService.GetContents(c.Context(), reqEntity)
	if err != nil {
		code := "[HANDLER] GetContentWithQuery - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"

	respContents := []response.ContentResponse{}
	for _, content := range results {
		respContent := response.ContentResponse{
			ID:           content.ID,
			Title:        content.Title,
			Excerpt:      content.Excerpt,
			Description:  content.Description,
			Image:        content.Image,
			Tags:         content.Tags,
			Status:       content.Status,
			CategoryID:   content.CategoryID,
			CreatedByID:  content.CreatedByID,
			CreatedAt:    content.CreatedAt.In(jakartaTZ).Format("02 January 2006"),
			CategoryName: content.Category.Title,
			Author:       content.User.Name,
		}

		respContents = append(respContents, respContent)
	}

	defaultSuccessResponse.Data = respContents
	defaultSuccessResponse.Pagination = &response.PaginationResponse{
		TotalRecords: int(totalData),
		Page:         page,
		PerPage:      limit,
		TotalPages:   int(totalPages),
	}

	return c.JSON(defaultSuccessResponse)
}

func (c2 *contentHandler) GetContentDetail(c *fiber.Ctx) error {
	idParam := c.Params("contentID")
	contentID, err := conv.StringToInt64(idParam)
	if err != nil {
		code := "[HANDLER] GetContentByID - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	result, err := c2.contentService.GetContentByID(context.Background(), contentID)
	if err != nil {
		code := "[HANDLER] GetContentByID - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"

	respContent := response.ContentResponse{
		ID:           result.ID,
		Title:        result.Title,
		Excerpt:      result.Excerpt,
		Description:  result.Description,
		Image:        result.Image,
		Tags:         result.Tags,
		Status:       result.Status,
		CategoryID:   result.CategoryID,
		CreatedByID:  result.CreatedByID,
		CreatedAt:    result.CreatedAt.In(jakartaTZ).Format("02 January 2006"),
		CategoryName: result.Category.Title,
		Author:       result.User.Name,
	}

	defaultSuccessResponse.Data = respContent
	defaultSuccessResponse.Pagination = nil
	return c.JSON(defaultSuccessResponse)
}

func (c2 *contentHandler) GetContents(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code = "[Handler] GetContents - 1"
		log.Errorw(code, err)
		errorResp.Status = false
		errorResp.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	page := 1
	if c.Query("page") != "" {
		page, err = conv.StringToInt(c.Query("page"))
		if err != nil {
			code = "[Handler] GetContents - 2"
			log.Errorw(code, err)
			errorResp.Status = false
			errorResp.Message = "Invalid page number"

			return c.Status(fiber.StatusBadRequest).JSON(errorResp)
		}
	}

	limit := 10
	if c.Query("limit") != "" {
		limit, err = conv.StringToInt(c.Query("limit"))
		if err != nil {
			code = "[Handler] GetContents - 3"
			log.Errorw(code, err)
			errorResp.Status = false
			errorResp.Message = "Invalid limit number"

			return c.Status(fiber.StatusBadRequest).JSON(errorResp)
		}
	}

	orderBy := "created_at"
	if c.Query("orderBy") != "" {
		orderBy = c.Query("orderBy")
	}

	orderType := "desc"
	if c.Query("orderType") != "" {
		orderType = c.Query("orderType")
	}

	search := ""
	if c.Query("search") != "" {
		search = c.Query("search")
	}

	categoryID := 0
	if c.Query("categoryID") != "" {
		categoryID, err = conv.StringToInt(c.Query("categoryID"))
		if err != nil {
			code = "[Handler] GetContents - 4"
			log.Errorw(code, err)
			errorResp.Status = false
			errorResp.Message = "Invalid categoryID"

			return c.Status(fiber.StatusBadRequest).JSON(errorResp)
		}
	}

	reqEntity := entity.QueryString{
		Limit:      limit,
		Page:       page,
		OrderBy:    orderBy,
		OrderType:  orderType,
		Search:     search,
		CategoryID: int64(categoryID),
	}

	results, totalData, totalPages, err := c2.contentService.GetContents(c.Context(), reqEntity)
	if err != nil {
		code = "[Handler] GetContents - 5"
		log.Errorw(code, err)
		errorResp.Status = false
		errorResp.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "success"

	respContents := []response.ContentResponse{}
	for _, content := range results {
		respContent := response.ContentResponse{
			ID:           content.ID,
			Title:        content.Title,
			Excerpt:      content.Excerpt,
			Description:  content.Description,
			Image:        content.Image,
			Tags:         content.Tags,
			Status:       content.Status,
			CategoryID:   content.CategoryID,
			CreatedByID:  content.CreatedByID,
			CreatedAt:    content.CreatedAt.In(jakartaTZ).Format("02 January 2006"),
			CategoryName: content.Category.Title,
			Author:       content.User.Name,
		}

		respContents = append(respContents, respContent)
	}

	defaultSuccessResponse.Data = respContents
	defaultSuccessResponse.Pagination = &response.PaginationResponse{
		TotalRecords: int(totalData),
		Page:         page,
		PerPage:      limit,
		TotalPages:   int(totalPages),
	}

	return c.JSON(defaultSuccessResponse)
}

func NewContentHandler(contentService service.ContentService) ContentHandler {
	return &contentHandler{contentService: contentService}
}
