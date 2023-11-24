package repository

import (
	interfaces "Zhooze/pkg/repository/interface"
	"Zhooze/pkg/utils/models"
	"errors"

	"gorm.io/gorm"
)

type cartRepository struct {
	DB *gorm.DB
}

func NewCartRepository(DB *gorm.DB) interfaces.CartRepository {
	return &cartRepository{
		DB: DB,
	}
}

func (cr *cartRepository) DisplayCart(userID int) ([]models.Cart, error) {

	var count int
	if err := cr.DB.Raw("SELECT COUNT(*) FROM carts WHERE user_id = ? ", userID).Scan(&count).Error; err != nil {
		return []models.Cart{}, err
	}

	if count == 0 {
		return []models.Cart{}, nil
	}

	var cartResponse []models.Cart

	if err := cr.DB.Raw("SELECT carts.user_id,users.firstname as user_name,carts.product_id,products.name as product_name,carts.quantity,carts.total_price from carts INNER JOIN users ON carts.user_id = users.id INNER JOIN products ON carts.product_id = products.id WHERE user_id = ?", userID).First(&cartResponse).Error; err != nil {
		return []models.Cart{}, err
	}
	return cartResponse, nil

}
func (cr *cartRepository) GetTotalPrice(userID int) (models.CartTotal, error) {

	var cartTotal models.CartTotal
	err := cr.DB.Raw("SELECT COALESCE(SUM(total_price), 0) FROM carts WHERE user_id = ?", userID).Scan(&cartTotal.TotalPrice).Error
	if err != nil {
		return models.CartTotal{}, err
	}
	err = cr.DB.Raw("SELECT COALESCE(SUM(total_price), 0) FROM carts WHERE user_id = ?", userID).Scan(&cartTotal.FinalPrice).Error
	if err != nil {
		return models.CartTotal{}, err
	}
	err = cr.DB.Raw("SELECT firstname as user_name FROM users WHERE id = ?", userID).Scan(&cartTotal.UserName).Error
	if err != nil {
		return models.CartTotal{}, err
	}

	return cartTotal, nil

}
func (cr *cartRepository) CartExist(userID int) (bool, error) {
	var count int
	if err := cr.DB.Raw("SELECT COUNT(*) FROM carts WHERE user_id = ? ", userID).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil

}
func (cr *cartRepository) EmptyCart(userID int) error {
	if err := cr.DB.Exec("DELETE FROM carts WHERE  user_id = ?", userID).Error; err != nil {
		return err
	}
	return nil
}
func (cr *cartRepository) CheckProduct(product_id int) (bool, string, error) {
	var count int
	err := cr.DB.Raw("SELECT COUNT(*) FROM products WHERE id = ?", product_id).Scan(&count).Error
	if err != nil {
		return false, "", err
	}
	if count > 0 {
		var category string
		err := cr.DB.Raw("SELECT categories.category FROM categories INNER JOIN products ON products.category_id = categories.id WHERE products.id = ?", product_id).Scan(&category).Error

		if err != nil {
			return false, "", err
		}
		return true, category, nil
	}
	return false, "", nil
}
func (cr *cartRepository) QuantityOfProductInCart(userId int, productId int) (int, error) {
	var productQty int
	err := cr.DB.Raw("SELECT quantity FROM carts WHERE user_id = ? AND product_id = ?", userId, productId).Scan(&productQty).Error
	if err != nil {
		return 0, err
	}
	return productQty, nil
}
func (cr *cartRepository) AddItemIntoCart(userId int, productId int, Quantity int, productprice float64) error {
	if err := cr.DB.Exec("insert into carts (user_id,product_id,quantity,total_price) values(?,?,?,?)", userId, productId, Quantity, productprice).Error; err != nil {
		return err
	}
	return nil

}

