package handlers

import (
	services "Zhooze/pkg/usecase/interface"
	"Zhooze/pkg/utils/models"
	"Zhooze/pkg/utils/response"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminUseCase services.AdminUseCase
}

func NewAdminHandler(usecase services.AdminUseCase) *AdminHandler {
	return &AdminHandler{
		adminUseCase: usecase,
	}
}

// @Summary		Admin Login
// @Description	Login handler for Zhooze admins
// @Tags			Admin
// @Accept			json
// @Produce		json
// @Param			admin	body		models.AdminLogin	true	"Admin login details"
// @Success		200		{object}	response.Response{}
// @Failure		500		{object}	response.Response{}
// @Router			/admin/adminlogin [POST]
func (ad *AdminHandler) LoginHandler(c *gin.Context) {
	var adminDetails models.AdminLogin
	if err := c.ShouldBindJSON(&adminDetails); err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	admin, err := ad.adminUseCase.LoginHandler(adminDetails)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "Cannot authenticate user", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Admin authenticated successfully", admin, nil)
	c.JSON(http.StatusOK, success)
}

// @Summary		Admin Dashboard
// @Description	Retrieve admin dashboard
// @Tags			Admin
// @Accept			json
// @Produce		json
// @Security		Bearer
// @Success		200		{object}	response.Response{}
// @Failure		500		{object}	response.Response{}
// @Router			/admin/dashboard [GET]
func (ad *AdminHandler) DashBoard(c *gin.Context) {
	adminDashboard, err := ad.adminUseCase.DashBoard()
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "Dashboard could not be displayed", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Admin dashboard displayed", adminDashboard, nil)
	c.JSON(http.StatusOK, success)
}

// @Summary		Get Users
// @Description	Retrieve users with pagination
// @Tags			Admin User Management
// @Accept			json
// @Produce		    json
// @Security		Bearer
// @Param page query string false "Page number"
// @Param count query string false "Page size"
// @Success		200		{object}	response.Response{}
// @Failure		500		{object}	response.Response{}
// @Router			/admin/users   [GET]
func (ad *AdminHandler) GetUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "page number not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	countStr := c.DefaultQuery("count", "10")
	pageSize, err := strconv.Atoi(countStr)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "user count in a page not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	users, err := ad.adminUseCase.ShowAllUsers(page, pageSize)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "couldn't retrieve users", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully Retrieved all Users", users, nil)
	c.JSON(http.StatusOK, success)
}

// @Summary		Block User
// @Description	using this handler admins can block an user
// @Tags			Admin User Management
// @Accept			json
// @Produce		json
// @Security		Bearer
// @Param			id	query		string	true	"user-id"
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/admin/users/block   [PUT]
func (ad *AdminHandler) BlockUser(c *gin.Context) {
	id := c.Query("id")
	err := ad.adminUseCase.BlockedUser(id)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "user could not be blocked", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	sucess := response.ClientResponse(http.StatusOK, "Successfully blocked the user", nil, nil)
	c.JSON(http.StatusOK, sucess)

}

// @Summary		UnBlock an existing user
// @Description	UnBlock user
// @Tags			Admin User Management
// @Accept			json
// @Produce		    json
// @Security		Bearer
// @Param			id	query		string	true	"user-id"
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/admin/users/unblock    [PUT]
func (ad *AdminHandler) UnBlockUser(c *gin.Context) {
	id := c.Query("id")
	err := ad.adminUseCase.UnBlockedUser(id)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "user could not be unblocked", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	sucess := response.ClientResponse(http.StatusOK, "Successfully unblocked the user", nil, nil)
	c.JSON(http.StatusOK, sucess)
}

// @Summary Filtered Sales Report
// @Description Get Filtered sales report by week, month and year
// @Tags Admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param period query string true "sales report"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /admin/sales-report    [GET]
func (ad *AdminHandler) FilteredSalesReport(c *gin.Context) {
	timePeriod := c.Query("period")
	salesReport, err := ad.adminUseCase.FilteredSalesReport(timePeriod)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "sales report could not be retrieved", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}
	success := response.ClientResponse(http.StatusOK, "sales report retrieved successfully", salesReport, nil)
	c.JSON(http.StatusOK, success)

}

