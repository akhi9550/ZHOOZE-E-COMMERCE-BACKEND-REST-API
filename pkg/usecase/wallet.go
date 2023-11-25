package usecase

import (
	interfaces "Zhooze/pkg/repository/interface"
	services "Zhooze/pkg/usecase/interface"
	"Zhooze/pkg/utils/models"
)

type walletUseCase struct {
	walletRepository interfaces.WalletRepository
}

func NewWalletUseCase(repository interfaces.WalletRepository) services.WalletUseCase {
	return &walletUseCase{
		walletRepository: repository,
	}
}

func (wt *walletUseCase) GetWallet(userID int) (models.WalletAmount, error) {
	return wt.walletRepository.GetWallet(userID)
}
func (wt *walletUseCase) GetWalletHistory(userID int) ([]models.WalletHistory, error) {
	return wt.walletRepository.GetWalletHistory(userID)

}
