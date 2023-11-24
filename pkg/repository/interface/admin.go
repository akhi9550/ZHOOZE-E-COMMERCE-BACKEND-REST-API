package interfaces

import (
	"Zhooze/pkg/domain"
	"Zhooze/pkg/utils/models"
	"time"
)

type AdminRepository interface {
	LoginHandler(adminDetails models.AdminLogin) (domain.Admin, error)
	DashBoardUserDetails() (models.DashBoardUser, error)
	DashBoardProductDetails() (models.DashBoardProduct, error)
	ShowAllUsersIn(page, count int) ([]models.UserDetailsAtAdmin, error)
	GetUserByID(id string) (domain.User, error)
	UpdateBlockUserByID(user domain.User) error
	DashBoardOrder() (models.DashboardOrder, error)
	TotalRevenue() (models.DashboardRevenue, error)
	AmountDetails() (models.DashboardAmount, error)
	FilteredSalesReport(startTime time.Time, endTime time.Time) (models.SalesReport, error)
	AddPaymentMethod(pay models.NewPaymentMethod) (domain.PaymentMethod, error)
	ListPaymentMethods() ([]domain.PaymentMethod, error)
	CheckIfPaymentMethodAlreadyExists(payment string) (bool, error)
	DeletePaymentMethod(id int) error
}
