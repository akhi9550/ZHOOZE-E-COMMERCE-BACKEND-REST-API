package interfaces

import (
	"Zhooze/pkg/domain"
	"Zhooze/pkg/utils/models"
)

type OfferUseCase interface {
	AddProductOffer(model models.ProductOfferReceiver)
	GetOffers() ([]domain.ProductOffer, error)
	MakeOfferExpire(id int) error
	AddCategoryOffer(model models.CategoryOfferReceiver) error
	GetCategoryOffer() ([]domain.CategoryOffer, error)
	ExpireCategoryOffer(id int) error
}
