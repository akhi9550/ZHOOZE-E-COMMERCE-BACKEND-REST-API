package interfaces

import "Zhooze/pkg/utils/models"

type PaymentUseCase interface {
	PaymentAlreadyPaid(orderID int) (bool, error)
	MakePaymentRazorPay(orderID int) (models.CombinedOrderDetails, string, error)
	SavePaymentDetails(orderID int, paymentID string) error
}
