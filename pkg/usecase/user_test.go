package usecase

import (
	mockRepository "Zhooze/pkg/repository/mock"
	"Zhooze/pkg/utils/models"
	"errors"
	"testing"

	"github.com/go-playground/assert/v2"
	gomock "github.com/golang/mock/gomock"
)

func Test_AddAddress(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock implementations for the repositories
	userRepo := mockRepository.NewMockUserRepository(ctrl)
	orderRepo := mockRepository.NewMockOrderRepository(ctrl)

	userUseCase := NewUserUseCase(userRepo, orderRepo)

	testData := map[string]struct {
		input   models.AddressInfo
		stub    func(*mockRepository.MockUserRepository, *mockRepository.MockOrderRepository, models.AddressInfo)
		wantErr error
	}{
		"success": {
			input: models.AddressInfo{
				Name:      "akhil",
				HouseName: "chekkiyil house",
				Street:    "pallippuram",
				City:      "cherthala",
				State:     "kerala",
				Pin:       "688541",
			},
			stub: func(userRepo *mockRepository.MockUserRepository, orderRepo *mockRepository.MockOrderRepository, data models.AddressInfo) {
				userRepo.EXPECT().AddAddress(1, data).Return(nil).Times(1)
			},
			wantErr: nil,
		},
		"failure": {
			input: models.AddressInfo{
				Name:      "akhil",
				HouseName: "chekkiyil house",
				Street:    "pallippuram",
				City:      "cherthala",
				State:     "kerala",
				Pin:       "688541",
			},
			stub: func(userRepo *mockRepository.MockUserRepository, orderRepo *mockRepository.MockOrderRepository, data models.AddressInfo) {
				userRepo.EXPECT().AddAddress(1, data).Return(errors.New("could not add the address")).Times(1)
			},
			wantErr: errors.New("could not add the address"),
		},
	}

	for testName, test := range testData {
		t.Run(testName, func(t *testing.T) {
			test.stub(userRepo, orderRepo, test.input)
			err := userUseCase.AddAddress(1, test.input)
			assert.Equal(t, test.wantErr, err)
		})
	}
}
func Test_GetAllAddress(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mockRepository.NewMockUserRepository(ctrl)
	orderRepo := mockRepository.NewMockOrderRepository(ctrl)
	userUseCase := NewUserUseCase(userRepo, orderRepo)

	testData := map[string]struct {
		input   int
		stub    func(*mockRepository.MockUserRepository, *mockRepository.MockOrderRepository, int)
		want    models.AddressInfoResponse
		wantErr error
	}{
		"success": {
			input: 1,
			stub: func(userRepo *mockRepository.MockUserRepository, orderRepo *mockRepository.MockOrderRepository, data int) {
				userRepo.EXPECT().GetAllAddres(data).Times(1).Return(models.AddressInfoResponse{}, nil)
			},
			want:    models.AddressInfoResponse{},
			wantErr: nil,
		},
		"failed": {
			input: 1,
			stub: func(userRepo *mockRepository.MockUserRepository, orderRepo *mockRepository.MockOrderRepository, data int) {
				userRepo.EXPECT().GetAllAddres(data).Times(1).Return(models.AddressInfoResponse{}, errors.New("error"))
			},
			want:    models.AddressInfoResponse{},
			wantErr: errors.New("error"),
		},
	}

	for _, test := range testData {
		test.stub(userRepo, orderRepo, test.input)
		result, err := userUseCase.GetAllAddres(test.input)
		assert.Equal(t, test.want, result)
		assert.Equal(t, test.wantErr, err)
	}
}
