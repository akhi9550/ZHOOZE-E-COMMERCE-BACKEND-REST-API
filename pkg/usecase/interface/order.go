package interfaces

import (
	"Zhooze/pkg/domain"
	"Zhooze/pkg/utils/models"
)

type OrderUseCase interface {
	OrderItemsFromCart(orderFromCart models.OrderFromCart, userID int) (domain.OrderSuccessResponse, error)
	GetOrderDetails(userId int, page int, count int) ([]models.FullOrderDetails, error)
	CancelOrders(orderID int, userId int) error
	Checkout(userID int) (models.CheckoutDetails, error)
	PaymentMethodID(order_id int) (int, error)
	ExecutePurchaseCOD(orderID int) error
	GetAllOrderDetailsForAdmin(page, pagesize int) ([]models.CombinedOrderDetails, error)
	ApproveOrder(orderId int) error
	CancelOrderFromAdmin(order_id int) error
}
