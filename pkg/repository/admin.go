package repository

import (
	"Zhooze/pkg/domain"
	"Zhooze/pkg/helper"
	interfaces "Zhooze/pkg/repository/interface"
	"Zhooze/pkg/utils/models"
	"errors"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type adminRepository struct {
	DB *gorm.DB
}

func NewAdminRepository(DB *gorm.DB) interfaces.AdminRepository {
	return &adminRepository{
		DB: DB,
	}
}

func (ad *adminRepository) LoginHandler(adminDetails models.AdminLogin) (domain.Admin, error) {
	var details domain.Admin
	if err := ad.DB.Raw("SELECT * FROM users WHERE email=? AND isadmin= true", adminDetails.Email).Scan(&details).Error; err != nil {
		return domain.Admin{}, err
	}
	return details, nil
}
func (ad *adminRepository) DashBoardUserDetails() (models.DashBoardUser, error) {
	var userDetails models.DashBoardUser
	err := ad.DB.Raw("SELECT COUNT(*) FROM users WHERE isadmin='false'").Scan(&userDetails.TotalUsers).Error
	if err != nil {
		return models.DashBoardUser{}, nil
	}
	err = ad.DB.Raw("SELECT COUNT(*)  FROM users WHERE blocked=true").Scan(&userDetails.BlockedUser).Error
	if err != nil {
		return models.DashBoardUser{}, nil
	}
	return userDetails, nil
}

func (ad *adminRepository) DashBoardProductDetails() (models.DashBoardProduct, error) {
	var productDetails models.DashBoardProduct
	err := ad.DB.Raw("SELECT COUNT(*) FROM products").Scan(&productDetails.TotalProducts).Error
	if err != nil {
		return models.DashBoardProduct{}, nil
	}
	err = ad.DB.Raw("SELECT COUNT(*) FROM products WHERE stock<=0").Scan(&productDetails.OutofStockProduct).Error
	if err != nil {
		return models.DashBoardProduct{}, nil
	}
	return productDetails, nil
}
func (ad *adminRepository) ShowAllUsersIn(page, count int) ([]models.UserDetailsAtAdmin, error) {
	var user []models.UserDetailsAtAdmin
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * count
	err := ad.DB.Raw("SELECT id,firstname,lastname,email,phone,blocked FROM users WHERE isadmin='false' limit ? offset ?", count, offset).Scan(&user).Error
	if err != nil {
		return []models.UserDetailsAtAdmin{}, err
	}
	return user, nil
}
func (ad *adminRepository) GetUserByID(id string) (domain.User, error) {
	user_id, err := strconv.Atoi(id)
	if err != nil {
		return domain.User{}, err
	}
	var count int
	if err := ad.DB.Raw("SELECT COUNT(*) FROM users WHERE id=?", user_id).Scan(&count).Error; err != nil {

		return domain.User{}, err
	}
	if count < 1 {
		return domain.User{}, errors.New("user for the given id does not exist")

	}
	var userDetails domain.User
	if err := ad.DB.Raw("SELECT * FROM users WHERE id=?", user_id).Scan(&userDetails).Error; err != nil {
		return domain.User{}, err
	}
	return userDetails, nil
}

func (ad *adminRepository) UpdateBlockUserByID(user domain.User) error {
	err := ad.DB.Exec("UPDATE users SET blocked = ? WHERE id = ?", user.Blocked, user.ID).Error
	if err != nil {
		fmt.Println("Error updating user:", err)
		return err
	}
	return nil
}
func (ad *adminRepository) DashBoardOrder() (models.DashboardOrder, error) {
	var orderDetail models.DashboardOrder
	err := ad.DB.Raw("SELECT COUNT(*) FROM orders WHERE payment_status= 'paid' AND approval =true").Scan(&orderDetail.CompletedOrder).Error
	if err != nil {
		return models.DashboardOrder{}, err
	}
	err = ad.DB.Raw("SELECT COUNT(*) FROM orders WHERE shipment_status='pending' OR shipment_status = 'processing'").Scan(&orderDetail.PendingOrder).Error
	if err != nil {
		return models.DashboardOrder{}, err
	}
	err = ad.DB.Raw("select count(*) from orders where shipment_status = 'cancelled'").Scan(&orderDetail.CancelledOrder).Error
	if err != nil {
		return models.DashboardOrder{}, nil
	}

	err = ad.DB.Raw("select count(*) from orders").Scan(&orderDetail.TotalOrder).Error
	if err != nil {
		return models.DashboardOrder{}, nil
	}

	err = ad.DB.Raw("select sum(quantity) from carts").Scan(&orderDetail.TotalOrderItem).Error
	if err != nil {
		return models.DashboardOrder{}, nil
	}
	return orderDetail, nil
}
func (ad *adminRepository) TotalRevenue() (models.DashboardRevenue, error) {
	var revenueDetails models.DashboardRevenue
	startTime := time.Now().AddDate(0, 0, -1)
	endTime := time.Now()
	err := ad.DB.Raw("SELECT COALESCE(SUM(final_price),0) FROM orders WHERE payment_status = 'paid' AND approval = true AND created_at >=? AND created_at <=?", startTime, endTime).Scan(&revenueDetails.TodayRevenue).Error
	if err != nil {
		return models.DashboardRevenue{}, nil
	}
	startTime, endTime = helper.GetTimeFromPeriod("month")
	err = ad.DB.Raw("SELECT COALESCE (SUM(final_price),0) FROM orders WHERE payment_status = 'paid' AND approval = true AND created_at >=? AND created_at <=?", startTime, endTime).Scan(&revenueDetails.MonthRevenue).Error
	if err != nil {
		return models.DashboardRevenue{}, nil
	}
	startTime, endTime = helper.GetTimeFromPeriod("year")
	err = ad.DB.Raw("SELECT COALESCE (SUM(final_price),0) FROM orders WHERE payment_status = 'paid' AND approval = true AND created_at >=? AND created_at <=?", startTime, endTime).Scan(&revenueDetails.YearRevenue).Error
	if err != nil {
		return models.DashboardRevenue{}, nil
	}
	return revenueDetails, nil
}

func (ad *adminRepository) AmountDetails() (models.DashboardAmount, error) {
	var amountDetails models.DashboardAmount
	err := ad.DB.Raw("SELECT COALESCE (SUM(final_price),0) FROM orders WHERE payment_status = 'paid' AND approval = true").Scan(&amountDetails.CreditedAmount).Error
	if err != nil {
		return models.DashboardAmount{}, nil
	}
	err = ad.DB.Raw("SELECT COALESCE(SUM(final_price),0) FROM orders WHERE payment_status = 'not paid' AND shipment_status = 'processing' OR shipment_status = 'pending' OR shipment_status = 'order placed'").Scan(&amountDetails.PendingAmount).Error
	if err != nil {
		return models.DashboardAmount{}, nil
	}
	return amountDetails, nil
}
func (ad *adminRepository) FilteredSalesReport(startTime time.Time, endTime time.Time) (models.SalesReport, error) {
	var salesReport models.SalesReport
	result := ad.DB.Raw("SELECT COALESCE(SUM(final_price),0) FROM orders WHERE payment_status='paid' AND approval = true AND created_at >= ? AND created_at <= ?", startTime, endTime).Scan(&salesReport.TotalSales)
	if result.Error != nil {
		return models.SalesReport{}, result.Error
	}
	result = ad.DB.Raw("SELECT COUNT(*) FROM orders").Scan(&salesReport.TotalOrders)
	if result.Error != nil {
		return models.SalesReport{}, result.Error
	}
	result = ad.DB.Raw("SELECT COUNT(*) FROM orders WHERE payment_status = 'paid' and approval = true and created_at >= ? AND created_at <= ?", startTime, endTime).Scan(&salesReport.CompletedOrders)
	if result.Error != nil {
		return models.SalesReport{}, result.Error
	}
	result = ad.DB.Raw("SELECT COUNT(*) FROM orders WHERE shipment_status = 'processing' AND approval = false AND created_at >= ? AND created_at<=?", startTime, endTime).Scan(&salesReport.PendingOrders)
	if result.Error != nil {
		return models.SalesReport{}, result.Error
	}
	var productID int
	result = ad.DB.Raw("SELECT product_id FROM order_items GROUP BY product_id order by SUM(quantity) DESC LIMIT 1").Scan(&productID)
	if result.Error != nil {
		return models.SalesReport{}, result.Error
	}
	result = ad.DB.Raw("SELECT name FROM products WHERE id = ?", productID).Scan(&salesReport.TrendingProduct)
	if result.Error != nil {
		return models.SalesReport{}, result.Error
	}
	return salesReport, nil
}
func (ad *adminRepository) AddPaymentMethod(pay models.NewPaymentMethod) (domain.PaymentMethod, error) {
	var payment string
	if err := ad.DB.Raw("INSERT INTO payment_methods (payment_name) VALUES (?) RETURNING payment_name", pay.PaymentName).Scan(&payment).Error; err != nil {
		return domain.PaymentMethod{}, err
	}
	var paymentResponse domain.PaymentMethod
	err := ad.DB.Raw("SELECT id, payment_name FROM payment_methods WHERE payment_name = ?", payment).Scan(&paymentResponse).Error
	if err != nil {
		return domain.PaymentMethod{}, err
	}
	return paymentResponse, nil

}

func (ad *adminRepository) ListPaymentMethods() ([]domain.PaymentMethod, error) {
	var model []domain.PaymentMethod
	err := ad.DB.Raw("SELECT * FROM payment_methods").Scan(&model).Error
	if err != nil {
		return []domain.PaymentMethod{}, err
	}

	return model, nil
}

func (ad *adminRepository) CheckIfPaymentMethodAlreadyExists(payment string) (bool, error) {
	var count int64
	err := ad.DB.Raw("SELECT COUNT(*) FROM payment_methods WHERE payment_name = $1", payment).Scan(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (ad *adminRepository) DeletePaymentMethod(id int) error {
	var count int
	if err := ad.DB.Raw("SELECT COUNT(*) FROM payment_methods WHERE id=?", id).Scan(&count).Error; err != nil {
		return err
	}
	if count < 1 {
		return errors.New("payment for given id does not exist")
	}

	if err := ad.DB.Exec("DELETE FROM payment_methods WHERE id=?", id).Error; err != nil {
		return err
	}
	return nil
}
