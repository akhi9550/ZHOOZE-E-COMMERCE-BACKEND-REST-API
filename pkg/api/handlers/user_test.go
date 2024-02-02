package handlers

import (
	"Zhooze/pkg/mock/mockUseCase"
	"Zhooze/pkg/utils/models"
	"bytes"
	"encoding/json"
	"errors"
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
		"user couldnot sign up": {
			input: models.UserSignUp{
				Firstname:    "akhil",
				Lastname:     "c",
				Email:        "akhilc89@gmail.com",
				Password:     "908765",
				Phone:        "+919087675645",
				ReferralCode: "659823",
			},
			buildStub: func(useCaseMock *mockUseCase.MockUserUseCase, signupData models.UserSignUp) {
				// copying signupData to domain.user for pass to Mock usecase
				err := validator.New().Struct(signupData)
				if err != nil {
					fmt.Println("validation failed")
				}

				useCaseMock.EXPECT().UsersSignUp(signupData).Times(1).Return(&models.TokenUser{}, errors.New("cannot sign up"))
			},
			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)

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

func Test_LoginHandler(t *testing.T) {
	testCase := map[string]struct {
		input         models.LoginDetail
		buildStub     func(useCaseMock *mockUseCase.MockUserUseCase, login models.LoginDetail)
		checkResponse func(t *testing.T, responseRecorder *httptest.ResponseRecorder)
	}{
		"Success": {
			input: models.LoginDetail{
				Email:    "akhilc23@gmail.com",
				Password: "898989",
			},
			buildStub: func(useCaseMock *mockUseCase.MockUserUseCase, login models.LoginDetail) {
				err := validator.New().Struct(login)
				if err != nil {
					fmt.Println("validation failed")
				}
				useCaseMock.EXPECT().UsersLogin(login).Times(1).Return(&models.TokenUser{
					Users: models.UserDetailsResponse{
						Id:        1,
						Firstname: "akhil",
						Lastname:  "c",
						Email:     "akhilc23@gmail.com",
						Phone:     "+919856123585",
					},
					AccessToken:  "tyhddfgfh.djdudhfffdf.isjjrhfhfs",
					RefreshToken: "lpoiisjaf.pidfs9d8fsf.sddddffsff",
				}, nil)
			},
			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				if responseRecorder.Code != http.StatusOK && responseRecorder.Code != http.StatusCreated {
					t.Errorf("unexpected status code: %d", responseRecorder.Code)
				}
			},
		},
		"user couldn't login": {
			input: models.LoginDetail{
				Email:    "akhilc23@gmail.com",
				Password: "no password",
			},
			buildStub: func(useCaseMock *mockUseCase.MockUserUseCase, login models.LoginDetail) {
				err := validator.New().Struct(login)
				if err != nil {
					fmt.Println("validation failed")
				}
				useCaseMock.EXPECT().UsersLogin(login).Times(1).Return(&models.TokenUser{}, errors.New("cannot login up"))
			},
			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
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
			UserHandler := NewUserHandler(mockUseCase)
			server := gin.Default()
			server.POST("/login", UserHandler.Userlogin)
			jsonData, err := json.Marshal(test.input)
			assert.NoError(t, err)
			body := bytes.NewBuffer(jsonData)
			mockRequest, err := http.NewRequest(http.MethodPost, "/login", body)
			assert.NoError(t, err)
			responseRecorder := httptest.NewRecorder()
			server.ServeHTTP(responseRecorder, mockRequest)
			test.checkResponse(t, responseRecorder)
		})
	}
}
