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
		input         interface{}
		buildStub     func(useCaseMock *mockUseCase.MockUserUseCase, signup interface{})
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
			buildStub: func(useCaseMock *mockUseCase.MockUserUseCase, signupData interface{}) {
				err := validator.New().Struct(signupData)
				if err != nil {
					fmt.Println("validation failed")
				}
				useCaseMock.EXPECT().UsersSignUp(signupData).Times(1).Return(models.TokenUser{
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
