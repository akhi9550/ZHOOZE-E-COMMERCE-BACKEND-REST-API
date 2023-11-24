package repository

import (
	interfaces "Zhooze/pkg/repository/interface"
	"Zhooze/pkg/utils/models"
	"errors"

	"gorm.io/gorm"
)

type wishlistRepository struct {
	DB *gorm.DB
}

func NewWishlistRepository(DB *gorm.DB) interfaces.WishlistRepository {
	return &wishlistRepository{
		DB: DB,
	}
}

func (ws *wishlistRepository) GetWishList(userID int) ([]models.WishListResponse, error) {
	var wishList []models.WishListResponse
	err := ws.DB.Raw("SELECT products.id as product_id,products.name as product_name,products.price as product_price FROM products INNER JOIN wish_lists ON products.id = wish_lists.product_id WHERE wish_lists.user_id = ?", userID).Scan(&wishList).Error
	if err != nil {
		return []models.WishListResponse{}, err
	}
	return wishList, nil
}
func (ws *wishlistRepository) ProductExistInWishList(productID, userID int) (bool, error) {
	var count int
	err := ws.DB.Raw("SELECT COUNT(*) FROM wish_lists WHERE user_id = ? AND product_id = ?", userID, productID).Scan(&count).Error
	if err != nil {
		return false, errors.New("error checking user product already present")
	}
	return count > 0, nil
}
func (ws *wishlistRepository) AddToWishlist(userID, productID int) error {
	err := ws.DB.Exec("INSERT INTO wish_lists (user_id,product_id) VALUES (?,?)", userID, productID).Error
	if err != nil {
		return err
	}
	return nil
}
func (ws *wishlistRepository) RemoveFromWishList(userID, productID int) error {
	err := ws.DB.Exec("DELETE FROM wish_lists WHERE user_id = ? AND product_id = ?", userID, productID).Error
	if err != nil {
		return err
	}
	return nil
}
