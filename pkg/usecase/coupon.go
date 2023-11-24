package usecase

import (
	interfaces "Zhooze/pkg/repository/interface"
	services "Zhooze/pkg/usecase/interface"
	"Zhooze/pkg/utils/models"
	"errors"
)

type couponUseCase struct {
	couponRepository interfaces.CouponRepository
	cartRepository interfaces.CartRepository
	orderRepository interfaces.OrderRepository
}

func NewCouponUseCase(repository interfaces.CouponRepository,cartRepo interfaces.CartRepository,orderRepo interfaces.OrderRepository) services.CouponUseCase {
	return &couponUseCase{
		couponRepository: repository,
		cartRepository:cartRepo,
		orderRepository: orderRepo,

		
	}
}
func (co *couponUseCase) AddCoupon(coupon models.AddCoupon) (string, error) {

	// if coupon already exist and if it is expired revalidate it. else give back an error message saying the coupon already exist
	couponExist, err := co.couponRepository.CouponExist(coupon.Coupon)
	if err != nil {
		return "", err
	}

	if couponExist {
		alreadyValid, err := co.couponRepository.CouponRevalidateIfExpired(coupon.Coupon)
		if err != nil {
			return "", err
		}

		if alreadyValid {
			return "The coupon which is valid already exists", nil
		}

		return "Made the coupon valid", nil

	}

	err = co.couponRepository.AddCoupon(coupon)
	if err != nil {
		return "", err
	}
	return "successfully added the coupon", nil
}
func (co *couponUseCase) GetCoupon() ([]models.Coupon, error) {
	coupons, err := co.couponRepository.GetCoupon()
	if err != nil {
		return []models.Coupon{}, err
	}
	return coupons, nil
}
func (co *couponUseCase) ExpireCoupon(couponID int) error {
	couponExist, err := co.couponRepository.ExistCoupon(couponID)
	if err != nil {
		return err
	}
	// if it exists expire it, if already expired send back relevant message
	if couponExist {
		err = co.couponRepository.CouponAlreadyExpired(couponID)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("coupon does not exist")

}

func (co *couponUseCase) ApplyCoupon(coupon string, userID int) error {

	cartExist, err := co.orderRepository.DoesCartExist(userID)
	if err != nil {
		return err
	}

	if !cartExist {
		return errors.New("cart empty, can't apply coupon")
	}

	couponExist, err := co.couponRepository.CouponExist(coupon)
	if err != nil {
		return err
	}

	if !couponExist {
		return errors.New("coupon does not exist")
	}

	couponValidity, err :=co. couponRepository.CouponValidity(coupon)
	if err != nil {
		return err
	}

	if !couponValidity {
		return errors.New("coupon expired")
	}

	minDiscountPrice, err := co.couponRepository.GetCouponMinimumAmount(coupon)
	if err != nil {
		return err
	}

	totalPriceFromCarts, err :=co. cartRepository.GetTotalPriceFromCart(userID)
	if err != nil {
		return err
	}

	// if the total Price is less than minDiscount price don't allow coupon to be added
	if totalPriceFromCarts < minDiscountPrice {
		return errors.New("coupon cannot be added as the total amount is less than minimum amount for coupon")
	}

	userAlreadyUsed, err :=co. couponRepository.DidUserAlreadyUsedThisCoupon(coupon, userID)
	if err != nil {
		return err
	}

	if userAlreadyUsed {
		return errors.New("user already used this coupon")
	}

	couponStatus, err :=co. couponRepository.UpdateUsedCoupon(coupon, userID)
	if err != nil {
		return err
	}

	if couponStatus {
		return nil
	}
	return errors.New("could not add the coupon")

}
