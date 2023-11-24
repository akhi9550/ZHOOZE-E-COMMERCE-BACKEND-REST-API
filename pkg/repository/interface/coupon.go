package interfaces

import "Zhooze/pkg/utils/models"

type CouponRepository interface {
	CouponExist(couponName string) (bool, error)
	CouponValidity(couponName string) (bool, error)
	CouponRevalidateIfExpired(couponName string) (bool, error)
	AddCoupon(coupon models.AddCoupon) error
	GetCoupon() ([]models.Coupon, error)
	ExistCoupon(couponID int) (bool, error)
	CouponAlreadyExpired(couponID int) error
	GetCouponMinimumAmount(coupon string) (float64, error)
	DidUserAlreadyUsedThisCoupon(coupon string, userID int) (bool, error)
	UpdateUsedCoupon(coupon string, userID int) (bool, error)
}
