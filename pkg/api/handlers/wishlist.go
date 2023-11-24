package handlers

import (
	services"Zhooze/pkg/usecase"
	"Zhooze/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ImageHandler struct {
	ImageUseCase services.ImageUseCase
}

func NewCouponHandler(useCase services.ImageUseCase) *ImageHandler {
	return &ImageHandler{
		ImageUseCase: useCase,
	}
}

// @Summary Display Wishlist
// @Description Display wish List
// @Tags WishList Management
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /user/wishlist [GET]
func GetWishList(c *gin.Context) {
	userID, _ := c.Get("user_id")
	wishList, err := usecase.GetWishList(userID.(int))
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "failed to retrieve wishlist detailss", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "SuccessFully retrieved wishlist", wishList, nil)
	c.JSON(http.StatusOK, success)
}

// @Summary Add to Wishlist
// @Description Add To wish List
// @Tags WishList Management
// @Accept json
// @Produce json
// @Security Bearer
// @Param product_id query string true "product id"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /user/wishlist [POST]
func AddToWishlist(c *gin.Context) {
	userID, _ := c.Get("user_id")
	product_id := c.Query("product_id")
	productID, err := strconv.Atoi(product_id)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "product id is in wrong format", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	err = usecase.AddToWishlist(userID.(int), productID)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "failed to item to the wishlist", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "SuccessFully added product to the wishlist", nil, nil)
	c.JSON(http.StatusOK, success)
}

// @Summary Delete From Wishlist
// @Description Delete From wish List
// @Tags WishList Management
// @Accept json
// @Produce json
// @Security Bearer
// @Param id query string true "product id"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /user/wishlist [DELETE]
func RemoveFromWishlist(c *gin.Context) {
	userID, _ := c.Get("user_id")
	product_id := c.Query("id")
	productID, err := strconv.Atoi(product_id)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "product id is in wrong format", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	err = usecase.RemoveFromWishlist(productID, userID.(int))
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "failed to remove item from wishlist", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "SuccessFully remove product from  wishlist", nil, nil)
	c.JSON(http.StatusOK, success)
}
