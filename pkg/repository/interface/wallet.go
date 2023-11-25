package interfaces

import "Zhooze/pkg/utils/models"

type WalletRepository interface {
	GetWallet(userID int) (models.WalletAmount, error)
	GetWalletHistory(userID int) ([]models.WalletHistory, error)
}
