package interfaces

import "Zhooze/pkg/utils/models"

type CouponUseCase interface {
	AddCoupon(coupon models.AddCoupon) (string, error)
	GetCoupon() ([]models.Coupon, error)
	ExpireCoupon(couponID int) error
	ApplyCoupon(coupon string, userID int) error
}
