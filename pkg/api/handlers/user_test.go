package handlers

import (
	"Zhooze/pkg/mock/mockUseCase"
	"Zhooze/pkg/utils/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_UserSignup(t *testing.T) {
	testCase := map[string]struct {
		input         models.UserSignUp
		buildStub     func(useCaseMock *mockUseCase.MockUserUseCase, signupData models.UserSignUp)
		checkResponse func(t *testing.T, responseRecorder *httptest.ResponseRecorder)
	}{
		"Valid Signup": {
			input: models.UserSignUp{
				Firstname:    "akhil",
				Lastname:     "c",
				Email:        "akhilc89@gmail.com",
				Password:     "908765",
				Phone:        "+919087675645",
				ReferralCode: "659823",
			},
			buildStub: func(useCaseMock *mockUseCase.MockUserUseCase, signupData models.UserSignUp) {
				err := validator.New().Struct(signupData)
				if err != nil {
					fmt.Println("validation failed")
				}
				useCaseMock.EXPECT().UsersSignUp(signupData).Times(1).Return(&models.TokenUser{
					Users: models.UserDetailsResponse{
						Id:        1,
						Firstname: "akhil",
						Lastname:  "c",
						Email:     "akhilc89@gmail.com",
						Phone:     "+919087675645",
					},
					AccessToken:  "adfsae.thjjshahfiurhf.ajherkuefeu",
					RefreshToken: "fkdgker.jrijigsiejggj.rlisjgjisg3",
				}, nil)
			},
			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, responseRecorder.Code)
			},
		},
	}
	for testName, test := range testCase {
		test := test
		t.Run(testName, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			mockUseCase := mockUseCase.NewMockUserUseCase(ctrl)
			test.buildStub(mockUseCase, test.input)
			userHandler := NewUserHandler(mockUseCase)

			server := gin.Default()
			server.POST("/signup", userHandler.UserSignup)

			jsonData, err := json.Marshal(test.input)
			assert.NoError(t, err)
			body := bytes.NewBuffer(jsonData)

			mockRequest, err := http.NewRequest(http.MethodPost, "/signup", body)
			assert.NoError(t, err)
			responseRecorder := httptest.NewRecorder()
			server.ServeHTTP(responseRecorder, mockRequest)

			test.checkResponse(t, responseRecorder)

		})
	}
}

// func Test_AddAddress(t *testing.T) {
// 	testCase := map[string]struct {
// 		buildStub     func(useCaseMock *mockUseCase.MockUserUseCase)
// 		checkResponse func(t *testing.T, responseRecorder *httptest.ResponseRecorder)
// 		expected      []domain.Address
// 	}{
// 		"successfull": {
// 			buildStub: func(useCaseMock *mockUseCase.MockUserUseCase) {

// 				useCaseMock.EXPECT().GetAllAddress(1).Times(1).Return([]domain.Address{}, nil)
// 			},
// 			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
// 				assert.Equal(t, http.StatusOK, responseRecorder.Code)

// 			},
// 		},
// 		"parameter problem": {
// 			buildStub: func(useCaseMock *mockUseCase.MockUserUseCase) {
// 			},
// 			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
// 				assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)

// 			},
// 		},
// 		"error retrieving records": {
// 			buildStub: func(useCaseMock *mockUseCase.MockUserUseCase) {

// 				useCaseMock.EXPECT().GetAllAddress(1).Times(1).Return(nil, errors.New("error retrieving records"))
// 			},
// 			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
// 				assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)

// 			},
// 		},
// 	}

// 	for testName, test := range testCase {
// 		testName := testName
// 		test := test
// 		t.Run(testName, func(t *testing.T) {
// 			t.Parallel()
// 			ctrl := gomock.NewController(t)
// 			mockUseCase := mockUseCase.NewMockUserUseCase(ctrl)
// 			test.buildStub(mockUseCase)

// 			userHandler := NewUserHandler(mockUseCase)

// 			server := gin.Default()
// 			server.POST("/getAddresses", userHandler.GetAllAddress)

// 			mockRequest, err := http.NewRequest(http.MethodGet, "id=1", nil)
// 			assert.NoError(t, err)
// 			if testName == "parameter problem" {
// 				mockRequest, err = http.NewRequest(http.MethodGet, "id=invalid", nil)
// 				assert.NoError(t, err)
// 			}
// 			responseRecorder := httptest.NewRecorder()

// 			server.ServeHTTP(responseRecorder, mockRequest)

// 			test.checkResponse(t, responseRecorder)

// 		})

// 	}
// }
