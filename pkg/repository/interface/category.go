package interfaces

import (
	"Zhooze/pkg/domain"
	"Zhooze/pkg/utils/models"
)

type CategoryRepository interface {
	GetCategory() ([]domain.Category, error)
	CheckIfCategoryAlreadyExists(category string) (bool, error)
	AddCategory(category models.Category) (domain.Category, error)
	DeleteCategory(id int) error
	UpdateCategory(current string, new string) (domain.Category, error)
	CheckCategory(current string) (bool, error)
}
