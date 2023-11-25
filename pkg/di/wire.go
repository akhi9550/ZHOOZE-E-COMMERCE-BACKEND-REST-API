//go:build wireinject
// +build wireinject

package di

import (
	http "Zhooze/pkg/api"
	handler "Zhooze/pkg/api/handlers"
	config "Zhooze/pkg/config"
	db "Zhooze/pkg/db"
	repository "Zhooze/pkg/repository"
	usecase "Zhooze/pkg/usecase"

	"github.com/google/wire"
)

func InitializeAPI(cfg config.Config) (*http.ServerHTTP, error) {
	wire.Build(db.ConnectDatabase, repository.NewUserRepository, usecase.NewUserUseCase, handler.NewUserHandler,
		http.NewServerHTTP, repository.NewProductRepository, usecase.NewProductUseCase, handler.NewProductHandler,
		handler.NewOtpHandler, usecase.NewOtpUseCase, repository.NewOtpRepository,
		repository.NewAdminRepository, usecase.NewAdminUseCase, handler.NewAdminHandler,
		handler.NewCartHandler, usecase.NewCartUseCase, repository.NewCartRepository,
		repository.NewOfferRepository, usecase.NewOfferUseCase, handler.NewOfferHandler,
		repository.NewCategoryRepository, usecase.NewCategoryUseCase, handler.NewCategory,
		repository.NewWishlistRepository, usecase.NewWishListUseCase, handler.NewWishListHandler,
		repository.NewWalletRepository, usecase.NewWalletUseCase, handler.NewWalletHandler,
		repository.NewCouponRepository, usecase.NewCouponUseCase, handler.NewCouponHandler,
		repository.NewOrderRepository, usecase.NewOrderUseCase, handler.NewOrderHandler,
		repository.NewPaymentRepository, usecase.NewPaymentUseCase, handler.NeWPaymentHandler)
	return &http.ServerHTTP{}, nil
}
