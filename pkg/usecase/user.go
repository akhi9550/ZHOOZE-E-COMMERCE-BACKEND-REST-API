package usecase

import (
	"Zhooze/pkg/config"
	"Zhooze/pkg/helper"
	interfaces "Zhooze/pkg/repository/interface"
	services "Zhooze/pkg/usecase/interface"

	"Zhooze/pkg/utils/models"
	"strconv"

	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
)

type userUseCase struct {
	userRepository  interfaces.UserRepository
	orderRepository interfaces.OrderRepository
}

func NewUserUseCase(repository interfaces.UserRepository, orderRepo interfaces.OrderRepository) services.UserUseCase {
	return &userUseCase{
		userRepository:  repository,
		orderRepository: orderRepo,
	}
}

func (ur *userUseCase) UsersSignUp(user models.UserSignUp) (*models.TokenUser, error) {
	email, err := ur.userRepository.CheckUserExistsByEmail(user.Email)
	fmt.Println(email)
	if err != nil {
		return &models.TokenUser{}, errors.New("error with server")
	}
	if email != nil {
		return &models.TokenUser{}, errors.New("user with this email is already exists")
	}

	phone, err := ur.userRepository.CheckUserExistsByPhone(user.Phone)
	fmt.Println(phone, nil)
	if err != nil {
		return &models.TokenUser{}, errors.New("error with server")
	}
	if phone != nil {
		return &models.TokenUser{}, errors.New("user with this phone is already exists")
	}

	hashPassword, err := helper.PasswordHash(user.Password)
	if err != nil {
		return &models.TokenUser{}, errors.New("error in hashing password")
	}
	user.Password = hashPassword
	userData, err := ur.userRepository.UserSignUp(user)
	if err != nil {
		return &models.TokenUser{}, errors.New("could not add the user")
	}
	// create referral code for the user and send in details of referred id of user if it exist
	id := uuid.New().ID()
	str := strconv.Itoa(int(id))
	userReferral := str[:8]
	err = ur.userRepository.CreateReferralEntry(userData, userReferral)
	if err != nil {
		return &models.TokenUser{}, err
	}
	if user.ReferralCode != "" {
		// first check whether if a user with that referralCode exist
		referredUserId, err := ur.userRepository.GetUserIdFromReferrals(user.ReferralCode)
		if err != nil {
			return &models.TokenUser{}, err
		}
		if referredUserId != 0 {
			referralAmount := 150
			err := ur.userRepository.UpdateReferralAmount(float64(referralAmount), referredUserId, userData.Id)
			if err != nil {
				return &models.TokenUser{}, err
			}

		}
	}
	accessToken, err := helper.GenerateAccessToken(userData)
	if err != nil {
		return &models.TokenUser{}, errors.New("couldn't create access token due to error")
	}
	refreshToken, err := helper.GenerateRefreshToken(userData)
	if err != nil {
		return &models.TokenUser{}, errors.New("couldn't create refresh token due to error")
	}
	return &models.TokenUser{
		Users:        userData,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
func (ur *userUseCase) UsersLogin(user models.LoginDetail) (*models.TokenUser, error) {
	email, err := ur.userRepository.CheckUserExistsByEmail(user.Email)
	if err != nil {
		return &models.TokenUser{}, errors.New("error with server")
	}
	if email == nil {
		return &models.TokenUser{}, errors.New("email doesn't exist")
	}
	userdeatils, err := ur.userRepository.FindUserByEmail(user)
	if err != nil {
		return &models.TokenUser{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(userdeatils.Password), []byte(user.Password))
	if err != nil {
		return &models.TokenUser{}, errors.New("password not matching")
	}
	var userDetails models.UserDetailsResponse
	err = copier.Copy(&userDetails, &userdeatils)
	if err != nil {
		return &models.TokenUser{}, err
	}
	accessToken, err := helper.GenerateAccessToken(userDetails)
	if err != nil {
		return &models.TokenUser{}, errors.New("couldn't create accesstoken due to internal error")
	}
	refreshToken, err := helper.GenerateRefreshToken(userDetails)
	if err != nil {
		return &models.TokenUser{}, errors.New("counldn't create refreshtoken due to internal error")
	}
	return &models.TokenUser{
		Users:        userDetails,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
func (ur *userUseCase) AddAddress(userID int, address models.AddressInfo) error {
	err := ur.userRepository.AddAddress(userID, address)
	if err != nil {
		return errors.New("could not add the address")
	}
	return nil
}
func (ur *userUseCase) GetAllAddress(userId int) ([]models.AddressInfoResponse, error) {
	addressInfo, err := ur.userRepository.GetAllAddress(userId)
	if err != nil {
		return []models.AddressInfoResponse{}, err
	}
	return addressInfo, nil

}
func (ur *userUseCase) GetAllAddres(userId int) (models.AddressInfoResponse, error) {
	addressInfo, err := ur.userRepository.GetAllAddres(userId)
	if err != nil {
		return models.AddressInfoResponse{}, err
	}
	return addressInfo, nil

}
func (ur *userUseCase) UserDetails(userID int) (models.UsersProfileDetails, error) {
	return ur.userRepository.UserDetails(userID)
}
func (ur *userUseCase) UpdateUserDetails(userDetails models.UsersProfileDetails, userID int) (models.UsersProfileDetails, error) {
	userExist := ur.userRepository.CheckUserAvailabilityWithUserID(userID)
	if !userExist {
		return models.UsersProfileDetails{}, errors.New("user doesn't exist")
	}
	if userDetails.Email != "" {
		ur.userRepository.UpdateUserEmail(userDetails.Email, userID)
	}
	if userDetails.Firstname != "" {
		ur.userRepository.UpdateFirstName(userDetails.Firstname, userID)
	}
	if userDetails.Lastname != "" {
		ur.userRepository.UpdateLastName(userDetails.Lastname, userID)
	}

	if userDetails.Phone != "" {
		ur.userRepository.UpdateUserPhone(userDetails.Phone, userID)
	}
	return ur.userRepository.UserDetails(userID)
}

func (ur *userUseCase) UpdateAddress(addressDetails models.AddressInfo, addressID, userID int) (models.AddressInfoResponse, error) {
	addressExist := ur.userRepository.CheckAddressAvailabilityWithAddressID(addressID, userID)
	if !addressExist {
		return models.AddressInfoResponse{}, errors.New("address doesn't exist")
	}
	if addressDetails.Name != "" {
		ur.userRepository.UpdateName(addressDetails.Name, addressID)
	}
	if addressDetails.HouseName != "" {
		ur.userRepository.UpdateHouseName(addressDetails.HouseName, addressID)
	}
	if addressDetails.Street != "" {
		ur.userRepository.UpdateStreet(addressDetails.Street, addressID)
	}
	if addressDetails.City != "" {
		ur.userRepository.UpdateCity(addressDetails.City, addressID)
	}
	if addressDetails.State != "" {
		ur.userRepository.UpdateState(addressDetails.State, addressID)
	}
	if addressDetails.Pin != "" {
		ur.userRepository.UpdatePin(addressDetails.Pin, addressID)
	}
	return ur.userRepository.AddressDetails(addressID)
}

func (ur *userUseCase) DeleteAddress(addressID, userID int) error {
	addressExist, err := ur.userRepository.AddressExistInUserProfile(addressID, userID)
	if err != nil {
		return err
	}
	if !addressExist {
		return errors.New("address does not exist in user profile")
	}
	err = ur.userRepository.RemoveFromUserProfile(userID, addressID)
	if err != nil {
		return err
	}
	return nil
}

func (ur *userUseCase) ChangePassword(id int, old string, password string, repassword string) error {
	userPassword, err := ur.userRepository.GetPassword(id)
	if err != nil {
		return errors.New("internal error")
	}
	err = helper.CompareHashAndPassword(userPassword, old)
	if err != nil {
		return errors.New("password incorrect")
	}
	if password != repassword {
		return errors.New("password doesn't match")
	}
	newpassword, err := helper.PasswordHash(password)
	if err != nil {
		return errors.New("error in hashing password")
	}
	return ur.userRepository.ChangePassword(id, string(newpassword))
}
func (ur *userUseCase) UpdateQuantityAdd(id, productID int) error {
	productExist, err := ur.userRepository.ProductExistCart(id, productID)
	if !productExist {
		return errors.New("product doesnot exist cart")
	}
	if err != nil {
		return err
	}
	stock, err := ur.userRepository.ProductStock(productID)
	if err != nil {
		return err
	}
	if stock <= 0 {
		return errors.New("not available out of stock")
	}
	err = ur.userRepository.UpdateQuantityAdd(id, productID)
	if err != nil {
		return err
	}
	err = ur.userRepository.UpdateTotalPrice(id, productID)
	if err != nil {
		return err
	}
	return nil

}

func (ur *userUseCase) UpdateQuantityless(id, productID int) error {
	productExist, err := ur.userRepository.ProductExistCart(id, productID)
	if !productExist {
		return errors.New("product doesnot exist cart")
	}
	if err != nil {
		return err
	}
	stock, err := ur.userRepository.ExistStock(id, productID)
	if err != nil {
		return err
	}
	if stock <= 1 {
		return errors.New("its a maximum")
	}
	err = ur.userRepository.UpdateQuantityless(id, productID)
	if err != nil {
		return err
	}
	err = ur.userRepository.UpdateTotalPrice(id, productID)
	if err != nil {
		return err
	}

	return nil
}

func (ur *userUseCase) ForgotPasswordSend(phone string) error {
	cfg, _ := config.LoadConfig()
	ok := ur.userRepository.FindUserByMobileNumber(phone)
	if !ok {
		return errors.New("the user does not exist")
	}
	helper.TwilioSetup(cfg.ACCOUNTSID, cfg.AUTHTOKEN)
	_, err := helper.TwilioSendOTP(phone, cfg.SERVICESSID)
	if err != nil {
		return errors.New("error ocurred while generating OTP")
	}
	return nil
}

func (ur *userUseCase) ForgotPasswordVerifyAndChange(model models.ForgotVerify) error {
	cfg, _ := config.LoadConfig()
	helper.TwilioSetup(cfg.ACCOUNTSID, cfg.AUTHTOKEN)
	err := helper.TwilioVerifyOTP(cfg.SERVICESSID, model.Otp, model.Phone)
	if err != nil {
		return errors.New("error while verifying")
	}

	id, err := ur.userRepository.FindIdFromPhone(model.Phone)
	if err != nil {
		return errors.New("cannot find user from mobile number")
	}

	newpassword, err := helper.PasswordHashing(model.NewPassword)
	if err != nil {
		return errors.New("error in hashing password")
	}

	// if user is authenticated then change the password i the database
	if err := ur.userRepository.ChangePassword(id, string(newpassword)); err != nil {
		return errors.New("could not change password")
	}

	return nil
}
func (ur *userUseCase) GetCart(id, cart_id int) (models.GetCartResponse, error) {
	products, err := ur.orderRepository.GetProductsInCart(cart_id)
	if err != nil {
		return models.GetCartResponse{}, errors.New("internal error")
	}
	var product_names []string
	for i := range products {
		product_name, err := ur.orderRepository.FindProductNames(products[i])
		if err != nil {
			return models.GetCartResponse{}, errors.New("internal error")
		}
		product_names = append(product_names, product_name)
	}
	var quantity []int
	for i := range products {
		q, err := ur.orderRepository.FindCartQuantity(cart_id, products[i])
		if err != nil {
			return models.GetCartResponse{}, errors.New("internal error")
		}
		quantity = append(quantity, q)
	}

	var price []float64
	for i := range products {
		q, err := ur.orderRepository.FindPrice(products[i])
		if err != nil {
			return models.GetCartResponse{}, errors.New("internal error")
		}
		price = append(price, q)
	}
	var stocks []int

	for _, v := range products {
		stock, err := ur.orderRepository.FindStock(v)
		if err != nil {
			return models.GetCartResponse{}, errors.New("internal error")
		}
		stocks = append(stocks, stock)
	}
	var getcart []models.GetCart
	for i := range product_names {
		var get models.GetCart
		get.ID = products[i]
		get.ProductName = product_names[i]
		get.Quantity = float64(quantity[i])
		get.TotalPrice = price[i]
		get.Product.Stock = stocks[i]
		getcart = append(getcart, get)
	}
	var response models.GetCartResponse
	response.ID = cart_id
	response.Data = getcart
	return response, nil
}
func (ur *userUseCase) ApplyReferral(userID int) (string, error) {
	exist, err := ur.orderRepository.DoesCartExist(userID)
	if err != nil {
		return "", err
	}
	if !exist {
		return "", errors.New("cart does not exist, can't apply offer")
	}
	referralAmount, totalCartAmount, err := ur.userRepository.GetReferralAndTotalAmount(userID)
	if err != nil {
		return "", err
	}
	if totalCartAmount > referralAmount {
		totalCartAmount = totalCartAmount - referralAmount
		referralAmount = 0
	} else {
		referralAmount = referralAmount - totalCartAmount
		totalCartAmount = 0
	}

	err = ur.userRepository.UpdateSomethingBasedOnUserID("referrals", "referral_amount", referralAmount, userID)
	if err != nil {
		return "", err
	}

	err = ur.userRepository.UpdateSomethingBasedOnUserID("carts", "total_price", totalCartAmount, userID)
	if err != nil {
		return "", err
	}

	return "", nil
}
