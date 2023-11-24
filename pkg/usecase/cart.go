package usecase

import (
	interfaces "Zhooze/pkg/repository/interface"
	services "Zhooze/pkg/usecase/interface"
	"Zhooze/pkg/utils/models"
	"errors"
)

type cartUseCase struct {
	cartRepository    interfaces.CartRepository
	couponRepository  interfaces.CouponRepository
	productRepository interfaces.ProductRepository
	offerRepository   interfaces.OfferRepository
	orderRepository   interfaces.OrderRepository
}

func NewCartUseCase(repository interfaces.CartRepository, couponRepo interfaces.CouponRepository, productRepo interfaces.ProductRepository, offerRepo interfaces.OfferRepository, orderRepo interfaces.OrderRepository) services.CartUseCase {

	return &cartUseCase{
		cartRepository:    repository,
		couponRepository:  couponRepo,
		productRepository: productRepo,
		offerRepository:   offerRepo,
		orderRepository:   orderRepo,
	}

}
func (cr *cartUseCase) AddToCart(product_id int, user_id int) (models.CartResponse, error) {
	ok, _, err := cr.cartRepository.CheckProduct(product_id)
	if err != nil {
		return models.CartResponse{}, err
	}
	if !ok {
		return models.CartResponse{}, errors.New("product Does not exist")
	}
	QuantityOfProductInCart, err := cr.cartRepository.QuantityOfProductInCart(user_id, product_id)
	if err != nil {

		return models.CartResponse{}, err
	}
	quantityOfProduct, err := cr.productRepository.GetQuantityFromProductID(product_id)
	if err != nil {

		return models.CartResponse{}, err
	}
	if quantityOfProduct <= 0 {
		return models.CartResponse{}, errors.New("out of stock")
	}
	if quantityOfProduct == QuantityOfProductInCart {
		return models.CartResponse{}, errors.New("stock limit exceeded")
	}
	productPrice, err := cr.productRepository.GetPriceOfProductFromID(product_id)
	if err != nil {

		return models.CartResponse{}, err
	}
	/////////////
	discount_percentage, err := cr.offerRepository.FindDiscountPercentageForProduct(product_id)
	if err != nil {
		return models.CartResponse{}, errors.New("there was some error in finding the discounted prices")
	}
	var discount float64

	if discount_percentage > 0 {
		discount = (productPrice * float64(discount_percentage)) / 100
	}

	Price := productPrice - discount
	categoryID, err := cr.productRepository.FindCategoryID(product_id)
	if err != nil {
		return models.CartResponse{}, err
	}
	discount_percentageCategory, err := cr.offerRepository.FindDiscountPercentageForCategory(categoryID)
	if err != nil {
		return models.CartResponse{}, errors.New("there was some error in finding the discounted prices")
	}
	var discountcategory float64

	if discount_percentageCategory > 0 {
		discountcategory = (productPrice * float64(discount_percentageCategory)) / 100
	}

	FinalPrice := Price - discountcategory
	//////////////////////
	if QuantityOfProductInCart == 0 {
		err := cr.cartRepository.AddItemIntoCart(user_id, product_id, 1, FinalPrice)
		if err != nil {

			return models.CartResponse{}, err
		}

	} else {
		currentTotal, err := cr.cartRepository.TotalPriceForProductInCart(user_id, product_id)
		if err != nil {
			return models.CartResponse{}, err
		}
		err = cr.cartRepository.UpdateCart(QuantityOfProductInCart+1, currentTotal+productPrice, user_id, product_id)
		if err != nil {
			return models.CartResponse{}, err
		}
	}
	cartDetails, err := cr.cartRepository.DisplayCart(user_id)
	if err != nil {
		return models.CartResponse{}, err
	}
	cartTotal, err := cr.cartRepository.GetTotalPrice(user_id)
	if err != nil {

		return models.CartResponse{}, err
	}
	err = cr.orderRepository.ProductStockMinus(product_id, QuantityOfProductInCart)
	if err != nil {
		return models.CartResponse{}, err
	}
	return models.CartResponse{
		UserName:   cartTotal.UserName,
		TotalPrice: cartTotal.TotalPrice,
		Cart:       cartDetails,
	}, nil

}

func (cr *cartUseCase) RemoveFromCart(product_id, user_id int) (models.CartResponse, error) {
	ok, err := cr.cartRepository.ProductExist(user_id, product_id)
	if err != nil {
		return models.CartResponse{}, err
	}
	if !ok {
		return models.CartResponse{}, errors.New("product doesn't exist in the cart")
	}
	var cartDetails struct {
		Quantity   int
		TotalPrice float64
	}

	cartDetails, err = cr.cartRepository.GetQuantityAndProductDetails(user_id, product_id, cartDetails)
	if err != nil {
		return models.CartResponse{}, err
	}
	if err := cr.cartRepository.RemoveProductFromCart(user_id, product_id); err != nil {
		return models.CartResponse{}, err
	}

	if cartDetails.Quantity != 0 {

		product_price, err := cr.productRepository.GetPriceOfProductFromID(product_id)
		if err != nil {
			return models.CartResponse{}, err
		}
		cartDetails.TotalPrice = cartDetails.TotalPrice - product_price
		err = cr.cartRepository.UpdateCartDetails(cartDetails, user_id, product_id)
		if err != nil {
			return models.CartResponse{}, err
		}
	}
	updatedCart, err := cr.cartRepository.CartAfterRemovalOfProduct(user_id)
	if err != nil {
		return models.CartResponse{}, err
	}
	cartTotal, err := cr.cartRepository.GetTotalPrice(user_id)
	if err != nil {
		return models.CartResponse{}, err
	}
	return models.CartResponse{
		UserName:   cartTotal.UserName,
		TotalPrice: cartTotal.TotalPrice,
		Cart:       updatedCart,
	}, nil

}
func (cr *cartUseCase) DisplayCart(user_id int) (models.CartResponse, error) {
	cart, err := cr.cartRepository.DisplayCart(user_id)
	if err != nil {
		return models.CartResponse{}, err
	}
	cartTotal, err := cr.cartRepository.GetTotalPrice(user_id)
	if err != nil {
		return models.CartResponse{}, err
	}
	return models.CartResponse{
		UserName:   cartTotal.UserName,
		TotalPrice: cartTotal.TotalPrice,
		Cart:       cart,
	}, nil
}
func (cr *cartUseCase) EmptyCart(userID int) (models.CartResponse, error) {
	ok, err := cr.cartRepository.CartExist(userID)
	if err != nil {
		return models.CartResponse{}, err
	}
	if !ok {
		return models.CartResponse{}, errors.New("cart already empty")
	}
	if err := cr.cartRepository.EmptyCart(userID); err != nil {
		return models.CartResponse{}, err
	}

	cartTotal, err := cr.cartRepository.GetTotalPrice(userID)

	if err != nil {
		return models.CartResponse{}, err
	}

	cartResponse := models.CartResponse{
		UserName:   cartTotal.UserName,
		TotalPrice: cartTotal.TotalPrice,
		Cart:       []models.Cart{},
	}

	return cartResponse, nil

}
