package usecase

import (
	"Zhooze/pkg/config"
	"Zhooze/pkg/helper"
	interfaces "Zhooze/pkg/repository/interface"
	services "Zhooze/pkg/usecase/interface"

	"Zhooze/pkg/utils/models"
	"errors"

	"github.com/jinzhu/copier"
)

type otpUseCase struct {
	otpRepository interfaces.OtpRepository
}

func NewOtpUseCase(repo interfaces.OtpRepository) services.OtpUseCase {
	return &otpUseCase{
		otpRepository: repo,
	}
}

func (op *otpUseCase) SendOtp(phone string) error {
	cfg, err := config.LoadConfig()

	if err != nil {
		return err
	}
	user, err := op.otpRepository.FindUserByPhoneNumber(phone)
	if err != nil {
		return errors.New("error with server")
	}
	if user == nil {
		return errors.New("user with this phone is not exists")
	}
	helper.TwilioSetup(cfg.ACCOUNTSID, cfg.AUTHTOKEN)
	_, err = helper.TwilioSendOTP(phone, cfg.SERVICESSID)

	if err != nil {
		return errors.New("error occured while generating otp")
	}
	return nil
}
func (op *otpUseCase) VerifyOTP(code models.VerifyData) (models.TokenUser, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return models.TokenUser{}, err
	}
	helper.TwilioSetup(cfg.ACCOUNTSID, cfg.AUTHTOKEN)
	err = helper.TwilioVerifyOTP(cfg.SERVICESSID, code.Code, code.User.PhoneNumber)
	if err != nil {
		return models.TokenUser{}, errors.New("error while verifying")
	}
	userDetails, err := op.otpRepository.UserDetailsUsingPhone(code.User.PhoneNumber)
	if err != nil {
		return models.TokenUser{}, err
	}
	accessToken, err := helper.GenerateAccessToken(userDetails)
	if err != nil {
		return models.TokenUser{}, errors.New("couldn't create token due to some internal error")
	}
	refreshToken, err := helper.GenerateRefreshToken(userDetails)
	if err != nil {
		return models.TokenUser{}, errors.New("couldn't create token due to some internal error")
	}
	var user models.UserDetailsResponse
	err = copier.Copy(&user, &userDetails)
	if err != nil {
		return models.TokenUser{}, err
	}
	return models.TokenUser{
		Users:        user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
