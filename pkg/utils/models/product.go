package models

type Image struct {
	Url string `JSON:"url" `
}

type ProductBrief struct {
	ID              uint     `json:"id" gorm:"unique;not null"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	CategoryID      int      `json:"category_id"`
	Size            int      `json:"size"`
	Stock           int      `json:"stock"`
	Price           float64  `json:"price"`
	DiscountedPrice float64  `json:"discounted_price"`
	ProductStatus   string   `json:"product_status"`
	Image           []string `json:"image"`
}

type ProductReceiver struct {
	Name        string  `json:"name" `
	Description string  `json:"description"`
	CategoryID  uint    `json:"category_id"`
	Size        int     `json:"size"`
	Stock       int     `json:"stock"`
	Price       float64 `json:"price"`
}
type Product struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	CategoryID  uint    `json:"category_id" validate:"required"`
	Size        int     `json:"size" validate:"required"`
	Stock       int     `json:"stock" validate:"required"`
	Price       float64 `json:"price" validate:"required"`
}
type Category struct {
	Category string `json:"category"`
}
type SetNewName struct {
	Current string `json:"current"`
	New     string `json:"new"`
}
type ProductUpdate struct {
	ProductId int `json:"product_id"`
	Stock     int `json:"stock"`
}
type ProductUpdateReciever struct {
	ProductID int
	Stock     int
}
