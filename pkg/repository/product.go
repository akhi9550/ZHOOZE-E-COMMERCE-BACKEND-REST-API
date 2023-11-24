package repository

import (
	"Zhooze/pkg/domain"
	interfaces "Zhooze/pkg/repository/interface"
	"Zhooze/pkg/utils/models"
	"errors"
	"log"
	"strconv"

	"gorm.io/gorm"
)

type productRepository struct {
	DB *gorm.DB
}

func NewProductRepository(DB *gorm.DB) interfaces.ProductRepository {
	return &productRepository{
		DB: DB,
	}
}

func (pr *productRepository) ShowAllProducts(page int, count int) ([]models.ProductBrief, error) {
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * count
	var productBrief []models.ProductBrief
	err := pr.DB.Raw(`SELECT * FROM products limit ? offset ?`, count, offset).Scan(&productBrief).Error
	if err != nil {
		return nil, err
	}
	return productBrief, nil
}
func (pr *productRepository) ShowAllProductsFromAdmin(page int, count int) ([]models.ProductBrief, error) {
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * count
	var productBrief []models.ProductBrief
	err := pr.DB.Raw(`SELECT * FROM products limit ? offset ?`, count, offset).Scan(&productBrief).Error
	if err != nil {
		return nil, err
	}
	return productBrief, nil
}
func (pr *productRepository) CheckValidateCategory(data map[string]int) error {
	for _, id := range data {
		var count int
		err := pr.DB.Raw("SELECT COUNT(*) FROM categories WHERE id=?", id).Scan(&count).Error
		if err != nil {
			return err
		}
		if count < 1 {
			return errors.New("doesn't exist")
		}
	}
	return nil
}
func (pr *productRepository) GetProductFromCategory(id int) ([]models.ProductBrief, error) {
	var product []models.ProductBrief
	err := pr.DB.Raw(`SELECT * FROM products JOIN categories ON products.category_id=categories.id WHERE categories.id=?`, id).Scan(&product).Error
	if err != nil {
		return []models.ProductBrief{}, err
	}
	return product, nil
}
func (pr *productRepository) GetQuantityFromProductID(id int) (int, error) {
	var quantity int
	err := pr.DB.Raw("SELECT stock FROM products WHERE id= ?", id).Scan(&quantity).Error
	if err != nil {
		return 0.0, err
	}

	return quantity, nil
}
func (pr *productRepository) GetPriceOfProductFromID(prodcut_id int) (float64, error) {
	var productPrice float64

	if err := pr.DB.Raw("SELECT price FROM products WHERE id = ?", prodcut_id).Scan(&productPrice).Error; err != nil {
		return 0.0, err
	}
	return productPrice, nil
}

func (pr *productRepository) ProductAlreadyExist(Name string) bool {
	var count int
	if err := pr.DB.Raw("SELECT count(*) FROM products WHERE name = ?", Name).Scan(&count).Error; err != nil {
		return false
	}
	return count > 0
}
func (pr *productRepository) FindCategoryID(id int) (int, error) {
	var a int
	if err := pr.DB.Raw("SELECT category_id FROM products WHERE id = ?", id).Scan(&a).Error; err != nil {
		return 0.0, err
	}
	return a, nil
}
func (pr *productRepository) StockInvalid(Name string) bool {
	var count int
	if err := pr.DB.Raw("SELECT SUM(stock) FROM products WHERE name = ? AND stock >= 0", Name).Scan(&count).Error; err != nil {
		return false
	}
	return count > 0
}
func (pr *productRepository) AddProducts(product models.Product) (domain.Product, error) {
	var p domain.Product
	query := `
    INSERT INTO products (name, description, category_id, size, stock, price)
    VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING name, description, category_id, size, stock, price`
	err := pr.DB.Raw(query, product.Name, product.Description, product.CategoryID, product.Size, product.Stock, product.Price).Scan(&p).Error
	if err != nil {
		log.Println(err.Error())
		return domain.Product{}, err
	}
	var ProductResponses domain.Product
	err = pr.DB.Raw("SELECT * FROM products WHERE name = ?", p.Name).Scan(&ProductResponses).Error
	if err != nil {
		log.Println(err.Error())
		return domain.Product{}, err
	}
	return ProductResponses, nil
}

func (pr *productRepository) DeleteProducts(id string) error {
	product_id, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	var count int
	if err := pr.DB.Raw("SELECT COUNT(*) FROM products WHERE id=?", product_id).Scan(&count).Error; err != nil {
		return err
	}
	if count < 1 {
		return errors.New("product for given id does not exist")
	}
	if err := pr.DB.Exec("DELETE FROM products WHERE id=?", product_id).Error; err != nil {
		return err
	}
	if err := pr.DB.Exec("DELETE FROM images WHERE product_id = ?", product_id).Error; err != nil {
		return err
	}
	return nil
}
func (pr *productRepository) CheckProductExist(pid int) (bool, error) {
	var a int
	err := pr.DB.Raw("SELECT COUNT(*) FROM products WHERE id=?", pid).Scan(&a).Error
	if err != nil {
		return false, err
	}
	if a == 0 {
		return false, err
	}
	return true, err
}
func (pr *productRepository) UpdateProduct(pid int, stock int) (models.ProductUpdateReciever, error) {
	if stock <= 0 {
		return models.ProductUpdateReciever{}, errors.New("stock doesnot update invalid input")
	}
	if pr.DB == nil {
		return models.ProductUpdateReciever{}, errors.New("database connection is nil")
	}
	if err := pr.DB.Exec("UPDATE products SET stock = stock + $1 WHERE id = $2", stock, pid).Error; err != nil {
		return models.ProductUpdateReciever{}, err
	}
	var newdetails models.ProductUpdateReciever
	var newQuantity int
	if err := pr.DB.Raw("SELECT stock FROM products WHERE id =?", pid).Scan(&newQuantity).Error; err != nil {
		return models.ProductUpdateReciever{}, err
	}
	newdetails.ProductID = pid
	newdetails.Stock = newQuantity
	return newdetails, nil
}
func (pr *productRepository) DoesProductExist(productID int) (bool, error) {
	var count int
	err := pr.DB.Raw("SELECT COUNT(*) FROM products WHERE id = ?", productID).Scan(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (pr *productRepository) UpdateProductImage(productID int, url string) error {
	err := pr.DB.Exec("INSERT INTO images (product_id,url) VALUES ($1,$2) RETURNING * ", productID, url).Error
	if err != nil {
		return errors.New("error while insert image to database")
	}
	return nil
}
func (pr *productRepository) DisplayImages(productID int) (domain.Product, []domain.Image, error) {
	var product domain.Product
	var image []domain.Image
	err := pr.DB.Raw(`SELECT * FROM products WHERE product_id = $1`, productID).Scan(&product).Error
	if err != nil {
		return domain.Product{}, []domain.Image{}, err
	}
	err = pr.DB.Raw(`SELECT * FROM images WHERE product_id = $1`, productID).Scan(&image).Error
	if err != nil {
		return domain.Product{}, []domain.Image{}, err
	}
	return product, image, nil
}
func (pr *productRepository) ShowImages(productID int) ([]models.Image, error) {
	var image []models.Image
	err := pr.DB.Raw(`SELECT url FROM images  WHERE images.product_id = $1`, productID).Scan(&image).Error
	if err != nil {
		return nil, err
	}
	return image, nil
}
