package handlers

import (
	services "Zhooze/pkg/usecase/interface"
	"Zhooze/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	cartUseCase services.CartUseCase
	userUseCase services.UserUseCase
}

func NewCartHandler(usecase services.CartUseCase, usecaseUser services.UserUseCase) *CartHandler {

	return &CartHandler{
		cartUseCase: usecase,
		userUseCase: usecaseUser,
	}

}

// @Summary		Add To Cart
// @Description	Add products to carts  for the purchase
// @Tags			User Cart Management
// @Accept			json
// @Produce		    json
// @Param			product_id	query		string	true	"product-id"
// @Security		Bearer
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/user/cart  [post]
func (ct *CartHandler) AddToCart(c *gin.Context) {
	id := c.Query("product_id")
	product_id, err := strconv.Atoi(id)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadGateway, "Product id is given in the wrong format", nil, err.Error())
		c.JSON(http.StatusBadGateway, errs)
		return
	}
	user_ID, _ := c.Get("user_id")
	cartResponse, err := ct.cartUseCase.AddToCart(product_id, user_ID.(int))
	if err != nil {
		errs := response.ClientResponse(http.StatusBadGateway, "could not add product to the cart", nil, err.Error())
		c.JSON(http.StatusBadGateway, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Added porduct Successfully to the cart", cartResponse, nil)
	c.JSON(http.StatusOK, success)
}

// @Summary		Remove From Cart
// @Description	Remove products to carts  for the purchase
// @Tags			User Cart Management
// @Accept			json
// @Produce		    json
// @Param			id	query		string	true	"product-id"
// @Security		Bearer
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/user/cart    [DELETE]
func (ct *CartHandler) RemoveFromCart(c *gin.Context) {
	id := c.Query("id")
	product_id, err := strconv.Atoi(id)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "product not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	userID, _ := c.Get("user_id")
	updatedCart, err := ct.cartUseCase.RemoveFromCart(product_id, userID.(int))
	if err != nil {
		errs := response.ClientResponse(http.StatusBadGateway, "cannot remove product from the cart", nil, err.Error())
		c.JSON(http.StatusBadGateway, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "product removed successfully", updatedCart, nil)
	c.JSON(http.StatusOK, success)
}

// @Summary		Display Cart
// @Description	Display products to carts
// @Tags			User Cart Management
// @Accept			json
// @Produce		    json
// @Security		Bearer
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/user/cart  [GET]
func (ct *CartHandler) DisplayCart(c *gin.Context) {
	userID, _ := c.Get("user_id")
	cart, err := ct.cartUseCase.DisplayCart(userID.(int))
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "cannot display cart", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Cart items displayed successfully", cart, nil)
	c.JSON(http.StatusOK, success)
}

// @Summary		Empty Cart
// @Description	Empty products to carts  for the purchase
// @Tags			User Cart Management
// @Accept			json
// @Produce		    json
// @Security		Bearer
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/user/cart/empty   [DELETE]
func (ct *CartHandler) EmptyCart(c *gin.Context) {
	userID, _ := c.Get("user_id")
	cart, err := ct.cartUseCase.EmptyCart(userID.(int))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "cannot empty the cart", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Cart emptied successfully", cart, nil)
	c.JSON(http.StatusOK, successRes)

}

// @Summary		Add quantity in cart by one
// @Description	user can add 1 quantity of product to their cart
// @Tags			User Cart Management
// @Accept			json
// @Produce		    json
// @Param			product_id	query	string	true	"product_id"
// @Security		Bearer
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/user/cart/updatequantityadd   [PUT]
func (ct *CartHandler) UpdateQuantityAdd(c *gin.Context) {
	id, _ := c.Get("user_id")
	productID, err := strconv.Atoi(c.Query("product_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "check product id parameter properly", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	if err := ct.userUseCase.UpdateQuantityAdd(id.(int), productID); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not Add the quantity", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully added quantity", nil, nil)
	c.JSON(http.StatusOK, success)
}

// @Summary		Subtract quantity in cart by one
// @Description	user can subtract 1 quantity of product from their cart
// @Tags			User Cart Management
// @Accept			json
// @Produce		    json
// @Param			product_id	query	string	true	"product_id"
// @Security		Bearer
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/user/cart/updatequantityless     [PUT]
func (ct *CartHandler) UpdateQuantityless(c *gin.Context) {
	id, _ := c.Get("user_id")
	productID, err := strconv.Atoi(c.Query("product_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "check parameters properly", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	if err := ct.userUseCase.UpdateQuantityless(id.(int), productID); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not Add the quantity", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	success := response.ClientResponse(http.StatusOK, "Successfully less quantity", nil, nil)
	c.JSON(http.StatusOK, success)
}
