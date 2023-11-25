package handlers

import (
	services "Zhooze/pkg/usecase/interface"
	"Zhooze/pkg/utils/models"
	"Zhooze/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderUseCase services.OrderUseCase
}

func NewOrderHandler(useCase services.OrderUseCase) *OrderHandler {
	return &OrderHandler{
		orderUseCase: useCase,
	}
}

// @Summary Approve Order
// @Description Approve Order from admin side which is in processing state
// @Tags Admin Order Management
// @Accept   json
// @Produce  json
// @Security Bearer
// @Param    order_id   query   string   true    "Order ID"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /admin/order/approve [GET]
func (or *OrderHandler) ApproveOrder(c *gin.Context) {
	orderId, err := strconv.Atoi(c.Query("order_id"))
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "error from orderID", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	err = or.orderUseCase.ApproveOrder(orderId)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "Couldn't approve the order", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Order Approved Successfully", nil, nil)
	c.JSON(http.StatusOK, success)
}

// @Summary Cancel Order Admin
// @Description Cancel Order from admin side
// @Tags Admin Order Management
// @Accept   json
// @Produce  json
// @Security Bearer
// @Param order_id query string true "Order ID"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /admin/order/cancel   [GET]
func (or *OrderHandler) CancelOrderFromAdmin(c *gin.Context) {
	order_id, err := strconv.Atoi(c.Query("order_id"))
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "error from orderID", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	err = or.orderUseCase.CancelOrderFromAdmin(order_id)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "Couldn't cancel the order", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Order Cancel Successfully", nil, nil)
	c.JSON(http.StatusOK, success)
}

// @Summary Get All order details for admin
// @Description Get all order details to the admin side
// @Tags Admin Order Management
// @Accept   json
// @Produce  json
// @Security Bearer
// @Param page query string false "Page number"
// @Param count query string false "Page Count"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /admin/order   [GET]
func (or *OrderHandler) GetAllOrderDetailsForAdmin(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "page number not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	countStr := c.DefaultQuery("count", "20")
	pageSize, err := strconv.Atoi(countStr)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "page count not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	allOrderDetails, err := or.orderUseCase.GetAllOrderDetailsForAdmin(page, pageSize)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "Could not retrieve order details", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Order Details Retrieved successfully", allOrderDetails, nil)
	c.JSON(http.StatusOK, success)
}

// @Summary Order Items from cart
// @Description Order all products which is currently present inside  the cart
// @Tags User Order Management
// @Accept json
// @Produce json
// @Security Bearer
// @Param orderBody body models.OrderFromCart true "Order details"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /user/order [POST]
func (or *OrderHandler) OrderItemsFromCart(c *gin.Context) {
	id, _ := c.Get("user_id")
	userID := id.(int)
	var orderFromCart models.OrderFromCart
	if err := c.ShouldBindJSON(&orderFromCart); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "bad request", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	orderSuccessResponse, err := or.orderUseCase.OrderItemsFromCart(orderFromCart, userID)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "Could not do the order", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully created the order", orderSuccessResponse, nil)
	c.JSON(http.StatusOK, successRes)

}

// @Summary Get Order Details to user side
// @Description Get all order details done by user to user side
// @Tags User Order Management
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query string false "Page"
// @Param count query string false "Count"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /user/order   [GET]
func (or *OrderHandler) GetOrderDetails(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "page number not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("count", "10"))
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "page count not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	id, _ := c.Get("user_id")
	UserID := id.(int)
	OrderDetails, err := or.orderUseCase.GetOrderDetails(UserID, page, pageSize)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "Could not do the order", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Full Order Details", OrderDetails, nil)
	c.JSON(http.StatusOK, successRes)
}

// @Summary Cancel order
// @Description Cancel order by the user using order ID
// @Tags User Order Management
// @Accept json
// @Produce json
// @Security Bearer
// @Param id query string true "Order ID"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /user/order   [PUT]
func (or *OrderHandler) CancelOrder(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "error from orderID", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	id, _ := c.Get("user_id")
	userID := id.(int)
	err = or.orderUseCase.CancelOrders(orderID, userID)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "Could not cancel the order", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Cancel Successfull", nil, nil)
	c.JSON(http.StatusOK, success)
}

// @Summary		Checkout section
// @Description	Add products to carts  for the purchase
// @Tags			User Order Management
// @Accept			json
// @Produce		    json
// @Security		Bearer
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/user/order/checkout    [GET]
func (or *OrderHandler) CheckOut(c *gin.Context) {
	userID, _ := c.Get("user_id")
	checkoutDetails, err := or.orderUseCase.Checkout(userID.(int))

	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "failed to retrieve details", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Checkout Page loaded successfully", checkoutDetails, nil)
	c.JSON(http.StatusOK, successRes)
}

// @Summary		Checkout section
// @Description	Add products to carts  for the purchase
// @Tags			User Order Management
// @Accept			json
// @Produce		    json
// @Param    order_id    query    int    true    "address id"
// @Security		Bearer
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/user/order/place-order     [GET]
func (or *OrderHandler) PlaceOrderCOD(c *gin.Context) {
	order_id, err := strconv.Atoi(c.Query("order_id"))
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "error from orderID", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	paymentMethodID, err := or.orderUseCase.PaymentMethodID(order_id)
	if err != nil {
		err := response.ClientResponse(http.StatusInternalServerError, "error from paymentId ", nil, err.Error())
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if paymentMethodID == 1 {
		err := or.orderUseCase.ExecutePurchaseCOD(order_id)
		if err != nil {
			errorRes := response.ClientResponse(http.StatusInternalServerError, "error in cash on delivery ", nil, err.Error())
			c.JSON(http.StatusInternalServerError, errorRes)
			return
		}
		success := response.ClientResponse(http.StatusOK, "Placed Order with cash on delivery", nil, nil)
		c.JSON(http.StatusOK, success)
	}
	if paymentMethodID == 2 {
		success := response.ClientResponse(http.StatusOK, "Placed Order with razor pay", nil, nil)
		c.JSON(http.StatusOK, success)
	}
}
