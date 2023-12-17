package interfaces

import (
	"Zhooze/pkg/domain"
	"Zhooze/pkg/utils/models"
)

type OrderRepository interface {
	DoesCartExist(userID int) (bool, error)
	AddressExist(orderBody models.OrderIncoming) (bool, error)
	PaymentExist(orderBody models.OrderIncoming) (bool, error)
	PaymentStatus(orderID int) (string, error)
	TotalAmountFromOrder(orderID int) (float64, error)
	UserIDFromOrder(orderID int) (int, error)
	CheckOrderID(orderId int) (bool, error)
	GetOrderDetails(userId int, page int, count int) ([]models.FullOrderDetails, error)
	GetOrderDetail(orderId int) (models.OrderDetails, error)
	GetShipmentStatus(orderID int) (string, error)
	UserOrderRelationship(orderID int, userID int) (int, error)
	ApproveOrder(orderID int) error
	CancelOrders(orderID int) error
	PaymentMethodID(orderID int) (int, error)
	GetProductDetailsFromOrders(orderID int) ([]models.OrderProducts, error)
	UpdateQuantityOfProduct(orderProducts []models.OrderProducts) error
	PaymentAlreadyPaid(orderID int) (bool, error)
	GetOrderDetailsByOrderId(orderID int) (models.CombinedOrderDetails, error)
	OrderItems(ob models.OrderIncoming, price float64) (int, error)
	OrderExist(orderID int) error
	UpdateOrder(orderID int) error
	AddOrderProducts(order_id int, cart []models.Cart) error
	UpdateCartAfterOrder(userID, productID int, quantity float64) error
	GetBriefOrderDetails(orderID int) (domain.OrderSuccessResponse, error)
	UpdateStockOfProduct(orderProducts []models.OrderProducts) error
	GetAllOrderDetailsBrief(page, count int) ([]models.CombinedOrderDetails, error)
	AddpaymentMethod(paymentID int, orderID uint) error
	CheckAddressAvailabilityWithID(addressID, userID int) bool
	CheckCartAvailabilityWithID(cartID, UserID int) bool
	FindOrderStock(cartID int) (int, error)
	AddAmountToOrder(Price float64, orderID uint) error
	GetOrder(orderID int) (domain.Order, error)
	FindProductFromCart(cartID int) (int, error)
	CartEmpty(cartID int) error
	ProductStockMinus(productID, stock int) error
	GetPaymentId(paymentID int) bool
	TotalAmountInCart(userID int) (float64, error)
	GetCouponDiscountPrice(UserID int, Total float64) (float64, error)
	UpdateCouponDetails(discount_price float64, UserID int) error
	GetAllAddresses(userID int) ([]models.AddressInfoResponse, error)
	GetAllPaymentOption() ([]models.PaymentDetails, error)
	GetAddressFromOrderId(orderID int) (models.AddressInfoResponse, error)
	GetOrderDetailOfAproduct(orderID int) (models.OrderDetails, error)
	GetProductsInCart(cart_id int) ([]int, error)
	FindProductNames(product_id int) (string, error)
	FindCartQuantity(cart_id, product_id int) (int, error)
	FindPrice(product_id int) (float64, error)
	FindStock(id int) (int, error)
	UpdateHistory(userID, orderID int, amount float64, reason string) error
	UpdateAmountToWallet(userID int, amount float64) error
	GetDetailedOrderThroughId(orderId int) (models.CombinedOrderDetails, error)
	GetItemsByOrderId(orderId int) ([]models.Invoice, error)
}
