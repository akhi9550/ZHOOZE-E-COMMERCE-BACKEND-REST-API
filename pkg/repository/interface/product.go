package interfaces

import (
	"Zhooze/pkg/domain"
	"Zhooze/pkg/utils/models"
)

type ProductRepository interface {
	ShowAllProducts(page int, count int) ([]models.ProductBrief, error)
	ShowAllProductsFromAdmin(page int, count int) ([]models.ProductBrief, error)
	GetImage(productID int) ([]string, error)
	CheckValidateCategory(data map[string]int) error
	GetProductFromCategory(id int) ([]models.ProductBrief, error)
	GetQuantityFromProductID(id int) (int, error)
	GetPriceOfProductFromID(prodcut_id int) (float64, error)
	ProductAlreadyExist(Name string) bool
	FindCategoryID(id int) (int, error)
	StockInvalid(Name string) bool
	AddProducts(product models.Product) (domain.Product, error)
	DeleteProducts(id string) error
	CheckProductExist(pid int) (bool, error)
	UpdateProduct(pid int, stock int) (models.ProductUpdateReciever, error)
	DoesProductExist(productID int) (bool, error)
	UpdateProductImage(productID int, url string) error
	DisplayImages(productID int) (domain.Product, []domain.Image, error)
	ShowImages(productID int) ([]models.Image, error)
}
