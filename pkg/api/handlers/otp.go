package handlers

import (
	services "Zhooze/pkg/usecase/interface"
	"Zhooze/pkg/utils/models"
	"Zhooze/pkg/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OtpHandler struct {
	OtpUseCase services.OtpUseCase
}

func NewOtpHandler(useCase services.OtpUseCase) *OtpHandler {
	return &OtpHandler{
		OtpUseCase: useCase,
	}
}

// @Summary  OTP login
// @Description Send OTP to Authenticate user
// @Tags User OTP Login
// @Accept json
// @Produce json
// @Param phone body models.OTPData true "phone number details"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /user/send-otp   [POST]
func (op *OtpHandler) SendOtp(c *gin.Context) {
	var phone models.OTPData
	if err := c.ShouldBindJSON(&phone); err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	err := op.OtpUseCase.SendOtp(phone.PhoneNumber)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "OTP sent successfully", nil, nil)
	c.JSON(http.StatusOK, success)
}

// @Summary Verify OTP
// @Description Verify OTP by passing the OTP in order to authenticate user
// @Tags User OTP Login
// @Accept json
// @Produce json
// @Param phone body models.VerifyData true "Verify OTP Details"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /user/verify-otp      [POST]
func (op *OtpHandler) VerifyOtp(c *gin.Context) {
	var code models.VerifyData
	if err := c.ShouldBindJSON(&code); err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	user, err := op.OtpUseCase.VerifyOTP(code)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "Could not verify OTP", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	sucess := response.ClientResponse(http.StatusOK, "Successfully verified OTP", user, nil)
	c.JSON(http.StatusOK, sucess)
}
