package interfaces

import "Zhooze/pkg/utils/models"

type WishlistRepository interface {
	GetWishList(userID int) ([]models.WishListResponse, error)
	ProductExistInWishList(productID, userID int) (bool, error)
	AddToWishlist(userID, productID int) error
	RemoveFromWishList(userID, productID int) error
}
