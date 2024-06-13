package db

import (
	"Zhooze/pkg/config"
	"Zhooze/pkg/domain"
	"Zhooze/pkg/helper"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase(confg config.Config) (*gorm.DB, error) {
	connectTo := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s", confg.DBHost, confg.DBUser, confg.DBName, confg.DBPort, confg.DBPassword)
	db, err := gorm.Open(postgres.Open(connectTo), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database:%w", err)
	}
	DB = db
	db.AutoMigrate(&domain.Admin{})
	db.AutoMigrate(&domain.User{})
	db.AutoMigrate(&domain.Product{})
	db.AutoMigrate(&domain.Category{})
	db.AutoMigrate(&domain.Address{})
	db.AutoMigrate(&domain.Cart{})
	db.AutoMigrate(&domain.Order{})
	db.AutoMigrate(&domain.OrderItem{})
	db.AutoMigrate(&domain.RazerPay{})
	db.AutoMigrate(&domain.PaymentMethod{})
	db.AutoMigrate(&domain.Coupons{})
	db.AutoMigrate(&domain.UsedCoupon{})
	db.AutoMigrate(&domain.OrderCoupon{})
	db.AutoMigrate(&domain.ProductOffer{})
	db.AutoMigrate(&domain.CategoryOffer{})
	db.AutoMigrate(&domain.Referral{})
	db.AutoMigrate(&domain.WishList{})
	db.AutoMigrate(&domain.Image{})
	db.AutoMigrate(&domain.Wallet{})
	db.AutoMigrate(&domain.WalletHistory{})
	CheckAndCreateAdmin(db)
	return DB, nil
}
func CheckAndCreateAdmin(db *gorm.DB) {
	var count int64
	db.Model(&domain.User{}).Count(&count)
	if count == 0 {
		password := "admin@123"
		hashPassword, err := helper.PasswordHash(password)
		if err != nil {
			return
		}
		admin := domain.User{
			ID:        1,
			Firstname: "Zhooze",
			Lastname:  "Admin",
			Email:     "admin@zhooze.com",
			Password:  hashPassword,
			Phone:     "+919061757507",
			Blocked:   false,
			Isadmin:   true,
		}
		db.Create(&admin)
	}
}
