package repository

import (
	interfaces "Zhooze/pkg/repository/interface"
	"Zhooze/pkg/utils/models"

	"gorm.io/gorm"
)

type walletRepository struct {
	DB *gorm.DB
}

func NewWalletRepository(DB *gorm.DB) interfaces.WalletRepository {
	return &walletRepository{
		DB: DB,
	}
}
func (wt *walletRepository) GetWallet(userID int) (models.WalletAmount, error) {
	var walletAmount models.WalletAmount
	err := wt.DB.Raw("SELECT amount FROM wallets WHERE user_id = ?", userID).Scan(&walletAmount).Error
	if err != nil {
		return models.WalletAmount{}, err
	}
	return walletAmount, nil
}
func (wt *walletRepository) GetWalletHistory(userID int) ([]models.WalletHistory, error) {
	var history []models.WalletHistory
	err := wt.DB.Raw("SELECT id,order_id,description,amount,is_credited FROM wallet_histories WHERE user_id = ?", userID).Scan(&history).Error
	if err != nil {
		return []models.WalletHistory{}, err
	}
	return history, nil
}
