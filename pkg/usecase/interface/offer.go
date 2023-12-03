package interfaces

import (
	"Zhooze/pkg/domain"
	"Zhooze/pkg/utils/models"
)

type OfferUseCase interface {
	AddProductOffer(models.ProductOfferReceiver)error
	GetOffers() ([]domain.ProductOffer, error)
	// MakeOfferExpire(id int) error
	AddCategoryOffer(models.CategoryOfferReceiver) error
	GetCategoryOffer() ([]domain.CategoryOffer, error)
	// ExpireCategoryOffer(id int) error
}
