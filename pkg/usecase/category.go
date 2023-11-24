package usecase

import (
	"Zhooze/pkg/domain"
	interfaces "Zhooze/pkg/repository/interface"
	services "Zhooze/pkg/usecase/interface"
	"Zhooze/pkg/utils/models"
	"errors"
)

type categoryUseCase struct {
	categoryRepository interfaces.CategoryRepository
}

func NewCategoryUseCase(repo interfaces.CategoryRepository) services.CategoryUseCase {
	return &categoryUseCase{
		categoryRepository: repo,
	}
}

func (cr *categoryUseCase) GetCategory() ([]domain.Category, error) {
	category, err := cr.categoryRepository.GetCategory()
	if err != nil {
		return []domain.Category{}, err
	}
	return category, nil

}
func (cr *categoryUseCase) AddCategory(category models.Category) (domain.Category, error) {
	exists, err := cr.categoryRepository.CheckIfCategoryAlreadyExists(category.Category)
	if err != nil {
		return domain.Category{}, err
	}

	if exists {
		return domain.Category{}, errors.New("category already exists")
	}
	categories, err := cr.categoryRepository.AddCategory(category)
	if err != nil {
		return domain.Category{}, err
	}
	return categories, nil
}
func (cr *categoryUseCase) UpdateCategory(current string, new string) (domain.Category, error) {
	categries, err := cr.categoryRepository.CheckCategory(current)
	if err != nil {
		return domain.Category{}, err
	}
	if !categries {
		return domain.Category{}, errors.New("category doesn't exist")
	}
	newCate, err := cr.categoryRepository.UpdateCategory(current, new)
	if err != nil {
		return domain.Category{}, err
	}
	return newCate, nil
}
func (cr *categoryUseCase) DeleteCategory(id int) error {
	err := cr.categoryRepository.DeleteCategory(id)
	if err != nil {
		return err
	}
	return nil
}
