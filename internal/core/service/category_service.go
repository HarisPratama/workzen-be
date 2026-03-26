package service

import (
	"context"
	"workzen-be/internal/adapter/repository"
	"workzen-be/internal/core/domain/entity"
	"workzen-be/lib/conv"

	"github.com/gofiber/fiber/v2/log"
)

type CategoryService interface {
	GetCategories(ctx context.Context) ([]entity.CategoryEntity, error)
	GetCategoryByID(ctx context.Context, id int64) (*entity.CategoryEntity, error)
	CreateCategory(ctx context.Context, req entity.CategoryEntity) error
	EditCategoryByID(ctx context.Context, req entity.CategoryEntity) error
	DeleteCategory(ctx context.Context, id int64) error
}

type categoryService struct {
	categoryRepository repository.CategoryRepository
}

func (c *categoryService) GetCategories(ctx context.Context) ([]entity.CategoryEntity, error) {
	results, err := c.categoryRepository.GetCategories(ctx)
	if err != nil {
		code = "[SERVICE] GetCategories - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return results, nil
}

func (c *categoryService) GetCategoryByID(ctx context.Context, id int64) (*entity.CategoryEntity, error) {
	result, err := c.categoryRepository.GetCategoryByID(ctx, id)
	if err != nil {
		code = "[SERVICE] GetCategoryByID - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return result, nil
}

func (c *categoryService) CreateCategory(ctx context.Context, req entity.CategoryEntity) error {
	slug := conv.GenerateSlug(req.Title)
	req.Slug = slug

	err := c.categoryRepository.CreateCategory(ctx, req)

	if err != nil {
		code = "[SERVICE] CreateCategory - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (c *categoryService) EditCategoryByID(ctx context.Context, req entity.CategoryEntity) error {
	category, err := c.categoryRepository.GetCategoryByID(ctx, req.ID)
	if err != nil {
		code = "[SERVICE] EditCategory - 1"
		log.Errorw(code, err)
		return err
	}

	slug := conv.GenerateSlug(req.Title)
	if category.Title == req.Title {
		slug = category.Slug
	}

	req.Slug = slug

	err = c.categoryRepository.EditCategoryByID(ctx, req)
	if err != nil {
		code = "[SERVICE] EditCategory - 2"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (c *categoryService) DeleteCategory(ctx context.Context, id int64) error {
	err = c.categoryRepository.DeleteCategory(ctx, id)
	if err != nil {
		code = "[SERVICE] DeleteCategory - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func NewCategoryService(categoryRepository repository.CategoryRepository) CategoryService {
	return &categoryService{categoryRepository: categoryRepository}
}
