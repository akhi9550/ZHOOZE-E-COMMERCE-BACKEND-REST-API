package usecase

import (
	"Zhooze/pkg/config"
	interfaces "Zhooze/pkg/repository/interface"
	services "Zhooze/pkg/usecase/interface"

	"Zhooze/pkg/utils/models"
	"errors"
	"fmt"

	"github.com/razorpay/razorpay-go"
)

type paymentUseCase struct {
	paymentRepository interfaces.PaymentRepository
	orderRepository   interfaces.OrderRepository
}

func NewPaymentUseCase(repository interfaces.PaymentRepository, orderRepo interfaces.OrderRepository) services.PaymentUseCase {
	return &paymentUseCase{
		paymentRepository: repository,
		orderRepository:   orderRepo,
	}
}
func (pt *paymentUseCase) PaymentAlreadyPaid(orderID int) (bool, error) {
	AlreadyPayed, err := pt.orderRepository.PaymentAlreadyPaid(orderID)
	if err != nil {
		return false, err
	}
	return AlreadyPayed, nil
}

func (pt *paymentUseCase) MakePaymentRazorPay(orderID int) (models.CombinedOrderDetails, string, error) {
	cfg, _ := config.LoadConfig()
	combinedOrderDetails, err := pt.orderRepository.GetOrderDetailsByOrderId(orderID)
	if err != nil {
		return models.CombinedOrderDetails{}, "", err
	}

	client := razorpay.NewClient(cfg.KEY_ID_FOR_PAY, cfg.SECRET_KEY_FOR_PAY)

	data := map[string]interface{}{
		"amount":   int(combinedOrderDetails.FinalPrice) * 100,
		"currency": "INR",
		"receipt":  "some_receipt_id",
	}
	body, err := client.Order.Create(data, nil)
	if err != nil {

		return models.CombinedOrderDetails{}, "", err
	}

	razorPayOrderID := body["id"].(string)

	err = pt.paymentRepository.AddRazorPayDetails(orderID, razorPayOrderID)
	if err != nil {
		return models.CombinedOrderDetails{}, "", err
	}

	return combinedOrderDetails, razorPayOrderID, nil

}

func (pt *paymentUseCase) SavePaymentDetails(orderID int, paymentID string) error {
	status, err := pt.paymentRepository.CheckPaymentStatus(orderID)
	if err != nil {
		return err
	}
	if status == "not paid" {
		err = pt.paymentRepository.UpdatePaymentDetails(orderID, paymentID)
		if err != nil {
			return err
		}
		err := pt.paymentRepository.UpdateShipmentAndPaymentByOrderID("processing", "paid", orderID)
		if err != nil {
			return err
		}
		return nil
	}
	fmt.Println("❌already paid❌")
	return errors.New("already paid")
}
