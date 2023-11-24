package interfaces

import "Zhooze/pkg/utils/models"

type OtpUseCase interface {
	SendOtp(phone string) error
	VerifyOTP(code models.VerifyData) (models.TokenUser, error)
}
