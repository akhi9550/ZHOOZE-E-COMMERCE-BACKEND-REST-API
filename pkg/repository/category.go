package repository

import (
	"Zhooze/pkg/domain"
	interfaces "Zhooze/pkg/repository/interface"
	"Zhooze/pkg/utils/models"
	"errors"

	"gorm.io/gorm"
)

type categoryRepository struct {
	DB *gorm.DB
}

func NewCategoryRepository(DB *gorm.DB) interfaces.CartRepository {
	return &cartRepository{
		DB: DB,
	}
}
func (car *categoryRepository) GetCategory() ([]domain.Category, error) {
	var category []domain.Category
	err := car.DB.Raw("SELECT * FROM categories").Scan(&category).Error
	if err != nil {
		return nil, err
	}
	return category, nil
}
func (car *categoryRepository) CheckIfCategoryAlreadyExists(category string) (bool, error) {
	var count int64
	err := car.DB.Raw("SELECT COUNT(*) FROM categories WHERE category = $1", category).Scan(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (car *categoryRepository) AddCategory(category models.Category) (domain.Category, error) {
	var categore string
	err := car.DB.Raw("INSERT INTO categories (category) VALUES (?) RETURNING category", category.Category).Scan(&categore).Error
	if err != nil {
		return domain.Category{}, err
	}
	var categoriesResponse domain.Category
	err = car.DB.Raw("SELECT id , category FROM categories  WHERE category = ?", categore).Scan(&categoriesResponse).Error
	if err != nil {
		return domain.Category{}, err
	}
	return categoriesResponse, nil
}
func (car *categoryRepository) DeleteCategory(id int) error {
	var count int
	if err := car.DB.Raw("SELECT COUNT(*) FROM categories WHERE id=?", id).Scan(&count).Error; err != nil {
		return err
	}
	if count < 1 {
		return errors.New("category for given id does not exist")
	}

	if err := car.DB.Exec("DELETE FROM categories WHERE id=?", id).Error; err != nil {
		return err
	}
	return nil
}
func (car *categoryRepository) UpdateCategory(current string, new string) (domain.Category, error) {
	if car.DB == nil {
		return domain.Category{}, errors.New("database connection is nil")
	}
	if err := car.DB.Exec("UPDATE categories SET category=? WHERE category = ?", new, current).Error; err != nil {
		return domain.Category{}, err
	}
	var newcat domain.Category
	if err := car.DB.Raw("SELECT id,category FROM categories WHERE category = ?", new).Scan(&newcat).Error; err != nil {
		return domain.Category{}, nil
	}
	return newcat, nil
}
func (car *categoryRepository) CheckCategory(current string) (bool, error) {
	var count int
	err := car.DB.Raw("SELECT COUNT(*) FROM categories WHERE category=?", current).Scan(&count).Error
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, err
	}
	return true, err
}
