package interfaces

import "Zhooze/pkg/utils/models"

type WishListUseCase interface {
	GetWishList(userID int) ([]models.WishListResponse, error)
	AddToWishlist(userID, productID int) error
	RemoveFromWishlist(productID, userID int) error
}
