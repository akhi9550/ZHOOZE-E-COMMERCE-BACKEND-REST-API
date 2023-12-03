package usecase

import (
	"Zhooze/pkg/domain"
	interfaces "Zhooze/pkg/repository/interface"
	services "Zhooze/pkg/usecase/interface"

	"Zhooze/pkg/utils/models"
)

type offerUseCase struct {
	offerRepository interfaces.OfferRepository
}

func NewOfferUseCase(repo interfaces.OfferRepository) services.OfferUseCase {
	return &offerUseCase{
		offerRepository: repo,
	}
}

func (of *offerUseCase) AddProductOffer(model models.ProductOfferReceiver) error {
	if err := of.offerRepository.AddProductOffer(model); err != nil {
		return err
	}

	return nil
}
func (of *offerUseCase) GetOffers() ([]domain.ProductOffer, error) {

	offers, err := of.offerRepository.GetOffers()
	if err != nil {
		return []domain.ProductOffer{}, err
	}
	return offers, nil

}
// func (of *offerUseCase) MakeOfferExpire(id int) error {
// 	if err := of.offerRepository.MakeOfferExpire(id); err != nil {
// 		return err
// 	}

// 	return nil
// }

func (of *offerUseCase) AddCategoryOffer(model models.CategoryOfferReceiver) error {
	if err := of.offerRepository.AddCategoryOffer(model); err != nil {
		return err
	}

	return nil
}
func (of *offerUseCase) GetCategoryOffer() ([]domain.CategoryOffer, error) {

	offers, err := of.offerRepository.GetCategoryOffer()
	if err != nil {
		return []domain.CategoryOffer{}, err
	}
	return offers, nil

}
// func (of *offerUseCase) ExpireCategoryOffer(id int) error {
// 	if err := of.offerRepository.ExpireCategoryOffer(id); err != nil {
// 		return err
// 	}

// 	return nil
// }