//	@Summary		Sales report by date
//	@Description	Showing the sales report with respect to the given date
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//
// @Security        Bearer
//
//	@Param			start	query	string		true	"start date DD-MM-YYYY"
//	@Param			end		query	string		true	"end   date DD-MM-YYYY"
//	@Success		200		body	entity.SalesReport	"report"
//	@Router			/admin/sales-report-date   [GET]
func (ad *AdminHandler) SalesReportByDate(c *gin.Context) {
	startDateStr := c.Query("start")
	endDateStr := c.Query("end")
	if startDateStr == "" || endDateStr == "" {
		err := response.ClientResponse(http.StatusBadRequest, "start or end date is empty", nil, "Empty date string")
		c.JSON(http.StatusBadRequest, err)
		return
	}
	startDate, err := time.Parse("2-1-2006", startDateStr)
	if err != nil {
		err := response.ClientResponse(http.StatusBadRequest, "start date conversion failed", nil, err.Error())
		c.JSON(http.StatusBadRequest, err)
		return
	}
	endDate, err := time.Parse("2-1-2006", endDateStr)
	if err != nil {
		err := response.ClientResponse(http.StatusBadRequest, "end date conversion failed", nil, err.Error())
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if startDate.After(endDate) {
		err := response.ClientResponse(http.StatusBadRequest, "start date is after end date", nil, "Invalid date range")
		c.JSON(http.StatusBadRequest, err)
		return
	}

	report, err := ad.adminUseCase.ExecuteSalesReportByDate(startDate, endDate)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "sales report could not be retrieved", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	success := response.ClientResponse(http.StatusOK, "sales report retrieved successfully", report, nil)
	c.JSON(http.StatusOK, success)
}

// @Summary		Add Payment Method
// @Description	Admin can add new payment methods
// @Tags			Admin Payment Method
// @Accept			json
// @Produce		    json
// @Security		Bearer
// @Param			payment 	body		models.NewPaymentMethod	 true	"payment method"
// @Success		200		{object}	response.Response{}
// @Failure		500		{object}	response.Response{}
// @Router			/admin/payment-method  [POST]
func (ad *AdminHandler) AddPaymentMethod(c *gin.Context) {
	var method models.NewPaymentMethod
	if err := c.ShouldBindJSON(&method); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	pay, err := ad.adminUseCase.AddPaymentMethod(method)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not add the payment method", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully added Payment Method", pay, nil)
	c.JSON(http.StatusOK, success)

}

// @Summary		Get Payment Method
// @Description	Admin can add new payment methods
// @Tags			Admin Payment Method
// @Accept			json
// @Produce		    json
// @Security		Bearer
// @Success		200		{object}	response.Response{}
// @Failure		500		{object}	response.Response{}
// @Router			/admin/payment-method  [GET]
func (ad *AdminHandler) ListPaymentMethods(c *gin.Context) {
	categories, err := ad.adminUseCase.ListPaymentMethods()
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	success := response.ClientResponse(http.StatusOK, "Successfully got all payment methods", categories, nil)
	c.JSON(http.StatusOK, success)

}

// @Summary		Delete Payment Method
// @Description	Admin can add new payment methods
// @Tags			Admin Payment Method
// @Accept			json
// @Produce		    json
// @Security		Bearer
// @Param			id	query	string	true	"id"
// @Success		200		{object}	response.Response{}
// @Failure		500		{object}	response.Response{}
// @Router			/admin/payment-method  [DELETE]
func (ad *AdminHandler) DeletePaymentMethod(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		error := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, error)
		return
	}

	err = ad.adminUseCase.DeletePaymentMethod(id)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "error in deleting data", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully Deleted the PaymentMethod", nil, nil)
	c.JSON(http.StatusOK, success)

}
