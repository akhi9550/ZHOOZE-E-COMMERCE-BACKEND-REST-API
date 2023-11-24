package repository

import (
	interfaces "Zhooze/pkg/repository/interface"

	"gorm.io/gorm"
)

type paymentRepository struct {
	DB *gorm.DB
}

func NewPaymentRepository(DB *gorm.DB) interfaces.PaymentRepository {
	return &paymentRepository{
		DB: DB,
	}
}
func (pt *paymentRepository) CheckPaymentStatus(orderID int) (string, error) {
	var paymentStatus string
	err := pt.DB.Raw(`SELECT payment_status FROM orders WHERE id = $1`, orderID).Scan(&paymentStatus).Error
	if err != nil {
		return "", err
	}
	return paymentStatus, nil
}
func (pt *paymentRepository) UpdatePaymentDetails(orderID int, paymentID string) error {
	err := pt.DB.Exec("UPDATE razer_pays set payment_id = ? WHERE order_id= ?", paymentID, orderID).Error
	if err != nil {
		return err
	}
	return nil
}
func (pt *paymentRepository) AddRazorPayDetails(orderID int, razorPayOrderID string) error {
	err := pt.DB.Exec("INSERT INTO razer_pays (order_id,razor_id) VALUES (?,?)", orderID, razorPayOrderID).Error
	if err != nil {
		return err
	}
	return nil
}
func (pt *paymentRepository) UpdateShipmentAndPaymentByOrderID(shipmentStatus string, paymentStatus string, orderID int) error {
	err := pt.DB.Exec("UPDATE orders SET payment_status = ?,shipment_status = ?  WHERE id = ?", paymentStatus, shipmentStatus, orderID).Error
	if err != nil {
		return err
	}
	return nil
}
