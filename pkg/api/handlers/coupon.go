package handlers

import (
	services"Zhooze/pkg/usecase"
	"Zhooze/pkg/utils/models"
	"Zhooze/pkg/utils/response"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CouponHandler struct {
	couponUseCase services.CouponUseCase
}

func NewCouponHandler(useCase services.CouponUseCase) *CouponHandler {
	return &CouponHandler{
		couponUseCase: useCase,
	}
}

// @Summary Add  a new coupon by Admin
// @Description Add A new Coupon which can be used by the users from the checkout section
// @Tags Admin Coupon Management
// @Accept json
// @Produce json
// @Security Bearer
// @Param coupon body models.AddCoupon true "Add new Coupon"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /admin/coupons [POST]
func AddCoupon(c *gin.Context) {

	var coupon models.AddCoupon

	if err := c.ShouldBindJSON(&coupon); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not bind the coupon details", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	err := validator.New().Struct(coupon)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "constraints not satisfied", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	message, err := usecase.AddCoupon(coupon)

	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "Could not add coupon", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusCreated, "Coupon Added", message, nil)
	c.JSON(http.StatusCreated, successRes)

}

// @Summary Get coupon details
// @Description Get Available coupon details for admin side
// @Tags Admin Coupon Management
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /admin/coupons [get]
func GetCoupon(c *gin.Context) {

	coupons, err := usecase.GetCoupon()

	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "Could not get coupon details", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Coupon Retrieved successfully", coupons, nil)
	c.JSON(http.StatusOK, successRes)

}

// @Summary Expire Coupon
// @Description Expire Coupon by admin which are already present by passing coupon id
// @Tags Admin Coupon Management
// @Accept json
// @Produce json
// @Security Bearer
// @Param coupon_id query string true "Coupon id"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /admin/coupons   [PATCH]
func ExpireCoupon(c *gin.Context) {

	id := c.Query("coupon_id")
	couponID, err := strconv.Atoi(id)

	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "coupon id not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	err = usecase.ExpireCoupon(couponID)
	if err != nil {
		if errors.Is(err, errors.New("the offer already exists")) {
			errorRes := response.ClientResponse(http.StatusForbidden, "could not expire coupon", nil, err.Error())
			c.JSON(http.StatusForbidden, errorRes)
			return
		}
		errorRes := response.ClientResponse(http.StatusInternalServerError, "could not expire coupon", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Coupon expired successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)

}

// @Summary Apply coupon on Checkout Section
// @Description Add coupon to get discount on Checkout section
// @Tags User Checkout
// @Accept json
// @Produce json
// @Security Bearer
// @Param couponDetails body models.CouponAddUser true "Add coupon to order"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /user/coupon/apply [POST]
func ApplyCoupon(c *gin.Context) {

	userID, _ := c.Get("user_id")
	var couponDetails models.CouponAddUser

	if err := c.ShouldBindJSON(&couponDetails); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not bind the coupon", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	err := usecase.ApplyCoupon(couponDetails.CouponName, userID.(int))

	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "coupon could not be added", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusCreated, "Coupon added successfully", nil, nil)
	c.JSON(http.StatusCreated, successRes)

}
