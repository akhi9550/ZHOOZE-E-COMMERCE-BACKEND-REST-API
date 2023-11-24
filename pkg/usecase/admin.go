package usecase

import (
	"Zhooze/pkg/domain"
	"Zhooze/pkg/helper"
	interfaces "Zhooze/pkg/repository/interface"
	services "Zhooze/pkg/usecase/interface"
	"Zhooze/pkg/utils/models"
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
)

type adminUseCase struct {
	adminRepository interfaces.AdminRepository
	orderRepository interfaces.OrderRepository
	paymentRepository interfaces.PaymentRepository
	
}

func NewAdminUseCase(repository interfaces.AdminRepository, orderRepo interfaces.OrderRepository,paymentRepo interfaces.PaymentRepository) services.AdminUseCase {
	return &adminUseCase{
		adminRepository: repository,
		orderRepository: orderRepo,
		paymentRepository: paymentRepo,
	}
}

func (ad *adminUseCase) LoginHandler(adminDetails models.AdminLogin) (domain.TokenAdmin, error) {
	adminCompareDetails, err := ad.adminRepository.LoginHandler(adminDetails)
	if err != nil {
		return domain.TokenAdmin{}, err
	}
	// compare password from database and that provided from admins
	err = bcrypt.CompareHashAndPassword([]byte(adminCompareDetails.Password), []byte(adminDetails.Password))
	if err != nil {
		return domain.TokenAdmin{}, err
	}
	var adminDetailsResponse models.AdminDetailsResponse
	//  copy all details except password and sent it back to the front end
	err = copier.Copy(&adminDetailsResponse, &adminCompareDetails)
	if err != nil {
		return domain.TokenAdmin{}, err
	}

	tokenString, err := helper.GenerateTokenAdmin(adminDetailsResponse)

	if err != nil {
		return domain.TokenAdmin{}, err
	}

	return domain.TokenAdmin{
		Admin: adminDetailsResponse,
		Token: tokenString,
	}, nil
}

func (ad *adminUseCase) DashBoard() (models.CompleteAdminDashboard, error) {
	userDetails, err := ad.adminRepository.DashBoardUserDetails()
	if err != nil {
		return models.CompleteAdminDashboard{}, err
	}
	productDetails, err := ad.adminRepository.DashBoardProductDetails()
	if err != nil {
		return models.CompleteAdminDashboard{}, err
	}
	orderDetails, err := ad.adminRepository.DashBoardOrder()
	if err != nil {
		return models.CompleteAdminDashboard{}, err
	}
	totalRevenue, err := ad.adminRepository.TotalRevenue()
	if err != nil {
		return models.CompleteAdminDashboard{}, err
	}
	amountDetails, err := ad.adminRepository.AmountDetails()
	if err != nil {
		return models.CompleteAdminDashboard{}, err
	}
	return models.CompleteAdminDashboard{
		DashboardUser:    userDetails,
		DashboardProduct: productDetails,
		DashboardOrder:   orderDetails,
		DashboardRevenue: totalRevenue,
		DashboardAmount:  amountDetails,
	}, nil
}
func (ad *adminUseCase) ShowAllUsers(page, count int) ([]models.UserDetailsAtAdmin, error) {
	users, err := ad.adminRepository.ShowAllUsersIn(page, count)
	if err != nil {
		return []models.UserDetailsAtAdmin{}, err
	}
	return users, nil
}
func (ad *adminUseCase) BlockedUser(id string) error {
	user, err := ad.adminRepository.GetUserByID(id)
	if err != nil {
		return err
	}
	if user.Blocked {
		return errors.New("already blocked")
	} else {
		user.Blocked = true
	}
	err = ad.adminRepository.UpdateBlockUserByID(user)
	if err != nil {
		return err
	}
	return nil
}
func (ad *adminUseCase) UnBlockedUser(id string) error {
	user, err := ad.adminRepository.GetUserByID(id)
	if err != nil {
		return err
	}
	if user.Blocked {
		user.Blocked = false
	} else {
		return errors.New("already unblocked")
	}
	err = ad.adminRepository.UpdateBlockUserByID(user)
	if err != nil {
		return err
	}
	return nil
}

