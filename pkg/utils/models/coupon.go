package models

type Coupon struct {
	ID                 uint    `json:"id"`
	Coupon             string  `json:"coupon"`
	DiscountPercentage int     `json:"discount_percentage"`
	MinimumPrice       float64 `json:"minimum_price"`
	Validity           bool    `json:"validity"`
}

type AddCoupon struct {
	Coupon             string  `json:"coupon" binding:"required" validate:"required"`
	DiscountPercentage int     `json:"discount_percentage" binding:"required"`
	MinimumPrice       float64 `json:"minimum_price" binding:"required"`
	Validity           bool    `json:"validity" binding:"required"`
}

type CouponAddUser struct {
	CouponName string `json:"coupon_name" binding:"required"`
}

type ProductOfferReceiver struct {
	ProductID          uint   `json:"product_id" binding:"required"`
	OfferName          string `json:"offer_name" binding:"required"`
	DiscountPercentage int    `json:"discount_percentage" binding:"required"`
	// OfferLimit         int    `json:"offer_limit" binding:"required"`
}
type CategoryOfferReceiver struct {
	CategoryID         uint   `json:"category_id" binding:"required"`
	OfferName          string `json:"offer_name" binding:"required"`
	DiscountPercentage int    `json:"discount_percentage" binding:"required"`
	// OfferLimit         int    `json:"offer_limit" binding:"required"`
}
type ReferralAmount struct {
	ReferralAmount float64 `json:"referral_amount"`
}
