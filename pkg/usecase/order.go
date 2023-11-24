package usecase

import (
	"Zhooze/pkg/domain"
	interfaces "Zhooze/pkg/repository/interface"
	services "Zhooze/pkg/usecase/interface"
	"Zhooze/pkg/utils/models"
	"errors"
	"fmt"

	"github.com/jinzhu/copier"
)

type orderUseCase struct {
	orderRepository interfaces.OrderRepository
	cartRepository interfaces.CartRepository
}

func NewOrderUseCase(repository interfaces.OrderRepository,cartRepo interfaces.CartRepository) services.OrderUseCase {
	return &orderUseCase{
		orderRepository: repository,
		cartRepository: cartRepo,
	}
}
func (or *orderUseCase) OrderItemsFromCart(orderFromCart models.OrderFromCart, userID int) (domain.OrderSuccessResponse, error) {
	var orderBody models.OrderIncoming
	err := copier.Copy(&orderBody, &orderFromCart)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	orderBody.UserID = userID
	cartExist, err := or.orderRepository.DoesCartExist(userID)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	if !cartExist {
		return domain.OrderSuccessResponse{}, errors.New("cart empty can't order")
	}

	addressExist, err := or.orderRepository.AddressExist(orderBody)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}

	if !addressExist {
		return domain.OrderSuccessResponse{}, errors.New("address does not exist")
	}
	PaymentExist, err := or.orderRepository.PaymentExist(orderBody)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}

	if !PaymentExist {
		return domain.OrderSuccessResponse{}, errors.New("paymentmethod does not exist")
	}
	cartItems, err := or.cartRepository.GetAllItemsFromCart(orderBody.UserID)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}

	total, err := or.orderRepository.TotalAmountInCart(orderBody.UserID)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	discount_price, err := or.orderRepository.GetCouponDiscountPrice(int(orderBody.UserID), total)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}

	err = or.orderRepository.UpdateCouponDetails(discount_price, orderBody.UserID)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	FinalPrice := total - discount_price
	order_id, err := or.orderRepository.OrderItems(orderBody, FinalPrice)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}

	if err := or.orderRepository.AddOrderProducts(order_id, cartItems); err != nil {
		return domain.OrderSuccessResponse{}, err
	}

	orderSuccessResponse, err := or.orderRepository.GetBriefOrderDetails(order_id)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	var orderItemDetails domain.OrderItem
	for _, c := range cartItems {
		orderItemDetails.ProductID = c.ProductID
		orderItemDetails.Quantity = c.Quantity
		err := or.orderRepository.UpdateCartAfterOrder(userID, int(orderItemDetails.ProductID), orderItemDetails.Quantity)
		if err != nil {
			return domain.OrderSuccessResponse{}, err
		}
	}
	return orderSuccessResponse, nil
}
func (or *orderUseCase) GetOrderDetails(userId int, page int, count int) ([]models.FullOrderDetails, error) {

	fullOrderDetails, err := or.orderRepository.GetOrderDetails(userId, page, count)
	if err != nil {
		return []models.FullOrderDetails{}, err
	}
	return fullOrderDetails, nil

}

