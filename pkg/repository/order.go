package repository

import (
	"Zhooze/pkg/domain"
	"Zhooze/pkg/helper"
	interfaces "Zhooze/pkg/repository/interface"
	"fmt"

	"Zhooze/pkg/utils/models"
	"errors"

	"gorm.io/gorm"
)

type orderRepository struct {
	DB *gorm.DB
}

func NewOrderRepository(DB *gorm.DB) interfaces.OrderRepository {
	return &orderRepository{
		DB: DB,
	}
}

func (or *orderRepository) DoesCartExist(userID int) (bool, error) {

	var exist bool
	err := or.DB.Raw("select exists(select 1 from carts where user_id = ?)", userID).Scan(&exist).Error
	if err != nil {
		return false, err
	}

	return exist, nil
}
func (or *orderRepository) AddressExist(orderBody models.OrderIncoming) (bool, error) {

	var count int
	if err := or.DB.Raw("SELECT COUNT(*) FROM addresses WHERE user_id = ? AND id = ?", orderBody.UserID, orderBody.AddressID).Scan(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil

}
func (or *orderRepository) PaymentExist(orderBody models.OrderIncoming) (bool, error) {
	var count int
	if err := or.DB.Raw("SELECT count(*) FROM payment_methods WHERE id = ?", orderBody.PaymentID).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil

}
func (or *orderRepository) CheckOrderID(orderId int) (bool, error) {
	var count int
	err := or.DB.Raw("SELECT COUNT(*) FROM orders WHERE id = ?", orderId).Scan(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (or *orderRepository) GetOrderDetails(userId int, page int, count int) ([]models.FullOrderDetails, error) {
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * count
	var orderDetails []models.OrderDetails
	or.DB.Raw("SELECT id as order_id,final_price,shipment_status,payment_status FROM orders WHERE user_id = ? LIMIT ? OFFSET ? ", userId, count, offset).Scan(&orderDetails)
	var fullOrderDetails []models.FullOrderDetails
	for _, od := range orderDetails {
		var orderProductDetails []models.OrderProductDetails
		or.DB.Raw(`SELECT
		order_items.product_id,
		products.name AS product_name,
		order_items.quantity,
		order_items.total_price
	    FROM
		order_items
	    INNER JOIN
		products ON order_items.product_id = products.id
	    WHERE
		order_items.order_id = $1 `, od.OrderId).Scan(&orderProductDetails)
		fullOrderDetails = append(fullOrderDetails, models.FullOrderDetails{OrderDetails: od, OrderProductDetails: orderProductDetails})
	}
	return fullOrderDetails, nil
}

func (or *orderRepository) GetOrderDetail(orderId int) (models.OrderDetails, error) {
	var OrderDetails models.OrderDetails

	if err := or.DB.Raw("select id,final_price,shipment_status,payment_status from orders where id = ?", orderId).Scan(&OrderDetails).Error; err != nil {
		return models.OrderDetails{}, err
	}
	return OrderDetails, nil
}

func (or *orderRepository) GetShipmentStatus(orderID int) (string, error) {
	var status string
	err := or.DB.Raw("SELECT shipment_status FROM orders WHERE id= ?", orderID).Scan(&status).Error
	if err != nil {
		return "", err
	}
	return status, nil
}
func (or *orderRepository) UserOrderRelationship(orderID int, userID int) (int, error) {

	var testUserID int
	err := or.DB.Raw("select user_id from orders where id = ?", orderID).Scan(&testUserID).Error
	if err != nil {
		return -1, err
	}
	return testUserID, nil
}

func (or *orderRepository) ApproveOrder(orderID int) error {
	err := or.DB.Exec("UPDATE orders SET shipment_status = 'order placed' , approval = 'true' WHERE id = ?", orderID).Error
	if err != nil {
		return err
	}
	return nil
}

func (or *orderRepository) CancelOrders(orderID int) error {
	status := "cancelled"
	err := or.DB.Exec("UPDATE orders SET shipment_status = ? , approval='false' WHERE id = ? ", status, orderID).Error
	if err != nil {
		return err
	}
	var paymentMethod int
	err = or.DB.Raw("SELECT payment_method_id FROM orders WHERE id = ? ", orderID).Scan(&paymentMethod).Error
	if err != nil {
		return err
	}
	if paymentMethod == 3 || paymentMethod == 2 {
		err = or.DB.Exec("UPDATE orders SET payment_status = 'refunded' WHERE id = ?", orderID).Error
		if err != nil {
			return err
		}
	}
	return nil
}
func (or *orderRepository) PaymentMethodID(orderID int) (int, error) {
	var a int
	err := or.DB.Raw("SELECT payment_method_id FROM orders WHERE id = ?", orderID).Scan(&a).Error
	if err != nil {
		return 0, err
	}
	return a, nil
}
func (or *orderRepository) GetProductDetailsFromOrders(orderID int) ([]models.OrderProducts, error) {
	var OrderProductDetails []models.OrderProducts
	if err := or.DB.Raw("SELECT product_id,quantity as stock FROM order_items WHERE order_id = ?", orderID).Scan(&OrderProductDetails).Error; err != nil {
		return []models.OrderProducts{}, err
	}
	return OrderProductDetails, nil
}

func (or *orderRepository) UpdateQuantityOfProduct(orderProducts []models.OrderProducts) error {

	for _, od := range orderProducts {

		var quantity int
		if err := or.DB.Raw("SELECT stock FROM products WHERE id = ?", od.ProductId).Scan(&quantity).Error; err != nil {
			return err
		}

		od.Stock += quantity
		if err := or.DB.Exec("UPDATE products SET stock = ? WHERE id = ?", od.Stock, od.ProductId).Error; err != nil {
			return err
		}
	}
	return nil

}
func (or *orderRepository) PaymentAlreadyPaid(orderID int) (bool, error) {
	var a bool
	err := or.DB.Raw("SELECT shipment_status = 'processing' AND payment_status = 'paid' FROM orders WHERE id = ?", orderID).Scan(&a).Error
	if err != nil {
		return false, err
	}
	return a, nil
}

func (or *orderRepository) GetOrderDetailsByOrderId(orderID int) (models.CombinedOrderDetails, error) {

	var orderDetails models.CombinedOrderDetails
	err := or.DB.Raw(`SELECT
    orders.id as order_id,
    orders.final_price,
    orders.shipment_status,
    orders.payment_status,
    users.firstname,
    users.email,
    users.phone,
    addresses.house_name,
    addresses.state,
    addresses.street,
    addresses.city,
    addresses.pin
FROM
    orders
INNER JOIN
    users ON orders.user_id = users.id
INNER JOIN
    addresses ON users.id = addresses.user_id
WHERE
    orders.id = ?`, orderID).Scan(&orderDetails).Error
	if err != nil {
		return models.CombinedOrderDetails{}, nil
	}

	return orderDetails, nil
}

func (or *orderRepository) OrderItems(ob models.OrderIncoming, price float64) (int, error) {
	var id int
	query := `
    INSERT INTO orders (created_at , user_id , address_id , payment_method_id , final_price)
    VALUES (NOW(),?, ?, ?, ?)
    RETURNING id`
	or.DB.Raw(query, ob.UserID, ob.AddressID, ob.PaymentID, price).Scan(&id)
	return id, nil
}

func (or *orderRepository) OrderExist(orderID int) error {
	err := or.DB.Raw("SELECT id FROM orders WHERE id = ?", orderID).Error
	if err != nil {
		return err
	}
	return nil
}
func (or *orderRepository) UpdateOrder(orderID int) error {
	err := or.DB.Exec("UPDATE orders SET Shipment_status = 'processing' WHERE id = ?", orderID).Error
	if err != nil {
		return err
	}
	return nil
}

func (or *orderRepository) AddOrderProducts(order_id int, cart []models.Cart) error {
	query := `
    INSERT INTO order_items (order_id,product_id,quantity,total_price)
    VALUES (?, ?, ?, ?) `
	for _, v := range cart {
		var productID int
		if err := or.DB.Raw("SELECT id FROM products WHERE name = $1", v.ProductName).Scan(&productID).Error; err != nil {
			return err
		}
		if err := or.DB.Exec(query, order_id, productID, v.Quantity, v.TotalPrice).Error; err != nil {
			return err
		}
	}
	return nil
}
func (or *orderRepository) UpdateCartAfterOrder(userID, productID int, quantity float64) error {
	err := or.DB.Exec("DELETE FROM carts WHERE user_id = ? and product_id = ?", userID, productID).Error
	if err != nil {
		return err
	}

	err = or.DB.Exec("UPDATE products SET stock = stock - ? WHERE id = ?", quantity, productID).Error
	if err != nil {
		return err
	}

	return nil
}
func (or *orderRepository) GetBriefOrderDetails(orderID int) (domain.OrderSuccessResponse, error) {
	var orderSuccessResponse domain.OrderSuccessResponse
	err := or.DB.Raw(`SELECT id as order_id,shipment_status FROM orders WHERE id = ?`, orderID).Scan(&orderSuccessResponse).Error
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	return orderSuccessResponse, nil
}
func (or *orderRepository) UpdateStockOfProduct(orderProducts []models.OrderProducts) error {
	for _, ok := range orderProducts {
		var quantity int
		if err := or.DB.Raw("SELECT stock FROM products WHERE id = ?", ok.ProductId).Scan(&quantity).Error; err != nil {
			return err
		}
		ok.Stock += quantity
		if err := or.DB.Exec("UPDATE products SET stock  = ? WHERE id = ?", ok.Stock, ok.ProductId).Error; err != nil {
			return err
		}
	}
	return nil
}
func (or *orderRepository) GetAllOrderDetailsBrief(page, count int) ([]models.CombinedOrderDetails, error) {
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * count
	var orderDatails []models.CombinedOrderDetails
	err := or.DB.Raw("SELECT orders.id as order_id,orders.final_price,orders.shipment_status,orders.payment_status,users.firstname,users.email,users.phone,addresses.house_name,addresses.street,addresses.city,addresses.state,addresses.pin FROM orders INNER JOIN users ON orders.user_id = users.id INNER JOIN addresses ON orders.address_id = addresses.id limit ? offset ?", count, offset).Scan(&orderDatails).Error
	if err != nil {
		return []models.CombinedOrderDetails{}, nil
	}
	return orderDatails, nil

}

func (or *orderRepository) AddpaymentMethod(paymentID int, orderID uint) error {
	fmt.Println("payment id : ", orderID)
	err := or.DB.Exec(`UPDATE orders SET payment_method_id = $1 WHERE id = $2`, paymentID, orderID).Error
	if err != nil {
		return err
	}
	return nil
}

func (or *orderRepository) CheckAddressAvailabilityWithID(addressID, userID int) bool {
	var count int
	if err := or.DB.Raw("SELECT COUNT(*) FROM addresses WHERE id = ? AND user_id = ?", addressID, userID).Scan(&count).Error; err != nil {
		return false
	}
	return count > 0
}
func (or *orderRepository) CheckCartAvailabilityWithID(cartID, UserID int) bool {

	var count int
	if err := or.DB.Raw("SELECT COUNT(*) FROM cart_items JOIN carts ON cart_items.cart_id = carts.id WHERE cart_items.cart_id = ? AND carts.user_id = ?", cartID, UserID).Scan(&count).Error; err != nil {
		return false
	}
	return count > 0
}
func (or *orderRepository) FindOrderStock(cartID int) (int, error) {
	var count int
	if err := or.DB.Raw("SELECT quantity FROM cart_items WHERE cart_id = ?", cartID).Scan(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
func (or *orderRepository) AddAmountToOrder(Price float64, orderID uint) error {
	err := or.DB.Exec("UPDATE orders SET final_price = ? WHERE id = ?", Price, orderID).Error
	if err != nil {
		return err
	}
	return nil
}
func (or *orderRepository) GetOrder(orderID int) (domain.Order, error) {
	var order domain.Order
	err := or.DB.Raw("SELECT * FROM orders WHERE id = ?", orderID).Scan(&order).Error
	if err != nil {
		return domain.Order{}, err
	}
	return order, nil
}

func (or *orderRepository) FindProductFromCart(cartID int) (int, error) {
	var p int
	if err := or.DB.Raw("SELECT product_id FROM cart_items WHERE cart_id = ?", cartID).Scan(&p).Error; err != nil {
		return 0, err
	}
	return p, nil
}
func (or *orderRepository) CartEmpty(cartID int) error {
	if err := or.DB.Exec("DELETE FROM cart_items WHERE cart_id = ?", cartID).Error; err != nil {
		return err
	}
	return nil

}
func (or *orderRepository) ProductStockMinus(productID, stock int) error {
	err := or.DB.Exec("UPDATE products SET stock = stock - ? WHERE id = ?", stock, productID).Error
	if err != nil {
		return err
	}
	return nil
}
func (or *orderRepository) GetPaymentId(paymentID int) bool {
	var count int
	if err := or.DB.Raw("SELECT COUNT(*) FROM payment_methods WHERE id = ? ", paymentID).Scan(&count).Error; err != nil {
		return false
	}
	return count > 0
}
func (or *orderRepository) TotalAmountInCart(userID int) (float64, error) {
	var price float64
	if err := or.DB.Raw("SELECT sum(total_price) FROM carts WHERE  user_id= $1", userID).Scan(&price).Error; err != nil {
		return 0, err
	}
	return price, nil
}
func (or *orderRepository) GetCouponDiscountPrice(UserID int, Total float64) (float64, error) {
	discountPrice, err := helper.GetCouponDiscountPrice(UserID, Total, or.DB)
	if err != nil {
		return 0.0, err
	}

	return discountPrice, nil

}
func (or *orderRepository) UpdateCouponDetails(discount_price float64, UserID int) error {

	if discount_price != 0.0 {
		err := or.DB.Exec("update used_coupons set used = true where user_id = ?", UserID).Error
		if err != nil {
			return err
		}
	}
	return nil
}
func (or *orderRepository) GetAllAddresses(userID int) ([]models.AddressInfoResponse, error) {
	var addressResponse []models.AddressInfoResponse
	err := or.DB.Raw(`SELECT * FROM addresses WHERE user_id = $1`, userID).Scan(&addressResponse).Error
	if err != nil {
		return []models.AddressInfoResponse{}, err
	}

	return addressResponse, nil
}
func (or *orderRepository) GetAllPaymentOption() ([]models.PaymentDetails, error) {
	var paymentMethods []models.PaymentDetails
	err := or.DB.Raw("SELECT * FROM payment_methods").Scan(&paymentMethods).Error
	if err != nil {
		return []models.PaymentDetails{}, err
	}

	return paymentMethods, nil

}
func (or *orderRepository) GetAddressFromOrderId(orderID int) (models.AddressInfoResponse, error) {
	var addressInfoResponse models.AddressInfoResponse
	var addressId int
	if err := or.DB.Raw("SELECT address_id FROM orders WHERE id =?", orderID).Scan(&addressId).Error; err != nil {
		return models.AddressInfoResponse{}, errors.New("first in orders")
	}
	if err := or.DB.Raw("SELECT * FROM addresses WHERE id=?", addressId).Scan(&addressInfoResponse).Error; err != nil {
		return models.AddressInfoResponse{}, errors.New("second  in address")
	}
	return addressInfoResponse, nil
}
func (or *orderRepository) GetOrderDetailOfAproduct(orderID int) (models.OrderDetails, error) {
	var OrderDetails models.OrderDetails

	if err := or.DB.Raw("SELECT id,final_price,shipment_status,payment_status FROM orders WHERE id = ?", orderID).Scan(&OrderDetails).Error; err != nil {
		return models.OrderDetails{}, err
	}
	return OrderDetails, nil
}

func (or *orderRepository) GetProductsInCart(cart_id int) ([]int, error) {

	var cart_products []int

	if err := or.DB.Raw("select product_id from cart_items where cart_id=?", cart_id).Scan(&cart_products).Error; err != nil {
		return []int{}, err
	}

	return cart_products, nil

}
func (or *orderRepository) FindProductNames(product_id int) (string, error) {

	var product_name string

	if err := or.DB.Raw("select name from products where id=?", product_id).Scan(&product_name).Error; err != nil {
		return "", err
	}

	return product_name, nil

}

func (or *orderRepository) FindCartQuantity(cart_id, product_id int) (int, error) {

	var quantity int

	if err := or.DB.Raw("select quantity from cart_items where cart_id=$1 and product_id=$2", cart_id, product_id).Scan(&quantity).Error; err != nil {
		return 0, err
	}

	return quantity, nil

}

func (or *orderRepository) FindPrice(product_id int) (float64, error) {

	var price float64

	if err := or.DB.Raw("select price from products where id=?", product_id).Scan(&price).Error; err != nil {
		return 0, err
	}

	return price, nil

}
func (or *orderRepository) FindStock(id int) (int, error) {
	var stock int
	err := or.DB.Raw("SELECT stock FROM prodcuts WHERE id = ?", id).Scan(&stock).Error
	if err != nil {
		return 0, err
	}

	return stock, nil
}
