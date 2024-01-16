package interfaces

import (
	"Zhooze/pkg/domain"
	"Zhooze/pkg/utils/models"
	"mime/multipart"
)

type ProductUseCase interface {
	ShowAllProducts(page int, count int) ([]models.ProductBrief, error)
	ShowAllProductsFromAdmin(page int, count int) ([]models.ProductBrief, error)
	SearchProductsOnPrefix(prefix string) ([]models.ProductBrief, error)
	FilterCategory(data map[string]int) ([]models.ProductBrief, error)
	AddProducts(product models.Product) (domain.Product, error)
	DeleteProducts(id string) error
	UpdateProduct(pid int, stock int) (models.ProductUpdateReciever, error)
	UpdateProductImage(id int, file *multipart.FileHeader) error
	ShowImages(productID int) ([]models.Image, error)
}