func (or *orderUseCase) CancelOrders(orderID int, userId int) error {
	userTest, err := or.orderRepository.UserOrderRelationship(orderID, userId)
	if err != nil {
		return err
	}
	if userTest != userId {
		return errors.New("the order is not done by this user")
	}
	orderProductDetails, err := or.orderRepository.GetProductDetailsFromOrders(orderID)
	if err != nil {
		return err
	}
	shipmentStatus, err := or.orderRepository.GetShipmentStatus(orderID)
	if err != nil {
		return err
	}
	if shipmentStatus == "delivered" {
		return errors.New("item already delivered, cannot cancel")
	}

	if shipmentStatus == "pending" || shipmentStatus == "returned" || shipmentStatus == "return" {
		message := fmt.Sprint(shipmentStatus)
		return errors.New("the order is in" + message + ", so no point in cancelling")
	}

	if shipmentStatus == "cancelled" {
		return errors.New("the order is already cancelled, so no point in cancelling")
	}
	err = or.orderRepository.CancelOrders(orderID)
	if err != nil {
		return err
	}
	err = or.orderRepository.UpdateQuantityOfProduct(orderProductDetails)
	if err != nil {
		return err
	}
	return nil

}
func (or *orderUseCase) Checkout(userID int) (models.CheckoutDetails, error) {
	allUserAddress, err := or.orderRepository.GetAllAddresses(userID)
	if err != nil {
		return models.CheckoutDetails{}, err
	}
	paymentDetails, err := or.orderRepository.GetAllPaymentOption()
	if err != nil {
		return models.CheckoutDetails{}, err
	}
	cartItems, err := or.cartRepository.DisplayCart(userID)
	if err != nil {
		return models.CheckoutDetails{}, err
	}
	grandTotal, err := or.cartRepository.GetTotalPrice(userID)
	if err != nil {
		return models.CheckoutDetails{}, err
	}

	return models.CheckoutDetails{
		AddressInfoResponse: allUserAddress,
		Payment_Method:      paymentDetails,
		Cart:                cartItems,
		Total_Price:         grandTotal.FinalPrice,
	}, nil
}
func (or *orderUseCase) PaymentMethodID(order_id int) (int, error) {
	id, err := or.orderRepository.PaymentMethodID(order_id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
func (or *orderUseCase) ExecutePurchaseCOD(orderID int) error {
	err := or.orderRepository.OrderExist(orderID)
	if err != nil {
		return err
	}
	shipmentStatus, err := or.orderRepository.GetShipmentStatus(orderID)
	if err != nil {
		return err
	}
	if shipmentStatus == "delivered" {
		return errors.New("item  delivered, cannot pay")
	}
	if shipmentStatus == "order placed" {
		return errors.New("item placed, cannot pay")
	}
	if shipmentStatus == "cancelled" || shipmentStatus == "returned" || shipmentStatus == "return" {
		message := fmt.Sprint(shipmentStatus)
		return errors.New("the order is in" + message + "so can't paid")
	}
	if shipmentStatus == "processing" {
		return errors.New("the order is already paid")
	}
	err = or.orderRepository.UpdateOrder(orderID)
	if err != nil {
		return err
	}

	return nil

}
func (or *orderUseCase) GetAllOrderDetailsForAdmin(page, pagesize int) ([]models.CombinedOrderDetails, error) {
	orderDetail, err := or.orderRepository.GetAllOrderDetailsBrief(page, pagesize)
	if err != nil {
		return []models.CombinedOrderDetails{}, err
	}
	return orderDetail, nil
}
func (or *orderUseCase) ApproveOrder(orderId int) error {
	ShipmentStatus, err := or.orderRepository.GetShipmentStatus(orderId)
	if err != nil {
		return err
	}
	if ShipmentStatus == "cancelled" {
		return errors.New("the order is cancelled,cannot approve it")
	}
	if ShipmentStatus == "pending" {
		return errors.New("the order is pending,cannot approve it")
	}
	if ShipmentStatus == "processing" {
		err := or.orderRepository.ApproveOrder(orderId)
		if err != nil {
			return err
		}
		return nil
	}
	// if the shipment status is not processing or cancelled. Then it is defenetely cancelled
	return nil
}
func (or *orderUseCase) CancelOrderFromAdmin(order_id int) error {
	ok, err := or.orderRepository.CheckOrderID(order_id)
	fmt.Println(err)
	if !ok {
		return err
	}
	orderProduct, err := or.orderRepository.GetProductDetailsFromOrders(order_id)
	if err != nil {
		return err
	}
	err = or.orderRepository.CancelOrders(order_id)
	if err != nil {
		return err
	}
	err = or.orderRepository.UpdateStockOfProduct(orderProduct)
	if err != nil {
		return err
	}
	return nil
}
