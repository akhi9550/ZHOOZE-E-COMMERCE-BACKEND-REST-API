package repository

import (
	"Zhooze/pkg/domain"
	interfaces "Zhooze/pkg/repository/interface"
	"Zhooze/pkg/utils/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

type offerRepository struct {
	DB *gorm.DB
}

func NewOfferRepository(DB *gorm.DB) interfaces.OfferRepository {
	return &offerRepository{
		DB: DB,
	}
}

func (of *offerRepository) AddProductOffer(productOffer models.ProductOfferReceiver) error {
	// check if the offer with the offer name already exist in the db
	var count int
	err := of.DB.Raw("SELECT COUNT(*) FROM product_offers WHERE offer_name = ? AND product_id = ?", productOffer.OfferName, productOffer.ProductID).Scan(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("the offer already exists")
	}

	// if there is any other offer for this product delete that before adding this one
	count = 0
	err = of.DB.Raw("SELECT COUNT(*) FROM product_offers WHERE product_id = ?", productOffer.ProductID).Scan(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		err = of.DB.Exec("DELETE FROM product_offers WHERE product_id = ?", productOffer.ProductID).Error
		if err != nil {
			return err
		}
	}

	startDate := time.Now()
	endDate := time.Now().Add(time.Hour * 24 * 5)
	err = of.DB.Exec("INSERT INTO product_offers (product_id, offer_name, discount_percentage, start_date, end_date) VALUES (?, ?, ?, ?, ?)", productOffer.ProductID, productOffer.OfferName, productOffer.DiscountPercentage, startDate, endDate).Error
	if err != nil {
		return err
	}

	return nil

}
func (of *offerRepository) GetOffers() ([]domain.ProductOffer, error) {
	var model []domain.ProductOffer
	err := of.DB.Raw("SELECT * FROM product_offers").Scan(&model).Error
	if err != nil {
		return []domain.ProductOffer{}, err
	}

	return model, nil
}
func (of *offerRepository) MakeOfferExpire(id int) error {
	if err := of.DB.Exec("DELETE FROM product_offers WHERE id = $1", id).Error; err != nil {
		return err
	}

	return nil
}
func (of *offerRepository) FindDiscountPercentageForProduct(id int) (int, error) {
	var percentage int
	err := of.DB.Raw("SELECT discount_percentage FROM product_offers WHERE product_id= $1 ", id).Scan(&percentage).Error
	if err != nil {
		return 0, err
	}

	return percentage, nil
}
func (of *offerRepository) AddCategoryOffer(categoryOffer models.CategoryOfferReceiver) error {

	// check if the offer with the offer name already exist in the db
	var count int
	err := of.DB.Raw("SELECT COUNT(*) FROM category_offers WHERE offer_name = ?", categoryOffer.OfferName).Scan(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("the offer already exists")
	}
	// if there is any other offer for this category delete that before adding this one
	count = 0
	err = of.DB.Raw("SELECT COUNT(*) FROM category_offers WHERE category_id = ?", categoryOffer.CategoryID).Scan(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		err = of.DB.Exec("DELETE FROM category_offers WHERE category_id = ?", categoryOffer.CategoryID).Error
		if err != nil {
			return err
		}
	}
	startDate := time.Now()
	endDate := time.Now().Add(time.Hour * 24 * 5)
	err = of.DB.Exec("INSERT INTO category_offers (category_id, offer_name, discount_percentage, start_date, end_date) VALUES (?, ?, ?, ?, ?)", categoryOffer.CategoryID, categoryOffer.OfferName, categoryOffer.DiscountPercentage, startDate, endDate).Error
	if err != nil {
		return err
	}
	return nil
}
func (of *offerRepository) GetCategoryOffer() ([]domain.CategoryOffer, error) {
	var model []domain.CategoryOffer
	err := of.DB.Raw("SELECT * FROM category_offers").Scan(&model).Error
	if err != nil {
		return []domain.CategoryOffer{}, err
	}
	return model, nil
}
func (of *offerRepository) ExpireCategoryOffer(id int) error {
	if err := of.DB.Exec("DELETE FROM category_offers WHERE id = $1", id).Error; err != nil {
		return err
	}

	return nil
}
func (of *offerRepository) FindDiscountPercentageForCategory(id int) (int, error) {
	var percentage int
	err := of.DB.Raw("SELECT discount_percentage FROM category_offers WHERE category_id= $1 ", id).Scan(&percentage).Error
	if err != nil {
		return 0, err
	}

	return percentage, nil
}
