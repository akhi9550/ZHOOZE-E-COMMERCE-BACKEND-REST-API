package repository

import (
	"Zhooze/pkg/domain"
	interfaces "Zhooze/pkg/repository/interface"
	"Zhooze/pkg/utils/models"
	"errors"

	"gorm.io/gorm"
)



type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) interfaces.UserRepository {
	return &userRepository{
		DB: DB,
	}
}

func (ur *userRepository) CheckUserExistsByEmail(email string) (*domain.User, error) {
	var user domain.User
	res := ur.DB.Where(&domain.User{Email: email}).First(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, res.Error
	}
	return &user, nil
}
func (ur *userRepository) CheckUserExistsByPhone(phone string) (*domain.User, error) {
	var user domain.User
	res := ur.DB.Where(&domain.User{Phone: phone}).First(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, res.Error
	}
	return &user, nil
}
func (ur *userRepository) UserSignUp(user models.UserSignUp) (models.UserDetailsResponse, error) {
	var SignupDetail models.UserDetailsResponse
	err := ur.DB.Raw("INSERT INTO users(firstname,lastname,email,password,phone)VALUES(?,?,?,?,?)RETURNING id,firstname,lastname,email,password,phone", user.Firstname, user.Lastname, user.Email, user.Password, user.Phone).Scan(&SignupDetail).Error
	if err != nil {
		return models.UserDetailsResponse{}, err
	}
	return SignupDetail, nil
}
func (ur *userRepository) FindUserByEmail(user models.LoginDetail) (models.UserLoginResponse, error) {
	var userDetails models.UserLoginResponse
	err := ur.DB.Raw("SELECT * FROM users WHERE email=? and blocked=false and isadmin=false", user.Email).Scan(&userDetails).Error
	if err != nil {
		return models.UserLoginResponse{}, errors.New("error checking user details")
	}
	return userDetails, nil
}
func (ur *userRepository) AddAddress(userID int, address models.AddressInfo) error {
	err := ur.DB.Exec("INSERT INTO addresses(user_id,name,house_name,street,city,state,pin)VALUES(?,?,?,?,?,?,?)", userID, address.Name, address.HouseName, address.Street, address.City, address.State, address.Pin).Error
	if err != nil {
		return errors.New("could not add address")
	}
	return nil
}
func (ur *userRepository) GetAllAddress(userId int) ([]models.AddressInfoResponse, error) {
	var addressInfoResponse []models.AddressInfoResponse
	if err := ur.DB.Raw("SELECT * FROM addresses WHERE user_id = ?", userId).Scan(&addressInfoResponse).Error; err != nil {
		return []models.AddressInfoResponse{}, err
	}
	return addressInfoResponse, nil
}
func (ur *userRepository) GetAllAddres(userId int) (models.AddressInfoResponse, error) {
	var addressInfoResponse models.AddressInfoResponse
	if err := ur.DB.Raw("SELECT * FROM addresses WHERE user_id = ?", userId).Scan(&addressInfoResponse).Error; err != nil {
		return models.AddressInfoResponse{}, err
	}
	return addressInfoResponse, nil
}
func (ur *userRepository) UserDetails(userID int) (models.UsersProfileDetails, error) {
	var userDetails models.UsersProfileDetails
	err := ur.DB.Raw("SELECT u.firstname,u.lastname,u.email,u.phone FROM users u WHERE u.id = ?", userID).Row().Scan(&userDetails.Firstname, &userDetails.Lastname, &userDetails.Email, &userDetails.Phone)
	if err != nil {
		return models.UsersProfileDetails{}, err
	}
	return userDetails, nil
}
func (ur *userRepository) CheckUserAvailabilityWithUserID(userID int) bool {
	var count int
	if err := ur.DB.Raw("SELECT COUNT(*) FROM users WHERE id= ?", userID).Scan(&count).Error; err != nil {
		return false
	}
	return count > 0
}
func (ur *userRepository) UpdateUserEmail(email string, userID int) error {
	err := ur.DB.Exec("UPDATE users SET email= ? WHERE id = ?", email, userID).Error
	if err != nil {
		return err
	}
	return nil
}
func (ur *userRepository) UpdateUserPhone(phone string, userID int) error {
	if err := ur.DB.Exec("UPDATE users SET phone = ? WHERE id = ?", phone, userID).Error; err != nil {
		return err
	}
	return nil
}
func (ur *userRepository) UpdateFirstName(name string, userID int) error {

	err := ur.DB.Exec("UPDATE users SET firstname = ? WHERE id = ?", name, userID).Error
	if err != nil {
		return err
	}
	return nil

}
func (ur *userRepository) UpdateLastName(name string, userID int) error {

	err := ur.DB.Exec("UPDATE users SET lastname = ? WHERE id = ?", name, userID).Error
	if err != nil {
		return err
	}
	return nil

}
func (ur *userRepository) CheckAddressAvailabilityWithAddressID(addressID, userID int) bool {
	var count int
	if err := ur.DB.Raw("SELECT COUNT(*) FROM addresses WHERE id = ? AND user_id = ?", addressID, userID).Scan(&count).Error; err != nil {
		return false
	}
	return count > 0
}
func (ur *userRepository) UpdateName(name string, addressID int) error {
	err := ur.DB.Exec("UPDATE addresses SET name= ? WHERE id = ?", name, addressID).Error
	if err != nil {
		return err
	}
	return nil
}
func (ur *userRepository) UpdateHouseName(HouseName string, addressID int) error {
	err := ur.DB.Exec("UPDATE addresses SET house_name= ? WHERE id = ?", HouseName, addressID).Error
	if err != nil {
		return err
	}
	return nil
}
func (ur *userRepository) UpdateStreet(street string, addressID int) error {
	err := ur.DB.Exec("UPDATE addresses SET street= ? WHERE id = ?", street, addressID).Error
	if err != nil {
		return err
	}
	return nil
}
func (ur *userRepository) UpdateCity(city string, addressID int) error {
	err := ur.DB.Exec("UPDATE addresses SET city= ? WHERE id = ?", city, addressID).Error
	if err != nil {
		return err
	}
	return nil
}
func (ur *userRepository) UpdateState(state string, addressID int) error {
	err := ur.DB.Exec("UPDATE addresses SET state= ? WHERE id = ?", state, addressID).Error
	if err != nil {
		return err
	}
	return nil
}
func (ur *userRepository) UpdatePin(pin string, addressID int) error {
	err := ur.DB.Exec("UPDATE addresses SET pin= ? WHERE id = ?", pin, addressID).Error
	if err != nil {
		return err
	}
	return nil
}
func (ur *userRepository) AddressDetails(addressID int) (models.AddressInfoResponse, error) {
	var addressDetails models.AddressInfoResponse
	err := ur.DB.Raw("SELECT a.id, a.name, a.house_name, a.street, a.city, a.state, a.pin FROM addresses a WHERE a.id = ?", addressID).Row().Scan(&addressDetails.ID, &addressDetails.Name, &addressDetails.HouseName, &addressDetails.Street, &addressDetails.City, &addressDetails.State, &addressDetails.Pin)
	if err != nil {
		return models.AddressInfoResponse{}, err
	}
	return addressDetails, nil
}

