package usecase

import (
	"Zhooze/pkg/domain"
	"Zhooze/pkg/helper"
	interfaces "Zhooze/pkg/repository/interface"
	services "Zhooze/pkg/usecase/interface"
	"strings"

	"Zhooze/pkg/utils/models"
	"errors"
	"mime/multipart"
)

type productUseCase struct {
	productRepository interfaces.ProductRepository
	offerRepository   interfaces.OfferRepository
}

func NewProductUseCase(repository interfaces.ProductRepository, offerRepo interfaces.OfferRepository) services.ProductUseCase {
	return &productUseCase{
		productRepository: repository,
		offerRepository:   offerRepo,
	}
}
func (pr *productUseCase) ShowAllProducts(page int, count int) ([]models.ProductBrief, error) {
	productDetails, err := pr.productRepository.ShowAllProducts(page, count)
	if err != nil {
		return []models.ProductBrief{}, err
	}

	for i := range productDetails {
		p := &productDetails[i]
		if p.Stock <= 0 {
			p.ProductStatus = "out of stock"
		} else {
			p.ProductStatus = "in stock"
		}
	}
	//loop inside products and then calculate discounted price of each then return
	for j := range productDetails {
		discount_percentage, err := pr.offerRepository.FindDiscountPercentageForProduct(int(productDetails[j].ID))
		if err != nil {
			return []models.ProductBrief{}, errors.New("there was some error in finding the discounted prices")
		}
		var discount float64
		if discount_percentage > 0 {
			discount = (productDetails[j].Price * float64(discount_percentage)) / 100
		}
		productDetails[j].DiscountedPrice = productDetails[j].Price - discount

		discount_percentageCategory, err := pr.offerRepository.FindDiscountPercentageForCategory(int(productDetails[j].CategoryID))
		if err != nil {
			return []models.ProductBrief{}, errors.New("there was some error in finding the discounted prices")
		}
		var categorydiscount float64
		if discount_percentageCategory > 0 {
			categorydiscount = (productDetails[j].Price * float64(discount_percentageCategory)) / 100
		}

		productDetails[j].DiscountedPrice = productDetails[j].DiscountedPrice - categorydiscount
	}
	var updatedproductDetails []models.ProductBrief
	for _, p := range productDetails {
		img, err := pr.productRepository.GetImage(int(p.ID))
		if err != nil {
			return nil, err
		}
		p.Image = img
		updatedproductDetails = append(updatedproductDetails, p)
	}

	return updatedproductDetails, nil
}

