package interfaces

import (
	"Zhooze/pkg/domain"
	"Zhooze/pkg/utils/models"
)

type CategoryUseCase interface {
	GetCategory() ([]domain.Category, error)
	AddCategory(category models.Category) (domain.Category, error)
	UpdateCategory(current string, new string) (domain.Category, error)
	DeleteCategory(id int) error
}
