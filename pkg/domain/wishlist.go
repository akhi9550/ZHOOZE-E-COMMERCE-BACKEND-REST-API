package domain

type WishList struct {
	ID        uint    `json:"id" gorm:"uniquekey; not null"`
	UserID    uint    `json:"user_id"`
	Users     User    `json:"-" gorm:"foreignkey:UserID"`
	ProductID uint    `json:"product_id"`
	Products  Product `json:"-" gorm:"foreignkey:ProductID"`
}