func (ad *adminUseCase) FilteredSalesReport(timePeriod string) (models.SalesReport, error) {
	startTime, endTime := helper.GetTimeFromPeriod(timePeriod)
	saleReport, err := ad.adminRepository.FilteredSalesReport(startTime, endTime)

	if err != nil {
		return models.SalesReport{}, err
	}
	return saleReport, nil

}
func (ad *adminUseCase) ExecuteSalesReportByDate(startDate, endDate time.Time) (models.SalesReport, error) {
	orders, err := ad.adminRepository.FilteredSalesReport(startDate, endDate)
	if err != nil {
		return models.SalesReport{}, errors.New("report fetching failed")
	}
	return orders, nil
}

func (ad *adminUseCase) AddPaymentMethod(payment models.NewPaymentMethod) (domain.PaymentMethod, error) {
	exists, err := ad.adminRepository.CheckIfPaymentMethodAlreadyExists(payment.PaymentName)
	if err != nil {
		return domain.PaymentMethod{}, err
	}
	if exists {
		return domain.PaymentMethod{}, errors.New("payment method already exists")
	}
	paymentadd, err := ad.adminRepository.AddPaymentMethod(payment)
	if err != nil {
		return domain.PaymentMethod{}, err
	}
	return paymentadd, nil
}

func (ad *adminUseCase) ListPaymentMethods() ([]domain.PaymentMethod, error) {

	categories, err := ad.adminRepository.ListPaymentMethods()
	if err != nil {
		return []domain.PaymentMethod{}, err
	}
	return categories, nil

}

func (ad *adminUseCase) DeletePaymentMethod(id int) error {

	err := ad.adminRepository.DeletePaymentMethod(id)
	if err != nil {
		return err
	}
	return nil

}
func (ad *adminUseCase) GetAllOrderDetailsForAdmin(page, pagesize int) ([]models.CombinedOrderDetails, error) {
	orderDetail, err := ad.orderRepository.GetAllOrderDetailsBrief(page, pagesize)
	if err != nil {
		return []models.CombinedOrderDetails{}, err
	}
	return orderDetail, nil
}
func (ad *adminUseCase) ApproveOrder(orderId int) error {
	ShipmentStatus, err := ad.orderRepository.GetShipmentStatus(orderId)
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
		err := ad.orderRepository.ApproveOrder(orderId)
		if err != nil {
			return err
		}
		return nil
	}
	// if the shipment status is not processing or cancelled. Then it is defenetely cancelled
	return nil
}
func (ad *adminUseCase) CancelOrderFromAdmin(orderID int) error {
	ok, err := ad.orderRepository.CheckOrderID(orderID)
	fmt.Println(err)
	if !ok {
		return err
	}
	orderProduct, err := ad.orderRepository.GetProductDetailsFromOrders(orderID)
	if err != nil {
		return err
	}
	err = ad.orderRepository.CancelOrders(orderID)
	if err != nil {
		return err
	}
	err = ad.orderRepository.UpdateStockOfProduct(orderProduct)
	if err != nil {
		return err
	}
	payment_status, err := ad.orderRepository.PaymentStatus(orderID)
	if err != nil {
		return err
	}
	amount, err := ad.orderRepository.TotalAmountFromOrder(orderID)
	if err != nil {
		return err
	}
	userID, err := ad.orderRepository.UserIDFromOrder(orderID)
	if err != nil {
		return err
	}
	if payment_status == "refunded" {
		err = ad.adminRepository.UpdateAmountToWallet(userID, amount)
		if err != nil {
			return err
		}
		reason := "Amount credited for  cancellation of order by admin"
		err := ad.adminRepository.UpdateHistory(userID, orderID, amount, reason)
		if err != nil {
			return err
		}
	}
	return nil

}
