package interfaces

import "Zhooze/pkg/utils/models"

type WalletUseCase interface {
	GetWallet(userID int) (models.WalletAmount, error)
	GetWalletHistory(userID int) ([]models.WalletHistory, error)
}