func (ur *userRepository) ChangePassword(id int, password string) error {
	err := ur.DB.Exec("UPDATE users SET password = $1 WHERE id = $2", password, id).Error
	if err != nil {
		return err
	}
	return nil
}
func (ur *userRepository) GetPassword(id int) (string, error) {
	var userPassword string
	err := ur.DB.Raw("SELECT password FROM users WHERE id = ?", id).Scan(&userPassword).Error
	if err != nil {
		return "", err
	}
	return userPassword, nil
}
func (ur *userRepository) ProductStock(productID int) (int, error) {
	var a int
	err := ur.DB.Raw("SELECT stock FROM products WHERE id = ?", productID).Scan(&a).Error
	if err != nil {
		return 0, err
	}
	return a, nil
}
func (ur *userRepository) ProductExistCart(userID, productID int) (bool, error) {
	var count int
	err := ur.DB.Raw("SELECT COUNT(*) FROM carts WHERE user_id = ? AND product_id = ?", userID, productID).Scan(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (ur *userRepository) UpdateQuantityAdd(id, prdt_id int) error {
	err := ur.DB.Exec(`	UPDATE carts
	SET quantity = quantity + 1
	WHERE user_id=$1 AND product_id=$2`, id, prdt_id).Error
	if err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) UpdateTotalPrice(ID, productID int) error {
	err := ur.DB.Exec(` UPDATE carts SET total_price = quantity * (select price from products where id = $1) WHERE user_id =$2 AND product_id = $3`, productID, ID, productID).Error
	if err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) UpdateQuantityless(id, prdt_id int) error {
	err := ur.DB.Exec(`	UPDATE carts
	SET quantity = quantity - 1
	WHERE user_id=$1 AND product_id=$2`, id, prdt_id).Error
	if err != nil {
		return err
	}
	return nil
}
func (ur *userRepository) ExistStock(id, productID int) (int, error) {
	var a int
	err := ur.DB.Raw("SELECT quantity FROM carts WHERE user_id = ? AND product_id = ?", id, productID).Scan(&a).Error
	if err != nil {
		return 0, err
	}
	return a, nil
}
func (ur *userRepository) FindUserByMobileNumber(phone string) bool {

	var count int
	if err := ur.DB.Raw("SELECT count(*) FROM users WHERE phone = ?", phone).Scan(&count).Error; err != nil {
		return false
	}

	return count > 0

}
func (ur *userRepository) FindIdFromPhone(phone string) (int, error) {
	var id int
	if err := ur.DB.Raw("SELECT id FROM users WHERE phone=?", phone).Scan(&id).Error; err != nil {
		return id, err
	}
	return id, nil
}
func (ur *userRepository) AddressExistInUserProfile(addressID, userID int) (bool, error) {
	var count int
	err := ur.DB.Raw("SELECT COUNT (*) FROM addresses WHERE user_id = $1 AND id = $2", userID, addressID).Scan(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (ur *userRepository) RemoveFromUserProfile(userID, addressID int) error {
	err := ur.DB.Exec("DELETE FROM addresses WHERE user_id = ? AND  id= ?", userID, addressID).Error
	if err != nil {
		return err
	}
	return nil
}
func (ur *userRepository) GetReferralAndTotalAmount(userID int) (float64, float64, error) {
	// first check whether the cart is empty -- do this for coupon too
	var cartDetails struct {
		ReferralAmount  float64
		TotalCartAmount float64
	}

	err := ur.DB.Raw("SELECT (SELECT referral_amount FROM referrals WHERE user_id = ?) AS referral_amount, COALESCE(SUM(total_price), 0) AS total_cart_amount FROM carts WHERE user_id = ?", userID, userID).Scan(&cartDetails).Error
	if err != nil {
		return 0.0, 0.0, err
	}

	return cartDetails.ReferralAmount, cartDetails.TotalCartAmount, nil

}
func (ur *userRepository) UpdateSomethingBasedOnUserID(tableName string, columnName string, updateValue float64, userID int) error {

	err := ur.DB.Exec("UPDATE "+tableName+" SET "+columnName+" = ? WHERE user_id = ?", updateValue, userID).Error
	if err != nil {
		ur.DB.Rollback()
		return err
	}
	return nil

}
func (ur *userRepository) CreateReferralEntry(userDetails models.UserDetailsResponse, userReferral string) error {

	err := ur.DB.Exec("INSERT INTO referrals (user_id,referral_code,referral_amount) VALUES (?,?,?)", userDetails.Id, userReferral, 0).Error
	if err != nil {
		return err
	}

	return nil

}
func (ur *userRepository) GetUserIdFromReferrals(ReferralCode string) (int, error) {

	var referredUserId int
	err := ur.DB.Raw("SELECT user_id FROM referrals WHERE referral_code = ?", ReferralCode).Scan(&referredUserId).Error
	if err != nil {
		return 0, nil
	}

	return referredUserId, nil
}

func (ur *userRepository) UpdateReferralAmount(referralAmount float64, referredUserId int, currentUserID int) error {

	err := ur.DB.Exec("UPDATE referrals SET referral_amount = ? , referred_user_id = ? WHERE user_id = ? ", referralAmount, referredUserId, currentUserID).Error
	if err != nil {
		return err
	}

	// find the current amount in referred users referral table and add 100 with that
	err = ur.DB.Exec("UPDATE referrals SET referral_amount = referral_amount + ? WHERE user_id = ? ", referralAmount, referredUserId).Error
	if err != nil {
		return err
	}

	return nil

}
