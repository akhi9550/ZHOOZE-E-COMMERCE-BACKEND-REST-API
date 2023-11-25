package handlers

import (
	services "Zhooze/pkg/usecase/interface"
	"Zhooze/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WishListHandler struct {
	WishListUsecase services.WishListUseCase
}

func NewWishListHandler(useCase services.WishListUseCase) *WishListHandler {
	return &WishListHandler{
		WishListUsecase: useCase,
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
func (wl *WishListHandler) GetWishList(c *gin.Context) {
	userID, _ := c.Get("user_id")
	wishList, err := wl.WishListUsecase.GetWishList(userID.(int))
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
func (wl *WishListHandler) AddToWishlist(c *gin.Context) {
	userID, _ := c.Get("user_id")
	product_id := c.Query("product_id")
	productID, err := strconv.Atoi(product_id)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "product id is in wrong format", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	err = wl.WishListUsecase.AddToWishlist(userID.(int), productID)
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
func (wl *WishListHandler) RemoveFromWishlist(c *gin.Context) {
	userID, _ := c.Get("user_id")
	product_id := c.Query("id")
	productID, err := strconv.Atoi(product_id)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "product id is in wrong format", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	err = wl.WishListUsecase.RemoveFromWishlist(productID, userID.(int))
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "failed to remove item from wishlist", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "SuccessFully remove product from  wishlist", nil, nil)
	c.JSON(http.StatusOK, success)
}
