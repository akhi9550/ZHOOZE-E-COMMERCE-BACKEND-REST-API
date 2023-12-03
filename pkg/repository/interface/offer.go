package interfaces

import (
	"Zhooze/pkg/domain"
	"Zhooze/pkg/utils/models"
)

type OfferRepository interface {
	AddProductOffer(productOffer models.ProductOfferReceiver) error
	GetOffers() ([]domain.ProductOffer, error)
	// MakeOfferExpire(id int) error
	FindDiscountPercentageForProduct(id int) (int, error)
	AddCategoryOffer(categoryOffer models.CategoryOfferReceiver) error
	GetCategoryOffer() ([]domain.CategoryOffer, error)
	// ExpireCategoryOffer(id int) error
	FindDiscountPercentageForCategory(id int) (int, error)
}