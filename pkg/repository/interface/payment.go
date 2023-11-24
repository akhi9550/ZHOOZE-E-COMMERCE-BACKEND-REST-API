package interfaces

type PaymentRepository interface {
	CheckPaymentStatus(orderID int) (string, error)
	UpdatePaymentDetails(orderID int, paymentID string) error
	AddRazorPayDetails(orderID int, razorPayOrderID string) error
	UpdateShipmentAndPaymentByOrderID(shipmentStatus string, paymentStatus string, orderID int) error
}
