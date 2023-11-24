package interfaces

import (
	"Zhooze/pkg/domain"
	"Zhooze/pkg/utils/models"
)

type OtpRepository interface {
	FindUserByPhoneNumber(phone string) (*domain.User, error)
	UserDetailsUsingPhone(phone string) (models.UserDetailsResponse, error)
	FindUsersByEmail(email string) (bool, error)
	GetUserPhoneByEmail(email string) (string, error)
}