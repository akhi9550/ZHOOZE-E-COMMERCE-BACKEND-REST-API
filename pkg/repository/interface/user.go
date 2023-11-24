package interfaces

import (
	"Zhooze/pkg/utils/models"
	"Zhooze/pkg/domain"
)

type UserRepository interface {
	CheckUserExistsByEmail(email string) (*domain.User, error)
	CheckUserExistsByPhone(phone string) (*domain.User, error)
	UserSignUp(user models.UserSignUp) (models.UserDetailsResponse, error)
	FindUserByEmail(user models.LoginDetail) (models.UserLoginResponse, error)
	AddAddress(userID int, address models.AddressInfo) error
	GetAllAddress(userId int) ([]models.AddressInfoResponse, error)
	GetAllAddres(userId int) (models.AddressInfoResponse, error)
	UserDetails(userID int) (models.UsersProfileDetails, error)
	CheckUserAvailabilityWithUserID(userID int) bool
	UpdateUserEmail(email string, userID int) error
	UpdateUserPhone(phone string, userID int) error
	UpdateFirstName(name string, userID int) error
	UpdateLastName(name string, userID int) error
	CheckAddressAvailabilityWithAddressID(addressID, userID int) bool
	UpdateName(name string, addressID int) error
	UpdateHouseName(HouseName string, addressID int) error
	UpdateStreet(street string, addressID int) error
	UpdateCity(city string, addressID int) error
	UpdateState(state string, addressID int) error
	UpdatePin(pin string, addressID int) error
	AddressDetails(addressID int) (models.AddressInfoResponse, error)
	ChangePassword(id int, password string) error
	GetPassword(id int) (string, error)
	ProductStock(productID int) (int, error)
	ProductExistCart(userID, productID int) (bool, error)
	UpdateQuantityAdd(id, prdt_id int) error
	UpdateTotalPrice(ID, productID int) error
	UpdateQuantityless(id, prdt_id int) error
	ExistStock(id, productID int) (int, error)
	FindUserByMobileNumber(phone string) bool
	FindIdFromPhone(phone string) (int, error)
	AddressExistInUserProfile(addressID, userID int) (bool, error)
	RemoveFromUserProfile(userID, addressID int) error
	GetReferralAndTotalAmount(userID int) (float64, float64, error)
	UpdateSomethingBasedOnUserID(tableName string, columnName string, updateValue float64, userID int) error
	CreateReferralEntry(userDetails models.UserDetailsResponse, userReferral string) error
	GetUserIdFromReferrals(ReferralCode string) (int, error)
	UpdateReferralAmount(referralAmount float64, referredUserId int, currentUserID int) error
}
