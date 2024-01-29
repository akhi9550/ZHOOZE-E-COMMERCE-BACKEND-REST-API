package mockusecase

import (
	"Zhooze/pkg/utils/models"
	"reflect"

	"github.com/golang/mock/gomock"
)

type MockUserUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockUserUseCaseMockRecorder
}
type MockUserUseCaseMockRecorder struct {
	mock *MockUserUseCase
}

func NewMockUserUseCase(ctrl *gomock.Controller) *MockUserUseCase {
	mock := &MockUserUseCase{ctrl: ctrl}
	mock.recorder = &MockUserUseCaseMockRecorder{mock}
	return mock
}
func (m *MockUserUseCase) EXPECT() *MockUserUseCaseMockRecorder {
	return m.recorder
}
func (m *MockUserUseCase) UserSignUp(user models.UserSignUp, ref string) (models.TokenUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserSignUp", user, ref)
	ret0, _ := ret[0].(models.TokenUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserSignUp indicates an expected call of UserSignUp.
func (mr *MockUserUseCaseMockRecorder) UserSignUp(user, ref interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserSignUp", reflect.TypeOf((*MockUserUseCase)(nil).UserSignUp), user, ref)
}