func (pr *productUseCase) ShowAllProductsFromAdmin(page int, count int) ([]models.ProductBrief, error) {
	productDetails, err := pr.productRepository.ShowAllProductsFromAdmin(page, count)
	if err != nil {
		return []models.ProductBrief{}, err
	}
	for i := range productDetails {
		p := &productDetails[i]
		if p.Stock <= 0 {
			p.ProductStatus = "out of stock"
		} else {
			p.ProductStatus = "in stock"
		}
	}
	for j := range productDetails {
		discount_percentage, err := pr.offerRepository.FindDiscountPercentageForProduct(int(productDetails[j].ID))
		if err != nil {
			return []models.ProductBrief{}, errors.New("there was some error in finding the discounted prices")
		}
		var discount float64
		if discount_percentage > 0 {
			discount = (productDetails[j].Price * float64(discount_percentage)) / 100
		}
		productDetails[j].DiscountedPrice = productDetails[j].Price - discount

		discount_percentageCategory, err := pr.offerRepository.FindDiscountPercentageForCategory(int(productDetails[j].CategoryID))
		if err != nil {
			return []models.ProductBrief{}, errors.New("there was some error in finding the discounted prices")
		}
		var categorydiscount float64
		if discount_percentageCategory > 0 {
			categorydiscount = (productDetails[j].Price * float64(discount_percentageCategory)) / 100
		}

		productDetails[j].DiscountedPrice = productDetails[j].DiscountedPrice - categorydiscount
	}
	var updatedproductDetails []models.ProductBrief
	for _, p := range productDetails {
		img, err := pr.productRepository.GetImage(int(p.ID))
		if err != nil {
			return nil, err
		}
		p.Image = img
		updatedproductDetails = append(updatedproductDetails, p)
	}

	return updatedproductDetails, nil

}
func (pr *productUseCase) FilterCategory(data map[string]int) ([]models.ProductBrief, error) {
	err := pr.productRepository.CheckValidateCategory(data)
	if err != nil {
		return []models.ProductBrief{}, err
	}
	var ProductFromCategory []models.ProductBrief
	for _, id := range data {
		product, err := pr.productRepository.GetProductFromCategory(id)
		if err != nil {
			return []models.ProductBrief{}, err
		}
		for _, products := range product {
			stock, err := pr.productRepository.GetQuantityFromProductID(int(products.ID))
			if err != nil {
				return []models.ProductBrief{}, err
			}
			if stock <= 0 {
				products.ProductStatus = "out of stock"
			} else {
				products.ProductStatus = "in stock"
			}
			if products.ID != 0 {
				ProductFromCategory = append(ProductFromCategory, products)
			}
		}
	}
	for j := range ProductFromCategory {
		discount_percentage, err := pr.offerRepository.FindDiscountPercentageForProduct(int(ProductFromCategory[j].ID))
		if err != nil {
			return []models.ProductBrief{}, errors.New("there was some error in finding the discounted prices")
		}
		var discount float64
		if discount_percentage > 0 {
			discount = (ProductFromCategory[j].Price * float64(discount_percentage)) / 100
		}
		ProductFromCategory[j].DiscountedPrice = ProductFromCategory[j].Price - discount

		discount_percentageCategory, err := pr.offerRepository.FindDiscountPercentageForCategory(int(ProductFromCategory[j].CategoryID))
		if err != nil {
			return []models.ProductBrief{}, errors.New("there was some error in finding the discounted prices")
		}
		var categorydiscount float64
		if discount_percentageCategory > 0 {
			categorydiscount = (ProductFromCategory[j].Price * float64(discount_percentageCategory)) / 100
		}

		ProductFromCategory[j].DiscountedPrice = ProductFromCategory[j].DiscountedPrice - categorydiscount
	}
	var updatedproductDetails []models.ProductBrief
	for _, p := range ProductFromCategory {
		img, err := pr.productRepository.GetImage(int(p.ID))
		if err != nil {
			return nil, err
		}
		p.Image = img
		updatedproductDetails = append(updatedproductDetails, p)
	}

	return updatedproductDetails, nil

}

func (pr *productUseCase) AddProducts(product models.Product) (domain.Product, error) {
	exist := pr.productRepository.ProductAlreadyExist(product.Name)
	if exist {
		return domain.Product{}, errors.New("product already exist")
	}
	productResponse, err := pr.productRepository.AddProducts(product)
	if err != nil {
		return domain.Product{}, err
	}
	stock := pr.productRepository.StockInvalid(productResponse.Name)
	if !stock {
		return domain.Product{}, errors.New("stock is invalid input")
	}
	return productResponse, nil
}
func (pr *productUseCase) DeleteProducts(id string) error {
	err := pr.productRepository.DeleteProducts(id)
	if err != nil {
		return err
	}
	return nil
}
func (pr *productUseCase) UpdateProduct(pid int, stock int) (models.ProductUpdateReciever, error) {
	if stock <= 0 {
		return models.ProductUpdateReciever{}, errors.New("stock doesnot update invalid input")
	}
	result, err := pr.productRepository.CheckProductExist(pid)
	if err != nil {
		return models.ProductUpdateReciever{}, err
	}
	if !result {
		return models.ProductUpdateReciever{}, errors.New("there is no product as you mentioned")
	}
	newcat, err := pr.productRepository.UpdateProduct(pid, stock)
	if err != nil {
		return models.ProductUpdateReciever{}, err
	}
	return newcat, err

}
func (pr *productUseCase) UpdateProductImage(id int, file *multipart.FileHeader) error {

	url, err := helper.AddImageToS3(file)
	if err != nil {
		return err
	}
	err = pr.productRepository.UpdateProductImage(id, url)
	if err != nil {
		return err
	}
	return nil
}
func (pr *productUseCase) ShowImages(productID int) ([]models.Image, error) {
	image, err := pr.productRepository.ShowImages(productID)
	if err != nil {
		return nil, err
	}
	return image, nil
}
func (pr *productUseCase) SearchProductsOnPrefix(prefix string) ([]models.ProductBrief, error) {

	inventoryList, err := pr.productRepository.GetInventory(prefix)

	if err != nil {
		return nil, err
	}

	var filteredProducts []models.ProductBrief

	for _, product := range inventoryList {
		if strings.HasPrefix(strings.ToLower(product.Name), strings.ToLower(prefix)) {
			filteredProducts = append(filteredProducts, product)
		}
	}

	if len(filteredProducts) == 0 {
		return nil, errors.New("no items matching your keyword")
	}

	return filteredProducts, nil
}
