package usecase

import (
	interfaces "Zhooze/pkg/repository/interface"
	services "Zhooze/pkg/usecase/interface"

	"Zhooze/pkg/utils/models"
	"errors"
)

type wishlistUseCase struct {
	wishlistRepository interfaces.WishlistRepository
	productRepository  interfaces.ProductRepository
}

func NewWishListUseCase(repository interfaces.WishlistRepository, productRepo interfaces.ProductRepository) services.WishListUseCase {
	return &wishlistUseCase{
		wishlistRepository: repository,
		productRepository:  productRepo,
	}
}
func (wh *wishlistUseCase) GetWishList(userID int) ([]models.WishListResponse, error) {
	wishList, err := wh.wishlistRepository.GetWishList(userID)
	if err != nil {
		return []models.WishListResponse{}, err
	}
	return wishList, err
}

func (wh *wishlistUseCase) AddToWishlist(userID, productID int) error {
	productExist, err := wh.productRepository.DoesProductExist(productID)
	if err != nil {
		return err
	}
	if !productExist {
		return errors.New("product does not exist")
	}
	productExistInWishList, err := wh.wishlistRepository.ProductExistInWishList(productID, userID)
	if err != nil {
		return err
	}
	if productExistInWishList {
		return errors.New("product already exist in wishlist")
	}
	err = wh.wishlistRepository.AddToWishlist(userID, productID)
	if err != nil {
		return err
	}
	return nil
}
func (wh *wishlistUseCase) RemoveFromWishlist(productID, userID int) error {
	productExistInWishList, err := wh.wishlistRepository.ProductExistInWishList(productID, userID)
	if err != nil {
		return err
	}
	if !productExistInWishList {
		return errors.New("product does not exist in wishlist")
	}
	err = wh.wishlistRepository.RemoveFromWishList(userID, productID)
	if err != nil {
		return err
	}
	return nil
}
