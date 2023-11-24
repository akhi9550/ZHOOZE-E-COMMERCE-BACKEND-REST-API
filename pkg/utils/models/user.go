package models

type UserSignUp struct {
	Firstname    string `json:"firstname" validate:"gte=3"`
	Lastname     string `json:"lastname" validate:"gte=1"`
	Email        string `json:"email" validate:"email"`
	Password     string `json:"password" validate:"min=6,max=20"`
	Phone        string `json:"phone" validate:"e164"`
	ReferralCode string `json:"referral_code"`
}
type UserDetailsResponse struct {
	Id        int    `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}
type UserDetailsAtAdmin struct {
	Id          int    `json:"id"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	BlockStatus bool   `json:"block_status"`
}
type TokenUser struct {
	Users        UserDetailsResponse
	AccessToken  string
	RefreshToken string
}
type LoginDetail struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	Id        uint   `json:"id"`
	UserId    uint   `json:"user_id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
}
type AddressInfoResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name" validate:"required"`
	HouseName string `json:"house_name" validate:"required"`
	Street    string `json:"street"`
	City      string `json:"city"`
	State     string `json:"state" validate:"required"`
	Pin       string `json:"pin" validate:"required"`
}
type AddressInfo struct {
	Name      string `json:"name" validate:"required"`
	HouseName string `json:"house_name" validate:"required"`
	Street    string `json:"street"`
	City      string `json:"city"`
	State     string `json:"state" validate:"required"`
	Pin       string `json:"pin" validate:"required"`
}
type UsersProfileDetails struct {
	Firstname string `json:"firstname" `
	Lastname  string `json:"lastname" `
	Email     string `json:"email" `
	Phone     string `json:"phone" `
}

type UpdatePassword struct {
	OldPassword        string `json:"old_password" binding:"required"`
	NewPassword        string `json:"new_password" binding:"required"`
	ConfirmNewPassword string `json:"confirm_new_password" binding:"required"`
}

type PaymentDetails struct {
	ID           uint   `json:"id"`
	Payment_Name string `json:"payment_name"`
}

type CheckoutDetails struct {
	AddressInfoResponse []AddressInfoResponse
	Payment_Method      []PaymentDetails
	Cart                []Cart
	Total_Price         float64
}
type ChangePassword struct {
	Oldpassword string `json:"old_password"`
	Password    string `json:"password"`
	Repassword  string `json:"re_password"`
}

type ForgotPasswordSend struct {
	Phone string `json:"phone"`
}
type ForgotVerify struct {
	Phone       string `json:"phone" binding:"required" validate:"required"`
	Otp         string `json:"otp" binding:"required"`
	NewPassword string `json:"newpassword" binding:"required" validate:"min=6,max=20"`
}
