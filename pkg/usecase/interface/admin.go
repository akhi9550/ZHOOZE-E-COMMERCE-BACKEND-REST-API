package interfaces

import (
	"Zhooze/pkg/domain"
	"Zhooze/pkg/utils/models"
	"time"
)

type AdminUseCase interface {
	LoginHandler(adminDetails models.AdminLogin) (domain.TokenAdmin, error)
	DashBoard() (models.CompleteAdminDashboard, error)
	ShowAllUsers(page, count int) ([]models.UserDetailsAtAdmin, error)
	BlockedUser(id string) error
	UnBlockedUser(id string) error
	GetAllOrderDetailsForAdmin(page, pagesize int) ([]models.CombinedOrderDetails, error)
	ApproveOrder(orderId int) error
	CancelOrderFromAdmin(order_id int) error
	FilteredSalesReport(timePeriod string) (models.SalesReport, error)
	ExecuteSalesReportByDate(startDate, endDate time.Time) (models.SalesReport, error)
	AddPaymentMethod(payment models.NewPaymentMethod) (domain.PaymentMethod, error)
	ListPaymentMethods() ([]domain.PaymentMethod, error)
	DeletePaymentMethod(id int) error
}
