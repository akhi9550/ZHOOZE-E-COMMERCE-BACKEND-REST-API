package interfaces

import "Zhooze/pkg/utils/models"

type CartRepository interface {
	DisplayCart(userID int) ([]models.Cart, error)
	GetTotalPrice(userID int) (models.CartTotal, error)
	CartExist(userID int) (bool, error)

	EmptyCart(userID int) error
	CheckProduct(product_id int) (bool, string, error)
	QuantityOfProductInCart(userId int, productId int) (int, error)
	AddItemIntoCart(userId int, productId int, Quantity int, productprice float64) error
	TotalPriceForProductInCart(userID int, productID int) (float64, error)
	UpdateCart(quantity int, price float64, userID int, product_id int) error
	ProductExist(userID int, productID int) (bool, error)
	GetQuantityAndProductDetails(userId int, productId int, cartDetails struct {
		Quantity   int
		TotalPrice float64
	}) (struct {
		Quantity   int
		TotalPrice float64
	}, error)
	RemoveProductFromCart(userID int, product_id int) error
	UpdateCartDetails(cartDetails struct {
		Quantity   int
		TotalPrice float64
	}, userId int, productId int) error
	CartAfterRemovalOfProduct(user_id int) ([]models.Cart, error)
	GetAllItemsFromCart(userID int) ([]models.Cart, error)
	GetTotalPriceFromCart(userID int) (float64, error)
}