func (cr *cartRepository) TotalPriceForProductInCart(userID int, productID int) (float64, error) {

	var totalPrice float64
	if err := cr.DB.Raw("SELECT SUM(total_price) as total_price FROM carts  WHERE user_id = ? AND product_id = ?", userID, productID).Scan(&totalPrice).Error; err != nil {
		return 0.0, err
	}
	return totalPrice, nil
}
func (cr *cartRepository) UpdateCart(quantity int, price float64, userID int, product_id int) error {
	if err := cr.DB.Exec(`UPDATE carts
	SET quantity = ?, total_price = ? 
	WHERE user_id = ? and product_id = ?`, quantity, price, product_id, userID).Error; err != nil {
		return err
	}

	return nil

}
func (cr *cartRepository) ProductExist(userID int, productID int) (bool, error) {
	var count int
	if err := cr.DB.Raw("SELECT count(*) FROM carts  WHERE carts.user_id = ? AND product_id = ?", userID, productID).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil

}
func (cr *cartRepository) GetQuantityAndProductDetails(userId int, productId int, cartDetails struct {
	Quantity   int
	TotalPrice float64
}) (struct {
	Quantity   int
	TotalPrice float64
}, error) {
	if err := cr.DB.Raw("SELECT quantity,total_price FROM carts WHERE user_id = ? AND product_id = ?", userId, productId).Scan(&cartDetails).Error; err != nil {
		return struct {
			Quantity   int
			TotalPrice float64
		}{}, err
	}
	return cartDetails, nil
}
func (cr *cartRepository) RemoveProductFromCart(userID int, product_id int) error {

	if err := cr.DB.Exec("DELETE FROM carts WHERE user_id = ? AND product_id = ?", uint(userID), uint(product_id)).Error; err != nil {
		return err
	}

	return nil
}
func (cr *cartRepository) UpdateCartDetails(cartDetails struct {
	Quantity   int
	TotalPrice float64
}, userId int, productId int) error {
	if err := cr.DB.Raw("UPDATE carts SET quantity = ? , total_price = ? WHERE user_id = ? AND product_id = ? ", cartDetails.Quantity, cartDetails.TotalPrice, userId, productId).Scan(&cartDetails).Error; err != nil {
		return err
	}
	return nil

}
func (cr *cartRepository) CartAfterRemovalOfProduct(user_id int) ([]models.Cart, error) {
	var cart []models.Cart
	if err := cr.DB.Raw("SELECT carts.product_id,products.name as product_name,carts.quantity,carts.total_price FROM carts INNER JOIN products on carts.product_id = products.id WHERE carts.user_id = ?", user_id).Scan(&cart).Error; err != nil {
		return []models.Cart{}, err
	}
	return cart, nil
}
func (cr *cartRepository) GetAllItemsFromCart(userID int) ([]models.Cart, error) {
	var count int
	var cartResponse []models.Cart
	err := cr.DB.Raw("SELECT COUNT(*) FROM carts WHERE user_id = ?", userID).Scan(&count).Error
	if err != nil {
		return []models.Cart{}, err
	}
	if count == 0 {
		return []models.Cart{}, nil
	}
	err = cr.DB.Raw("SELECT carts.user_id,users.firstname as user_name,carts.product_id,products.name as product_name,carts.quantity,carts.total_price from carts INNER JOIN users on carts.user_id = users.id INNER JOIN products ON carts.product_id = products.id where user_id = ?", userID).First(&cartResponse).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if len(cartResponse) == 0 {
				return []models.Cart{}, nil
			}
			return []models.Cart{}, err
		}
		return []models.Cart{}, err
	}
	return cartResponse, nil
}
func (cr *cartRepository) GetTotalPriceFromCart(userID int) (float64, error) {
	var totalPrice float64
	err := cr.DB.Raw("SELECT COALESCE(SUM(total_price), 0) FROM carts WHERE user_id = ?", userID).Scan(&totalPrice).Error
	if err != nil {
		return 0.0, err
	}
	return totalPrice, nil

}
