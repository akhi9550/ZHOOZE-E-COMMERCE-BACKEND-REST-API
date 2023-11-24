package handlers

import (
	"Zhooze/pkg/usecase"
	services "Zhooze/pkg/usecase"
	"Zhooze/pkg/utils/models"
	"Zhooze/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ImageHandler struct {
	ImageUseCase services.ImageUseCase
}

func NewCouponHandler(useCase services.ImageUseCase) *ImageHandler {
	return &ImageHandler{
		ImageUseCase: useCase,
	}
}

// @Summary Add  Product Offer
// @Description Add a new Offer for a product by specifying a limit
// @Tags Admin Offer Management
// @Accept json
// @Produce json
// @Security Bearer
// @Param coupon body models.ProductOfferReceiver true "Add new Product Offer"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /admin/offer/product-offer [POST]
func AddProdcutOffer(c *gin.Context) {

	var productOffer models.ProductOfferReceiver

	if err := c.ShouldBindJSON(&productOffer); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "request fields in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	err := validator.New().Struct(productOffer)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "constraints not satisfied", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	err = usecase.AddProductOffer(productOffer)

	if err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "could not add offer", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusCreated, "Successfully added offer", nil, nil)
	c.JSON(http.StatusCreated, successRes)

}

// @Summary Add  Product Offer
// @Description Add a new Offer for a product by specifying a limit
// @Tags Admin Offer Management
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /admin/offer/product-offer [GET]
func GetProductOffer(c *gin.Context) {

	categories, err := usecase.GetOffers()
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully got all offers", categories, nil)
	c.JSON(http.StatusOK, successRes)

}

// @Summary Add  Product Offer
// @Description Add a new Offer for a product by specifying a limit
// @Tags Admin Offer Management
// @Accept json
// @Produce json
// @Security Bearer
// @Param	id	query	string	true	"id"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /admin/offer/product-offer   [DELETE]
func ExpireProductOffer(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	if err := usecase.MakeOfferExpire(id); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Coupon cannot be made invalid", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully made Coupon as invalid", nil, nil)
	c.JSON(http.StatusOK, successRes)

}

// @Summary Add Category Offer
// @Description Add a new Offer for a product by specifying a limit
// @Tags Admin Offer Management
// @Accept json
// @Produce json
// @Security Bearer
// @Param coupon body models.CategoryOfferReceiver  true "Add new category Offer"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /admin/offer/category-offer [POST]
func AddCategoryOffer(c *gin.Context) {

	var categoryOffer models.CategoryOfferReceiver

	if err := c.ShouldBindJSON(&categoryOffer); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "request fields in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	err := validator.New().Struct(categoryOffer)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "constraints not satisfied", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	err = usecase.AddCategoryOffer(categoryOffer)

	if err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "could not add offer", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusCreated, "Successfully added offer", nil, nil)
	c.JSON(http.StatusCreated, successRes)

}

// @Summary Add  Category Offer
// @Description Add a new Offer for a category by specifying a limit
// @Tags Admin Offer Management
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /admin/offer/category-offer [GET]
func GetCategoryOffer(c *gin.Context) {

	categories, err := usecase.GetCategoryOffer()
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully got all offers", categories, nil)
	c.JSON(http.StatusOK, successRes)

}

// @Summary Add  Category Offer
// @Description Add a new Offer for a category by specifying a limit
// @Tags Admin Offer Management
// @Accept json
// @Produce json
// @Security Bearer
// @Param	id	query	string	true	"id"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /admin/offer/category-offer   [DELETE]
func ExpireCategoryOffer(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	if err := usecase.ExpireCategoryOffer(id); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Coupon cannot be made invalid", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully made Coupon as invaid", nil, nil)
	c.JSON(http.StatusOK, successRes)

}
