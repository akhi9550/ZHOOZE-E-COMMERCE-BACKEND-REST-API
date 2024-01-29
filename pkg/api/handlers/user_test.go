package handlers

import (
	"Zhooze/pkg/mock/mockusecase"
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

// func Test_UserSignup(t *testing.T) {
// 	testCase := map[string]struct {
// 		input         interface{}
// 		buildStub     func(useCaseMock *mockusecase.MockUserUseCase, signup interface{})
// 		checkResponse func(t *testing.T, responseRecorder *httptest.ResponseRecorder)
// 	}{
// 		"Valid Signup": {
// 			input: models.UserSignUp{
// 				Firstname: "akhil",
// 				Lastname:  "c",
// 				Email:     "akhilc89@gmail.com",
// 				Password:  "908765",
// 				Phone:     "9087675645",
// 			},
// 			buildStub: func(useCaseMock *mockusecase.MockUserUseCase, signupData interface{}) {
// 				err := validator.New().Struct(signupData)
// 				if err != nil {
// 					fmt.Println("validation failed")
// 				}
// 				useCaseMock.EXPECT().UserSignUp(signupData, "").Times(1).Return(models.TokenUser{
// 					Users: models.UserDetailsResponse{
// 						Id:        1,
// 						Firstname: "akhil",
// 						Lastname:  "c",
// 						Email:     "akhil89@gmail.com",
// 						Phone:     "9087675645",
// 					},
// 					AccessToken:  "adfsaethjjshahfiurhfajherkuefeu",
// 					RefreshToken: "fkdgkerjrijigsiejggjrlisjgjisg",
// 				}, nil)
// 			},
// 			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
// 				assert.Equal(t, http.StatusCreated, responseRecorder.Code)
// 			},
// 		},
// 	}
// 	for testName, test := range testCase {
// 		test := test
// 		t.Run(testName, func(t *testing.T) {
// 			t.Parallel()
// 			ctrl := gomock.NewController(t)
// 			mockUseCase := mockusecase.NewMockUserUseCase(ctrl)
// 			test.buildStub(mockUseCase, test.input)
// 			userHandler := NewUserHandler(mockUseCase)

// 			server := gin.Default()
// 			server.POST("/signup", userHandler.UserSignup)

// 			jsonData, err := json.Marshal(test.input)
// 			assert.NoError(t, err)
// 			body := bytes.NewBuffer(jsonData)

// 			mockRequest, err := http.NewRequest(http.MethodPost, "/signup", body)
// 			assert.NoError(t, err)
// 			responseRecorder := httptest.NewRecorder()
// 			server.ServeHTTP(responseRecorder, mockRequest)

// 			test.checkResponse(t, responseRecorder)

// 		})
// 	}
// }

func TestUserSignup(t *testing.T) {

	testCase := map[string]struct {
		input         interface{}
		buildStub     func(useCaseMock *mockusecase.MockUserUseCase, signup interface{})
		checkResponse func(t *testing.T, responseRecorder *httptest.ResponseRecorder)
	}{
		"Valid Signup": {
			input: models.UserSignUp{
				Firstname: "akhil",
				Lastname:  "c",
				Email:     "akhilc89@gmail.com",
				Password:  "908765",
				Phone:     "9087675645",
			},
			buildStub: func(useCaseMock *mockusecase.MockUserUseCase, signupData interface{}) {
				// copying signupData to domain.user for pass to Mock usecase
				err := validator.New().Struct(signupData)
				if err != nil {
					fmt.Println("validation failed")
				}

				useCaseMock.EXPECT().UserSignUp(signupData, "").Times(1).Return(models.TokenUser{
					Users: models.UserDetailsResponse{
						Id:        1,
						Firstname: "akhil",
						Lastname:  "c",
						Email:     "akhil89@gmail.com",
						Phone:     "9087675645",
					},
					AccessToken:  "adfsaethjjshahfiurhfajherkuefeu",
					RefreshToken: "fkdgkerjrijigsiejggjrlisjgjisg",
				}, nil)
			},
			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, responseRecorder.Code)

			},
		},

		"fields provided in wrong format": {
			input: "",
			buildStub: func(useCaseMock *mockusecase.MockUserUseCase, signupData interface{}) {
				// copying signupData to domain.user for pass to Mock usecase
				err := validator.New().Struct(signupData)
				if err != nil {
					fmt.Println("validation failed")
				}
			},
			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)

			},
		},

		"struct validation fails": {
			input: models.SetNewName{},
			buildStub: func(useCaseMock *mockusecase.MockUserUseCase, signupData interface{}) {
				// copying signupData to domain.user for pass to Mock usecase
				err := validator.New().Struct(signupData)
				if err != nil {
					fmt.Println("validation failed")
				}
			},
			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)

			},
		},

		"user couldnot sign up": {
			input: models.UserSignUp{
				Firstname: "akhil",
				Lastname:  "c",
				Email:     "akhilc89@gmail.com",
				Password:  "908765",
				Phone:     "9087675645",
			},
			buildStub: func(useCaseMock *mockusecase.MockUserUseCase, signupData interface{}) {
				// copying signupData to domain.user for pass to Mock usecase
				err := validator.New().Struct(signupData)
				if err != nil {
					fmt.Println("validation failed")
				}

				useCaseMock.EXPECT().UserSignUp(signupData, "").Times(1).Return(models.TokenUser{}, errors.New("cannot sign up"))
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
			mockUseCase := mockusecase.NewMockUserUseCase(ctrl)
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
