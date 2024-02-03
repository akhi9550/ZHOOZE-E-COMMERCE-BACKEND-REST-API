package interfaces

import "Zhooze/pkg/utils/models"

type UserUseCase interface {
	UsersSignUp(user models.UserSignUp) (*models.TokenUser, error)
	UsersLogin(user models.LoginDetail) (*models.TokenUser, error)
	AddAddress(userID int, address models.AddressInfo) error
	GetAllAddress(userId int) ([]models.AddressInfoResponse, error)
	GetAllAddres(userId int) (models.AddressInfoResponse, error)
	UserDetails(userID int) (models.UsersProfileDetails, error)
	UpdateUserDetails(userDetails models.UsersProfileDetails, userID int) (models.UsersProfileDetails, error)
	UpdateAddress(addressDetails models.AddressInfo, addressID, userID int) (models.AddressInfoResponse, error)
	DeleteAddress(addressID, userID int) error
	ChangePassword(id int, old string, password string, repassword string) error
	UpdateQuantityAdd(id, productID int) error
	UpdateQuantityless(id, productID int) error
	ForgotPasswordSend(phone string) error
	ForgotPasswordVerifyAndChange(model models.ForgotVerify) error
	GetCart(id, cart_id int) (models.GetCartResponse, error)
	ApplyReferral(userID int) (string, error)
}
