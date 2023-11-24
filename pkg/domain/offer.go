package domain

import (
	"time"
)

type ProductOffer struct {
	ID                 uint      `json:"id" gorm:"unique; not null"`
	ProductID          uint      `json:"product_id"`
	Products           Product   `json:"-" gorm:"foreignkey:ProductID"`
	OfferName          string    `json:"offer_name"`
	DiscountPercentage int       `json:"discount_percentage"`
	StartDate          time.Time `json:"start_date"`
	EndDate            time.Time `json:"end_date"`
	// OfferLimit         int       `json:"offer_limit"`
	// OfferUsed          int       `json:"offer_used"`
	// Valid bool `json:"valid" gorm:"default:True"`
}

type CategoryOffer struct {
	ID                 uint      `json:"id" gorm:"unique; not null"`
	CategoryID         uint      `json:"category_id"`
	Category           Category  `json:"-" gorm:"foreignkey:CategoryID"`
	OfferName          string    `json:"offer_name"`
	DiscountPercentage int       `json:"discount_percentage"`
	StartDate          time.Time `json:"start_date"`
	EndDate            time.Time `json:"end_date"`
}

// type ProductOfferUsed struct {
// 	gorm.Model
// 	UserID         uint         `json:"user_id"`
// 	Users          User         `json:"-" gorm:"foreignkey:UserID"`
// 	ProductOfferID uint         `json:"product_offer_id"`
// 	ProductOffer   ProductOffer `json:"-" gorm:"constraint:OnDelete:CASCADE;foreignKey:ProductOfferID"`
// 	OfferAmount    float64      `json:"offer_amount"`
// 	OfferCount     int          `json:"offer_count"`
// 	Used           bool         `json:"used"`
// }

// type CategoryOfferUsed struct {
// 	gorm.Model
// 	UserID          uint          `json:"user_id"`
// 	Users           User          `json:"-" gorm:"foreignkey:UserID"`
// 	CategoryOfferID uint          `json:"product_offer_id"`
// 	CategoryOffer   CategoryOffer `json:"-" gorm:"constraint:OnDelete:CASCADE;foreignKey:CategoryOfferID"`
// 	OfferAmount     float64       `json:"offer_amount"`
// 	OfferCount      int           `json:"offer_count"`
// 	Used            bool          `json:"used"`
// }
